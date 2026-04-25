package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"ai-anonymizer/internal/catalog"
	"ai-anonymizer/internal/config"
	"ai-anonymizer/internal/preflight"
	"ai-anonymizer/internal/scanner"
)

var doctorDeps = preflight.DoctorDeps{}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: anonym <doctor|scan|list>")
	}

	switch args[0] {
	case "doctor":
		return runDoctor(args[1:], stdout)
	case "scan":
		return runScan(args[1:], stdout)
	case "list":
		return runList(args[1:], stdout)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func runDoctor(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(stdout)

	cfgPath := fs.String("config", "config.example.yaml", "Path to config file")
	rawPath := fs.String("raw-path", "", "Absolute path to raw files")
	outputPath := fs.String("output-path", "", "Absolute output path")
	workspacePath := fs.String("workspace-path", "", "Absolute IDE workspace path")
	baseURL := fs.String("api-base-url", "", "API base URL")
	partnerID := fs.String("api-partner-id", "", "API partner ID")
	allowedHosts := fs.String("allowed-hosts", "", "Comma separated list of allowed hosts")
	requestTimeoutSec := fs.Int("request-timeout-sec", 5, "API preflight request timeout in seconds")
	maxFileSizeMB := fs.Int("max-file-size-mb", 25, "Maximum input file size in MB")
	maxParallelRuns := fs.Int("max-parallel-runs", 1, "Maximum parallel runs for v1")

	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.Load(config.LoadOptions{
		ConfigPath:        *cfgPath,
		RawPath:           *rawPath,
		OutputPath:        *outputPath,
		WorkspacePath:     *workspacePath,
		APIBaseURL:        *baseURL,
		APIPartnerID:      *partnerID,
		AllowedHostsCSV:   *allowedHosts,
		RequestTimeoutSec: *requestTimeoutSec,
		MaxFileSizeMB:     *maxFileSizeMB,
		MaxParallelRuns:   *maxParallelRuns,
	})
	if err != nil {
		return err
	}

	cfg.APIAllowedHosts = normalizeCSV(cfg.APIAllowedHosts)
	if err := preflight.ValidateDoctor(cfg, doctorDeps); err != nil {
		return err
	}

	fmt.Fprintln(stdout, "doctor: ok")
	return nil
}

func runScan(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(stdout)

	cfgPath := fs.String("config", "config.example.yaml", "Path to config file")
	rawPath := fs.String("raw-path", "", "Absolute path to raw files")
	outputPath := fs.String("output-path", "", "Absolute output path")
	workspacePath := fs.String("workspace-path", "", "Absolute IDE workspace path")
	baseURL := fs.String("api-base-url", "", "API base URL")
	partnerID := fs.String("api-partner-id", "", "API partner ID")
	allowedHosts := fs.String("allowed-hosts", "", "Comma separated list of allowed hosts")
	requestTimeoutSec := fs.Int("request-timeout-sec", 5, "API preflight request timeout in seconds")
	maxFileSizeMB := fs.Int("max-file-size-mb", 25, "Maximum input file size in MB")
	maxParallelRuns := fs.Int("max-parallel-runs", 1, "Maximum parallel runs for v1")

	if err := fs.Parse(args); err != nil {
		return err
	}
	cfg, err := config.Load(config.LoadOptions{
		ConfigPath:        *cfgPath,
		RawPath:           *rawPath,
		OutputPath:        *outputPath,
		WorkspacePath:     *workspacePath,
		APIBaseURL:        *baseURL,
		APIPartnerID:      *partnerID,
		AllowedHostsCSV:   *allowedHosts,
		RequestTimeoutSec: *requestTimeoutSec,
		MaxFileSizeMB:     *maxFileSizeMB,
		MaxParallelRuns:   *maxParallelRuns,
	})
	if err != nil {
		return err
	}
	cfg.APIAllowedHosts = normalizeCSV(cfg.APIAllowedHosts)
	if err := preflight.Validate(cfg); err != nil {
		return err
	}

	limitBytes := int64(cfg.MaxFileSizeMB) * 1024 * 1024
	items, err := scanner.DiscoverWithLimits(cfg.RawPath, scanner.Limits{
		MaxFileSizeBytes: limitBytes,
	})
	if err != nil {
		return err
	}
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save(items); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "scan: indexed %d file(s)\n", len(items))
	return nil
}

func runList(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(stdout)

	cfgPath := fs.String("config", "config.example.yaml", "Path to config file")
	rawPath := fs.String("raw-path", "", "Absolute path to raw files")
	outputPath := fs.String("output-path", "", "Absolute output path")
	workspacePath := fs.String("workspace-path", "", "Absolute IDE workspace path")
	baseURL := fs.String("api-base-url", "", "API base URL")
	partnerID := fs.String("api-partner-id", "", "API partner ID")
	allowedHosts := fs.String("allowed-hosts", "", "Comma separated list of allowed hosts")
	requestTimeoutSec := fs.Int("request-timeout-sec", 5, "API preflight request timeout in seconds")
	maxFileSizeMB := fs.Int("max-file-size-mb", 25, "Maximum input file size in MB")
	maxParallelRuns := fs.Int("max-parallel-runs", 1, "Maximum parallel runs for v1")

	if err := fs.Parse(args); err != nil {
		return err
	}
	cfg, err := config.Load(config.LoadOptions{
		ConfigPath:        *cfgPath,
		RawPath:           *rawPath,
		OutputPath:        *outputPath,
		WorkspacePath:     *workspacePath,
		APIBaseURL:        *baseURL,
		APIPartnerID:      *partnerID,
		AllowedHostsCSV:   *allowedHosts,
		RequestTimeoutSec: *requestTimeoutSec,
		MaxFileSizeMB:     *maxFileSizeMB,
		MaxParallelRuns:   *maxParallelRuns,
	})
	if err != nil {
		return err
	}
	cfg.APIAllowedHosts = normalizeCSV(cfg.APIAllowedHosts)
	if err := preflight.Validate(cfg); err != nil {
		return err
	}

	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	items, err := c.Load()
	if err != nil {
		return err
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Path < items[j].Path })

	if len(items) == 0 {
		fmt.Fprintln(stdout, "list: empty")
		return nil
	}
	for _, item := range items {
		fmt.Fprintf(stdout, "%s\t%s\t%s\n", item.ID, item.Path, item.Status)
	}
	return nil
}

func normalizeCSV(items []string) []string {
	var result []string
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
