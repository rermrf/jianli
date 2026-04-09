# PDF 导出开关设计

## 目标

为管理员提供一个可持久化的“允许导出 PDF”开关。

行为要求：

- 管理员可以在 `/edit` 页面直接控制是否允许导出 PDF
- 当开关关闭时，公开简历页隐藏“导出 PDF”按钮
- 当开关关闭时，打印页隐藏“下载 PDF”按钮
- 当开关关闭时，后端 `GET /api/resume/pdf` 必须拒绝导出，不能只依赖前端隐藏

## 非目标

- 不把该开关混入简历正文数据
- 不通过手改 `config.json` 控制该功能
- 不新增复杂的站点设置中心，只先支持本次需要的单个开关

## 方案概览

采用独立的站点设置记录，而不是把开关塞进 `resume` 数据。

原因：

- `allowPdfExport` 属于站点行为策略，不属于简历内容
- 后续如果再加“允许访客统计”“允许公开展示手机号”之类的开关，可以继续放进同一份设置记录
- 前后端职责更清晰，简历数据和站点设置互不污染

## 数据模型

新增 `site_settings` 持久化记录，当前只保存一个字段：

- `allowPdfExport: boolean`

建议保持单记录模式，主键固定为 `1`。

后端 store 层负责：

- 首次为空时创建默认记录
- 默认值为 `true`
- 提供读取和更新方法

## API 设计

### 公开接口

保留 `GET /api/resume`，但返回结构扩展为：

```json
{
  "code": 0,
  "data": {
    "resume": { "...": "现有简历数据" },
    "siteSettings": {
      "allowPdfExport": true
    }
  }
}
```

这样公开页和打印页都能一次请求拿到按钮显示所需的状态。

### 管理接口

新增管理员接口：

- `GET /api/settings`
- `PUT /api/settings`

请求与响应只暴露当前需要的字段：

```json
{
  "allowPdfExport": true
}
```

这两个接口都走现有管理员鉴权中间件。

### PDF 导出接口

保留 `GET /api/resume/pdf` 路由。

新增行为：

- 当 `allowPdfExport=true` 时，维持当前导出逻辑
- 当 `allowPdfExport=false` 时，直接返回拒绝状态，不生成 PDF

建议返回：

- HTTP `403 Forbidden`
- 业务错误码单独分配
- 文案类似：`pdf export is disabled`

## 前端改动

### 数据访问

前端现有 `loadResumeDraft()` 依赖 `/api/resume` 直接返回简历数据。

需要把公开返回结构调整为：

- `resume`
- `siteSettings`

前端新增对应类型，避免把设置字段混进 `ResumeData`。

### 编辑页 `/edit`

在编辑页顶部区域新增一个轻量设置卡片，放在现有“保存主简历 / 保存为草稿”区域附近。

建议文案：

- 标题：`站点设置`
- 开关标签：`允许访客导出 PDF`
- 辅助说明：`关闭后，公开页和打印页都会隐藏导出按钮，后端接口也会拒绝 PDF 导出。`

交互建议：

- 编辑页加载时同时获取当前设置
- 切换开关后直接调用 `PUT /api/settings`
- 成功后显示轻量 toast
- 失败时回退开关状态并显示错误

### 公开简历页 `/`

当前“导出 PDF”入口在 `ResumePage`。

改为：

- 只有 `siteSettings.allowPdfExport === true` 时显示按钮
- 关闭时不渲染按钮

### 打印页 `/print`

打印页仍可访问，但：

- 当 `allowPdfExport=true` 时显示“下载 PDF”
- 当 `allowPdfExport=false` 时隐藏“下载 PDF”

这样 `/print` 还能作为只读打印预览存在，但不能再导出文件。

真正的权限边界仍以后端 `/api/resume/pdf` 拒绝为准。

## 后端结构建议

新增模块：

- `internal/model/site_settings.go`
- `internal/store/site_settings_store.go`
- `internal/handler/settings.go`

现有 `resume` handler 调整：

- `GET /api/resume` 组装 `resume + siteSettings`
- `GET /api/resume/pdf` 导出前先读取 `siteSettings`

`cmd/server/main.go` 调整：

- 初始化 `siteSettingsStore`
- 注册 `GET /api/settings`
- 注册 `PUT /api/settings`
- 把 `siteSettingsStore` 注入需要它的 handler

## 错误处理

前端：

- 设置保存失败时恢复开关显示状态
- 对管理员显示明确错误提示

后端：

- 设置记录缺失时自动补默认值，不把“无记录”暴露给用户
- PDF 导出关闭时返回明确的 `403`

## 测试策略

### 后端

新增或修改测试覆盖：

- `site_settings` store 空库默认创建并返回 `allowPdfExport=true`
- `PUT /api/settings` 能持久化 `allowPdfExport=false`
- `GET /api/resume` 返回 `resume + siteSettings`
- `GET /api/resume/pdf` 在开关关闭时返回 `403`

### 前端

新增或修改测试覆盖：

- 编辑页加载后显示“允许访客导出 PDF”开关
- 管理员切换开关时会调用 `PUT /api/settings`
- 首页在 `allowPdfExport=false` 时不显示“导出 PDF”
- 打印页在 `allowPdfExport=false` 时不显示“下载 PDF”

## 兼容性与迁移

- 已有数据库无需手工迁移复杂结构，只需在启动或首次读取时自动创建 `site_settings` 默认记录
- `/api/resume` 返回结构会变化，前端必须同步更新后一起发布
- 由于本仓库前后端同源部署，适合一次性联动上线

## 验收标准

- 管理员能在 `/edit` 页面切换“允许访客导出 PDF”
- 设置刷新后仍然保留
- 公开简历页在关闭时不显示“导出 PDF”
- 打印页在关闭时不显示“下载 PDF”
- 手工请求 `GET /api/resume/pdf` 在关闭时返回拒绝
- 开关重新打开后上述能力恢复
