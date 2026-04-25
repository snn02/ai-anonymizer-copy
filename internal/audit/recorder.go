package audit

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Entry struct {
	Timestamp time.Time `json:"ts"`
	Event     string    `json:"event"`
	Status    string    `json:"status"`
	FileID    string    `json:"file_id"`
	Error     string    `json:"error,omitempty"`
}

type Recorder struct {
	Path          string
	RetentionDays int
	HMACKey       string
	Now           func() time.Time
}

func (r Recorder) Record(event string, status string, rawFileID string, errText string) error {
	if strings.TrimSpace(r.Path) == "" {
		return errors.New("audit: path is required")
	}
	if r.RetentionDays <= 0 {
		return fmt.Errorf("audit: retention_days must be positive, got %d", r.RetentionDays)
	}
	if strings.TrimSpace(r.HMACKey) == "" {
		return errors.New("audit: hmac key is required")
	}
	now := time.Now().UTC()
	if r.Now != nil {
		now = r.Now().UTC()
	}

	keepAfter := now.AddDate(0, 0, -r.RetentionDays)
	existing, err := readEntries(r.Path)
	if err != nil {
		return err
	}
	filtered := make([]Entry, 0, len(existing)+1)
	for _, e := range existing {
		if e.Timestamp.IsZero() || e.Timestamp.Before(keepAfter) {
			continue
		}
		filtered = append(filtered, e)
	}

	filtered = append(filtered, Entry{
		Timestamp: now,
		Event:     strings.TrimSpace(event),
		Status:    strings.TrimSpace(status),
		FileID:    hashFileID(rawFileID, r.HMACKey),
		Error:     sanitizeError(errText),
	})
	return writeEntries(r.Path, filtered)
}

func readEntries(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, fmt.Errorf("audit: open: %w", err)
	}
	defer f.Close()

	entries := make([]Entry, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("audit: scan: %w", err)
	}
	return entries, nil
}

func writeEntries(path string, entries []Entry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("audit: mkdir: %w", err)
	}

	lines := make([]string, 0, len(entries))
	for _, e := range entries {
		payload, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("audit: marshal: %w", err)
		}
		lines = append(lines, string(payload))
	}

	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n"
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return fmt.Errorf("audit: write: %w", err)
	}
	return nil
}

func hashFileID(rawFileID string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	_, _ = h.Write([]byte(strings.TrimSpace(rawFileID)))
	return hex.EncodeToString(h.Sum(nil))
}

func sanitizeError(errText string) string {
	text := strings.TrimSpace(errText)
	if text == "" {
		return ""
	}
	if len(text) > 120 {
		return text[:120]
	}
	return text
}
