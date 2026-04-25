package preflight

import (
	"testing"

	"ai-anonymizer/internal/config"
)

func TestValidatePassesForSafeConfiguration(t *testing.T) {
	cfg := config.Config{
		RawPath:         "D:/secure-raw",
		OutputPath:      "C:/work/project/anonymized",
		WorkspacePath:   "C:/work/project",
		APIBaseURL:      "https://api.company.local",
		APIAllowedHosts: []string{"api.company.local"},
		MaxParallelRuns: 1,
	}

	if err := Validate(cfg); err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
}

func TestValidateFailsWhenRawInsideWorkspace(t *testing.T) {
	cfg := config.Config{
		RawPath:         "C:/work/project/raw",
		OutputPath:      "C:/work/project/anonymized",
		WorkspacePath:   "C:/work/project",
		APIBaseURL:      "https://api.company.local",
		APIAllowedHosts: []string{"api.company.local"},
		MaxParallelRuns: 1,
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for raw_path inside workspace")
	}
}

func TestValidateFailsWhenBaseURLNotHTTPS(t *testing.T) {
	cfg := config.Config{
		RawPath:         "D:/secure-raw",
		OutputPath:      "C:/work/project/anonymized",
		WorkspacePath:   "C:/work/project",
		APIBaseURL:      "http://api.company.local",
		APIAllowedHosts: []string{"api.company.local"},
		MaxParallelRuns: 1,
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for non-https api.base_url")
	}
}

func TestValidateFailsWhenAllowedHostsEmpty(t *testing.T) {
	cfg := config.Config{
		RawPath:         "D:/secure-raw",
		OutputPath:      "C:/work/project/anonymized",
		WorkspacePath:   "C:/work/project",
		APIBaseURL:      "https://api.company.local",
		APIAllowedHosts: []string{},
		MaxParallelRuns: 1,
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for empty allowed_hosts")
	}
}

func TestValidateFailsWhenMaxParallelRunsNotOne(t *testing.T) {
	cfg := config.Config{
		RawPath:         "D:/secure-raw",
		OutputPath:      "C:/work/project/anonymized",
		WorkspacePath:   "C:/work/project",
		APIBaseURL:      "https://api.company.local",
		APIAllowedHosts: []string{"api.company.local"},
		MaxParallelRuns: 2,
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for max_parallel_runs != 1")
	}
}
