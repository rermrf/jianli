# 打印预览与 PDF 样式优化设计

- 日期：2026-05-03
- 状态：草稿，待评审

## 背景与问题

当前简历项目存在三套互相独立的渲染路径：

1. **在线展示页 (`/`)** —— `ResumePage` 通过 `ResumeDesktopLayout` / `ResumeMobileLayout` 渲染，使用 `SectionCard` / `ProfileCard` / `SkillSection` / `TimelineSection` 组件，桌面端为「左侧窄栏 + 右侧主栏」的双栏布局，视觉成熟。
2. **打印预览页 (`/print`)** —— `PrintPage` 通过 `PrintResume` 单文件渲染，单栏、无卡片、章节顺序与在线不同（技能 → 教育 → 工作 → 项目 → 荣誉），样式朴素。
3. **PDF 导出 (`GET /api/resume/pdf`)** —— `internal/pdf/export.go` 完全独立的 Go HTML 模板 + 内联 CSS，章节顺序又不同（技能 → 工作 → 项目 → 教育 → 荣誉），头像通过 base64 内嵌。

三套实现没有任何代码共享，必然出现「漂移」：每改一处就要在另外两处同步，每次改动都有遗漏风险。用户反馈：

- 打印预览与 PDF 输出和在线页不对齐（章节顺序、视觉风格不一致）
- 打印预览与 PDF 样式简陋（黑白 + 朴素列表，缺少视觉层级）

## 目标

- **结构对齐**：打印预览和 PDF 与在线展示包含相同的章节、相同的数据。章节出现顺序按打印场景优化（不要求与在线页逐字相同）。
- **打印优化样式**：在 A4 纸面上呈现专业、面向 HR 的视觉效果。
- **架构统一**：消除三套渲染路径，使打印预览和 PDF 共用同一份 React 组件。PDF 由后端 chromedp 打开本地 `/print` 路由 + 浏览器原生「打印为 PDF」生成。

## 非目标

- 不做客户端 PDF 生成（jsPDF / html2pdf 等）。
- 不做多主题/多模板切换，本次只做一套版式。
- 不要求在线展示页和打印页视觉 100% 一致，二者按媒介各自演化。
- 不做服务端 PDF 缓存。
- 不修复"现有 `export_test.go` 在无 Chrome 的 CI 上跑不动"这一既存问题。

## 决策摘要

| 维度 | 选择 |
|---|---|
| 对齐目标 | 结构对齐 + 打印优化样式 |
| 架构 | 一套实现：后端 chromedp 打开本地 `/print` React 页打印为 PDF |
| 版式 | 单栏顺序阅读 |
| 章节顺序 | 个人信息 → 个人技能 → 工作经历 → 项目经历 → 教育经历 → 荣誉奖项 |
| 视觉调性 | 黑白为主，少量品牌色点缀 |
| 分页策略 | 智能分页（避免单个条目跨页拆分） |

## 整体架构

```
浏览器访问 /print   ─────────────►  React 渲染最终样式
                                          ▲
                                          │ 同一份组件
后端 GET /api/resume/pdf                   │
   │                                       │
   ├─► chromedp 打开 http://127.0.0.1:<port>/print?pdfMode=1
   │       │
   │       ▼
   │   React 加载 → fetch /api/resume → 渲染完毕 → window.__printReady = true
   │       │
   │       ▼
   │   chromedp 等待 sentinel（Poll window.__printReady）
   │       │
   │       ▼
   │   chromedp PrintToPDF（A4 / printBackground / preferCSSPageSize）
   ▼
返回 PDF 字节流
```

可行性依据（基于现状）：

- `GET /api/resume` 已是公开接口，无需鉴权。
- `/print` 路由已是公开路由（`web/src/app/router.tsx`）。
- Go server 已通过 `router.NoRoute(serveFrontend)` 提供 SPA 静态文件，前后端同进程。
- 头像 `/uploads/...` 是公开静态资源，chromedp 可直接 HTTP 取，无需 base64 内嵌。

## `/print` 页视觉设计

### 页面骨架

单栏布局，章节由上至下：

1. Header（个人信息）
2. 个人技能
3. 工作经历
4. 项目经历
5. 教育经历
6. 荣誉奖项

### 样式规则

**页面容器**

- 屏幕预览：`max-width: 210mm`，居中，外层浅灰背景模拟纸感。
- 打印 (`@media print`)：去掉外层灰背景、占满整页。
- 字体：`"Inter", "PingFang SC", "Microsoft YaHei", "Noto Sans CJK SC", sans-serif`。

**Header（个人信息块）**

- 左 64px 圆形头像；无头像时使用中性灰底占位（不显示 emoji）。
- 右侧：
  - 姓名 24px / 600
  - 职位 14px / 500 / `text-slate-700`
  - 联系方式行 11px / `text-slate-500`：所在地、手机、邮箱
  - 个人事实行 11px / `text-slate-500`：年龄 / 性别 / 学历 / 经验 / 籍贯
  - 意向城市行 11px / `text-slate-500`（如有）
- 底部 1px 浅灰分隔线 (`border-slate-200`)。

**章节标题（统一样式）**

- 4px 宽 16px 高品牌色竖条 (`bg-brand-500`) + 14px / 700 / `text-slate-900` 标题。
- 章节标题上方 16px 间距，下方 8px 间距。

**条目（工作 / 项目 / 教育 / 荣誉）**

- 第一行：左侧加粗主标题（14px / 600）+ 右侧灰色日期（11px / `text-slate-400`），`flex justify-between`。
- 第二行（如有）：副标题（角色 / 专业，11px / `text-slate-600`）。
- 描述列表：`<ul>` 11px / 行高 1.6 / `text-slate-700`，自定义品牌色小圆点（不用默认 `disc`）。
- 项目超链接：标题用 `text-brand-600`，末尾追加 `↗` 上标小图标。

**技能区域**

- 不使用 Tag 圆角胶囊（打印浪费墨水），改为 `Go · MySQL · Redis · K8s` 中点分隔的纯文本行。

**品牌色用量**

| 位置 | 色值 |
|---|---|
| 章节标题左侧竖条 | `brand-500` |
| 项目链接标题 | `brand-600` |
| 描述列表小圆点 | `brand-400` |
| 其他文字、分隔线、日期 | `slate-900 / slate-700 / slate-500 / slate-200 / slate-400` |

**工具栏**

- 顶部"打印预览 / 下载 PDF"工具栏继续用 `print:hidden`。chromedp PrintToPDF 走 print 媒介查询，工具栏自动隐藏。

**空数据处理**

- 章节内容为空时（如 `awards` 数组为空），整个章节连标题一起不渲染，避免出现"荣誉奖项"标题但下方为空。

## 智能分页

通过 CSS 控制，单个条目不允许跨页拆分；章节标题不能孤立在页底。

```css
@media print {
  @page {
    size: A4;
    margin: 12mm 14mm;
  }

  /* 单个条目不允许跨页拆分 */
  .print-item {
    break-inside: avoid;
    page-break-inside: avoid;
  }

  /* 章节标题不能孤立在页底 */
  .print-section-title {
    break-after: avoid;
    page-break-after: avoid;
  }

  /* 章节内的第一个条目跟随标题 */
  .print-section-title + .print-item {
    break-before: avoid;
  }
}
```

落到代码：

- 每个工作 / 项目 / 教育 / 荣誉条目包裹在 `<div class="print-item">`。
- 每个章节标题加 `class="print-section-title"`。

## PDF 后端实现

### Exporter 重构

`Exporter` 构造函数从 `NewExporter(browserPath)` 改为 `NewExporter(browserPath, port)`，`port` 复用 `cfg.Port`，不引入新环境变量。

`ExportResume` 签名从 `ExportResume(ctx, resume json.RawMessage)` 改为 `ExportResume(ctx)` —— 不再需要传 resume 数据，因为 React 会自己 fetch。

```go
// 简化伪代码
func (e Exporter) ExportResume(ctx context.Context) ([]byte, error) {
    browserPath, err := resolveBrowserPath(e.browserPath, fallbacks)
    if err != nil { return nil, err }

    allocCtx, cancelA := chromedp.NewExecAllocator(ctx,
        append(chromedp.DefaultExecAllocatorOptions[:], chromedp.ExecPath(browserPath))...)
    defer cancelA()
    browserCtx, cancelB := chromedp.NewContext(allocCtx)
    defer cancelB()

    var pdfBytes []byte
    url := fmt.Sprintf("http://127.0.0.1:%s/print?pdfMode=1", e.port)

    err = chromedp.Run(browserCtx,
        chromedp.Navigate(url),
        chromedp.WaitVisible(`#print-root`, chromedp.ByID),
        chromedp.Poll(`window.__printReady === true`, nil,
            chromedp.WithPollingTimeout(10*time.Second)),
        chromedp.ActionFunc(func(ctx context.Context) error {
            buf, _, err := page.PrintToPDF().
                WithPrintBackground(true).
                WithPaperWidth(8.27).   // A4 inches
                WithPaperHeight(11.69).
                WithMarginTop(0).WithMarginBottom(0).
                WithMarginLeft(0).WithMarginRight(0).
                WithPreferCSSPageSize(true).
                Do(ctx)
            if err != nil { return err }
            pdfBytes = buf
            return nil
        }),
    )
    return pdfBytes, err
}
```

关键点：

- **`PreferCSSPageSize(true)` + 后端边距设 0**：让 `@page { margin: 12mm 14mm }` 生效，由 CSS 单点控制边距，避免后端和 CSS 各设一份导致漂移。
- **`Poll(window.__printReady === true)`**：等 React fetch 完数据、DOM 渲染稳定、图片加载完成后才打印。
- **超时 10 秒**：避免后端 `/api/resume` 抖动时无限挂起。
- **`?pdfMode=1`**：仅作语义信号，目前前端不强依赖（print CSS 已能区分），保留扩展位。

### React 端的 ready 信号

`PrintPage.tsx`（或新拆出的 `PrintResume.tsx`）在数据 + 图片都加载完成后设置 `window.__printReady`：

```tsx
const { draft, loading } = useResumeDraft()

useEffect(() => {
  if (loading) return

  Promise.all(
    Array.from(document.images).map((img) =>
      img.complete
        ? Promise.resolve()
        : new Promise<void>((resolve) => {
            img.onload = img.onerror = () => resolve()
          }),
    ),
  ).then(() => {
    ;(window as any).__printReady = true
  })
}, [loading])
```

避坑：头像图片若未加载完就打印，PDF 可能缺图。`Promise.all + img.complete` 是兜底方案。

### Handler 调整

`internal/handler/resume.go` 的 `ExportPDF`：

- 不再 `h.store.Get()` 取 resume。
- 只检查 `AllowPDFExport` 开关，然后调 `h.exporter.ExportResume(c.Request.Context())`。

### 配置注入

`cmd/server/main.go`：

```go
resumeHandler := handler.NewResumeHandler(
    resumeStore,
    siteSettingsStore,
    pdf.NewExporter(cfg.BrowserPath, cfg.Port),
)
```

## 迁移与代码清单

| 文件 | 操作 |
|---|---|
| `internal/pdf/export.go` | 删除 `resumeView` 结构体、`buildResumeView`、`resolveAvatarDataURL`、`renderResumeHTML`；重写 `Exporter` 与 `ExportResume` |
| `internal/pdf/export_test.go` | 重写：基于 `httptest.Server` 模拟 `/print` 返回最小占位 HTML，跑 chromedp 拉 PDF |
| `internal/handler/resume.go` | `ExportPDF` 不再取 resume；调用 `ExportResume(ctx)` |
| `internal/handler/resume_test.go` | 跟随新 PDFExporter 接口调整 mock |
| `cmd/server/main.go` | `pdf.NewExporter(cfg.BrowserPath, cfg.Port)` |
| `web/src/pages/PrintPage.tsx` | 工具栏样式微调；引入 `__printReady` 信号 |
| `web/src/components/resume/PrintResume.tsx` | **完全重写**：新视觉、章节顺序、`print-item` / `print-section-title` 类、空数据隐藏、品牌色点缀、`@media print` + `@page` 块 |

`ProfileCard` / `SkillSection` / `TimelineSection` 不动 —— 它们是在线展示专用，强行复用反而拖累两边。在线和打印两种媒介视觉语言天然不同，共享数据类型 `ResumeData` 而不共享组件，这是有意的解耦。

## 测试策略

### 前端单测（Vitest + RTL）

`PrintResume.test.tsx`：

- 基于固定 `ResumeData` 渲染 DOM，断言关键结构。
- 断言章节顺序：技能 → 工作 → 项目 → 教育 → 荣誉。
- 断言空数据章节不渲染（传入空 `awards` 数组，不应出现"荣誉奖项"标题）。
- 断言 `print-item` / `print-section-title` 类名正确挂载（智能分页依赖这些 hook）。

`PrintPage.test.tsx`：

- 断言数据加载完成后 `window.__printReady` 被置为 `true`。

### 后端单测（Go）

`internal/pdf/export_test.go` 重写：

- 起一个内嵌 `httptest.Server` 模拟 `/print` 返回最小占位 HTML：
  ```html
  <div id="print-root">test</div>
  <script>window.__printReady = true</script>
  ```
- 跑 chromedp 拉 PDF。
- 断言 PDF 字节流非空、以 `%PDF-` 开头。
- 断言 `Poll(__printReady)` 不超时。

不在单元测试里跑真实 React 构建产物（成本高、依赖前端构建产物存在）；端到端 PDF 视觉验证靠手动验收，暂不做截图比对。

### 验收标准（人眼）

- 打开 `/print` 屏幕上的版式 = 下载的 PDF 版式（除工具栏外像素级一致）。
- 一份典型简历（当前 mock 数据）单页打得下；A4 上字号舒适（正文 11px）。
- 头像在 PDF 中正确显示（不丢图）。
- 中文字符无方框、无字体缺失。
- 空数据章节自动消失。

## 风险与应对

| 风险 | 概率 | 应对 |
|---|---|---|
| chromedp `Poll` 一直拿不到 `__printReady`（前端报错挂死） | 中 | 10 秒超时 → 返回 504，前端 toast 报错 |
| 后端访问 `127.0.0.1:<port>` 在容器/k8s 里失败 | 低 | Docker 部署里前后端同进程，环回访问无问题；如有反例可加配置项覆盖 host |
| 头像未加载完就 print（PDF 丢图） | 中 | `Promise.all + img.complete` 兜底 |
| Chrome 不同版本 `@page` + `PreferCSSPageSize` 行为差异 | 低 | 锁定 chromedp + 现有 Chrome 路径检测；CI 不需要改 |
| 现有 `export_test.go` 在 CI 上无 Chrome 跑不动 | 已存在 | 现状如此，本次不修 |
| 字段后续新增（如「自我评价」） | 低 | 在 `PrintResume` 加一节即可；类型驱动 |

## 开放问题

无。所有架构和视觉决策已敲定。
