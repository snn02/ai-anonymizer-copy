package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRecordPrunesByRetentionAndKeepsSafePayload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".anonym", "audit.log")
	old := `{"ts":"2026-01-01T00:00:00Z","event":"run","status":"failed","file_id":"old"}`
	recent := `{"ts":"2026-04-25T00:00:00Z","event":"run","status":"succeeded","file_id":"recent"}`
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(old+"\n"+recent+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	rec := Recorder{
		Path:          path,
		RetentionDays: 7,
		HMACKey:       "audit-secret",
		Now: func() time.Time {
			return time.Date(2026, 4, 26, 12, 0, 0, 0, time.UTC)
		},
	}
	if err := rec.Record("run", "succeeded", "abc123", "api upstream timeout"); err != nil {
		t.Fatalf("record failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read audit failed: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 entries after retention, got %d", len(lines))
	}
	if strings.Contains(lines[0], `"file_id":"old"`) {
		t.Fatalf("old entry should be pruned, got %q", lines[0])
	}
	if !strings.Contains(lines[0], `"file_id":"recent"`) {
		t.Fatalf("recent entry should stay, got %q", lines[0])
	}
	if strings.Contains(lines[1], "abc123") {
		t.Fatalf("raw file id leaked in audit line: %q", lines[1])
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(lines[1]), &payload); err != nil {
		t.Fatalf("invalid json payload: %v", err)
	}
	if payload["event"] != "run" || payload["status"] != "succeeded" {
		t.Fatalf("unexpected event payload: %#v", payload)
	}
	if payload["error"] != "api upstream timeout" {
		t.Fatalf("expected sanitized error text, got %#v", payload["error"])
	}
	if got, _ := payload["file_id"].(string); len(got) != 64 {
		t.Fatalf("expected sha256 hmac hex (64 chars), got %q", got)
	}
}
