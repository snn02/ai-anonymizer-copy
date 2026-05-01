package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-anonymizer/internal/anonymizer"
	"ai-anonymizer/internal/catalog"
	"ai-anonymizer/internal/preflight"
)

const testPartnerUUID = "70bd3a91-0000-0000-0000-000000000000"

type alwaysTrueSecretChecker struct{}

func (alwaysTrueSecretChecker) GetSecret(partnerID string) (string, error) {
	return "token", nil
}

type staticSecretGetter struct {
	value string
}

func (s staticSecretGetter) GetSecret(partnerID string) (string, error) {
	return s.value, nil
}

type failingSecretGetter struct{}

func (failingSecretGetter) GetSecret(partnerID string) (string, error) {
	return "", os.ErrNotExist
}

func TestRunDoctorSuccess(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tasks/anonymization_fields" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	prevDeps := doctorDeps
	doctorDeps = preflight.DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: alwaysTrueSecretChecker{},
	}
	defer func() { doctorDeps = prevDeps }()

	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "anonymized")
	if err := os.MkdirAll(raw, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	cfgPath := writeRuntimeConfig(t, runtimeConfigValues{
		RawPath:       raw,
		OutputPath:    output,
		WorkspacePath: workspace,
		BaseURL:       server.URL,
		PartnerID:     testPartnerUUID,
		AllowedHosts:  []string{"127.0.0.1", "localhost"},
	})

	var out bytes.Buffer
	err := run([]string{
		"doctor",
		"--config", cfgPath,
	}, &out, &out)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(out.String(), "doctor: ok") {
		t.Fatalf("expected success output, got %q", out.String())
	}
}

func TestRunScanAndListFlow(t *testing.T) {
	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "anonymized")

	if err := os.MkdirAll(filepath.Join(raw, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(raw, "in", "sample.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfgPath := writeRuntimeConfig(t, runtimeConfigValues{
		RawPath:       raw,
		OutputPath:    output,
		WorkspacePath: workspace,
		BaseURL:       "https://api.company.local",
		PartnerID:     testPartnerUUID,
		AllowedHosts:  []string{"api.company.local"},
	})

	var scanOut bytes.Buffer
	scanErr := run([]string{
		"scan",
		"--config", cfgPath,
	}, &scanOut, &scanOut)
	if scanErr != nil {
		t.Fatalf("scan failed: %v", scanErr)
	}
	if !strings.Contains(scanOut.String(), "scan: indexed 1 file(s)") {
		t.Fatalf("unexpected scan output: %q", scanOut.String())
	}

	var listOut bytes.Buffer
	listErr := run([]string{
		"list",
		"--config", cfgPath,
	}, &listOut, &listOut)
	if listErr != nil {
		t.Fatalf("list failed: %v", listErr)
	}
	if !strings.Contains(listOut.String(), "in/sample.txt\tscanned") {
		t.Fatalf("unexpected list output: %q", listOut.String())
	}
}

func TestRunExecutesExplicitSendByID(t *testing.T) {
	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "anonymized")

	if err := os.MkdirAll(filepath.Join(raw, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}

	inputPath := filepath.Join(raw, "in", "sample.txt")
	if err := os.WriteFile(inputPath, []byte("PII"), 0o644); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/tasks/file_anonymization" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if got := r.Header.Get("Authorization"); got == "" {
			t.Fatalf("expected Authorization header")
		}
		if got := r.Header.Get("partner-id"); got != testPartnerUUID {
			t.Fatalf("unexpected partner-id header: %q", got)
		}
		if got := r.Header.Get("fields"); got == "" {
			t.Fatalf("expected fields header")
		}
		if got := r.Header.Get("user-id"); got != "70bd3a91-8888-8888-8888-888888888888" {
			t.Fatalf("unexpected user-id header: %q", got)
		}

		if err := r.ParseMultipartForm(2 << 20); err != nil {
			t.Fatalf("parse multipart failed: %v", err)
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("expected file in multipart: %v", err)
		}
		defer file.Close()
		body, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("read multipart file failed: %v", err)
		}
		if string(body) != "PII" {
			t.Fatalf("unexpected uploaded content: %q", string(body))
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": "ANON",
		})
	}))
	defer server.Close()
	prevRunDeps := runDeps
	runDeps = runDepsForTest(server.Client())
	defer func() { runDeps = prevRunDeps }()
	cfgPath := writeRuntimeConfig(t, runtimeConfigValues{
		RawPath:       raw,
		OutputPath:    output,
		WorkspacePath: workspace,
		BaseURL:       server.URL,
		PartnerID:     testPartnerUUID,
		UserID:        "70bd3a91-8888-8888-8888-888888888888",
		AllowedHosts:  []string{"127.0.0.1", "localhost"},
	})

	var scanOut bytes.Buffer
	if err := run([]string{
		"scan",
		"--config", cfgPath,
	}, &scanOut, &scanOut); err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	var runOut bytes.Buffer
	err := run([]string{
		"run",
		catalog.BuildID("in/sample.txt"),
		"--config", cfgPath,
	}, &runOut, &runOut)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !strings.Contains(runOut.String(), "run: succeeded") {
		t.Fatalf("unexpected run output: %q", runOut.String())
	}

	files, err := os.ReadDir(output)
	if err != nil {
		t.Fatalf("read output dir failed: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected one output file, got %d", len(files))
	}
	if strings.Contains(files[0].Name(), "sample") {
		t.Fatalf("output filename should not contain raw name, got %q", files[0].Name())
	}
	if !strings.Contains(files[0].Name(), catalog.BuildID("in/sample.txt")) {
		t.Fatalf("output filename should contain technical id, got %q", files[0].Name())
	}
}

func TestRunFuzzyAmbiguousRequiresExplicitID(t *testing.T) {
	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "anonymized")
	if err := os.MkdirAll(filepath.Join(raw, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"invoice-april.txt", "invoice-may.txt"} {
		if err := os.WriteFile(filepath.Join(raw, "in", name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()
	prevRunDeps := runDeps
	runDeps = runDepsForTest(server.Client())
	defer func() { runDeps = prevRunDeps }()
	cfgPath := writeRuntimeConfig(t, runtimeConfigValues{
		RawPath:       raw,
		OutputPath:    output,
		WorkspacePath: workspace,
		BaseURL:       server.URL,
		PartnerID:     testPartnerUUID,
		AllowedHosts:  []string{"127.0.0.1", "localhost"},
	})

	var scanOut bytes.Buffer
	if err := run([]string{
		"scan",
		"--config", cfgPath,
	}, &scanOut, &scanOut); err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	var runOut bytes.Buffer
	err := run([]string{
		"run",
		"invoice",
		"--config", cfgPath,
	}, &runOut, &runOut)
	if err == nil {
		t.Fatal("expected ambiguity error")
	}
	if !strings.Contains(err.Error(), "specify explicit id") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(runOut.String(), "multiple matches") {
		t.Fatalf("expected candidates list in output, got %q", runOut.String())
	}
}

func TestRunMarksCatalogFailedWhenAPIFails(t *testing.T) {
	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "anonymized")

	if err := os.MkdirAll(filepath.Join(raw, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(raw, "in", "sample.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("upstream error"))
	}))
	defer server.Close()
	prevRunDeps := runDeps
	runDeps = runDepsForTest(server.Client())
	defer func() { runDeps = prevRunDeps }()
	cfgPath := writeRuntimeConfig(t, runtimeConfigValues{
		RawPath:       raw,
		OutputPath:    output,
		WorkspacePath: workspace,
		BaseURL:       server.URL,
		PartnerID:     testPartnerUUID,
		AllowedHosts:  []string{"127.0.0.1", "localhost"},
	})

	var scanOut bytes.Buffer
	if err := run([]string{
		"scan",
		"--config", cfgPath,
	}, &scanOut, &scanOut); err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	var runOut bytes.Buffer
	runErr := run([]string{
		"run",
		catalog.BuildID("in/sample.txt"),
		"--config", cfgPath,
	}, &runOut, &runOut)
	if runErr == nil {
		t.Fatal("expected run failure")
	}

	var listOut bytes.Buffer
	if err := run([]string{
		"list",
		"--config", cfgPath,
	}, &listOut, &listOut); err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(listOut.String(), "\tfailed") {
		t.Fatalf("expected failed status in catalog, got %q", listOut.String())
	}
}

func TestRunDoctorWithConfigOnly(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tasks/anonymization_fields" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	prevDeps := doctorDeps
	doctorDeps = preflight.DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: alwaysTrueSecretChecker{},
	}
	defer func() { doctorDeps = prevDeps }()

	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "anonymized")
	if err := os.MkdirAll(raw, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}

	cfgPath := writeRuntimeConfig(t, runtimeConfigValues{
		RawPath:       raw,
		OutputPath:    output,
		WorkspacePath: workspace,
		BaseURL:       server.URL,
		PartnerID:     testPartnerUUID,
		AllowedHosts:  []string{"127.0.0.1", "localhost"},
	})

	var out bytes.Buffer
	err := run([]string{
		"doctor",
		"--config", cfgPath,
	}, &out, &out)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(out.String(), "doctor: ok") {
		t.Fatalf("expected success output, got %q", out.String())
	}
}

func TestRunCompactFlowWithConfigOnly(t *testing.T) {
	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "anonymized")
	if err := os.MkdirAll(filepath.Join(raw, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(raw, "in", "sample.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/tasks/file_anonymization":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"result": "ANON"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	prevRunDeps := runDeps
	runDeps = runDepsForTest(server.Client())
	defer func() { runDeps = prevRunDeps }()

	cfgPath := writeRuntimeConfig(t, runtimeConfigValues{
		RawPath:       raw,
		OutputPath:    output,
		WorkspacePath: workspace,
		BaseURL:       server.URL,
		PartnerID:     testPartnerUUID,
		AllowedHosts:  []string{"127.0.0.1", "localhost"},
	})

	var scanOut bytes.Buffer
	if err := run([]string{"scan", "--config", cfgPath}, &scanOut, &scanOut); err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if !strings.Contains(scanOut.String(), "scan: indexed 1 file(s)") {
		t.Fatalf("unexpected scan output: %q", scanOut.String())
	}

	var listOut bytes.Buffer
	if err := run([]string{"list", "--config", cfgPath}, &listOut, &listOut); err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(listOut.String(), "in/sample.txt\tscanned") {
		t.Fatalf("unexpected list output: %q", listOut.String())
	}

	var runOut bytes.Buffer
	err := run([]string{
		"run",
		catalog.BuildID("in/sample.txt"),
		"--config", cfgPath,
	}, &runOut, &runOut)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !strings.Contains(runOut.String(), "run: succeeded") {
		t.Fatalf("unexpected run output: %q", runOut.String())
	}
}

func TestRunDoctorFailsOnMissingExplicitConfig(t *testing.T) {
	missingCfg := filepath.Join(t.TempDir(), "missing.yaml")
	var out bytes.Buffer
	err := run([]string{
		"doctor",
		"--config", missingCfg,
	}, &out, &out)
	if err == nil {
		t.Fatal("expected error for missing explicit config")
	}
	if !strings.Contains(err.Error(), "config: file") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDoctorConfigEnvCLIConflictUsesCLIMode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tasks/anonymization_fields" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	prevDeps := doctorDeps
	doctorDeps = preflight.DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: alwaysTrueSecretChecker{},
	}
	defer func() { doctorDeps = prevDeps }()

	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "anonymized")
	if err := os.MkdirAll(raw, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}

	cfgPath := writeRuntimeConfig(t, runtimeConfigValues{
		RawPath:       raw,
		OutputPath:    output,
		WorkspacePath: workspace,
		BaseURL:       server.URL,
		PartnerID:     testPartnerUUID,
		AllowedHosts:  []string{"127.0.0.1", "localhost"},
	})

	t.Setenv("ANON_MODE", "prod")
	t.Setenv("ANON_API_AUTH_TOKEN", "env-token")
	t.Setenv("ANON_AUDIT_HMAC_SECRET", "env-audit-secret")

	var out bytes.Buffer
	err := run([]string{
		"doctor",
		"--config", cfgPath,
		"--mode", "mvp",
	}, &out, &out)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(out.String(), "doctor: ok") {
		t.Fatalf("expected success output, got %q", out.String())
	}
}

func TestRunDoctorPrintsWarningForMVPMode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tasks/anonymization_fields" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	prevDeps := doctorDeps
	doctorDeps = preflight.DoctorDeps{
		HTTPClient: server.Client(),
	}
	defer func() { doctorDeps = prevDeps }()

	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "anonymized")
	if err := os.MkdirAll(raw, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}

	cfgPath := writeRuntimeConfigV11WS(t, runtimeConfigV11WSValues{
		RawPath:        raw,
		OutputPath:     output,
		WorkspacePath:  workspace,
		BaseURL:        server.URL,
		PartnerID:      testPartnerUUID,
		AllowedHosts:   []string{"127.0.0.1", "localhost"},
		RuntimeMode:    "mvp",
		APIAuthToken:   "cfg-token",
		AuditHMACSecret: "cfg-audit-secret",
	})
	var out bytes.Buffer
	err := run([]string{"doctor", "--config", cfgPath}, &out, &out)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(out.String(), "warning: mvp insecure mode is active") {
		t.Fatalf("expected insecure profile warning, got %q", out.String())
	}
}

func TestRunScanAndListPrintWarningForMVPMode(t *testing.T) {
	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "anonymized")

	if err := os.MkdirAll(filepath.Join(raw, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(raw, "in", "sample.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfgPath := writeRuntimeConfigV11WS(t, runtimeConfigV11WSValues{
		RawPath:         raw,
		OutputPath:      output,
		WorkspacePath:   workspace,
		BaseURL:         "https://api.company.local",
		PartnerID:       testPartnerUUID,
		AllowedHosts:    []string{"api.company.local"},
		RuntimeMode:     "mvp",
		APIAuthToken:    "cfg-token",
		AuditHMACSecret: "cfg-audit-secret",
	})
	var scanOut bytes.Buffer
	if err := run([]string{"scan", "--config", cfgPath}, &scanOut, &scanOut); err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if !strings.Contains(scanOut.String(), "warning: mvp insecure mode is active") {
		t.Fatalf("expected insecure profile warning in scan, got %q", scanOut.String())
	}

	var listOut bytes.Buffer
	if err := run([]string{"list", "--config", cfgPath}, &listOut, &listOut); err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(listOut.String(), "warning: mvp insecure mode is active") {
		t.Fatalf("expected insecure profile warning in list, got %q", listOut.String())
	}
}

func TestRunMVPUsesEnvOverConfigWithoutSecretStore(t *testing.T) {
	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "anonymized")
	if err := os.MkdirAll(filepath.Join(raw, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(raw, "in", "sample.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	var gotAuth string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/tasks/file_anonymization":
			gotAuth = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"result": "ANON"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	prevRunDeps := runDeps
	runDeps = anonymizer.Deps{
		HTTPClient:   server.Client(),
		SecretGetter: failingSecretGetter{},
	}
	defer func() { runDeps = prevRunDeps }()

	cfgPath := writeRuntimeConfigV11WS(t, runtimeConfigV11WSValues{
		RawPath:         raw,
		OutputPath:      output,
		WorkspacePath:   workspace,
		BaseURL:         server.URL,
		PartnerID:       testPartnerUUID,
		AllowedHosts:    []string{"127.0.0.1", "localhost"},
		RuntimeMode:     "mvp",
		APIAuthToken:    "cfg-token",
		AuditHMACSecret: "cfg-audit-secret",
	})

	t.Setenv("ANON_MODE", "mvp")
	t.Setenv("ANON_API_AUTH_TOKEN", "env-token")
	t.Setenv("ANON_AUDIT_HMAC_SECRET", "env-audit-secret")

	var scanOut bytes.Buffer
	if err := run([]string{"scan", "--config", cfgPath}, &scanOut, &scanOut); err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	var runOut bytes.Buffer
	err := run([]string{
		"run",
		catalog.BuildID("in/sample.txt"),
		"--config", cfgPath,
	}, &runOut, &runOut)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !strings.Contains(runOut.String(), "run: succeeded") {
		t.Fatalf("unexpected run output: %q", runOut.String())
	}
	if gotAuth != "env-token" {
		t.Fatalf("expected Authorization from env token, got %q", gotAuth)
	}
}

func runDepsForTest(client *http.Client) anonymizer.Deps {
	return anonymizer.Deps{
		HTTPClient:   client,
		SecretGetter: staticSecretGetter{value: "token"},
	}
}

type runtimeConfigValues struct {
	RawPath       string
	OutputPath    string
	WorkspacePath string
	BaseURL       string
	PartnerID     string
	UserID        string
	AllowedHosts  []string
}

func writeRuntimeConfig(t *testing.T, values runtimeConfigValues) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	lines := []string{
		"paths:",
		"  raw_path: \"" + normalizeYAMLPath(values.RawPath) + "\"",
		"  output_path: \"" + normalizeYAMLPath(values.OutputPath) + "\"",
		"  workspace_path: \"" + normalizeYAMLPath(values.WorkspacePath) + "\"",
		"api:",
		"  base_url: \"" + values.BaseURL + "\"",
		"  partner_id: \"" + values.PartnerID + "\"",
	}
	if strings.TrimSpace(values.UserID) != "" {
		lines = append(lines, "  user_id: \""+values.UserID+"\"")
	}
	lines = append(lines,
		"limits:",
		"  max_file_size_mb: 25",
		"  max_pages: 300",
		"  request_timeout_sec: 5",
		"  max_parallel_runs: 1",
		"security:",
		"  allowed_hosts:",
	)
	for _, host := range values.AllowedHosts {
		trimmed := strings.TrimSpace(host)
		if trimmed != "" {
			lines = append(lines, "    - \""+trimmed+"\"")
		}
	}
	lines = append(lines,
		"audit:",
		"  retention_days: 30",
		"  hmac_key_id: \"audit-hmac-v1\"",
	)
	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}
	return path
}

func normalizeYAMLPath(value string) string {
	return strings.ReplaceAll(value, "\\", "/")
}

type runtimeConfigV11WSValues struct {
	RawPath           string
	OutputPath        string
	WorkspacePath     string
	BaseURL           string
	PartnerID         string
	AllowedHosts      []string
	RuntimeMode       string
	APIAuthToken      string
	AuditHMACSecret   string
}

func writeRuntimeConfigV11WS(t *testing.T, values runtimeConfigV11WSValues) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config-v11ws.yaml")
	lines := []string{
		"paths:",
		"  raw_path: \"" + normalizeYAMLPath(values.RawPath) + "\"",
		"  output_path: \"" + normalizeYAMLPath(values.OutputPath) + "\"",
		"  workspace_path: \"" + normalizeYAMLPath(values.WorkspacePath) + "\"",
		"api:",
		"  base_url: \"" + values.BaseURL + "\"",
		"  partner_id: \"" + values.PartnerID + "\"",
		"  auth_token: \"" + values.APIAuthToken + "\"",
		"limits:",
		"  max_file_size_mb: 25",
		"  max_pages: 300",
		"  request_timeout_sec: 5",
		"  max_parallel_runs: 1",
		"security:",
		"  allowed_hosts:",
	}
	for _, host := range values.AllowedHosts {
		trimmed := strings.TrimSpace(host)
		if trimmed != "" {
			lines = append(lines, "    - \""+trimmed+"\"")
		}
	}
	lines = append(lines,
		"audit:",
		"  retention_days: 30",
		"  hmac_key_id: \"audit-hmac-v1\"",
		"  hmac_secret: \""+values.AuditHMACSecret+"\"",
		"runtime:",
		"  mode: \""+values.RuntimeMode+"\"",
	)
	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}
	return path
}
