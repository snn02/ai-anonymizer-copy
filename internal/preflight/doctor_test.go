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
	err    error
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
		APIUserID:          "70bd3a91-9999-9999-9999-999999999999",
		APIFieldsPage:      1,
		APIFieldsPerPage:   10,
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
	var gotUserID string
	var gotPage string
	var gotPerPage string
	var gotMethod string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tasks/anonymization_fields" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		gotMethod = r.Method
		gotAuth = r.Header.Get("Authorization")
		gotPartnerID = r.Header.Get("partner-id")
		gotUserID = r.Header.Get("user-id")
		gotPage = r.URL.Query().Get("page")
		gotPerPage = r.URL.Query().Get("per_page")
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
	if gotMethod != http.MethodGet {
		t.Fatalf("unexpected method %q", gotMethod)
	}
	if gotPartnerID != cfg.APIPartnerID {
		t.Fatalf("unexpected partner-id header %q", gotPartnerID)
	}
	if gotUserID != cfg.APIUserID {
		t.Fatalf("unexpected user-id header %q", gotUserID)
	}
	if gotPage != "1" || gotPerPage != "10" {
		t.Fatalf("unexpected pagination query page=%q per_page=%q", gotPage, gotPerPage)
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

func TestValidateDoctorRemovesControlCharsFromSecretHeader(t *testing.T) {
	var gotAuth string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)
	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: fakeSecretChecker{secret: "\x00tok\ren\n"},
	})
	if err != nil {
		t.Fatalf("expected doctor preflight to pass, got: %v", err)
	}
	if gotAuth != "token" {
		t.Fatalf("unexpected Authorization header value %q", gotAuth)
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

func TestValidateDoctorMVPUsesTokenWithoutSecretStore(t *testing.T) {
	var gotAuth string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)
	cfg.RuntimeMode = "mvp"
	cfg.APIAuthToken = "session-token"
	cfg.AuditHMACSecret = "audit-secret"

	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient: server.Client(),
	})
	if err != nil {
		t.Fatalf("expected success for mvp token mode, got %v", err)
	}
	if gotAuth != "session-token" {
		t.Fatalf("unexpected auth header %q", gotAuth)
	}
}

func TestValidateDoctorMVPFailsWithoutAuditSecret(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)
	cfg.RuntimeMode = "mvp"
	cfg.APIAuthToken = "session-token"

	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient: server.Client(),
	})
	if err == nil || !strings.Contains(err.Error(), "audit.hmac_secret") {
		t.Fatalf("expected missing audit secret error, got %v", err)
	}
}

func TestValidateRuntimeMVPAllowsProductionHost(t *testing.T) {
	cfg := buildDoctorConfig(t, "https://production-retrievals.ai.rarus-cloud.ru")
	cfg.APIAllowedHosts = []string{"production-retrievals.ai.rarus-cloud.ru"}
	cfg.RuntimeMode = "mvp"
	cfg.APIAuthToken = "session-token"
	cfg.AuditHMACSecret = "audit-secret"

	err := ValidateRuntime(cfg)
	if err != nil {
		t.Fatalf("expected production host to be allowed in mvp mode, got %v", err)
	}
}

func TestValidateDoctorWritesHTTPDebugLogWhenEnabled(t *testing.T) {
	t.Setenv("ANON_HTTP_DEBUG", "true")
	var gotAuth string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := buildDoctorConfig(t, server.URL)
	err := ValidateDoctor(cfg, DoctorDeps{
		HTTPClient:    server.Client(),
		SecretChecker: fakeSecretChecker{secret: "secret-token-value"},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if gotAuth == "" {
		t.Fatal("expected auth header to be sent")
	}
	logPath := filepath.Join(cfg.WorkspacePath, ".anonym", "http-debug.log")
	data, readErr := os.ReadFile(logPath)
	if readErr != nil {
		t.Fatalf("expected debug log file, got %v", readErr)
	}
	logText := string(data)
	if !strings.Contains(logText, "scope=doctor") || !strings.Contains(logText, "status=200") {
		t.Fatalf("unexpected debug log content: %q", logText)
	}
	if strings.Contains(logText, "secret-token-value") {
		t.Fatalf("authorization must be masked in debug log: %q", logText)
	}
}
