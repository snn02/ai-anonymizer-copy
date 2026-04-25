package preflight

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ai-anonymizer/internal/config"
)

type SecretChecker interface {
	HasSecret(partnerID string) (bool, error)
}

type DoctorDeps struct {
	HTTPClient    *http.Client
	SecretChecker SecretChecker
}

func Validate(cfg config.Config) error {
	return ValidateRuntime(cfg)
}

func ValidateRuntime(cfg config.Config) error {
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
	if checker == nil {
		return errors.New("preflight: secret checker is not configured")
	}
	ok, err := checker.HasSecret(cfg.APIPartnerID)
	if err != nil {
		return fmt.Errorf("preflight: secret store check failed: %w", err)
	}
	if !ok {
		return errors.New("preflight: api secret was not found in OS secret store")
	}

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
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("preflight: api availability check failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("preflight: api availability check returned status %d", resp.StatusCode)
	}

	return nil
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
