# Docker 单容器 + 复用现有 Nginx 部署设计

## 概述

当前仓库中已经引入了 Docker + Jenkins + 项目内 Nginx 的全容器化部署方案，但对当前项目规模来说，这会带来不必要的复杂度。与此同时，服务器上已经有一个长期运行的 `nginx:alpine` 容器，负责统一处理 80/443 入口和证书。因此本次调整将部署方案简化为：本项目只运行一个 `app` 容器，前端静态资源与后端 API 统一由该容器提供，再通过服务器上现有的 Nginx 容器反向代理到该 `app` 容器。

## 目标

- 项目通过 Docker 以单 `app` 容器方式运行
- 不再引入 Jenkins
- 不再引入项目内单独的 Nginx 容器
- 保留 Docker 化部署能力
- 复用服务器上现有的 Nginx 容器作为统一入口
- 通过手动 `git pull` + `docker compose up -d --build` 完成发布
- 数据、上传文件、配置文件继续通过挂载持久化
- 正式环境仍支持 PDF 导出能力

## 非目标

- 不做 CI/CD
- 不做 Jenkins Pipeline
- 不新增第二个 Nginx 容器
- 不做多环境编排
- 不做自动回滚
- 不接入镜像仓库
- 不引入 Kubernetes

## 总体架构

部署后结构为：

- 服务器现有 `nginx:alpine` 容器
  - 继续监听 `80/443`
  - 负责 TLS、统一入口与反代
- 本项目新增 `app` 容器
  - 运行 Go 后端
  - 提供前端静态文件
  - 提供 `/api` 和 `/uploads`
  - 提供 PDF 导出能力

### 请求路径

外部访问路径：

- `https://wenemoji.com/` -> 现有 Nginx -> `app`
- `https://wenemoji.com/api/*` -> 现有 Nginx -> `app`
- `https://wenemoji.com/uploads/*` -> 现有 Nginx -> `app`

## 容器化策略

### 应用容器 `app`

单容器承载：

- Go 编译后的后端二进制
- 前端 `dist` 静态文件
- PDF 导出所需浏览器（推荐 `chromium`）

原因：

- 项目体量较小
- 前端请求本身使用相对路径 `/api`、`/uploads`
- 单容器部署更直接，降低维护复杂度

### 为什么不再新增项目内 Nginx

因为服务器已经有一个长期运行的 Nginx 容器：

- 已经对外占用了 `80/443`
- 已经挂载证书与全局配置
- 已经承担统一入口职责

继续为本项目再起一个 Nginx 只会增加：

- 容器数量
- 配置重复
- 调试复杂度

因此本项目内部不再保留 `nginx` 服务定义，只保留 `app`。

## 运行目录设计

建议服务器侧使用以下目录：

- `/srv/jianli/app-src`
  - 项目代码目录
- `/srv/jianli/config`
  - 生产版 `config.production.json`
- `/srv/jianli/data`
  - `resume.db`
  - `uploads/avatars`
- `/srv/jianli/deploy.env`
  - Compose 所需环境变量

## 配置设计

### 生产配置文件

`config.production.json` 中保留：

- `browserPath`
- `port`
- `dbPath`
- `frontendOrigin`

并要求：

- `AUTH_KEY` 通过环境变量提供，不再写入文件
- `browserPath` 使用 Linux 路径，例如 `/usr/bin/chromium`
- `dbPath` 使用容器内路径，例如 `/app/data/resume.db`
- `frontendOrigin` 使用正式域名 `https://wenemoji.com`

### 环境变量

至少包含：

- `AUTH_KEY`

可选保留：

- `APP_CONTAINER_NAME`
- `APP_PORT`
- `APP_CONFIG_PATH`
- `APP_DATA_DIR`

## Docker 设计

### `Dockerfile`

保留多阶段构建，但只保留一个最终运行目标：

1. 前端构建阶段
2. 后端构建阶段
3. `app` 运行阶段

删除项目内 `nginx` 运行目标。

### `docker-compose.yml`

只保留一个 `app` 服务。

该服务需要：

- 构建或运行 `app` 镜像
- 挂载：
  - `config.production.json`
  - `data/`
- 暴露内部端口，如 `8080`
- 如需宿主机调试，可映射 `127.0.0.1:8080:8080`

推荐只绑定到本机回环地址，避免绕过现有 Nginx：

- `127.0.0.1:8080:8080`

## 手动部署流程

部署流程改为：

1. 登录服务器
2. 进入项目目录
3. `git pull`
4. `docker compose up -d --build`
5. 检查健康状态

也就是：

```bash
cd /srv/jianli/app-src
git pull
docker compose --env-file /srv/jianli/deploy.env up -d --build
```

## 现有 Nginx 的接入方式

你已有的 Nginx 容器挂载了：

- `./nginx/nginx.conf:/etc/nginx/nginx.conf`
- 以及多个静态目录和证书目录

因此本项目只需要在你现有的 Nginx 配置中新增一个针对 `wenemoji.com` 的反代段，转发到本项目容器。

建议反代目标：

- `http://host.docker.internal:8080`
  或
- 若你的 Nginx 容器和本项目容器接入同一个 Docker 网络，则直接用容器名

首版建议优先选择你当前环境最容易接通的一种，不在本项目里额外引入网络复杂度。

## 文档与运维影响

部署文档需要改为：

- 删除 Jenkins 安装与 Jenkins Job 配置内容
- 删除项目内 Nginx 服务说明
- 改为说明如何：
  - 准备 `deploy.env`
  - 准备 `config.production.json`
  - 启动 `app` 容器
  - 修改现有 Nginx 配置进行反代

## 测试与验收边界

### 自动化验证

至少验证：

- Go 测试通过
- 前端测试通过
- 前端构建通过
- Docker Compose 配置静态校验通过

### 人工联调

部署后需要确认：

- `wenemoji.com` 首页可访问
- `/api/resume` 返回 200
- 登录可用
- 草稿功能可用
- 上传头像可用
- `/uploads/...` 可访问
- PDF 导出可用

## 风险与约束

### 现有 Nginx 与项目容器网络连通性

因为复用的是外部现有 Nginx 容器，而不是同一份 Compose 中的 `nginx` 服务，所以真正部署时必须确认：

- 现有 Nginx 容器能访问到本项目 `app` 容器

这可能需要：

- 共用 Docker 网络
- 或绑定宿主机回环端口再由 Nginx 转发

### PDF 导出浏览器依赖

如果容器内浏览器缺失，或 `browserPath` 配置错误，则 PDF 导出仍会失败。

## 验收标准

满足以下条件视为完成：

- 仓库部署资产已简化为单 `app` 容器
- 不再依赖 Jenkins
- 不再引入项目内独立 Nginx 容器
- 可通过手动 `git pull` + `docker compose up -d --build` 发布
- 能通过现有服务器 Nginx 容器访问 `wenemoji.com`
- 数据与上传文件持久化正常
- PDF 导出可用
