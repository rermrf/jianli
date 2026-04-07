# config.example.json 配置策略设计

## 概述

当前仓库仍在跟踪根目录 `config.json`，这会让部署者直接在受版本控制的文件上填写自己的配置，既容易造成误提交，也与“敏感项不应进入仓库”的目标相冲突。本次改动将配置策略统一为：仓库只保留 `config.example.json` 作为模板文件，`config.json` 改为本地生成且不纳入 git 跟踪；应用运行时仍默认读取项目根目录的 `config.json`。部署相关文件和文档也同步到这套约定。

## 目标

- 仓库不再跟踪 `config.json`
- 仓库新增或保留 `config.example.json` 作为配置模板
- 用户通过 `cp config.example.json config.json` 生成自己的配置文件
- `internal/config/config.go` 继续默认读取项目根目录 `config.json`
- `.gitignore` 忽略 `config.json`
- 部署文件与手动部署文档同步使用这套配置约定

## 非目标

- 不修改配置文件读取入口为其他路径
- 不引入环境变量化的普通配置项
- 不引入多套配置模板（dev/staging/prod）
- 不引入 secrets manager

## 核心设计

### 仓库内文件

改动后配置文件分工如下：

- `config.example.json`
  - 受 git 跟踪
  - 提供普通配置模板
  - 不包含任何私人或敏感真实值
- `config.json`
  - 不受 git 跟踪
  - 由用户从模板复制生成
  - 作为实际运行配置文件

### 配置读取行为

`internal/config/config.go` 不改变默认读取入口：

- 仍然读取项目根目录 `config.json`

这样可以避免扩大这次改动范围，也不需要修改所有部署脚本和运行入口的默认行为。

### 认证密钥

`AUTH_KEY` 仍保持环境变量读取，不进入 `config.example.json`。

## 文件策略

### `config.example.json`

模板中保留：

- `browserPath`
- `port`
- `dbPath`
- `frontendOrigin`

不包含：

- `authKey`

### `.gitignore`

增加：

- `config.json`

保证后续用户不会把自己的运行配置误提交进仓库。

### `config.json`

从 git 跟踪中移除，但保留文件在本地工作区中的存在方式由用户自己控制。

## 部署文件影响

### `deploy/docker-compose.yml`

继续保持：

- `../config.json:/app/config.json:ro`

因为应用仍默认读取项目根目录 `config.json`。

### `deploy/env/app.env.example`

继续只保留：

- `AUTH_KEY`
- `APP_BIND`
- `APP_HEALTHCHECK_URL`

不再引入额外配置文件路径变量。

### `deploy/config/config.production.json.example`

建议删除。

原因：

- 仓库中同时存在 `config.example.json` 和 `deploy/config/config.production.json.example` 会形成两套配置入口
- 用户容易困惑到底应该编辑哪一个
- 既然应用最终读取的是根目录 `config.json`，应只保留一套模板来源

## 文档影响

手动部署文档需要改为明确说明：

1. 进入项目目录
2. 执行：

```bash
cp config.example.json config.json
```

3. 编辑生成后的 `config.json`
4. 通过环境变量提供 `AUTH_KEY`
5. 运行 `docker compose`

同时删除任何仍指向 `deploy/config/config.production.json.example` 的说明。

## 本地开发影响

本地开发流程变为：

```bash
cp config.example.json config.json
```

然后编辑：

- `browserPath`
- `port`
- `dbPath`
- `frontendOrigin`

最后设置：

```bash
AUTH_KEY=your-local-key
```

再启动应用。

## 测试边界

至少验证：

- `internal/config/config.go` 仍能正确读取根目录 `config.json`
- 在 `AUTH_KEY` 存在时，删除 `config.json` 会导致启动失败（保持原行为）
- 部署文档中不再引用旧的 production config example
- `docker compose ... config` 仍能通过

## 验收标准

满足以下条件视为完成：

- 仓库已不再跟踪 `config.json`
- 仓库存在 `config.example.json` 作为唯一配置模板
- `.gitignore` 已忽略 `config.json`
- 应用仍默认读取根目录 `config.json`
- 部署文件与文档已同步到新的配置策略
- `AUTH_KEY` 仍保持环境变量读取
