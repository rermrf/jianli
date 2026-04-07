# Docker + Jenkins 全容器化部署设计

## 概述

当前项目运行方式仍偏开发态：前端通过 Vite 开发服务器访问，后端直接读取本地 `config.json` 和 `data/` 目录，没有正式的服务器部署编排。本次改动将该项目升级为“全 Docker 化部署”，并使用同机 Jenkins 负责持续集成与手动触发发布，使代码推送到 `master` 后可以自动测试、自动构建容器镜像，在人工确认后完成正式部署。

## 目标

- 项目可以通过 Docker 在 Linux 服务器上一键启动
- 前端、后端、Nginx、Jenkins 全部容器化
- Jenkins 与生产应用部署在同一台 Linux 服务器
- 推送到 `master` 后自动执行测试与构建
- 发布动作由 Jenkins 中的人工确认步骤触发
- 应用运行所需的数据库、上传文件、配置文件通过挂载持久化
- 正式环境支持 PDF 导出能力

## 非目标

- 不做蓝绿发布
- 不做多环境（dev / staging / prod）编排
- 不做回滚自动化
- 不接入镜像仓库（首版使用本机 Docker 构建与部署）
- 不引入 Kubernetes
- 不做监控、告警、日志平台接入

## 总体架构

在同一台 Linux 服务器上部署以下容器：

- `jenkins`
  - 负责拉取代码、运行测试、构建镜像、手动确认发布
- `app`
  - 运行 Go 后端
  - 提供前端构建后的静态资源
  - 处理 `/api`、`/uploads`、PDF 导出
- `nginx`
  - 对外暴露 80 / 443
  - 反向代理到 `app`
- 可选：`jenkins-agent` 不单独拆出，首版直接使用 Jenkins 容器内的执行环境

### 对外流量

用户请求路径：

- `/` -> `nginx` -> `app` 中的前端静态资源
- `/api/*` -> `nginx` -> `app` 中的 Gin 接口
- `/uploads/*` -> `nginx` -> `app` 中的上传文件目录

### CI/CD 路径

Jenkins Pipeline：

1. 拉取 `master`
2. 运行 Go 与前端测试
3. 构建应用镜像
4. 停在人工确认步骤
5. 人工确认后执行 `docker compose up -d --build`

## 容器化策略

### 应用容器 `app`

单容器承载：

- Go 编译后的后端二进制
- 前端 `dist` 静态文件
- 用于 PDF 导出的浏览器（推荐 `chromium`）

原因：

- 项目是单体应用
- 前端请求使用相对路径 `/api`、`/uploads`
- 统一容器可以减少前后端分离部署时的额外配置复杂度

### Nginx 容器

职责：

- 对外提供统一域名入口
- 反向代理 `/api` 和 `/uploads`
- 为前端 SPA 配置 `try_files ... /index.html`

### Jenkins 容器

职责：

- 承载 Jenkins UI
- 执行 Pipeline
- 挂载宿主机 Docker Socket，直接控制本机 Docker

风险说明：

- 挂载 Docker Socket 安全性较低
- 但对当前“单机部署、快速落地”的需求最直接
- 首版接受该权衡，后续再考虑隔离

## 运行目录设计

建议服务器目录结构：

- `/srv/jianli`
  - 项目仓库工作目录
- `/srv/jianli/deploy`
  - `docker-compose.yml`
  - `nginx.conf`
  - `.env`
- `/srv/jianli/config`
  - 生产版 `config.json`
- `/srv/jianli/data`
  - `resume.db`
  - `uploads/avatars`
- `/srv/jenkins_home`
  - Jenkins 持久化目录

## 配置设计

### 生产配置文件

`config.json` 仍作为应用配置入口，但内容改为适配 Linux 容器环境：

- `authKey`
  - 必须改为强随机值
- `browserPath`
  - 改为 Linux 浏览器路径，例如 `/usr/bin/chromium`
- `port`
  - 容器内建议固定为 `8080`
- `dbPath`
  - 改为容器内路径，例如 `/app/data/resume.db`
- `frontendOrigin`
  - 改为正式域名，如 `https://resume.example.com`

### 为什么必须处理 `browserPath`

当前 PDF 导出实现依赖 `chromedp`，并通过配置传入浏览器路径。开发环境的 Windows Edge 路径在 Linux 容器中必然失效，因此生产配置必须显式切换到 Linux 浏览器路径，否则 PDF 导出不可用。

## Docker 设计

### `Dockerfile`

建议使用多阶段构建：

1. 前端构建阶段
   - `node` 镜像
   - 执行 `npm ci`、`npm run build`
2. 后端构建阶段
   - `golang` 镜像
   - 执行 `go build`
3. 运行阶段
   - 基于 Debian / Ubuntu slim
   - 安装 `chromium`
   - 拷贝后端二进制和前端 `dist`
   - 设置工作目录 `/app`

运行阶段需要包含：

- `/app/bin/jianli-server`
- `/app/web/dist`
- `/app/data` 作为挂载目录
- `/app/config.json` 作为挂载文件

### `docker-compose.yml`

首版至少定义：

- `jenkins`
- `app`
- `nginx`

持久化挂载：

- Jenkins 数据目录
- 应用 `data/`
- 应用 `config.json`
- Nginx 配置文件

端口：

- `nginx` 暴露 `80:80`
- `jenkins` 暴露 `8081:8080` 或类似端口，避免和业务冲突
- `app` 不对公网直接暴露，可仅在 compose 内部网络使用

## Nginx 配置

Nginx 需要满足：

- `/` 走静态前端
- `/api/` 反代到 `app:8080`
- `/uploads/` 反代到 `app:8080`
- SPA 回退到 `/index.html`

可选后续增强：

- HTTPS
- gzip
- 静态缓存

但首版不是必须项。

## Jenkins Pipeline 设计

### 触发条件

- 监听 `master` 分支变更
- 代码推送后自动触发 Pipeline

### 流程

#### 阶段 1：Checkout

- 拉取 `master` 最新代码

#### 阶段 2：Test

执行：

- `go test ./... -count=1`
- `cd web && npm ci && npm run test`

目的：

- 保证部署前测试通过

#### 阶段 3：Build

执行：

- 构建应用镜像
- 镜像打上本地 tag，如：
  - `jianli:latest`
  - `jianli:build-<BUILD_NUMBER>`

#### 阶段 4：Manual Approval

通过 Jenkins `input` 步骤暂停：

- 测试与构建通过后，等待人工点击“Deploy”

#### 阶段 5：Deploy

执行：

- `docker compose up -d --build`
- 或更明确地先 `docker compose build app` 再 `docker compose up -d`

#### 阶段 6：Post-check

至少执行一个健康检查：

- `curl http://127.0.0.1/api/resume`
- 若失败则标记 Pipeline 失败

## 部署行为

### 为什么选择“测试自动、部署手动”

你明确要求：

- 推送代码后自动执行流程
- 但正式上线前需要人工确认

这适合当前单机部署场景：

- 既减少每次上线手工执行命令的负担
- 又避免把 `master` 上每次提交都立即放到正式环境

## 数据持久化

必须持久化的内容：

- SQLite 数据库
- 上传头像目录
- 生产配置文件
- Jenkins 数据目录

否则容器重建后会丢失：

- 简历数据
- 草稿版本
- 访客统计
- 上传头像
- Jenkins Job 配置与构建历史

## 测试与验收边界

### 自动化验证

Jenkins Pipeline 至少要验证：

- Go 测试通过
- 前端测试通过
- 镜像构建成功
- 部署后 `/api/resume` 可访问

### 人工联调

上线后需要确认：

- 首页可访问
- 管理后台登录正常
- 保存主简历正常
- 草稿列表正常
- 头像上传正常
- PDF 导出正常
- 容器重启后数据仍存在

## 风险与约束

### Docker Socket 风险

Jenkins 容器若直接挂载 Docker Socket，则拥有宿主机 Docker 控制权。

当前阶段接受该方案，因为：

- 目标是尽快实现单机 CI/CD
- 基础设施规模较小

但后续若走正式生产规范，应考虑：

- 独立 Jenkins agent
- 更细粒度权限控制

### PDF 导出浏览器依赖

如果容器内浏览器缺失，或 `browserPath` 配置错误，则 PDF 导出一定失败。

这是本项目部署时最容易被忽略的生产问题之一。

### 同机部署耦合

Jenkins、Nginx、业务应用都在同一台机器上，运维简单，但隔离性一般。

对当前项目规模可接受。

## 验收标准

满足以下条件视为完成：

- 项目可通过 Docker Compose 在 Linux 服务器上运行
- Jenkins、Nginx、业务应用均容器化
- 推送到 `master` 后 Jenkins 自动执行测试与构建
- Jenkins 在发布前要求人工确认
- 人工确认后 Jenkins 能成功部署最新版本
- 数据库、上传文件、配置文件、Jenkins 数据均持久化
- PDF 导出在生产容器内可正常工作
