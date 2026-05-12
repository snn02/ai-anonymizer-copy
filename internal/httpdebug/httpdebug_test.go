package httpdebug

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-anonymizer/internal/config"
)

func TestLogDoesNothingWhenDisabled(t *testing.T) {
	t.Setenv("ANON_HTTP_DEBUG", "false")
	workspace := t.TempDir()
	cfg := config.Config{WorkspacePath: workspace}
	req, _ := http.NewRequest(http.MethodGet, "https://api.company.local/v1/tasks/anonymization_fields", nil)
	req.Header.Set("Authorization", "token-value")
	Log(cfg, "doctor", req, 200, nil)
	if _, err := os.Stat(filepath.Join(workspace, ".anonym", "http-debug.log")); !os.IsNotExist(err) {
		t.Fatalf("expected no debug log when disabled, got err=%v", err)
	}
}

func TestLogMasksAuthorization(t *testing.T) {
	t.Setenv("ANON_HTTP_DEBUG", "true")
	workspace := t.TempDir()
	cfg := config.Config{WorkspacePath: workspace}
	req, _ := http.NewRequest(http.MethodGet, "https://api.company.local/v1/tasks/anonymization_fields", nil)
	req.Header.Set("Authorization", "token-value-123456")
	req.Header.Set("partner-id", "70bd3a91-0000-0000-0000-000000000000")
	Log(cfg, "doctor", req, 200, nil)
	data, err := os.ReadFile(filepath.Join(workspace, ".anonym", "http-debug.log"))
	if err != nil {
		t.Fatalf("expected debug log file, got %v", err)
	}
	text := string(data)
	if strings.Contains(text, "token-value-123456") {
		t.Fatalf("authorization should be masked: %q", text)
	}
	if !strings.Contains(text, "authorization=toke***3456") {
		t.Fatalf("masked authorization not found: %q", text)
	}
}
