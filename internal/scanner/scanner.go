package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"

	"ai-anonymizer/internal/catalog"
)

func Discover(rawRoot string) ([]catalog.Item, error) {
	rootAbs, err := filepath.Abs(rawRoot)
	if err != nil {
		return nil, fmt.Errorf("scan: invalid raw_path: %w", err)
	}
	rootAbs = filepath.Clean(rootAbs)

	items := make([]catalog.Item, 0)
	err = filepath.WalkDir(rootAbs, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			return nil
		}
		if d.Type()&fs.ModeSymlink != 0 {
			return nil
		}

		cleanPath := filepath.Clean(path)
		if !isInsideRoot(rootAbs, cleanPath) {
			return nil
		}
		if hasADS(cleanPath) {
			return nil
		}

		rel, err := filepath.Rel(rootAbs, cleanPath)
		if err != nil {
			return nil
		}
		if rel == "." || strings.HasPrefix(rel, "..") {
			return nil
		}
		rel = filepath.ToSlash(rel)
		items = append(items, catalog.Item{
			ID:     catalog.BuildID(rel),
			Path:   rel,
			Status: catalog.StatusScanned,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan: walk failed: %w", err)
	}
	return items, nil
}

func isInsideRoot(root, path string) bool {
	root = strings.ToLower(filepath.Clean(root))
	path = strings.ToLower(filepath.Clean(path))
	return path == root || strings.HasPrefix(path, root+string(filepath.Separator))
}

func hasADS(path string) bool {
	if runtime.GOOS != "windows" {
		return false
	}
	cleanPath := filepath.Clean(path)
	volume := filepath.VolumeName(cleanPath)
	withoutVolume := strings.TrimPrefix(cleanPath, volume)
	return strings.Contains(withoutVolume, ":")
}
