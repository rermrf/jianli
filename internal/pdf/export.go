package pdf

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"html/template"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type resumeView struct {
	Profile struct {
		Name     string `json:"name"`
		Title    string `json:"title"`
		Location string `json:"location"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	} `json:"profile"`
	Skills         []string `json:"skills"`
	WorkExperience []struct {
		Company     string   `json:"company"`
		Role        string   `json:"role"`
		StartDate   string   `json:"startDate"`
		EndDate     string   `json:"endDate"`
		Description []string `json:"description"`
	} `json:"workExperience"`
	Projects []struct {
		Name        string   `json:"name"`
		StartDate   string   `json:"startDate"`
		EndDate     string   `json:"endDate"`
		Description []string `json:"description"`
	} `json:"projects"`
}

type Exporter struct{}

func NewExporter() Exporter {
	return Exporter{}
}

func (Exporter) ExportResume(ctx context.Context, resume json.RawMessage) ([]byte, error) {
	var view resumeView
	if err := json.Unmarshal(resume, &view); err != nil {
		return nil, err
	}

	html, err := renderResumeHTML(view)
	if err != nil {
		return nil, err
	}

	browserCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var pdfBytes []byte
	encodedHTML := base64.StdEncoding.EncodeToString([]byte(html))
	if err := chromedp.Run(browserCtx,
		chromedp.Navigate("data:text/html;base64,"+encodedHTML),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			pdfBytes = buf
			return nil
		}),
	); err != nil {
		return nil, err
	}

	return pdfBytes, nil
}

func renderResumeHTML(view resumeView) (string, error) {
	const tpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<style>
body { font-family: Arial, sans-serif; padding: 32px; color: #1a1a1a; }
h1 { margin: 0 0 8px; font-size: 28px; }
h2 { margin: 28px 0 12px; font-size: 18px; }
p, li { font-size: 14px; line-height: 1.6; }
ul { padding-left: 18px; }
.meta { color: #666; margin-bottom: 16px; }
.block { margin-bottom: 16px; }
</style>
</head>
<body>
<h1>{{ .Profile.Name }}</h1>
<p class="meta">{{ .Profile.Title }} | {{ .Profile.Location }} | {{ .Profile.Phone }} | {{ .Profile.Email }}</p>
<h2>技能</h2>
<p>{{ range $index, $skill := .Skills }}{{ if $index }} / {{ end }}{{ $skill }}{{ end }}</p>
<h2>工作经历</h2>
{{ range .WorkExperience }}<div class="block"><strong>{{ .Company }} - {{ .Role }}</strong><p>{{ .StartDate }} - {{ .EndDate }}</p><ul>{{ range .Description }}<li>{{ . }}</li>{{ end }}</ul></div>{{ end }}
<h2>项目经历</h2>
{{ range .Projects }}<div class="block"><strong>{{ .Name }}</strong><p>{{ .StartDate }} - {{ .EndDate }}</p><ul>{{ range .Description }}<li>{{ . }}</li>{{ end }}</ul></div>{{ end }}
</body>
</html>`

	tmpl, err := template.New("resume").Parse(tpl)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, view); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
