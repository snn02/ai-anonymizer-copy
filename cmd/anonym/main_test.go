package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-anonymizer/internal/preflight"
)

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
		SecretChecker: preflight.AllowAllSecretChecker{},
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
