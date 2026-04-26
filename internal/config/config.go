package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
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
	ConfigPathExplicit bool
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

type fileConfig struct {
	Paths struct {
		RawPath       string `yaml:"raw_path"`
		OutputPath    string `yaml:"output_path"`
		WorkspacePath string `yaml:"workspace_path"`
	} `yaml:"paths"`
	API struct {
		BaseURL   string `yaml:"base_url"`
		PartnerID string `yaml:"partner_id"`
	} `yaml:"api"`
	Limits struct {
		MaxFileSizeMB     int `yaml:"max_file_size_mb"`
		MaxPages          int `yaml:"max_pages"`
		RequestTimeoutSec int `yaml:"request_timeout_sec"`
		MaxParallelRuns   int `yaml:"max_parallel_runs"`
	} `yaml:"limits"`
	Security struct {
		AllowedHosts []string `yaml:"allowed_hosts"`
	} `yaml:"security"`
	Audit struct {
		RetentionDays int    `yaml:"retention_days"`
		HMACKeyID     string `yaml:"hmac_key_id"`
	} `yaml:"audit"`
}

func Load(opts LoadOptions) (Config, error) {
	fileCfg, err := readFileConfig(opts.ConfigPath, opts.ConfigPathExplicit)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		RawPath:            firstNonEmpty(opts.RawPath, os.Getenv("ANON_RAW_PATH"), fileCfg.Paths.RawPath),
		OutputPath:         firstNonEmpty(opts.OutputPath, os.Getenv("ANON_OUTPUT_PATH"), fileCfg.Paths.OutputPath),
		WorkspacePath:      firstNonEmpty(opts.WorkspacePath, os.Getenv("ANON_WORKSPACE_PATH"), fileCfg.Paths.WorkspacePath),
		APIBaseURL:         firstNonEmpty(opts.APIBaseURL, os.Getenv("ANON_API_BASE_URL"), fileCfg.API.BaseURL),
		APIPartnerID:       firstNonEmpty(opts.APIPartnerID, os.Getenv("ANON_API_PARTNER_ID"), fileCfg.API.PartnerID),
		APIAllowedHosts:    resolveAllowedHosts(opts.AllowedHostsCSV, os.Getenv("ANON_ALLOWED_HOSTS"), fileCfg.Security.AllowedHosts),
		RequestTimeoutSec:  resolvePositiveInt(opts.RequestTimeoutSec, "ANON_REQUEST_TIMEOUT_SEC", fileCfg.Limits.RequestTimeoutSec, 5),
		MaxFileSizeMB:      resolvePositiveInt(opts.MaxFileSizeMB, "ANON_MAX_FILE_SIZE_MB", fileCfg.Limits.MaxFileSizeMB, 25),
		MaxPages:           resolvePositiveInt(opts.MaxPages, "ANON_MAX_PAGES", fileCfg.Limits.MaxPages, 300),
		MaxParallelRuns:    resolvePositiveInt(opts.MaxParallelRuns, "ANON_MAX_PARALLEL_RUNS", fileCfg.Limits.MaxParallelRuns, 1),
		AuditRetentionDays: resolvePositiveInt(opts.AuditRetentionDays, "ANON_AUDIT_RETENTION_DAYS", fileCfg.Audit.RetentionDays, 30),
		AuditHMACKeyID:     firstNonEmpty(opts.AuditHMACKeyID, os.Getenv("ANON_AUDIT_HMAC_KEY_ID"), fileCfg.Audit.HMACKeyID, "audit-hmac-v1"),
	}

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

func readFileConfig(path string, explicit bool) (fileConfig, error) {
	var cfg fileConfig
	if strings.TrimSpace(path) == "" {
		return cfg, nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) && !explicit {
			return cfg, nil
		}
		return fileConfig{}, fmt.Errorf("config: file %q: %w", path, err)
	}

	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("config: parse %q: %w", path, err)
	}
	return cfg, nil
}

func resolvePositiveInt(cliValue int, envName string, fileValue int, defaultValue int) int {
	if cliValue > 0 {
		return cliValue
	}
	if envValue, ok := parsePositiveInt(os.Getenv(envName)); ok {
		return envValue
	}
	if fileValue > 0 {
		return fileValue
	}
	return defaultValue
}

func parsePositiveInt(raw string) (int, bool) {
	envValue := strings.TrimSpace(raw)
	if envValue == "" {
		return 0, false
	}
	n, err := strconv.Atoi(envValue)
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}

func resolveAllowedHosts(cliValue, envValue string, fileValue []string) []string {
	if strings.TrimSpace(cliValue) != "" {
		return splitCSV(cliValue)
	}
	if strings.TrimSpace(envValue) != "" {
		return splitCSV(envValue)
	}
	return trimList(fileValue)
}

func splitCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}
	return strings.Split(raw, ",")
}

func trimList(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
