package anonymizer

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"ai-anonymizer/internal/catalog"
	"ai-anonymizer/internal/config"
)

type fakeSecretGetter struct {
	value string
	err   error
}

type mapSecretGetter struct {
	values map[string]string
	err    error
}

func (f fakeSecretGetter) GetSecret(partnerID string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.value, nil
}

func (m mapSecretGetter) GetSecret(partnerID string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.values[partnerID], nil
}

func TestResolveByID(t *testing.T) {
	items := []catalog.Item{
		{ID: "abc", Path: "in/a.txt", Status: catalog.StatusScanned},
		{ID: "def", Path: "in/b.txt", Status: catalog.StatusScanned},
	}

	idx, item, err := resolve(items, "def")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx != 1 || item.ID != "def" {
		t.Fatalf("unexpected match: idx=%d item=%+v", idx, item)
	}
}

func TestResolveFuzzyAmbiguous(t *testing.T) {
	items := []catalog.Item{
		{ID: "abc", Path: "in/invoice-april.txt", Status: catalog.StatusScanned},
		{ID: "def", Path: "in/invoice-may.txt", Status: catalog.StatusScanned},
	}

	_, _, err := resolve(items, "invoice")
	if err == nil {
		t.Fatal("expected ambiguous match error")
	}
	var ambiguous AmbiguousMatchError
	if !errors.As(err, &ambiguous) {
		t.Fatalf("expected AmbiguousMatchError, got %T", err)
	}
	if len(ambiguous.Candidates) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(ambiguous.Candidates))
	}
}

func TestAdaptResponseJSONVariants(t *testing.T) {
	got1, err := adaptResponse("application/json", []byte(`{"result":"ANON"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got1) != "ANON" {
		t.Fatalf("unexpected result %q", string(got1))
	}

	got2, err := adaptResponse("application/json", []byte(`{"result":{"content":"ANON2"}}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got2) != "ANON2" {
		t.Fatalf("unexpected result %q", string(got2))
	}
}

func TestExecuteRejectsNonHTTPS(t *testing.T) {
	cfg := baseConfig(t)
	cfg.APIBaseURL = "http://api.company.local"

	_, err := Execute(cfg, "x", Deps{
		SecretGetter: fakeSecretGetter{value: "token"},
	})
	if err == nil || !strings.Contains(err.Error(), "must use https") {
		t.Fatalf("expected https policy error, got %v", err)
	}
}

func TestExecuteRejectsHostOutsideAllowlist(t *testing.T) {
	cfg := baseConfig(t)
	cfg.APIAllowedHosts = []string{"allowed.company.local"}

	_, err := Execute(cfg, "x", Deps{
		SecretGetter: fakeSecretGetter{value: "token"},
	})
	if err == nil || !strings.Contains(err.Error(), "not in security.allowed_hosts") {
		t.Fatalf("expected allowlist policy error, got %v", err)
	}
}

func TestExecuteRejectsMissingSecretFallback(t *testing.T) {
	cfg := baseConfig(t)

	_, err := Execute(cfg, "x", Deps{
		SecretGetter: fakeSecretGetter{value: ""},
	})
	if err == nil || !strings.Contains(err.Error(), "was not found in OS secret store") {
		t.Fatalf("expected missing secret error, got %v", err)
	}
}

func TestExecuteRejectsMissingAuditHMACKey(t *testing.T) {
	cfg := baseConfig(t)
	cfg.AuditHMACKeyID = "audit-hmac-v1"

	_, err := Execute(cfg, "x", Deps{
		SecretGetter: mapSecretGetter{
			values: map[string]string{
				cfg.APIPartnerID: "api-secret",
			},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "audit hmac key was not found in OS secret store") {
		t.Fatalf("expected missing audit hmac key error, got %v", err)
	}
}

func TestExecuteRejectsWhenFileExceedsMaxSize(t *testing.T) {
	cfg := baseConfig(t)
	cfg.MaxFileSizeMB = 1

	if err := os.MkdirAll(filepath.Join(cfg.RawPath, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	largeRel := "in/large.bin"
	largePath := filepath.Join(cfg.RawPath, filepath.FromSlash(largeRel))
	if err := os.WriteFile(largePath, make([]byte, 2*1024*1024), 0o644); err != nil {
		t.Fatal(err)
	}
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save([]catalog.Item{
		{ID: catalog.BuildID(largeRel), Path: largeRel, Status: catalog.StatusScanned},
	}); err != nil {
		t.Fatal(err)
	}

	_, err := Execute(cfg, catalog.BuildID(largeRel), Deps{
		SecretGetter: fakeSecretGetter{value: "token"},
	})
	if err == nil || !strings.Contains(err.Error(), "exceeds max_file_size_mb") {
		t.Fatalf("expected size-limit error, got %v", err)
	}
}

func TestExecuteRejectsWhenPagesExceedLimitForText(t *testing.T) {
	cfg := baseConfig(t)
	cfg.MaxPages = 2

	if err := os.MkdirAll(filepath.Join(cfg.RawPath, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	rel := "in/multipage.txt"
	filePath := filepath.Join(cfg.RawPath, filepath.FromSlash(rel))
	if err := os.WriteFile(filePath, []byte("p1\f p2\f p3"), 0o644); err != nil {
		t.Fatal(err)
	}
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save([]catalog.Item{
		{ID: catalog.BuildID(rel), Path: rel, Status: catalog.StatusScanned},
	}); err != nil {
		t.Fatal(err)
	}

	_, err := Execute(cfg, catalog.BuildID(rel), Deps{
		SecretGetter: fakeSecretGetter{value: "token"},
	})
	if err == nil || !strings.Contains(err.Error(), "exceeds max_pages") {
		t.Fatalf("expected max_pages error, got %v", err)
	}
}

func TestExecuteRejectsWhenPagesExceedLimitForPDF(t *testing.T) {
	cfg := baseConfig(t)
	cfg.MaxPages = 2

	if err := os.MkdirAll(filepath.Join(cfg.RawPath, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	rel := "in/multipage.pdf"
	filePath := filepath.Join(cfg.RawPath, filepath.FromSlash(rel))
	if err := os.WriteFile(filePath, buildFakePDFWithPages(3), 0o644); err != nil {
		t.Fatal(err)
	}
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save([]catalog.Item{
		{ID: catalog.BuildID(rel), Path: rel, Status: catalog.StatusScanned},
	}); err != nil {
		t.Fatal(err)
	}

	_, err := Execute(cfg, catalog.BuildID(rel), Deps{
		SecretGetter: fakeSecretGetter{value: "token"},
	})
	if err == nil || !strings.Contains(err.Error(), "exceeds max_pages") {
		t.Fatalf("expected max_pages error for pdf, got %v", err)
	}
}

func TestExecuteRejectsWhenPagesExceedLimitForDOCX(t *testing.T) {
	cfg := baseConfig(t)
	cfg.MaxPages = 1

	if err := os.MkdirAll(filepath.Join(cfg.RawPath, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	rel := "in/multipage.docx"
	filePath := filepath.Join(cfg.RawPath, filepath.FromSlash(rel))
	if err := writeFakeDOCXWithPages(filePath, 3); err != nil {
		t.Fatal(err)
	}
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save([]catalog.Item{
		{ID: catalog.BuildID(rel), Path: rel, Status: catalog.StatusScanned},
	}); err != nil {
		t.Fatal(err)
	}

	_, err := Execute(cfg, catalog.BuildID(rel), Deps{
		SecretGetter: fakeSecretGetter{value: "token"},
	})
	if err == nil || !strings.Contains(err.Error(), "exceeds max_pages") {
		t.Fatalf("expected max_pages error for docx, got %v", err)
	}
}

func TestExecuteRejectsWindowsADSPath(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only ADS policy")
	}

	cfg := baseConfig(t)
	rel := "in/sample.txt:secret"
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save([]catalog.Item{
		{ID: catalog.BuildID(rel), Path: rel, Status: catalog.StatusScanned},
	}); err != nil {
		t.Fatal(err)
	}

	_, err := Execute(cfg, catalog.BuildID(rel), Deps{
		SecretGetter: fakeSecretGetter{value: "token"},
	})
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "ads") {
		t.Fatalf("expected ADS policy error, got %v", err)
	}
}

func TestExecuteRejectsWindowsReparsePointInPath(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only reparse policy")
	}

	cfg := baseConfig(t)
	if err := os.MkdirAll(filepath.Join(cfg.RawPath, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	rel := "in/sample.txt"
	filePath := filepath.Join(cfg.RawPath, filepath.FromSlash(rel))
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save([]catalog.Item{
		{ID: catalog.BuildID(rel), Path: rel, Status: catalog.StatusScanned},
	}); err != nil {
		t.Fatal(err)
	}

	prevDetector := reparsePointCheck
	reparsePointCheck = func(path string) (bool, error) {
		if strings.EqualFold(path, filepath.Join(cfg.RawPath, "in")) {
			return true, nil
		}
		return false, nil
	}
	defer func() { reparsePointCheck = prevDetector }()

	_, err := Execute(cfg, catalog.BuildID(rel), Deps{
		SecretGetter: fakeSecretGetter{value: "token"},
	})
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "reparse") {
		t.Fatalf("expected reparse policy error, got %v", err)
	}
}

func TestExecuteWritesAuditEntryWithHMACSafeFields(t *testing.T) {
	cfg := baseConfig(t)
	cfg.AuditRetentionDays = 30
	cfg.AuditHMACKeyID = "audit-hmac-v1"

	if err := os.MkdirAll(filepath.Join(cfg.RawPath, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	rel := "in/sample.txt"
	filePath := filepath.Join(cfg.RawPath, filepath.FromSlash(rel))
	if err := os.WriteFile(filePath, []byte("PII"), 0o644); err != nil {
		t.Fatal(err)
	}
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save([]catalog.Item{
		{ID: catalog.BuildID(rel), Path: rel, Status: catalog.StatusScanned},
	}); err != nil {
		t.Fatal(err)
	}

	serverURL, serverClient, cleanup := newSuccessAPIClient([]byte(`{"result":"ANON"}`), "application/json")
	defer cleanup()
	cfg.APIBaseURL = serverURL
	cfg.APIAllowedHosts = []string{"127.0.0.1", "localhost"}

	result, err := Execute(cfg, catalog.BuildID(rel), Deps{
		HTTPClient: serverClient,
		SecretGetter: mapSecretGetter{
			values: map[string]string{
				cfg.APIPartnerID:   "api-secret",
				cfg.AuditHMACKeyID: "audit-secret",
			},
		},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if result.Item.Status != catalog.StatusSucceeded {
		t.Fatalf("expected succeeded status, got %q", result.Item.Status)
	}

	data, err := os.ReadFile(cfg.AuditPath)
	if err != nil {
		t.Fatalf("read audit log failed: %v", err)
	}
	line := strings.TrimSpace(string(data))
	if line == "" {
		t.Fatal("expected non-empty audit line")
	}
	if strings.Contains(line, rel) || strings.Contains(line, filePath) || strings.Contains(strings.ToLower(line), "sample.txt") {
		t.Fatalf("audit line must not contain raw path/name, got %q", line)
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("invalid audit json: %v", err)
	}
	if payload["event"] != "run" {
		t.Fatalf("expected event=run, got %#v", payload["event"])
	}
	if payload["status"] != string(catalog.StatusSucceeded) {
		t.Fatalf("expected status=succeeded, got %#v", payload["status"])
	}
	fileID, _ := payload["file_id"].(string)
	if fileID == "" {
		t.Fatal("expected non-empty file_id")
	}
	if fileID == catalog.BuildID(rel) {
		t.Fatalf("expected HMAC file_id, got raw id %q", fileID)
	}
}

func TestExecuteMVPUsesEnvSecretsWithoutSecretGetter(t *testing.T) {
	cfg := baseConfig(t)
	cfg.RuntimeMode = "mvp"
	cfg.APIAuthToken = "session-token"
	cfg.AuditHMACSecret = "audit-secret"
	cfg.APIPartnerID = "70bd3a91-0000-0000-0000-000000000000"

	if err := os.MkdirAll(filepath.Join(cfg.RawPath, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	rel := "in/sample.txt"
	filePath := filepath.Join(cfg.RawPath, filepath.FromSlash(rel))
	if err := os.WriteFile(filePath, []byte("PII"), 0o644); err != nil {
		t.Fatal(err)
	}
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save([]catalog.Item{
		{ID: catalog.BuildID(rel), Path: rel, Status: catalog.StatusScanned},
	}); err != nil {
		t.Fatal(err)
	}

	serverURL, serverClient, cleanup := newSuccessAPIClient([]byte(`{"result":"ANON"}`), "application/json")
	defer cleanup()
	cfg.APIBaseURL = serverURL
	cfg.APIAllowedHosts = []string{"127.0.0.1", "localhost"}

	result, err := Execute(cfg, catalog.BuildID(rel), Deps{
		HTTPClient: serverClient,
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if result.Item.Status != catalog.StatusSucceeded {
		t.Fatalf("expected succeeded status, got %q", result.Item.Status)
	}
}

func TestExecuteMVPFailsWithoutAuditSecret(t *testing.T) {
	cfg := baseConfig(t)
	cfg.RuntimeMode = "mvp"
	cfg.APIAuthToken = "session-token"
	cfg.AuditHMACSecret = ""

	_, err := Execute(cfg, "x", Deps{})
	if err == nil || !strings.Contains(err.Error(), "audit.hmac_secret") {
		t.Fatalf("expected missing audit secret error, got %v", err)
	}
}

func TestExecuteTrimsAuthorizationHeaderValue(t *testing.T) {
	cfg := baseConfig(t)
	cfg.AuditHMACKeyID = "audit-hmac-v1"

	if err := os.MkdirAll(filepath.Join(cfg.RawPath, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	rel := "in/sample.txt"
	filePath := filepath.Join(cfg.RawPath, filepath.FromSlash(rel))
	if err := os.WriteFile(filePath, []byte("PII"), 0o644); err != nil {
		t.Fatal(err)
	}
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save([]catalog.Item{
		{ID: catalog.BuildID(rel), Path: rel, Status: catalog.StatusScanned},
	}); err != nil {
		t.Fatal(err)
	}

	var gotAuth string
	serverURL, serverClient, cleanup := newHeaderCaptureAPIClient(&gotAuth)
	defer cleanup()
	cfg.APIBaseURL = serverURL
	cfg.APIAllowedHosts = []string{"127.0.0.1", "localhost"}

	_, err := Execute(cfg, catalog.BuildID(rel), Deps{
		HTTPClient: serverClient,
		SecretGetter: mapSecretGetter{
			values: map[string]string{
				cfg.APIPartnerID:   "\r\napi-secret\r\n",
				cfg.AuditHMACKeyID: "audit-secret",
			},
		},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if gotAuth != "api-secret" {
		t.Fatalf("unexpected Authorization header value %q", gotAuth)
	}
}

func TestExecuteRemovesControlCharsFromAuthorizationHeaderValue(t *testing.T) {
	cfg := baseConfig(t)
	cfg.AuditHMACKeyID = "audit-hmac-v1"

	if err := os.MkdirAll(filepath.Join(cfg.RawPath, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	rel := "in/sample.txt"
	filePath := filepath.Join(cfg.RawPath, filepath.FromSlash(rel))
	if err := os.WriteFile(filePath, []byte("PII"), 0o644); err != nil {
		t.Fatal(err)
	}
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save([]catalog.Item{
		{ID: catalog.BuildID(rel), Path: rel, Status: catalog.StatusScanned},
	}); err != nil {
		t.Fatal(err)
	}

	var gotAuth string
	serverURL, serverClient, cleanup := newHeaderCaptureAPIClient(&gotAuth)
	defer cleanup()
	cfg.APIBaseURL = serverURL
	cfg.APIAllowedHosts = []string{"127.0.0.1", "localhost"}

	_, err := Execute(cfg, catalog.BuildID(rel), Deps{
		HTTPClient: serverClient,
		SecretGetter: mapSecretGetter{
			values: map[string]string{
				cfg.APIPartnerID:   "\x00api-\rse\ncret\x7f",
				cfg.AuditHMACKeyID: "audit-secret",
			},
		},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if gotAuth != "api-secret" {
		t.Fatalf("unexpected Authorization header value %q", gotAuth)
	}
}

func TestExecuteWritesHTTPDebugLogWhenEnabled(t *testing.T) {
	t.Setenv("ANON_HTTP_DEBUG", "true")
	cfg := baseConfig(t)
	cfg.AuditHMACKeyID = "audit-hmac-v1"

	if err := os.MkdirAll(filepath.Join(cfg.RawPath, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	rel := "in/sample.txt"
	filePath := filepath.Join(cfg.RawPath, filepath.FromSlash(rel))
	if err := os.WriteFile(filePath, []byte("PII"), 0o644); err != nil {
		t.Fatal(err)
	}
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save([]catalog.Item{
		{ID: catalog.BuildID(rel), Path: rel, Status: catalog.StatusScanned},
	}); err != nil {
		t.Fatal(err)
	}

	serverURL, serverClient, cleanup := newSuccessAPIClient([]byte(`{"result":"ANON"}`), "application/json")
	defer cleanup()
	cfg.APIBaseURL = serverURL
	cfg.APIAllowedHosts = []string{"127.0.0.1", "localhost"}

	_, err := Execute(cfg, catalog.BuildID(rel), Deps{
		HTTPClient: serverClient,
		SecretGetter: mapSecretGetter{
			values: map[string]string{
				cfg.APIPartnerID:   "very-secret-token",
				cfg.AuditHMACKeyID: "audit-secret",
			},
		},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	logPath := filepath.Join(cfg.WorkspacePath, ".anonym", "http-debug.log")
	data, readErr := os.ReadFile(logPath)
	if readErr != nil {
		t.Fatalf("expected debug log file, got %v", readErr)
	}
	logText := string(data)
	if !strings.Contains(logText, "scope=run") || !strings.Contains(logText, "status=200") {
		t.Fatalf("unexpected debug log content: %q", logText)
	}
	if strings.Contains(logText, "very-secret-token") {
		t.Fatalf("authorization must be masked in debug log: %q", logText)
	}
}

func baseConfig(t *testing.T) config.Config {
	t.Helper()

	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "out")
	if err := os.MkdirAll(raw, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}

	return config.Config{
		RawPath:            raw,
		OutputPath:         output,
		WorkspacePath:      workspace,
		CatalogPath:        filepath.Join(workspace, ".anonym", "catalog.json"),
		AuditPath:          filepath.Join(workspace, ".anonym", "audit.log"),
		APIBaseURL:         "https://api.company.local",
		APIPartnerID:       "70bd3a91-0000-0000-0000-000000000000",
		APIAllowedHosts:    []string{"api.company.local"},
		RequestTimeoutSec:  5,
		MaxFileSizeMB:      25,
		MaxPages:           300,
		MaxParallelRuns:    1,
		AuditRetentionDays: 30,
		AuditHMACKeyID:     "audit-hmac-v1",
	}
}

func newSuccessAPIClient(responseBody []byte, contentType string) (string, *http.Client, func()) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/tasks/file_anonymization" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(responseBody)
	}))
	return server.URL, server.Client(), server.Close
}

func newHeaderCaptureAPIClient(authHeader *string) (string, *http.Client, func()) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/tasks/file_anonymization" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		*authHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":"ANON"}`))
	}))
	return server.URL, server.Client(), server.Close
}

func buildFakePDFWithPages(pages int) []byte {
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	for i := 0; i < pages; i++ {
		buf.WriteString("<< /Type /Page >>\n")
	}
	buf.WriteString("%%EOF\n")
	return buf.Bytes()
}

func writeFakeDOCXWithPages(path string, pages int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := zip.NewWriter(file)
	app, err := writer.Create("docProps/app.xml")
	if err != nil {
		return err
	}
	_, err = app.Write([]byte("<Properties><Pages>" + strconv.Itoa(pages) + "</Pages></Properties>"))
	if err != nil {
		return err
	}

	doc, err := writer.Create("word/document.xml")
	if err != nil {
		return err
	}
	_, err = doc.Write([]byte("<w:document><w:body><w:p>content</w:p></w:body></w:document>"))
	if err != nil {
		return err
	}
	return writer.Close()
}
