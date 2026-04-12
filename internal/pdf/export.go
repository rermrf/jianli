package pdf

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type resumeView struct {
	Profile struct {
		AvatarDataURL template.URL `json:"avatarUrl"`
		Age           int          `json:"age"`
		Education     string       `json:"education"`
		Email         string       `json:"email"`
		Experience    string       `json:"experience"`
		Gender        string       `json:"gender"`
		Hometown      string       `json:"hometown"`
		Location      string       `json:"location"`
		Name          string       `json:"name"`
		Phone         string       `json:"phone"`
		Title         string       `json:"title"`
	} `json:"profile"`
	Awards []struct {
		Date  string `json:"date"`
		Title string `json:"title"`
	} `json:"awards"`
	Education []struct {
		Degree    string `json:"degree"`
		EndDate   string `json:"endDate"`
		Major     string `json:"major"`
		School    string `json:"school"`
		StartDate string `json:"startDate"`
	} `json:"education"`
	Projects []struct {
		Description []string `json:"description"`
		EndDate     string   `json:"endDate"`
		Name        string   `json:"name"`
		StartDate   string   `json:"startDate"`
		URL         string   `json:"url"`
	} `json:"projects"`
	Skills         []string `json:"skills"`
	WorkExperience []struct {
		Company     string   `json:"company"`
		Description []string `json:"description"`
		EndDate     string   `json:"endDate"`
		Role        string   `json:"role"`
		StartDate   string   `json:"startDate"`
	} `json:"workExperience"`
}

type Exporter struct {
	browserPath string
	uploadRoot  string
}

func NewExporter(browserPath string) Exporter {
	return Exporter{
		browserPath: browserPath,
		uploadRoot:  filepath.Clean("./data/uploads"),
	}
}

func (e Exporter) ExportResume(ctx context.Context, resume json.RawMessage) ([]byte, error) {
	view, err := buildResumeView(resume, e.uploadRoot)
	if err != nil {
		return nil, err
	}

	html, err := renderResumeHTML(view)
	if err != nil {
		return nil, err
	}

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

	var pdfBytes []byte
	if err := chromedp.Run(browserCtx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
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

func buildResumeView(resume json.RawMessage, uploadRoot string) (resumeView, error) {
	var view resumeView
	if err := json.Unmarshal(resume, &view); err != nil {
		return resumeView{}, err
	}

	if view.Profile.AvatarDataURL != "" {
		avatarDataURL, err := resolveAvatarDataURL(string(view.Profile.AvatarDataURL), uploadRoot)
		if err != nil {
			return resumeView{}, err
		}
		view.Profile.AvatarDataURL = template.URL(avatarDataURL)
	}

	return view, nil
}

func resolveAvatarDataURL(avatarURL, uploadRoot string) (string, error) {
	cleanURL := strings.TrimPrefix(avatarURL, "/")
	if !strings.HasPrefix(cleanURL, "uploads/") {
		return avatarURL, nil
	}

	filePath := filepath.Join(uploadRoot, strings.TrimPrefix(cleanURL, "uploads/"))
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(content)
	if contentType == "" {
		contentType = "image/png"
	}

	return "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(content), nil
}

func renderResumeHTML(view resumeView) (string, error) {
	const tpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<style>
body { font-family: "Inter", "Noto Sans CJK SC", "WenQuanYi Micro Hei", "PingFang SC", "Microsoft YaHei", Arial, sans-serif; padding: 32px; color: #1a1a1a; }
h1 { margin: 0; font-size: 30px; letter-spacing: 0.02em; }
h2 { margin: 28px 0 12px; font-size: 18px; }
p, li { font-size: 14px; line-height: 1.6; }
ul { padding-left: 18px; }
img.avatar { width: 84px; height: 84px; border-radius: 999px; object-fit: cover; }
.header { display: flex; gap: 20px; align-items: flex-start; }
.role { margin: 8px 0 0; font-size: 16px; font-weight: 600; color: #334155; }
.meta { color: #666; margin: 12px 0 0; }
.profile-facts { margin: 8px 0 0; color: #666; }
.block { margin-bottom: 16px; }
</style>
</head>
<body>
<div class="header">
{{ if .Profile.AvatarDataURL }}<img class="avatar" src="{{ .Profile.AvatarDataURL }}" />{{ end }}
<div>
<h1>{{ .Profile.Name }}</h1>
<p class="role">{{ .Profile.Title }}</p>
<p class="meta">所在地：{{ .Profile.Location }} | 手机号：{{ .Profile.Phone }} | 邮箱：{{ .Profile.Email }}</p>
<p class="profile-facts">{{ .Profile.Age }}岁 / {{ .Profile.Gender }} / {{ .Profile.Education }} / {{ .Profile.Experience }} / 籍贯：{{ .Profile.Hometown }}</p>
</div>
</div>
<h2>技能</h2>
<p>{{ range $index, $skill := .Skills }}{{ if $index }} / {{ end }}{{ $skill }}{{ end }}</p>
<h2>工作经历</h2>
{{ range .WorkExperience }}<div class="block"><strong>{{ .Company }} - {{ .Role }}</strong><p>{{ .StartDate }} - {{ .EndDate }}</p><ul>{{ range .Description }}<li>{{ . }}</li>{{ end }}</ul></div>{{ end }}
<h2>项目经历</h2>
{{ range .Projects }}<div class="block"><strong>{{ .Name }}</strong><p>{{ .StartDate }} - {{ .EndDate }}</p><ul>{{ range .Description }}<li>{{ . }}</li>{{ end }}</ul></div>{{ end }}
<h2>教育经历</h2>
{{ range .Education }}<div class="block"><strong>{{ .School }}</strong><p>{{ .StartDate }} - {{ .EndDate }}</p><p>{{ .Major }} · {{ .Degree }}</p></div>{{ end }}
<h2>荣誉奖项</h2>
{{ range .Awards }}<div class="block"><strong>{{ .Title }}</strong><p>{{ .Date }}</p></div>{{ end }}
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