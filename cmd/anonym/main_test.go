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

type alwaysTrueSecretChecker struct{}

func (alwaysTrueSecretChecker) HasSecret(partnerID string) (bool, error) {
	return true, nil
}

type staticSecretGetter struct {
	value string
}

func (s staticSecretGetter) GetSecret(partnerID string) (string, error) {
	return s.value, nil
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

	var out bytes.Buffer
	err := run([]string{
		"doctor",
		"--raw-path", raw,
		"--output-path", output,
		"--workspace-path", workspace,
		"--api-base-url", server.URL,
		"--api-partner-id", "partner-1",
		"--allowed-hosts", "127.0.0.1,localhost",
		"--max-parallel-runs", "1",
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

	var scanOut bytes.Buffer
	scanErr := run([]string{
		"scan",
		"--raw-path", raw,
		"--output-path", output,
		"--workspace-path", workspace,
		"--api-base-url", "https://api.company.local",
		"--api-partner-id", "partner-1",
		"--allowed-hosts", "api.company.local",
		"--max-parallel-runs", "1",
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
		"--raw-path", raw,
		"--output-path", output,
		"--workspace-path", workspace,
		"--api-base-url", "https://api.company.local",
		"--api-partner-id", "partner-1",
		"--allowed-hosts", "api.company.local",
		"--max-parallel-runs", "1",
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
		if got := r.Header.Get("partner-id"); got != "partner-1" {
			t.Fatalf("unexpected partner-id header: %q", got)
		}
		if got := r.Header.Get("fields"); got == "" {
			t.Fatalf("expected fields header")
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

	var scanOut bytes.Buffer
	if err := run([]string{
		"scan",
		"--raw-path", raw,
		"--output-path", output,
		"--workspace-path", workspace,
		"--api-base-url", server.URL,
		"--api-partner-id", "partner-1",
		"--allowed-hosts", "127.0.0.1,localhost",
		"--max-parallel-runs", "1",
	}, &scanOut, &scanOut); err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	var runOut bytes.Buffer
	err := run([]string{
		"run",
		catalog.BuildID("in/sample.txt"),
		"--raw-path", raw,
		"--output-path", output,
		"--workspace-path", workspace,
		"--api-base-url", server.URL,
		"--api-partner-id", "partner-1",
		"--allowed-hosts", "127.0.0.1,localhost",
		"--max-parallel-runs", "1",
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

	var scanOut bytes.Buffer
	if err := run([]string{
		"scan",
		"--raw-path", raw,
		"--output-path", output,
		"--workspace-path", workspace,
		"--api-base-url", server.URL,
		"--api-partner-id", "partner-1",
		"--allowed-hosts", "127.0.0.1,localhost",
		"--max-parallel-runs", "1",
	}, &scanOut, &scanOut); err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	var runOut bytes.Buffer
	err := run([]string{
		"run",
		"invoice",
		"--raw-path", raw,
		"--output-path", output,
		"--workspace-path", workspace,
		"--api-base-url", server.URL,
		"--api-partner-id", "partner-1",
		"--allowed-hosts", "127.0.0.1,localhost",
		"--max-parallel-runs", "1",
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

	var scanOut bytes.Buffer
	if err := run([]string{
		"scan",
		"--raw-path", raw,
		"--output-path", output,
		"--workspace-path", workspace,
		"--api-base-url", server.URL,
		"--api-partner-id", "partner-1",
		"--allowed-hosts", "127.0.0.1,localhost",
		"--max-parallel-runs", "1",
	}, &scanOut, &scanOut); err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	var runOut bytes.Buffer
	runErr := run([]string{
		"run",
		catalog.BuildID("in/sample.txt"),
		"--raw-path", raw,
		"--output-path", output,
		"--workspace-path", workspace,
		"--api-base-url", server.URL,
		"--api-partner-id", "partner-1",
		"--allowed-hosts", "127.0.0.1,localhost",
		"--max-parallel-runs", "1",
	}, &runOut, &runOut)
	if runErr == nil {
		t.Fatal("expected run failure")
	}

	var listOut bytes.Buffer
	if err := run([]string{
		"list",
		"--raw-path", raw,
		"--output-path", output,
		"--workspace-path", workspace,
		"--api-base-url", server.URL,
		"--api-partner-id", "partner-1",
		"--allowed-hosts", "127.0.0.1,localhost",
		"--max-parallel-runs", "1",
	}, &listOut, &listOut); err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(listOut.String(), "\tfailed") {
		t.Fatalf("expected failed status in catalog, got %q", listOut.String())
	}
}

func runDepsForTest(client *http.Client) anonymizer.Deps {
	return anonymizer.Deps{
		HTTPClient:   client,
		SecretGetter: staticSecretGetter{value: "token"},
	}
}
