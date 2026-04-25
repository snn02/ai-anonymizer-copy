package anonymizer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"ai-anonymizer/internal/catalog"
	"ai-anonymizer/internal/config"
	"ai-anonymizer/internal/preflight"
)

type SecretGetter interface {
	GetSecret(partnerID string) (string, error)
}

type Deps struct {
	HTTPClient   *http.Client
	SecretGetter SecretGetter
}

type Result struct {
	Item       catalog.Item
	OutputPath string
}

type AmbiguousMatchError struct {
	Query      string
	Candidates []catalog.Item
}

var reparsePointCheck = isReparsePoint

func (e AmbiguousMatchError) Error() string {
	return fmt.Sprintf("run: multiple matches for %q", e.Query)
}

func Execute(cfg config.Config, target string, deps Deps) (Result, error) {
	if err := preflight.ValidateRuntime(cfg); err != nil {
		return Result{}, err
	}
	if strings.TrimSpace(target) == "" {
		return Result{}, errors.New("run: target <id|path> is required")
	}
	if deps.SecretGetter == nil {
		return Result{}, errors.New("run: secret getter is not configured")
	}

	secret, err := deps.SecretGetter.GetSecret(cfg.APIPartnerID)
	if err != nil {
		return Result{}, fmt.Errorf("run: secret store read failed: %w", err)
	}
	if strings.TrimSpace(secret) == "" {
		return Result{}, errors.New("run: api secret was not found in OS secret store")
	}

	c := catalog.FileCatalog{Path: cfg.CatalogPath}
	items, err := c.Load()
	if err != nil {
		return Result{}, err
	}

	idx, selected, err := resolve(items, target)
	if err != nil {
		return Result{}, err
	}

	absPath := filepath.Join(cfg.RawPath, filepath.FromSlash(selected.Path))
	if err := validateRunFile(cfg.RawPath, absPath, cfg.MaxFileSizeMB, cfg.MaxPages); err != nil {
		return Result{}, err
	}

	items[idx].Status = catalog.StatusSent
	if err := c.Save(items); err != nil {
		return Result{}, err
	}

	sanitized, err := callFileAnonymization(cfg, absPath, secret, deps.HTTPClient)
	if err != nil {
		items[idx].Status = catalog.StatusFailed
		_ = c.Save(items)
		return Result{}, err
	}

	outPath, err := writeOutput(cfg.OutputPath, selected, sanitized)
	if err != nil {
		items[idx].Status = catalog.StatusFailed
		_ = c.Save(items)
		return Result{}, err
	}

	items[idx].Status = catalog.StatusSucceeded
	if err := c.Save(items); err != nil {
		return Result{}, err
	}
	return Result{
		Item:       items[idx],
		OutputPath: outPath,
	}, nil
}

func resolve(items []catalog.Item, target string) (int, catalog.Item, error) {
	normalizedTarget := normalizeMatchValue(target)

	for i, item := range items {
		if item.ID == strings.TrimSpace(target) {
			return i, item, nil
		}
	}

	for i, item := range items {
		if normalizeMatchValue(item.Path) == normalizedTarget {
			return i, item, nil
		}
	}

	matches := make([]int, 0)
	for i, item := range items {
		normPath := normalizeMatchValue(item.Path)
		if strings.Contains(normPath, normalizedTarget) || strings.Contains(strings.ToLower(item.ID), strings.ToLower(strings.TrimSpace(target))) {
			matches = append(matches, i)
		}
	}
	if len(matches) == 1 {
		i := matches[0]
		return i, items[i], nil
	}
	if len(matches) > 1 {
		candidates := make([]catalog.Item, 0, len(matches))
		for _, i := range matches {
			candidates = append(candidates, items[i])
		}
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Path < candidates[j].Path })
		return -1, catalog.Item{}, AmbiguousMatchError{
			Query:      target,
			Candidates: candidates,
		}
	}

	return -1, catalog.Item{}, fmt.Errorf("run: no matches for %q", target)
}

func normalizeMatchValue(s string) string {
	return strings.ToLower(filepath.ToSlash(filepath.Clean(strings.TrimSpace(s))))
}

func validateRunFile(rawPath, filePath string, maxFileSizeMB int, maxPages int) error {
	rootAbs, err := filepath.Abs(rawPath)
	if err != nil {
		return fmt.Errorf("run: invalid raw_path: %w", err)
	}
	fileAbs, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("run: invalid file path: %w", err)
	}

	root := strings.ToLower(filepath.Clean(rootAbs))
	file := strings.ToLower(filepath.Clean(fileAbs))
	if file != root && !strings.HasPrefix(file, root+string(filepath.Separator)) {
		return errors.New("run: selected file is outside raw_path")
	}
	if runtime.GOOS == "windows" {
		if hasADSPath(fileAbs) {
			return errors.New("run: windows ADS paths are not allowed")
		}
		if err := validateNoReparsePoints(rootAbs, fileAbs); err != nil {
			return err
		}
	}

	info, err := os.Lstat(fileAbs)
	if err != nil {
		return fmt.Errorf("run: selected file is not accessible: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return errors.New("run: symlink files are not allowed")
	}
	if !info.Mode().IsRegular() {
		return errors.New("run: selected path must be a regular file")
	}

	limitBytes := int64(maxFileSizeMB) * 1024 * 1024
	if maxFileSizeMB > 0 && info.Size() > limitBytes {
		return fmt.Errorf("run: file exceeds max_file_size_mb (%d)", maxFileSizeMB)
	}

	if err := validateMaxPages(fileAbs, maxPages); err != nil {
		return err
	}
	return nil
}

func validateMaxPages(filePath string, maxPages int) error {
	if maxPages <= 0 {
		return fmt.Errorf("run: max_pages must be positive, got %d", maxPages)
	}
	pages, supported, err := detectPages(filePath)
	if err != nil {
		return fmt.Errorf("run: page-count check failed: %w", err)
	}
	if !supported {
		return nil
	}
	if pages > maxPages {
		return fmt.Errorf("run: document exceeds max_pages (%d)", maxPages)
	}
	return nil
}

func detectPages(filePath string) (int, bool, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".txt", ".md", ".csv", ".json", ".xml", ".log":
		data, err := os.ReadFile(filePath)
		if err != nil {
			return 0, true, err
		}
		pages := 1
		if len(data) > 0 {
			pages += bytes.Count(data, []byte("\f"))
		}
		return pages, true, nil
	default:
		return 0, false, nil
	}
}

func hasADSPath(path string) bool {
	if runtime.GOOS != "windows" {
		return false
	}
	clean := filepath.Clean(path)
	volume := filepath.VolumeName(clean)
	withoutVolume := strings.TrimPrefix(clean, volume)
	return strings.Contains(withoutVolume, ":")
}

func validateNoReparsePoints(rootAbs, fileAbs string) error {
	rel, err := filepath.Rel(rootAbs, fileAbs)
	if err != nil {
		return fmt.Errorf("run: cannot evaluate path policy: %w", err)
	}
	current := rootAbs
	parts := strings.Split(rel, string(filepath.Separator))
	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		isReparse, err := reparsePointCheck(current)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("run: reparse-point check failed: %w", err)
		}
		if isReparse {
			return errors.New("run: reparse points are not allowed in selected path")
		}
	}
	return nil
}

func callFileAnonymization(cfg config.Config, filePath string, secret string, httpClient *http.Client) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("run: open file failed: %w", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("run: create multipart file failed: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("run: read file failed: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("run: finalize multipart failed: %w", err)
	}

	client := httpClient
	if client == nil {
		timeout := time.Duration(cfg.RequestTimeoutSec) * time.Second
		if timeout <= 0 {
			timeout = 5 * time.Second
		}
		client = &http.Client{Timeout: timeout}
	}

	endpoint := strings.TrimRight(cfg.APIBaseURL, "/") + "/v1/tasks/file_anonymization"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, &body)
	if err != nil {
		return nil, fmt.Errorf("run: build api request failed: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", secret)
	req.Header.Set("partner-id", cfg.APIPartnerID)
	req.Header.Set("fields", "anonymizer")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("run: api request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("run: read api response failed: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("run: api returned status %d", resp.StatusCode)
	}
	return adaptResponse(resp.Header.Get("Content-Type"), respBody)
}

func adaptResponse(contentType string, responseBody []byte) ([]byte, error) {
	if !strings.Contains(strings.ToLower(contentType), "application/json") {
		return responseBody, nil
	}

	var payload map[string]any
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return nil, fmt.Errorf("run: invalid json response: %w", err)
	}

	if value, ok := payload["result"].(string); ok {
		return []byte(value), nil
	}
	if resultMap, ok := payload["result"].(map[string]any); ok {
		if value, ok := resultMap["content"].(string); ok {
			return []byte(value), nil
		}
		if value, ok := resultMap["text"].(string); ok {
			return []byte(value), nil
		}
	}
	if value, ok := payload["content"].(string); ok {
		return []byte(value), nil
	}

	return nil, errors.New("run: unsupported api response format")
}

func writeOutput(outputRoot string, item catalog.Item, data []byte) (string, error) {
	ext := filepath.Ext(item.Path)
	filename := item.ID + ".anonymized" + ext
	outPath := filepath.Join(outputRoot, filename)

	if _, err := os.Stat(outPath); err == nil {
		return "", errors.New("run: output file already exists")
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("run: output file check failed: %w", err)
	}
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		return "", fmt.Errorf("run: write output failed: %w", err)
	}
	return outPath, nil
}
