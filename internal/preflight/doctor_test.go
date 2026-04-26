package preflight

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-anonymizer/internal/config"
	"ai-anonymizer/internal/secrets"
)

type fakeSecretChecker struct {
	secret string
	err error
}

func (f fakeSecretChecker) GetSecret(partnerID string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.secret, nil
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
		RawPath:            raw,
		OutputPath:         output,
		WorkspacePath:      workspace,
		APIBaseURL:         baseURL,
		APIAllowedHosts:    []string{"127.0.0.1", "localhost"},
		MaxPages:           300,
		MaxParallelRuns:    1,
		APIPartnerID:       "70bd3a91-0000-0000-0000-000000000000",
		RequestTimeoutSec:  3,
		AuditRetentionDays: 30,
		AuditHMACKeyID:     "audit-hmac-v1",
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
		SecretChecker: fakeSecretChecker{secret: ""},
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
		SecretChecker: fakeSecretChecker{secret: "token"},
	})
	if err == nil {
		t.Fatal("expected error when API preflight endpoint is unavailable")
	}
}

func TestValidateDoctorPassesForHealthySecretAndAPI(t *testing.T) {
	var gotAuth string
	var gotPartnerID string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tasks/anonymization_fields" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		gotAuth = r.Header.Get("Authorization")
		gotPartnerID = r.Header.Get("partner-id")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)

	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: fakeSecretChecker{secret: "token"},
	})
	if err != nil {
		t.Fatalf("expected doctor preflight to pass, got: %v", err)
	}
	if gotAuth != "token" {
		t.Fatalf("unexpected Authorization header value %q", gotAuth)
	}
	if gotPartnerID != cfg.APIPartnerID {
		t.Fatalf("unexpected partner-id header %q", gotPartnerID)
	}
}

func TestValidateDoctorTrimsSecretAndPartnerIDHeaders(t *testing.T) {
	var gotAuth string
	var gotPartnerID string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tasks/anonymization_fields" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		gotAuth = r.Header.Get("Authorization")
		gotPartnerID = r.Header.Get("partner-id")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)
	cfg.APIPartnerID = " 70bd3a91-0000-0000-0000-000000000000 "

	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: fakeSecretChecker{secret: "\r\ntoken\r\n"},
	})
	if err != nil {
		t.Fatalf("expected doctor preflight to pass, got: %v", err)
	}
	if gotAuth != "token" {
		t.Fatalf("unexpected Authorization header value %q", gotAuth)
	}
	if gotPartnerID != "70bd3a91-0000-0000-0000-000000000000" {
		t.Fatalf("unexpected partner-id header %q", gotPartnerID)
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

func TestValidateDoctorMapsSecretNotFoundToDeterministicError(t *testing.T) {
	serverCalled := false
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)
	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: fakeSecretChecker{err: secrets.ErrSecretNotFound},
	})
	if err == nil || !strings.Contains(err.Error(), "api secret was not found in OS secret store") {
		t.Fatalf("expected deterministic not-found error, got %v", err)
	}
	if serverCalled {
		t.Fatal("expected no api call when secret is not found")
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

func TestValidateDoctorFailsWhenPartnerIDIsNotUUID(t *testing.T) {
	serverCalled := false
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)
	cfg.APIPartnerID = "partner-1"

	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: fakeSecretChecker{secret: "token"},
	})
	if err == nil {
		t.Fatal("expected error for non-uuid api.partner_id")
	}
	if serverCalled {
		t.Fatal("expected no api call when api.partner_id is invalid")
	}
}
