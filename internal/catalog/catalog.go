package catalog

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	StatusScanned = "scanned"
)

type Item struct {
	ID     string `json:"id"`
	Path   string `json:"path"`
	Status string `json:"status"`
}

type FileCatalog struct {
	Path string
}

func (c FileCatalog) Save(items []Item) error {
	if err := os.MkdirAll(filepath.Dir(c.Path), 0o755); err != nil {
		return fmt.Errorf("catalog: mkdir: %w", err)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Path < items[j].Path })
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return fmt.Errorf("catalog: marshal: %w", err)
	}
	if err := os.WriteFile(c.Path, data, 0o644); err != nil {
		return fmt.Errorf("catalog: write: %w", err)
	}
	return nil
}

func (c FileCatalog) Load() ([]Item, error) {
	data, err := os.ReadFile(c.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Item{}, nil
		}
		return nil, fmt.Errorf("catalog: read: %w", err)
	}
	var items []Item
	if len(data) == 0 {
		return []Item{}, nil
	}
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("catalog: unmarshal: %w", err)
	}
	return items, nil
}

func BuildID(relativePath string) string {
	sum := sha256.Sum256([]byte(relativePath))
	return hex.EncodeToString(sum[:])[:12]
}
