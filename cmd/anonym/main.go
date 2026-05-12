package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"ai-anonymizer/internal/anonymizer"
	"ai-anonymizer/internal/catalog"
	"ai-anonymizer/internal/config"
	"ai-anonymizer/internal/preflight"
	"ai-anonymizer/internal/scanner"
	"ai-anonymizer/internal/secrets"
)

var doctorDeps = preflight.DoctorDeps{
	SecretChecker: secrets.NewOSSecretChecker(),
}

var runDeps = anonymizer.Deps{
	SecretGetter: secrets.NewOSSecretChecker(),
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: anonym <doctor|scan|list|run>")
	}

	switch args[0] {
	case "doctor":
		return runDoctor(args[1:], stdout)
	case "scan":
		return runScan(args[1:], stdout)
	case "list":
		return runList(args[1:], stdout)
	case "run":
		return runRun(args[1:], stdout)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func runDoctor(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(stdout)

	cfgPath := fs.String("config", "config.yaml", "Path to config file")
	outputPath := fs.String("output-path", "", "Absolute output path")
	workspacePath := fs.String("workspace-path", "", "Absolute IDE workspace path")
	mode := fs.String("mode", "", "Runtime mode (prod|mvp)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	configPathExplicit := wasFlagProvided(fs, "config")
	cfg, err := config.Load(config.LoadOptions{
		ConfigPath:         *cfgPath,
		ConfigPathExplicit: configPathExplicit,
		OutputPath:         *outputPath,
		WorkspacePath:      *workspacePath,
		RuntimeMode:        *mode,
	})
	if err != nil {
		return err
	}

	cfg.APIAllowedHosts = normalizeCSV(cfg.APIAllowedHosts)
	printProfileWarning(cfg, stdout)
	if err := preflight.ValidateDoctor(cfg, doctorDeps); err != nil {
		return err
	}

	fmt.Fprintln(stdout, "doctor: ok")
	return nil
}

func runScan(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(stdout)

	cfgPath := fs.String("config", "config.yaml", "Path to config file")
	outputPath := fs.String("output-path", "", "Absolute output path")
	workspacePath := fs.String("workspace-path", "", "Absolute IDE workspace path")
	mode := fs.String("mode", "", "Runtime mode (prod|mvp)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	configPathExplicit := wasFlagProvided(fs, "config")
	cfg, err := config.Load(config.LoadOptions{
		ConfigPath:         *cfgPath,
		ConfigPathExplicit: configPathExplicit,
		OutputPath:         *outputPath,
		WorkspacePath:      *workspacePath,
		RuntimeMode:        *mode,
	})
	if err != nil {
		return err
	}
	cfg.APIAllowedHosts = normalizeCSV(cfg.APIAllowedHosts)
	printProfileWarning(cfg, stdout)
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

	cfgPath := fs.String("config", "config.yaml", "Path to config file")
	outputPath := fs.String("output-path", "", "Absolute output path")
	workspacePath := fs.String("workspace-path", "", "Absolute IDE workspace path")
	mode := fs.String("mode", "", "Runtime mode (prod|mvp)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	configPathExplicit := wasFlagProvided(fs, "config")
	cfg, err := config.Load(config.LoadOptions{
		ConfigPath:         *cfgPath,
		ConfigPathExplicit: configPathExplicit,
		OutputPath:         *outputPath,
		WorkspacePath:      *workspacePath,
		RuntimeMode:        *mode,
	})
	if err != nil {
		return err
	}
	cfg.APIAllowedHosts = normalizeCSV(cfg.APIAllowedHosts)
	printProfileWarning(cfg, stdout)
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

func runRun(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(stdout)

	cfgPath := fs.String("config", "config.yaml", "Path to config file")
	outputPath := fs.String("output-path", "", "Absolute output path")
	workspacePath := fs.String("workspace-path", "", "Absolute IDE workspace path")
	mode := fs.String("mode", "", "Runtime mode (prod|mvp)")

	target := ""
	parseArgs := args
	if len(args) > 0 && !strings.HasPrefix(strings.TrimSpace(args[0]), "-") {
		target = strings.TrimSpace(args[0])
		parseArgs = args[1:]
	}

	if err := fs.Parse(parseArgs); err != nil {
		return err
	}
	if target == "" && fs.NArg() > 0 {
		target = fs.Arg(0)
	}
	if target == "" {
		return errors.New("run: usage anonym run <id|path>")
	}

	configPathExplicit := wasFlagProvided(fs, "config")
	cfg, err := config.Load(config.LoadOptions{
		ConfigPath:         *cfgPath,
		ConfigPathExplicit: configPathExplicit,
		OutputPath:         *outputPath,
		WorkspacePath:      *workspacePath,
		RuntimeMode:        *mode,
	})
	if err != nil {
		return err
	}
	cfg.APIAllowedHosts = normalizeCSV(cfg.APIAllowedHosts)
	printProfileWarning(cfg, stdout)

	result, err := anonymizer.Execute(cfg, target, runDeps)
	if err != nil {
		var ambiguous anonymizer.AmbiguousMatchError
		if errors.As(err, &ambiguous) {
			fmt.Fprintf(stdout, "run: multiple matches for %q:\n", target)
			for _, candidate := range ambiguous.Candidates {
				fmt.Fprintf(stdout, "%s\t%s\t%s\n", candidate.ID, candidate.Path, candidate.Status)
			}
			return errors.New("run: specify explicit id")
		}
		return err
	}

	fmt.Fprintf(stdout, "run: succeeded\t%s\t%s\n", result.Item.ID, result.OutputPath)
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

func wasFlagProvided(fs *flag.FlagSet, name string) bool {
	found := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func printProfileWarning(cfg config.Config, stdout io.Writer) {
	if strings.EqualFold(strings.TrimSpace(cfg.RuntimeMode), "mvp") {
		fmt.Fprintln(stdout, "warning: mvp insecure mode is active")
	}
	if strings.TrimSpace(cfg.LegacyModeSource) != "" {
		fmt.Fprintf(stdout, "warning: legacy profile mapping used (%s); please migrate to runtime.mode\n", cfg.LegacyModeSource)
	}
}
