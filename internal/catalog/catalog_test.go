package catalog

import (
	"path/filepath"
	"testing"
)

func TestFileCatalogSaveLoadRoundtrip(t *testing.T) {
	c := FileCatalog{Path: filepath.Join(t.TempDir(), ".anonym", "catalog.json")}
	input := []Item{
		{ID: "b", Path: "b.txt", Status: StatusScanned},
		{ID: "a", Path: "a.txt", Status: StatusScanned},
	}

	if err := c.Save(input); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	got, err := c.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
	if got[0].Path != "a.txt" || got[1].Path != "b.txt" {
		t.Fatalf("expected sorted items, got %#v", got)
	}
}
