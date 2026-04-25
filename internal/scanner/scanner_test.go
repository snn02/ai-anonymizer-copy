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
