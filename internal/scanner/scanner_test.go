package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverFindsRegularFiles(t *testing.T) {
	raw := t.TempDir()
	if err := os.MkdirAll(filepath.Join(raw, "nested"), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(raw, "nested", "doc.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	items, err := Discover(raw)
	if err != nil {
		t.Fatalf("discover failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Path != "nested/doc.txt" {
		t.Fatalf("unexpected path %q", items[0].Path)
	}
	if items[0].Status != "scanned" {
		t.Fatalf("unexpected status %q", items[0].Status)
	}
}

func TestHasADSWindowsOnly(t *testing.T) {
	if hasADS(`C:\safe\doc.txt`) {
		t.Fatal("expected no ADS for regular path")
	}
}

func TestDiscoverAppliesMaxFileSizeLimit(t *testing.T) {
	raw := t.TempDir()
	if err := os.WriteFile(filepath.Join(raw, "small.txt"), []byte("ok"), 0o644); err != nil {
		t.Fatalf("write small failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(raw, "big.txt"), []byte("1234"), 0o644); err != nil {
		t.Fatalf("write big failed: %v", err)
	}

	// 0 MB interpreted as explicit "disabled" is not used here; set hard tiny limit 3 bytes.
	items, err := DiscoverWithLimits(raw, Limits{MaxFileSizeBytes: 3})
	if err != nil {
		t.Fatalf("discover failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected only small file, got %d items", len(items))
	}
	if items[0].Path != "small.txt" {
		t.Fatalf("expected small.txt, got %q", items[0].Path)
	}
}
