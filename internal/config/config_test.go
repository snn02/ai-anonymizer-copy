package config

import "testing"

func TestLoadUsesFlagsOverEnv(t *testing.T) {
	t.Setenv("ANON_RAW_PATH", "D:/env-raw")
	t.Setenv("ANON_OUTPUT_PATH", "C:/env/output")
	t.Setenv("ANON_WORKSPACE_PATH", "C:/env/workspace")
	t.Setenv("ANON_API_BASE_URL", "https://env.company.local")
	t.Setenv("ANON_ALLOWED_HOSTS", "env.company.local")
	t.Setenv("ANON_MAX_PARALLEL_RUNS", "9")

	cfg, err := Load(LoadOptions{
		RawPath:         "D:/flag-raw",
		OutputPath:      "C:/flag/output",
		WorkspacePath:   "C:/flag/workspace",
		APIBaseURL:      "https://flag.company.local",
		AllowedHostsCSV: "flag.company.local",
		MaxParallelRuns: 1,
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
}

func TestLoadFailsWhenRequiredFieldsMissing(t *testing.T) {
	_, err := Load(LoadOptions{})
	if err == nil {
		t.Fatal("expected error for missing required fields")
	}
}
