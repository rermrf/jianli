package pdf

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// ErrPrintTimeout is returned when the /print page does not finish rendering
// (window.__printReady never becomes true) within the configured timeout.
// Handlers should map this to HTTP 504 Gateway Timeout.
var ErrPrintTimeout = errors.New("pdf: timed out waiting for print page")

const printReadyTimeout = 10 * time.Second

// Exporter renders the resume to a PDF by driving a headless Chrome instance
// to navigate to the locally-served /print React route and then printing the
// resulting page through Chrome's PrintToPDF.
type Exporter struct {
	browserPath string
	port        string
}

// NewExporter constructs an Exporter. browserPath is an optional explicit path
// to chrome/edge; an empty string falls back to a list of standard install
// locations. port is the local listening port of the same HTTP server that
// serves the SPA — chromedp will navigate to http://127.0.0.1:<port>/print.
func NewExporter(browserPath, port string) Exporter {
	return Exporter{
		browserPath: browserPath,
		port:        port,
	}
}

// ExportResume produces a PDF byte stream by navigating chromedp to the local
// /print route, waiting for window.__printReady === true, then invoking
// page.PrintToPDF with A4 settings and CSS-controlled margins.
func (e Exporter) ExportResume(ctx context.Context) ([]byte, error) {
	browserPath, err := resolveBrowserPath(e.browserPath, []string{
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
		`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
	})
	if err != nil {
		return nil, err
	}

	allocatorOptions := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(browserPath),
	)
	allocatorCtx, cancelAllocator := chromedp.NewExecAllocator(ctx, allocatorOptions...)
	defer cancelAllocator()
	browserCtx, cancelBrowser := chromedp.NewContext(allocatorCtx)
	defer cancelBrowser()

	host := net.JoinHostPort("127.0.0.1", e.port)
	url := fmt.Sprintf("http://%s/print", host)

	var pdfBytes []byte
	err = chromedp.Run(browserCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#print-root`, chromedp.ByID),
		chromedp.Poll(
			`window.__printReady === true`,
			nil,
			chromedp.WithPollingTimeout(printReadyTimeout),
		),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).
				WithPaperHeight(11.69).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithPreferCSSPageSize(true).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBytes = buf
			return nil
		}),
	)
	if err != nil {
		if isPollTimeout(err) {
			return nil, ErrPrintTimeout
		}
		return nil, err
	}

	return pdfBytes, nil
}

// isPollTimeout detects the timeout error returned by chromedp.Poll when the
// predicate never resolves within WithPollingTimeout. chromedp doesn't export
// a sentinel for this, so we string-match the variants we've observed:
//   - "waiting for function failed: timeout" (chromedp.Poll's normal timeout)
//   - "poll function timed out"              (older / alternate phrasing)
//   - "context deadline exceeded"            (outer ctx hit before Poll)
func isPollTimeout(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "waiting for function failed: timeout") ||
		strings.Contains(msg, "poll function timed out") ||
		strings.Contains(msg, "context deadline exceeded")
}

func resolveBrowserPath(configuredPath string, fallbacks []string) (string, error) {
	candidates := []string{}
	if configuredPath != "" {
		candidates = append(candidates, configuredPath)
	}
	candidates = append(candidates, fallbacks...)

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", errors.New("no supported browser executable found for pdf export")
}
