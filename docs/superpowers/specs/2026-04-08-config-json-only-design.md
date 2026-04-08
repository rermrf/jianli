# 全量 config.json 配置策略设计

## 概述

当前仓库中配置策略仍处于混合状态：一部分业务配置通过 `config.json` 提供，`authKey` 曾被迁移到环境变量，部署文件也引入了额外的环境变量输入。这与“项目规模较小，配置应尽量集中”的目标不一致。本次改动将配置策略进一步统一为：所有业务配置都由项目根目录的 `config.json` 提供，仓库只跟踪 `config.example.json` 模板文件，部署时用户通过 `cp config.example.json config.json` 生成自己的本地配置；Docker Compose 与部署文档不再依赖 `AUTH_KEY` 作为环境变量输入。

## 目标

- `authKey` 重新回到 `config.json` 中配置
- `internal/config/config.go` 不再读取 `AUTH_KEY` 环境变量
- 仓库只跟踪 `config.example.json`
- `config.json` 仍然本地生成且不纳入 git 跟踪
- 所有业务配置统一来自 `config.json`
- Docker 部署不再要求设置 `AUTH_KEY`
- 手动部署文档只要求复制并编辑 `config.example.json`

## 非目标

- 不引入多环境配置体系
- 不引入 secrets manager
- 不改动配置读取入口路径，仍默认读取项目根目录 `config.json`
- 不调整前端请求头格式或登录逻辑

## 核心设计

### 配置来源统一规则

最终统一为：

- `authKey`：从 `config.json` 读取
- `browserPath`：从 `config.json` 读取
- `port`：从 `config.json` 读取
- `dbPath`：从 `config.json` 读取
- `frontendOrigin`：从 `config.json` 读取

换句话说，运行时业务配置只有一个来源：

- 项目根目录 `config.json`

### 模板文件策略

仓库中只保留：

- `config.example.json`

用户本地通过：

```bash
cp config.example.json config.json
```

生成自己的运行配置。

## 文件策略

### `config.example.json`

模板文件中应包含完整业务配置字段：

- `authKey`
- `browserPath`
- `port`
- `dbPath`
- `frontendOrigin`

### `config.json`

- 不纳入 git 跟踪
- 作为实际运行配置
- 可因本地环境或服务器环境不同而自行修改

### `.gitignore`

继续忽略：

- `config.json`

## 代码边界

### `internal/config/config.go`

改动后应：

- 保留当前 `Config.AuthKey` 字段
- 通过 `json.Unmarshal` 直接从 `config.json` 读取 `authKey`
- 删除 `os.Getenv("AUTH_KEY")` 相关逻辑
- 若 `authKey` 缺失，则返回明确错误，例如：
  - `authKey is required in config.json`

### `internal/config/config_test.go`

测试应同步回退为：

- `authKey` 来自 `config.json`
- 缺少 `authKey` 时加载失败
- 不再需要设置环境变量 `AUTH_KEY`

## 部署文件影响

### `deploy/docker-compose.yml`

应调整为：

- 删除：
  - `environment.AUTH_KEY`
- 保留：
  - `../config.json:/app/config.json:ro`
  - `../data:/app/data`

如果你希望进一步减少部署输入，建议：

- 将 `ports` 直接写死为 `127.0.0.1:8088:8088`

而不再保留 `APP_BIND` 这样的外部变量。

### `deploy/env/app.env.example`

建议删除。

原因：

- 既然所有业务配置都统一回到 `config.json`
- 当前项目又不需要 CI/CD
- 保留单独 `.env` 文件只会继续制造第二套输入源

### `deploy/scripts/deploy.sh`

应调整为不依赖 env 文件：

```bash
docker compose -f deploy/docker-compose.yml up -d --build app
```

### `deploy/scripts/healthcheck.sh`

可保留，但默认端口应与配置示例一致，例如：

- `http://127.0.0.1:8088/api/resume`

## 文档影响

### 手动部署文档

部署文档应统一为：

1. `git pull`
2. `cp config.example.json config.json`
3. 编辑 `config.json`
4. `docker compose -f deploy/docker-compose.yml up -d --build app`

不再提：

- `.env`
- `AUTH_KEY`
- deployment env file

## 本地开发影响

本地运行方式将变得更简单：

```bash
cp config.example.json config.json
# 编辑 config.json

go run ./cmd/server
```

无需额外设置环境变量。

## 测试边界

至少覆盖：

- `internal/config/config.go` 正常从 `config.json` 读取 `authKey`
- 缺少 `authKey` 时配置加载失败
- BOM `config.json` 仍能正常解析
- Docker Compose 静态校验通过
- 现有后端测试继续通过

## 验收标准

满足以下条件视为完成：

- 所有业务配置统一回到 `config.json`
- `authKey` 不再使用环境变量
- 仓库只跟踪 `config.example.json`
- `config.json` 继续为本地生成文件，不纳入 git 跟踪
- 部署文件与文档不再依赖 `AUTH_KEY` 环境变量
- 配置和后端测试通过
