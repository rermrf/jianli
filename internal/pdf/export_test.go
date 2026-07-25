package pdf

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// findInstalledBrowser returns the path of an installed browser found via the
// same fallback list ExportResume uses, or "" if none is available. Used by
// integration tests to skip cleanly when chromedp can't run.
func findInstalledBrowser(t *testing.T) string {
	t.Helper()
	path, err := resolveBrowserPath("", []string{
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
		`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
	})
	if err != nil {
		return ""
	}
	return path
}

// startPrintStubServer stands up a local HTTP server that serves a minimal
// /print page satisfying the chromedp wait conditions used by ExportResume:
// a #print-root element and window.__printReady=true. Returns the listening
// port and a cleanup func.
func startPrintStubServer(t *testing.T) (string, func()) {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/print", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body>
<div id="print-root">stub print page</div>
<script>window.__printReady = true;</script>
</body>
</html>`))
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen() returned error: %v", err)
	}
	server := &http.Server{Handler: mux, ReadHeaderTimeout: 3 * time.Second}
	go func() { _ = server.Serve(listener) }()

	addr := listener.Addr().(*net.TCPAddr)
	port := strconv.Itoa(addr.Port)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}
	return port, cleanup
}

// startSlowPrintServer serves a /print page that never sets window.__printReady,
// so chromedp.Poll always times out. Used to assert the timeout is mapped to
// ErrPrintTimeout.
func startSlowPrintServer(t *testing.T) (string, func()) {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/print", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!DOCTYPE html>
<html><head><meta charset="utf-8"></head>
<body><div id="print-root">never ready</div></body>
</html>`))
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen() returned error: %v", err)
	}
	server := &http.Server{Handler: mux, ReadHeaderTimeout: 3 * time.Second}
	go func() { _ = server.Serve(listener) }()

	addr := listener.Addr().(*net.TCPAddr)
	port := strconv.Itoa(addr.Port)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}
	return port, cleanup
}

func TestExportResumeProducesPDFBytes(t *testing.T) {
	browserPath := findInstalledBrowser(t)
	if browserPath == "" {
		t.Skip("no chrome/edge installed; skipping chromedp integration test")
	}

	port, cleanup := startPrintStubServer(t)
	defer cleanup()

	exporter := NewExporter(browserPath, port)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pdfBytes, err := exporter.ExportResume(ctx)
	if err != nil {
		t.Fatalf("ExportResume() returned error: %v", err)
	}
	if len(pdfBytes) == 0 {
		t.Fatal("expected non-empty pdf bytes")
	}
	if !strings.HasPrefix(string(pdfBytes[:5]), "%PDF-") {
		end := 5
		if len(pdfBytes) < end {
			end = len(pdfBytes)
		}
		t.Fatalf("expected pdf bytes to start with %%PDF-, got %q", string(pdfBytes[:end]))
	}
}

func TestExportResumeReturnsErrPrintTimeoutWhenPageNeverReady(t *testing.T) {
	browserPath := findInstalledBrowser(t)
	if browserPath == "" {
		t.Skip("no chrome/edge installed; skipping chromedp integration test")
	}

	port, cleanup := startSlowPrintServer(t)
	defer cleanup()

	exporter := NewExporter(browserPath, port)

	// We need an outer ctx slightly longer than the internal 10s Poll timeout
	// so the Poll itself fires (not the outer context).
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := exporter.ExportResume(ctx)
	if !errors.Is(err, ErrPrintTimeout) {
		t.Fatalf("expected ErrPrintTimeout, got %v", err)
	}
}

func TestResolveBrowserPathPrefersConfiguredPath(t *testing.T) {
	configuredPath := filepath.Join(t.TempDir(), "chrome.exe")
	if err := os.WriteFile(configuredPath, []byte("stub"), 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	path, err := resolveBrowserPath(configuredPath, []string{"missing.exe"})
	if err != nil {
		t.Fatalf("resolveBrowserPath() returned error: %v", err)
	}

	if path != configuredPath {
		t.Fatalf("expected configured path %q, got %q", configuredPath, path)
	}
}

func TestResolveBrowserPathReturnsErrorWhenNoneFound(t *testing.T) {
	_, err := resolveBrowserPath("", []string{
		filepath.Join(t.TempDir(), "missing-1.exe"),
		filepath.Join(t.TempDir(), "missing-2.exe"),
	})
	if err == nil {
		t.Fatal("expected error when no browser found, got nil")
	}
}

func TestIsPollTimeoutDetectsKnownMessages(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"poll timeout (chromedp Poll)", errors.New("waiting for function failed: timeout"), true},
		{"poll timeout (legacy)", errors.New("poll function timed out"), true},
		{"context deadline", errors.New("context deadline exceeded"), true},
		{"unrelated", errors.New("something else"), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isPollTimeout(tc.err); got != tc.want {
				t.Fatalf("isPollTimeout(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}
