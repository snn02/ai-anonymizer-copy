package config

import "testing"

func TestLoadUsesFlagsOverEnv(t *testing.T) {
	t.Setenv("ANON_RAW_PATH", "D:/env-raw")
	t.Setenv("ANON_OUTPUT_PATH", "C:/env/output")
	t.Setenv("ANON_WORKSPACE_PATH", "C:/env/workspace")
	t.Setenv("ANON_API_BASE_URL", "https://env.company.local")
	t.Setenv("ANON_ALLOWED_HOSTS", "env.company.local")
	t.Setenv("ANON_MAX_PARALLEL_RUNS", "9")
	t.Setenv("ANON_MAX_FILE_SIZE_MB", "50")
	t.Setenv("ANON_MAX_PAGES", "999")
	t.Setenv("ANON_AUDIT_RETENTION_DAYS", "60")
	t.Setenv("ANON_AUDIT_HMAC_KEY_ID", "env-audit-key")

	cfg, err := Load(LoadOptions{
		RawPath:            "D:/flag-raw",
		OutputPath:         "C:/flag/output",
		WorkspacePath:      "C:/flag/workspace",
		APIBaseURL:         "https://flag.company.local",
		AllowedHostsCSV:    "flag.company.local",
		MaxFileSizeMB:      25,
		MaxPages:           300,
		MaxParallelRuns:    1,
		AuditRetentionDays: 30,
		AuditHMACKeyID:     "flag-audit-key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RawPath != "D:\\flag-raw" {
		t.Fatalf("expected flag raw_path, got %q", cfg.RawPath)
	}
	if cfg.APIBaseURL != "https://flag.company.local" {
		t.Fatalf("expected flag base_url, got %q", cfg.APIBaseURL)
	}
	if len(cfg.APIAllowedHosts) != 1 || cfg.APIAllowedHosts[0] != "flag.company.local" {
		t.Fatalf("expected flag allowed_hosts, got %#v", cfg.APIAllowedHosts)
	}
	if cfg.MaxParallelRuns != 1 {
		t.Fatalf("expected flag max_parallel_runs=1, got %d", cfg.MaxParallelRuns)
	}
	if cfg.MaxFileSizeMB != 25 {
		t.Fatalf("expected flag max_file_size_mb=25, got %d", cfg.MaxFileSizeMB)
	}
	if cfg.MaxPages != 300 {
		t.Fatalf("expected flag max_pages=300, got %d", cfg.MaxPages)
	}
	if cfg.AuditRetentionDays != 30 {
		t.Fatalf("expected flag audit_retention_days=30, got %d", cfg.AuditRetentionDays)
	}
	if cfg.AuditHMACKeyID != "flag-audit-key" {
		t.Fatalf("expected flag audit_hmac_key_id, got %q", cfg.AuditHMACKeyID)
	}
}

func TestLoadFailsWhenRequiredFieldsMissing(t *testing.T) {
	_, err := Load(LoadOptions{})
	if err == nil {
		t.Fatal("expected error for missing required fields")
	}
}
