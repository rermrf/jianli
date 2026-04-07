# AUTH_KEY 环境变量化设计

## 概述

当前后端认证密钥 `authKey` 仍从 `config.json` 读取，这使得部署到 Docker / Jenkins 环境时需要通过挂载配置文件来传递敏感信息，不利于在容器编排和 CI/CD 中管理机密。本次改动将 `authKey` 调整为仅从环境变量 `AUTH_KEY` 读取，而 `port`、`dbPath`、`browserPath`、`frontendOrigin` 等其他配置仍继续从 `config.json` 获取。

## 目标

- 后端认证密钥只从环境变量 `AUTH_KEY` 读取
- 启动时若未设置 `AUTH_KEY`，应用直接报错退出
- `config.json` 不再承载 `authKey`
- Docker Compose 可以直接为应用容器注入 `AUTH_KEY`
- Jenkins / 部署文档明确要求配置 `AUTH_KEY`

## 非目标

- 不把 `port` 改为环境变量
- 不把 `dbPath` 改为环境变量
- 不把 `browserPath` 改为环境变量
- 不改前端登录逻辑和请求头格式
- 不引入更复杂的 secrets manager

## 核心设计

### 配置读取规则

改动后配置规则为：

- `AUTH_KEY`：仅从环境变量读取
- `port`：继续从 `config.json` 读取
- `dbPath`：继续从 `config.json` 读取
- `browserPath`：继续从 `config.json` 读取
- `frontendOrigin`：继续从 `config.json` 读取

### 启动校验

在配置加载阶段：

- 若环境变量 `AUTH_KEY` 为空，返回明确错误
- 应用启动失败，不允许进入运行态

这和当前 `authKey is required in config.json` 的失败方式一致，只是来源改为环境变量。

## 代码边界

### `internal/config/config.go`

建议改成：

- 从 `config.json` 读取除 `authKey` 外的配置
- 在 `Load()` 内通过 `os.Getenv("AUTH_KEY")` 赋值给 `Config.AuthKey`
- 若为空则报错

实现上可以保留 `Config.AuthKey` 字段，因为后续业务层仍通过 `cfg.AuthKey` 使用它。

### `config.json`

- 删除 `authKey` 字段
- 保留其余字段

这样可以避免误导使用者继续把密钥写进文件。

### 测试

`internal/config/config_test.go` 需要同步调整：

- 原先依赖 `config.json` 中 `authKey` 的测试应改为设置环境变量 `AUTH_KEY`
- 增加未设置 `AUTH_KEY` 时返回错误的覆盖
- 增加“即使 `config.json` 缺少 `authKey`，但环境变量存在时仍可正常加载”的覆盖

## Docker / Compose 影响

### Compose

`deploy/docker-compose.yml` 中的 `app` 服务需要注入：

- `AUTH_KEY=${AUTH_KEY}`

或等价写法。

### 环境样例

`deploy/env/app.env.example` 需要增加：

- `AUTH_KEY=replace-with-strong-secret`

## Jenkins 影响

Jenkins Pipeline 本身不需要改动业务逻辑，但部署文档和运行环境必须保证：

- `deploy.env` 中存在 `AUTH_KEY`
- Jenkins 在部署容器时能把该变量传入应用容器

## 文档影响

以下文档需要同步更新：

- Docker + Jenkins 部署文档
- 若仓库中有本地运行说明，也应注明需要先设置 `AUTH_KEY`

## 本地开发影响

本地开发时，后端启动前必须设置环境变量，例如：

### Windows PowerShell

```powershell
$env:AUTH_KEY='your-local-key'
go run ./cmd/server
```

### Linux / macOS

```bash
AUTH_KEY=your-local-key go run ./cmd/server
```

## 错误处理

启动失败错误建议保持直白，例如：

- `AUTH_KEY is required`

避免继续提示 `config.json`，否则会误导使用者。

## 测试边界

至少覆盖：

- 设置 `AUTH_KEY` 且 `config.json` 正常时，配置加载成功
- 未设置 `AUTH_KEY` 时，配置加载失败
- `config.json` 缺少旧 `authKey` 字段时，只要环境变量存在仍能成功
- 现有使用 `cfg.AuthKey` 的 handler / middleware 测试继续通过

## 验收标准

满足以下条件视为完成：

- `authKey` 不再从 `config.json` 读取
- `AUTH_KEY` 为唯一认证密钥来源
- 未设置 `AUTH_KEY` 时应用无法启动
- Docker Compose 与部署文档已支持 `AUTH_KEY`
- 配置相关测试通过
