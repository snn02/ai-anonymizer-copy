package preflight

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"ai-anonymizer/internal/config"
)

func Validate(cfg config.Config) error {
	if err := validateBoundary(cfg.RawPath, cfg.WorkspacePath); err != nil {
		return err
	}
	if err := validateAPI(cfg.APIBaseURL, cfg.APIAllowedHosts); err != nil {
		return err
	}
	if cfg.MaxParallelRuns != 1 {
		return fmt.Errorf("preflight: max_parallel_runs must be 1 in v1, got %d", cfg.MaxParallelRuns)
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
