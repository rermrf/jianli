# AUTH_KEY 仅环境变量读取设计

## 概述

当前后端认证密钥 `authKey` 仍从 `config.json` 读取，这使得部署到 Docker / Jenkins 环境时需要通过挂载配置文件传递敏感信息，不利于在容器编排和 CI/CD 中统一管理机密。本次改动将认证密钥彻底调整为仅从环境变量 `AUTH_KEY` 读取，`config.json` 不再承载 `authKey` 字段；而 `port`、`dbPath`、`browserPath`、`frontendOrigin` 等其他配置继续保持从 `config.json` 读取。

## 目标

- 后端认证密钥只从环境变量 `AUTH_KEY` 读取
- 启动时若未设置 `AUTH_KEY`，应用直接报错退出
- `config.json` 不再包含 `authKey`
- Docker Compose 可直接为应用容器注入 `AUTH_KEY`
- Jenkins / 部署文档明确要求在部署环境提供 `AUTH_KEY`
- 本地开发仍可通过设置环境变量运行后端

## 非目标

- 不把 `port` 改为环境变量
- 不把 `dbPath` 改为环境变量
- 不把 `browserPath` 改为环境变量
- 不把 `frontendOrigin` 改为环境变量
- 不改前端登录逻辑和 `X-Auth-Key` 请求头格式
- 不接入更复杂的 secrets manager

## 核心设计

### 配置来源规则

改动后配置来源如下：

- `AUTH_KEY`：仅从环境变量读取
- `port`：继续从 `config.json` 读取
- `dbPath`：继续从 `config.json` 读取
- `browserPath`：继续从 `config.json` 读取
- `frontendOrigin`：继续从 `config.json` 读取

### 配置结构

`Config` 结构体可继续保留 `AuthKey` 字段，以便调用方继续通过 `cfg.AuthKey` 使用认证密钥；但该字段的赋值来源改为：

- `os.Getenv("AUTH_KEY")`

而不再来自 `config.json` 反序列化结果。

## 启动行为

在 `Load()` 过程中：

- 先读取并解析 `config.json`
- 再从环境变量读取 `AUTH_KEY`
- 若 `AUTH_KEY` 为空，则返回错误并阻止应用启动

建议错误文案明确为：

- `AUTH_KEY is required`

避免继续提示 `config.json`，否则会误导使用者。

## 代码边界

### `internal/config/config.go`

需要做的调整：

- `Config` 结构体仍保留 `AuthKey string`
- `json.Unmarshal` 之后，不再依赖 `cfg.AuthKey` 的 JSON 值
- 改为 `cfg.AuthKey = os.Getenv("AUTH_KEY")`
- 若为空则报错

### `config.json`

需要调整：

- 删除 `authKey` 字段
- 保留其他字段不变

这样可以防止开发者继续把认证密钥提交进仓库。

## Docker / Compose 影响

### `deploy/docker-compose.yml`

应用容器需要显式接收：

- `AUTH_KEY=${AUTH_KEY}`

### `deploy/env/app.env.example`

需要新增：

- `AUTH_KEY=replace-with-strong-secret`

这样部署环境可以直接通过 `deploy.env` 注入密钥。

## Jenkins 影响

Jenkins Pipeline 的主流程不需要改业务阶段，但要确保：

- `deploy.env` 中包含 `AUTH_KEY`
- 部署脚本和 `docker compose` 启动应用时会把该变量传给 `app` 容器

## 本地开发影响

本地运行后端时，开发者必须先设置环境变量：

### Windows PowerShell

```powershell
$env:AUTH_KEY='your-local-key'
go run ./cmd/server
```

### Linux / macOS

```bash
AUTH_KEY=your-local-key go run ./cmd/server
```

## 测试边界

### 配置测试

至少覆盖：

- `AUTH_KEY` 存在且 `config.json` 正常时，配置加载成功
- 未设置 `AUTH_KEY` 时，配置加载失败
- `config.json` 不包含旧 `authKey` 字段时，只要环境变量存在仍能成功加载
- BOM `config.json` 场景仍能正常加载

### 回归验证

需要确认：

- 现有使用 `cfg.AuthKey` 的 handler / middleware 测试继续通过
- Docker Compose 配置静态校验通过
- 部署文档与示例环境文件已同步更新

## 文档影响

需要同步更新：

- Docker + Jenkins 部署文档
- 任何仍写着 `config.json` 里包含 `authKey` 的说明

## 验收标准

满足以下条件视为完成：

- `AUTH_KEY` 成为唯一认证密钥来源
- `config.json` 不再承载 `authKey`
- 未设置 `AUTH_KEY` 时应用无法启动
- Docker Compose 与部署文档已支持 `AUTH_KEY`
- 配置与回归测试通过
