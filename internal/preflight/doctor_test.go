package preflight

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"ai-anonymizer/internal/config"
)

type fakeSecretChecker struct {
	ok  bool
	err error
}

func (f fakeSecretChecker) HasSecret(partnerID string) (bool, error) {
	return f.ok, f.err
}

func buildDoctorConfig(t *testing.T, baseURL string) config.Config {
	t.Helper()
	root := t.TempDir()
	workspace := filepath.Join(root, "workspace")
	output := filepath.Join(workspace, "anonymized")
	raw := filepath.Join(root, "secure-raw")

	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatalf("mkdir output failed: %v", err)
	}
	if err := os.MkdirAll(raw, 0o755); err != nil {
		t.Fatalf("mkdir raw failed: %v", err)
	}

	return config.Config{
		RawPath:           raw,
		OutputPath:        output,
		WorkspacePath:     workspace,
		APIBaseURL:        baseURL,
		APIAllowedHosts:   []string{"127.0.0.1", "localhost"},
		MaxParallelRuns:   1,
		APIPartnerID:      "partner-1",
		RequestTimeoutSec: 3,
	}
}

func TestValidateDoctorFailsWhenPartnerSecretMissing(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)

	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: fakeSecretChecker{ok: false},
	})
	if err == nil {
		t.Fatal("expected error when secret is missing")
	}
}

func TestValidateDoctorFailsWhenAPIUnavailable(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)

	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: fakeSecretChecker{ok: true},
	})
	if err == nil {
		t.Fatal("expected error when API preflight endpoint is unavailable")
	}
}

func TestValidateDoctorPassesForHealthySecretAndAPI(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tasks/anonymization_fields" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)

	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: fakeSecretChecker{ok: true},
	})
	if err != nil {
		t.Fatalf("expected doctor preflight to pass, got: %v", err)
	}
}

func TestValidateDoctorFailsWhenSecretCheckerErrors(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)

	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: fakeSecretChecker{err: errors.New("store unavailable")},
	})
	if err == nil {
		t.Fatal("expected error when secret store check fails")
	}
}

func TestValidateDoctorFailsWhenSecretCheckerIsNotConfigured(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)

	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient: server.Client(),
	})
	if err == nil {
		t.Fatal("expected error when secret checker is not configured")
	}
}
