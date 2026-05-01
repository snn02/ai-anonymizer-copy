package preflight

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"ai-anonymizer/internal/config"
	"ai-anonymizer/internal/httpdebug"
	"ai-anonymizer/internal/secrets"
)

type SecretChecker interface {
	GetSecret(partnerID string) (string, error)
}

type DoctorDeps struct {
	HTTPClient    *http.Client
	SecretChecker SecretChecker
}

var uuidPattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func Validate(cfg config.Config) error {
	return ValidateRuntime(cfg)
}

func ValidateRuntime(cfg config.Config) error {
	mode := runtimeMode(cfg)
	if err := validateBoundary(cfg.RawPath, cfg.WorkspacePath); err != nil {
		return err
	}
	if err := validatePathsAccess(cfg.RawPath, cfg.OutputPath); err != nil {
		return err
	}
	if err := validateAPI(cfg.APIBaseURL, cfg.APIAllowedHosts); err != nil {
		return err
	}
	if cfg.MaxParallelRuns != 1 {
		return fmt.Errorf("preflight: max_parallel_runs must be 1 in v1, got %d", cfg.MaxParallelRuns)
	}
	if cfg.MaxPages <= 0 {
		return fmt.Errorf("preflight: max_pages must be positive, got %d", cfg.MaxPages)
	}
	if cfg.AuditRetentionDays <= 0 {
		return fmt.Errorf("preflight: audit.retention_days must be positive, got %d", cfg.AuditRetentionDays)
	}
	if strings.TrimSpace(cfg.AuditHMACKeyID) == "" {
		return errors.New("preflight: audit.hmac_key_id is required")
	}
	if strings.TrimSpace(cfg.APIPartnerID) != "" && !uuidPattern.MatchString(strings.TrimSpace(cfg.APIPartnerID)) {
		return errors.New("preflight: api.partner_id must be a valid UUID")
	}
	if strings.TrimSpace(cfg.APIUserID) != "" && !uuidPattern.MatchString(strings.TrimSpace(cfg.APIUserID)) {
		return errors.New("preflight: api.user_id must be a valid UUID")
	}
	if mode != "prod" && mode != "mvp" {
		return fmt.Errorf("preflight: unsupported runtime mode %q", cfg.RuntimeMode)
	}
	return nil
}

func validatePathsAccess(rawPath, outputPath string) error {
	rawInfo, err := os.Stat(rawPath)
	if err != nil {
		return fmt.Errorf("preflight: raw_path is not accessible: %w", err)
	}
	if !rawInfo.IsDir() {
		return errors.New("preflight: raw_path must be a directory")
	}

	outputInfo, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("preflight: output_path is not accessible: %w", err)
	}
	if !outputInfo.IsDir() {
		return errors.New("preflight: output_path must be a directory")
	}

	probePath := filepath.Join(outputPath, ".anonym-write-probe")
	if err := os.WriteFile(probePath, []byte("ok"), 0o600); err != nil {
		return fmt.Errorf("preflight: output_path is not writable: %w", err)
	}
	_ = os.Remove(probePath)
	return nil
}

func ValidateDoctor(cfg config.Config, deps DoctorDeps) error {
	if err := ValidateRuntime(cfg); err != nil {
		return err
	}
	if strings.TrimSpace(cfg.APIPartnerID) == "" {
		return errors.New("preflight: api.partner_id is required for doctor")
	}
	checker := deps.SecretChecker
	normalizedSecret, err := resolveDoctorSecret(cfg, checker)
	if err != nil {
		return err
	}
	normalizedPartnerID := strings.TrimSpace(cfg.APIPartnerID)

	timeout := time.Duration(cfg.RequestTimeoutSec) * time.Second
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	client := deps.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}

	endpoint := strings.TrimRight(cfg.APIBaseURL, "/") + "/v1/tasks/anonymization_fields"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("preflight: build api request: %w", err)
	}
	req.Header.Set("Authorization", normalizedSecret)
	req.Header.Set("partner-id", normalizedPartnerID)
	if strings.TrimSpace(cfg.APIUserID) != "" {
		req.Header.Set("user-id", strings.TrimSpace(cfg.APIUserID))
	}
	q := req.URL.Query()
	page := cfg.APIFieldsPage
	if page <= 0 {
		page = 1
	}
	perPage := cfg.APIFieldsPerPage
	if perPage <= 0 {
		perPage = 10
	}
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("per_page", fmt.Sprintf("%d", perPage))
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		httpdebug.Log(cfg, "doctor", req, 0, err)
		return fmt.Errorf("preflight: api availability check failed: %w", err)
	}
	defer resp.Body.Close()
	httpdebug.Log(cfg, "doctor", req, resp.StatusCode, nil)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("preflight: api availability check returned status %d", resp.StatusCode)
	}

	return nil
}

func resolveDoctorSecret(cfg config.Config, checker SecretChecker) (string, error) {
	if runtimeMode(cfg) == "mvp" {
		token := sanitizeHeaderValue(cfg.APIAuthToken)
		if token == "" {
			return "", errors.New("preflight: api.auth_token is required in mvp mode")
		}
		if strings.TrimSpace(cfg.AuditHMACSecret) == "" {
			return "", errors.New("preflight: audit.hmac_secret is required in mvp mode")
		}
		return token, nil
	}
	if checker == nil {
		return "", errors.New("preflight: secret checker is not configured")
	}
	secret, err := checker.GetSecret(cfg.APIPartnerID)
	if err != nil {
		if errors.Is(err, secrets.ErrSecretNotFound) {
			return "", errors.New("preflight: api secret was not found in OS secret store")
		}
		return "", fmt.Errorf("preflight: secret store check failed: %w", err)
	}
	normalizedSecret := sanitizeHeaderValue(secret)
	if normalizedSecret == "" {
		return "", errors.New("preflight: api secret was not found in OS secret store")
	}
	return normalizedSecret, nil
}

func runtimeMode(cfg config.Config) string {
	mode := strings.TrimSpace(strings.ToLower(cfg.RuntimeMode))
	if mode == "" {
		return "prod"
	}
	return mode
}

func sanitizeHeaderValue(v string) string {
	trimmed := strings.TrimSpace(v)
	return strings.Map(func(r rune) rune {
		if r == unicode.ReplacementChar {
			return -1
		}
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, trimmed)
}

func validateBoundary(rawPath, workspacePath string) error {
	rawAbs, err := filepath.Abs(rawPath)
	if err != nil {
		return fmt.Errorf("preflight: invalid raw_path: %w", err)
	}
	workspaceAbs, err := filepath.Abs(workspacePath)
	if err != nil {
		return fmt.Errorf("preflight: invalid workspace_path: %w", err)
	}

	raw := strings.ToLower(filepath.Clean(rawAbs))
	workspace := strings.ToLower(filepath.Clean(workspaceAbs))

	if raw == workspace || strings.HasPrefix(raw, workspace+string(filepath.Separator)) {
		return errors.New("preflight: raw_path must be outside workspace_path")
	}

	return nil
}

func validateAPI(baseURL string, allowedHosts []string) error {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("preflight: invalid api.base_url: %w", err)
	}
	if !strings.EqualFold(parsed.Scheme, "https") {
		return errors.New("preflight: api.base_url must use https")
	}

	host := strings.ToLower(parsed.Hostname())
	if host == "" {
		return errors.New("preflight: api.base_url host is required")
	}
	if len(allowedHosts) == 0 {
		return errors.New("preflight: security.allowed_hosts must not be empty")
	}

	for _, allowed := range allowedHosts {
		if strings.EqualFold(strings.TrimSpace(allowed), host) {
			return nil
		}
	}

	return fmt.Errorf("preflight: api host %q is not in security.allowed_hosts", host)
}
