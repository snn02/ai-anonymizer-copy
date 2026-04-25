package preflight

import (
	"os"
	"path/filepath"
	"testing"

	"ai-anonymizer/internal/config"
)

func buildSafeRuntimeConfig(t *testing.T) config.Config {
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
		APIBaseURL:         "https://api.company.local",
		APIAllowedHosts:    []string{"api.company.local"},
		MaxPages:           300,
		MaxParallelRuns:    1,
		AuditRetentionDays: 30,
		AuditHMACKeyID:     "audit-hmac-v1",
	}
}

func TestValidatePassesForSafeConfiguration(t *testing.T) {
	cfg := buildSafeRuntimeConfig(t)

	if err := Validate(cfg); err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
}

func TestValidateFailsWhenRawPathDoesNotExist(t *testing.T) {
	workspace := t.TempDir()
	output := filepath.Join(workspace, "anonymized")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatalf("mkdir output failed: %v", err)
	}

	cfg := buildSafeRuntimeConfig(t)
	cfg.RawPath = filepath.Join(t.TempDir(), "missing-raw")
	cfg.OutputPath = output
	cfg.WorkspacePath = workspace

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for missing raw_path")
	}
}

func TestValidateFailsWhenOutputPathDoesNotExist(t *testing.T) {
	workspace := t.TempDir()
	raw := filepath.Join(t.TempDir(), "raw")
	if err := os.MkdirAll(raw, 0o755); err != nil {
		t.Fatalf("mkdir raw failed: %v", err)
	}

	cfg := buildSafeRuntimeConfig(t)
	cfg.RawPath = raw
	cfg.OutputPath = filepath.Join(workspace, "missing-output")
	cfg.WorkspacePath = workspace

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for missing output_path")
	}
}

func TestValidateFailsWhenRawInsideWorkspace(t *testing.T) {
	cfg := buildSafeRuntimeConfig(t)
	cfg.RawPath = filepath.Join(cfg.WorkspacePath, "raw")
	if err := os.MkdirAll(cfg.RawPath, 0o755); err != nil {
		t.Fatalf("mkdir workspace raw failed: %v", err)
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for raw_path inside workspace")
	}
}

func TestValidateFailsWhenBaseURLNotHTTPS(t *testing.T) {
	cfg := buildSafeRuntimeConfig(t)
	cfg.APIBaseURL = "http://api.company.local"

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for non-https api.base_url")
	}
}

func TestValidateFailsWhenAllowedHostsEmpty(t *testing.T) {
	cfg := buildSafeRuntimeConfig(t)
	cfg.APIAllowedHosts = []string{}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for empty allowed_hosts")
	}
}

func TestValidateFailsWhenMaxParallelRunsNotOne(t *testing.T) {
	cfg := buildSafeRuntimeConfig(t)
	cfg.MaxParallelRuns = 2

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for max_parallel_runs != 1")
	}
}

func TestValidateFailsWhenMaxPagesIsNotPositive(t *testing.T) {
	cfg := buildSafeRuntimeConfig(t)
	cfg.MaxPages = 0

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for max_pages <= 0")
	}
}

func TestValidateFailsWhenAuditRetentionIsNotPositive(t *testing.T) {
	cfg := buildSafeRuntimeConfig(t)
	cfg.AuditRetentionDays = 0

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for audit.retention_days <= 0")
	}
}

func TestValidateFailsWhenAuditHMACKeyIDIsEmpty(t *testing.T) {
	cfg := buildSafeRuntimeConfig(t)
	cfg.AuditHMACKeyID = " "

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for empty audit.hmac_key_id")
	}
}
