package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	RawPath            string
	OutputPath         string
	WorkspacePath      string
	CatalogPath        string
	AuditPath          string
	APIBaseURL         string
	APIPartnerID       string
	APIAllowedHosts    []string
	RequestTimeoutSec  int
	MaxFileSizeMB      int
	MaxPages           int
	MaxParallelRuns    int
	AuditRetentionDays int
	AuditHMACKeyID     string
}

type LoadOptions struct {
	ConfigPath         string
	RawPath            string
	OutputPath         string
	WorkspacePath      string
	APIBaseURL         string
	APIPartnerID       string
	AllowedHostsCSV    string
	RequestTimeoutSec  int
	MaxFileSizeMB      int
	MaxPages           int
	MaxParallelRuns    int
	AuditRetentionDays int
	AuditHMACKeyID     string
}

func Load(opts LoadOptions) (Config, error) {
	cfg := Config{
		RawPath:            firstNonEmpty(opts.RawPath, os.Getenv("ANON_RAW_PATH")),
		OutputPath:         firstNonEmpty(opts.OutputPath, os.Getenv("ANON_OUTPUT_PATH")),
		WorkspacePath:      firstNonEmpty(opts.WorkspacePath, os.Getenv("ANON_WORKSPACE_PATH")),
		APIBaseURL:         firstNonEmpty(opts.APIBaseURL, os.Getenv("ANON_API_BASE_URL")),
		APIPartnerID:       firstNonEmpty(opts.APIPartnerID, os.Getenv("ANON_API_PARTNER_ID")),
		APIAllowedHosts:    splitCSV(firstNonEmpty(opts.AllowedHostsCSV, os.Getenv("ANON_ALLOWED_HOSTS"))),
		RequestTimeoutSec:  resolveRequestTimeoutSec(opts.RequestTimeoutSec),
		MaxFileSizeMB:      resolveMaxFileSizeMB(opts.MaxFileSizeMB),
		MaxPages:           resolveMaxPages(opts.MaxPages),
		MaxParallelRuns:    1,
		AuditRetentionDays: resolveAuditRetentionDays(opts.AuditRetentionDays),
		AuditHMACKeyID:     firstNonEmpty(opts.AuditHMACKeyID, os.Getenv("ANON_AUDIT_HMAC_KEY_ID"), "audit-hmac-v1"),
	}

	cfg.MaxParallelRuns = resolveParallelRuns(opts.MaxParallelRuns)

	// Для первой итерации файл конфигурации резервируется как расширение.
	_ = opts.ConfigPath

	if cfg.RawPath == "" {
		return Config{}, errors.New("config: raw_path is required")
	}
	if cfg.OutputPath == "" {
		return Config{}, errors.New("config: output_path is required")
	}
	if cfg.WorkspacePath == "" {
		return Config{}, errors.New("config: workspace_path is required")
	}
	if cfg.APIBaseURL == "" {
		return Config{}, errors.New("config: api.base_url is required")
	}

	cfg.RawPath = filepath.Clean(cfg.RawPath)
	cfg.OutputPath = filepath.Clean(cfg.OutputPath)
	cfg.WorkspacePath = filepath.Clean(cfg.WorkspacePath)
	cfg.CatalogPath = filepath.Join(cfg.WorkspacePath, ".anonym", "catalog.json")
	cfg.AuditPath = filepath.Join(cfg.WorkspacePath, ".anonym", "audit.log")

	return cfg, nil
}

func resolveParallelRuns(flagValue int) int {
	if flagValue > 0 {
		return flagValue
	}
	envValue := strings.TrimSpace(os.Getenv("ANON_MAX_PARALLEL_RUNS"))
	if envValue == "" {
		return 1
	}
	n, err := strconv.Atoi(envValue)
	if err != nil || n <= 0 {
		return 1
	}
	return n
}

func resolveRequestTimeoutSec(flagValue int) int {
	if flagValue > 0 {
		return flagValue
	}
	envValue := strings.TrimSpace(os.Getenv("ANON_REQUEST_TIMEOUT_SEC"))
	if envValue == "" {
		return 5
	}
	n, err := strconv.Atoi(envValue)
	if err != nil || n <= 0 {
		return 5
	}
	return n
}

func resolveMaxFileSizeMB(flagValue int) int {
	if flagValue > 0 {
		return flagValue
	}
	envValue := strings.TrimSpace(os.Getenv("ANON_MAX_FILE_SIZE_MB"))
	if envValue == "" {
		return 25
	}
	n, err := strconv.Atoi(envValue)
	if err != nil || n <= 0 {
		return 25
	}
	return n
}

func resolveMaxPages(flagValue int) int {
	if flagValue > 0 {
		return flagValue
	}
	envValue := strings.TrimSpace(os.Getenv("ANON_MAX_PAGES"))
	if envValue == "" {
		return 300
	}
	n, err := strconv.Atoi(envValue)
	if err != nil || n <= 0 {
		return 300
	}
	return n
}

func resolveAuditRetentionDays(flagValue int) int {
	if flagValue > 0 {
		return flagValue
	}
	envValue := strings.TrimSpace(os.Getenv("ANON_AUDIT_RETENTION_DAYS"))
	if envValue == "" {
		return 30
	}
	n, err := strconv.Atoi(envValue)
	if err != nil || n <= 0 {
		return 30
	}
	return n
}

func splitCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}
	return strings.Split(raw, ",")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
