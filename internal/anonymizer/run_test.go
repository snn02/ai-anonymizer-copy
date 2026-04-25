package anonymizer

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-anonymizer/internal/catalog"
	"ai-anonymizer/internal/config"
)

type fakeSecretGetter struct {
	value string
	err   error
}

func (f fakeSecretGetter) GetSecret(partnerID string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.value, nil
}

func TestResolveByID(t *testing.T) {
	items := []catalog.Item{
		{ID: "abc", Path: "in/a.txt", Status: catalog.StatusScanned},
		{ID: "def", Path: "in/b.txt", Status: catalog.StatusScanned},
	}

	idx, item, err := resolve(items, "def")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx != 1 || item.ID != "def" {
		t.Fatalf("unexpected match: idx=%d item=%+v", idx, item)
	}
}

func TestResolveFuzzyAmbiguous(t *testing.T) {
	items := []catalog.Item{
		{ID: "abc", Path: "in/invoice-april.txt", Status: catalog.StatusScanned},
		{ID: "def", Path: "in/invoice-may.txt", Status: catalog.StatusScanned},
	}

	_, _, err := resolve(items, "invoice")
	if err == nil {
		t.Fatal("expected ambiguous match error")
	}
	var ambiguous AmbiguousMatchError
	if !errors.As(err, &ambiguous) {
		t.Fatalf("expected AmbiguousMatchError, got %T", err)
	}
	if len(ambiguous.Candidates) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(ambiguous.Candidates))
	}
}

func TestAdaptResponseJSONVariants(t *testing.T) {
	got1, err := adaptResponse("application/json", []byte(`{"result":"ANON"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got1) != "ANON" {
		t.Fatalf("unexpected result %q", string(got1))
	}

	got2, err := adaptResponse("application/json", []byte(`{"result":{"content":"ANON2"}}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got2) != "ANON2" {
		t.Fatalf("unexpected result %q", string(got2))
	}
}

func TestExecuteRejectsNonHTTPS(t *testing.T) {
	cfg := baseConfig(t)
	cfg.APIBaseURL = "http://api.company.local"

	_, err := Execute(cfg, "x", Deps{
		SecretGetter: fakeSecretGetter{value: "token"},
	})
	if err == nil || !strings.Contains(err.Error(), "must use https") {
		t.Fatalf("expected https policy error, got %v", err)
	}
}

func TestExecuteRejectsHostOutsideAllowlist(t *testing.T) {
	cfg := baseConfig(t)
	cfg.APIAllowedHosts = []string{"allowed.company.local"}

	_, err := Execute(cfg, "x", Deps{
		SecretGetter: fakeSecretGetter{value: "token"},
	})
	if err == nil || !strings.Contains(err.Error(), "not in security.allowed_hosts") {
		t.Fatalf("expected allowlist policy error, got %v", err)
	}
}

func TestExecuteRejectsMissingSecretFallback(t *testing.T) {
	cfg := baseConfig(t)

	_, err := Execute(cfg, "x", Deps{
		SecretGetter: fakeSecretGetter{value: ""},
	})
	if err == nil || !strings.Contains(err.Error(), "was not found in OS secret store") {
		t.Fatalf("expected missing secret error, got %v", err)
	}
}

func TestExecuteRejectsWhenFileExceedsMaxSize(t *testing.T) {
	cfg := baseConfig(t)
	cfg.MaxFileSizeMB = 1

	if err := os.MkdirAll(filepath.Join(cfg.RawPath, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	largeRel := "in/large.bin"
	largePath := filepath.Join(cfg.RawPath, filepath.FromSlash(largeRel))
	if err := os.WriteFile(largePath, make([]byte, 2*1024*1024), 0o644); err != nil {
		t.Fatal(err)
	}
	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	if err := c.Save([]catalog.Item{
		{ID: catalog.BuildID(largeRel), Path: largeRel, Status: catalog.StatusScanned},
	}); err != nil {
		t.Fatal(err)
	}

	_, err := Execute(cfg, catalog.BuildID(largeRel), Deps{
		SecretGetter: fakeSecretGetter{value: "token"},
	})
	if err == nil || !strings.Contains(err.Error(), "exceeds max_file_size_mb") {
		t.Fatalf("expected size-limit error, got %v", err)
	}
}

func baseConfig(t *testing.T) config.Config {
	t.Helper()

	raw := filepath.Join(t.TempDir(), "raw")
	workspace := filepath.Join(t.TempDir(), "workspace")
	output := filepath.Join(workspace, "out")
	if err := os.MkdirAll(raw, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}

	return config.Config{
		RawPath:           raw,
		OutputPath:        output,
		WorkspacePath:     workspace,
		CatalogPath:       filepath.Join(workspace, ".anonym", "catalog.json"),
		APIBaseURL:        "https://api.company.local",
		APIPartnerID:      "partner-1",
		APIAllowedHosts:   []string{"api.company.local"},
		RequestTimeoutSec: 5,
		MaxFileSizeMB:     25,
		MaxParallelRuns:   1,
	}
}
