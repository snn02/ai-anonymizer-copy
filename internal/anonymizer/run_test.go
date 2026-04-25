package anonymizer

import (
	"archive/zip"
	"bytes"
	"errors"
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

func (f fakeSecretGetter) GetSecret(partnerID string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.value, nil
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
		RawPath:           raw,
		OutputPath:        output,
		WorkspacePath:     workspace,
		CatalogPath:       filepath.Join(workspace, ".anonym", "catalog.json"),
		APIBaseURL:        "https://api.company.local",
		APIPartnerID:      "partner-1",
		APIAllowedHosts:   []string{"api.company.local"},
		RequestTimeoutSec: 5,
		MaxFileSizeMB:     25,
		MaxPages:          300,
		MaxParallelRuns:   1,
	}
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
