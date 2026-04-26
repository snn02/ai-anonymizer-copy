package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadUsesConfigFileWhenNoOverrides(t *testing.T) {
	cfgPath := writeTempConfig(t, `
paths:
  raw_path: "D:/file-raw"
  output_path: "C:/file/output"
  workspace_path: "C:/file/workspace"
api:
  base_url: "https://file.company.local"
  partner_id: "70bd3a91-0000-0000-0000-000000000000"
limits:
  max_file_size_mb: 41
  max_pages: 420
  request_timeout_sec: 17
  max_parallel_runs: 1
security:
  allowed_hosts:
    - "file.company.local"
audit:
  retention_days: 44
  hmac_key_id: "file-audit-key"
`)

	cfg, err := Load(LoadOptions{
		ConfigPath: cfgPath,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RawPath != filepath.Clean("D:/file-raw") {
		t.Fatalf("expected raw_path from file, got %q", cfg.RawPath)
	}
	if cfg.OutputPath != filepath.Clean("C:/file/output") {
		t.Fatalf("expected output_path from file, got %q", cfg.OutputPath)
	}
	if cfg.WorkspacePath != filepath.Clean("C:/file/workspace") {
		t.Fatalf("expected workspace_path from file, got %q", cfg.WorkspacePath)
	}
	if cfg.APIBaseURL != "https://file.company.local" {
		t.Fatalf("expected base_url from file, got %q", cfg.APIBaseURL)
	}
	if cfg.APIPartnerID != "70bd3a91-0000-0000-0000-000000000000" {
		t.Fatalf("expected partner_id from file, got %q", cfg.APIPartnerID)
	}
	if len(cfg.APIAllowedHosts) != 1 || cfg.APIAllowedHosts[0] != "file.company.local" {
		t.Fatalf("expected allowed_hosts from file, got %#v", cfg.APIAllowedHosts)
	}
	if cfg.RequestTimeoutSec != 17 {
		t.Fatalf("expected request_timeout_sec=17, got %d", cfg.RequestTimeoutSec)
	}
	if cfg.MaxFileSizeMB != 41 {
		t.Fatalf("expected max_file_size_mb=41, got %d", cfg.MaxFileSizeMB)
	}
	if cfg.MaxPages != 420 {
		t.Fatalf("expected max_pages=420, got %d", cfg.MaxPages)
	}
	if cfg.MaxParallelRuns != 1 {
		t.Fatalf("expected max_parallel_runs=1, got %d", cfg.MaxParallelRuns)
	}
	if cfg.AuditRetentionDays != 44 {
		t.Fatalf("expected audit_retention_days=44, got %d", cfg.AuditRetentionDays)
	}
	if cfg.AuditHMACKeyID != "file-audit-key" {
		t.Fatalf("expected audit_hmac_key_id=file-audit-key, got %q", cfg.AuditHMACKeyID)
	}
}

func TestLoadUsesPrecedenceCLIOverEnvOverFile(t *testing.T) {
	cfgPath := writeTempConfig(t, `
paths:
  raw_path: "D:/file-raw"
  output_path: "C:/file/output"
  workspace_path: "C:/file/workspace"
api:
  base_url: "https://file.company.local"
  partner_id: "70bd3a91-1111-1111-1111-111111111111"
limits:
  max_file_size_mb: 41
  max_pages: 420
  request_timeout_sec: 17
  max_parallel_runs: 1
security:
  allowed_hosts:
    - "file.company.local"
audit:
  retention_days: 44
  hmac_key_id: "file-audit-key"
`)

	t.Setenv("ANON_RAW_PATH", "D:/env-raw")
	t.Setenv("ANON_OUTPUT_PATH", "C:/env/output")
	t.Setenv("ANON_WORKSPACE_PATH", "C:/env/workspace")
	t.Setenv("ANON_API_BASE_URL", "https://env.company.local")
	t.Setenv("ANON_API_PARTNER_ID", "70bd3a91-2222-2222-2222-222222222222")
	t.Setenv("ANON_ALLOWED_HOSTS", "env.company.local")
	t.Setenv("ANON_REQUEST_TIMEOUT_SEC", "33")
	t.Setenv("ANON_MAX_FILE_SIZE_MB", "50")
	t.Setenv("ANON_MAX_PAGES", "500")
	t.Setenv("ANON_MAX_PARALLEL_RUNS", "1")
	t.Setenv("ANON_AUDIT_RETENTION_DAYS", "60")
	t.Setenv("ANON_AUDIT_HMAC_KEY_ID", "env-audit-key")

	cfg, err := Load(LoadOptions{
		ConfigPath:         cfgPath,
		RawPath:            "D:/cli-raw",
		AllowedHostsCSV:    "cli.company.local",
		MaxFileSizeMB:      25,
		MaxPages:           300,
		RequestTimeoutSec:  5,
		MaxParallelRuns:    1,
		AuditRetentionDays: 30,
		AuditHMACKeyID:     "cli-audit-key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RawPath != filepath.Clean("D:/cli-raw") {
		t.Fatalf("expected CLI raw_path, got %q", cfg.RawPath)
	}
	if cfg.OutputPath != filepath.Clean("C:/env/output") {
		t.Fatalf("expected env output_path, got %q", cfg.OutputPath)
	}
	if cfg.WorkspacePath != filepath.Clean("C:/env/workspace") {
		t.Fatalf("expected env workspace_path, got %q", cfg.WorkspacePath)
	}
	if cfg.APIBaseURL != "https://env.company.local" {
		t.Fatalf("expected env base_url, got %q", cfg.APIBaseURL)
	}
	if cfg.APIPartnerID != "70bd3a91-2222-2222-2222-222222222222" {
		t.Fatalf("expected env partner_id, got %q", cfg.APIPartnerID)
	}
	if len(cfg.APIAllowedHosts) != 1 || strings.TrimSpace(cfg.APIAllowedHosts[0]) != "cli.company.local" {
		t.Fatalf("expected CLI allowed_hosts, got %#v", cfg.APIAllowedHosts)
	}
	if cfg.RequestTimeoutSec != 5 {
		t.Fatalf("expected CLI request_timeout_sec=5, got %d", cfg.RequestTimeoutSec)
	}
	if cfg.MaxFileSizeMB != 25 {
		t.Fatalf("expected CLI max_file_size_mb=25, got %d", cfg.MaxFileSizeMB)
	}
	if cfg.MaxPages != 300 {
		t.Fatalf("expected CLI max_pages=300, got %d", cfg.MaxPages)
	}
	if cfg.MaxParallelRuns != 1 {
		t.Fatalf("expected CLI/env max_parallel_runs=1, got %d", cfg.MaxParallelRuns)
	}
	if cfg.AuditRetentionDays != 30 {
		t.Fatalf("expected CLI audit_retention_days=30, got %d", cfg.AuditRetentionDays)
	}
	if cfg.AuditHMACKeyID != "cli-audit-key" {
		t.Fatalf("expected CLI audit_hmac_key_id, got %q", cfg.AuditHMACKeyID)
	}
}

func TestLoadFailsOnInvalidConfigYAML(t *testing.T) {
	cfgPath := writeTempConfig(t, "paths:\n  raw_path: [bad")
	_, err := Load(LoadOptions{
		ConfigPath: cfgPath,
	})
	if err == nil {
		t.Fatal("expected config parse error")
	}
	if !strings.Contains(err.Error(), "config: parse") {
		t.Fatalf("expected parse error prefix, got %v", err)
	}
}

func TestLoadAllowsMissingDefaultConfigPath(t *testing.T) {
	tmp := t.TempDir()
	missingPath := filepath.Join(tmp, "missing-config.yaml")
	_, statErr := os.Stat(missingPath)
	if !os.IsNotExist(statErr) {
		t.Fatalf("expected missing test config path, got %v", statErr)
	}

	cfg, err := Load(LoadOptions{
		ConfigPath:      missingPath,
		RawPath:         "D:/cli-raw",
		OutputPath:      "C:/cli/output",
		WorkspacePath:   "C:/cli/workspace",
		APIBaseURL:      "https://cli.company.local",
		APIPartnerID:    "70bd3a91-3333-3333-3333-333333333333",
		AllowedHostsCSV: "cli.company.local",
		MaxParallelRuns: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error for missing default path: %v", err)
	}
	if cfg.RawPath != filepath.Clean("D:/cli-raw") {
		t.Fatalf("expected CLI raw path, got %q", cfg.RawPath)
	}
}

func TestLoadFailsOnMissingExplicitConfigPath(t *testing.T) {
	tmp := t.TempDir()
	missingPath := filepath.Join(tmp, "missing-explicit.yaml")
	_, err := Load(LoadOptions{
		ConfigPath:         missingPath,
		ConfigPathExplicit: true,
		RawPath:            "D:/cli-raw",
		OutputPath:         "C:/cli/output",
		WorkspacePath:      "C:/cli/workspace",
		APIBaseURL:         "https://cli.company.local",
		APIPartnerID:       "70bd3a91-4444-4444-4444-444444444444",
		AllowedHostsCSV:    "cli.company.local",
		MaxParallelRuns:    1,
	})
	if err == nil {
		t.Fatal("expected error for missing explicit config path")
	}
	if !strings.Contains(err.Error(), "config: file") {
		t.Fatalf("expected missing file error, got %v", err)
	}
}

func TestLoadFailsWhenRequiredFieldsMissing(t *testing.T) {
	_, err := Load(LoadOptions{})
	if err == nil {
		t.Fatal("expected error for missing required fields")
	}
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}
