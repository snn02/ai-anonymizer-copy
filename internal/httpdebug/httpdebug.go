package httpdebug

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ai-anonymizer/internal/config"
)

func Enabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("ANON_HTTP_DEBUG")))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func Log(cfg config.Config, scope string, req *http.Request, statusCode int, reqErr error) {
	if !Enabled() || req == nil {
		return
	}
	logPath := filepath.Join(cfg.WorkspacePath, ".anonym", "http-debug.log")
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()

	errText := "-"
	if reqErr != nil {
		errText = reqErr.Error()
	}
	auth := mask(req.Header.Get("Authorization"))
	partnerID := strings.TrimSpace(req.Header.Get("partner-id"))
	_, _ = fmt.Fprintf(
		f,
		"%s scope=%s method=%s url=%s status=%d authorization=%s partner-id=%s err=%s\n",
		time.Now().UTC().Format(time.RFC3339),
		scope,
		req.Method,
		req.URL.String(),
		statusCode,
		auth,
		partnerID,
		errText,
	)
}

func mask(value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return "<empty>"
	}
	if len(v) <= 8 {
		return "***"
	}
	return v[:4] + "***" + v[len(v)-4:]
}
