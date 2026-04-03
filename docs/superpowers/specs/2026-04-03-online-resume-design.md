# 在线简历系统设计文档

## 概述

将现有 Word 格式简历转化为在线版本，支持简历展示、编辑、PDF导出和访客记录功能。技术栈：Go + React，响应式设计（移动端 + PC端）。

## 技术架构

### 方案：单体全栈应用

- **后端**：Go + Gin，单个服务
- **前端**：React SPA，移动端优先，响应式
- **存储**：SQLite（轻量，单文件，免运维）
- **PDF导出**：chromedp 渲染 HTML 为 PDF
- **部署**：单个二进制 + 前端静态文件

### 选择理由

个人简历项目，YAGNI 原则——SQLite 完全够用，单体部署省心。复杂度集中在前端交互而非后端架构。

## API 设计

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/resume` | 获取简历 JSON | 否 |
| PUT | `/api/resume` | 更新简历 | 需要 |
| GET | `/api/resume/pdf` | 导出 PDF（A4，chromedp 渲染 `/print` 路由） | 否 |
| POST | `/api/auth/verify` | 验证 Key，返回 ` "valid": true }` | 否 |
| GET | `/api/visitors` | 获取访客记录，支持 `?page=&limit=&days=` 分页筛选 | 需要 |
| GET | `/api/visitors/stats` | 聚合统计（总量/今日/独立访客/平均停留），支持 `?days=` | 需要 |
| POST | `/api/visitors` | 记录访客（IP 去重，同 IP 10分钟内不重复记录） | 否 |
| PATCH | `/api/visitors/:id` | 更新停留时长（前端 `navigator.sendBeacon` 页面离开时上报） | 否 |

### 认证方式

请求头 `X-Auth-Key` 匹配环境变量 `AUTH_KEY`。访客打开简历页面无需认证，编辑和访客统计功能需要 Key 登录。

- 前端登录后将 Key 存入 `sessionStorage`，axios 拦截器自动注入请求头
- 收到 401 响应时自动跳转到登录页
- 关闭浏览器 Tab 后 session 失效

### 前端路由

| 路径 | 页面 | 认证 |
|------|------|------|
| `/` | 简历展示页 | 否 |
| `/login` | Key 登录页 | 否 |
| `/edit` | 简历编辑页 | 需要（未认证跳转 `/login`） |
| `/visitors` | 访客统计页 | 需要（未认证跳转 `/login`） |
| `/print` | PDF 打印专用页（隐藏导航，A4 排版） | 否 |

### API 响应格式

```json
// 成功
{ "code": 0, "data": { ... } }

// 错误
{ "code": 40001, "message": "invalid auth key" }
```

常用 HTTP 状态码：200（成功）、400（参数错误）、401（未认证）、500（服务端错误）。

## 数据模型

### Resume（JSON 存储于 SQLite）

Resume 表结构：

```sql
CREATE TABLE resume (
  id INTEGER PRIMARY KEY DEFAULT 1,
  data TEXT NOT NULL,           -- JSON 格式简历数据
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- 只有一行数据，id 固定为 1
```

JSON 结构：

```json
{
  "profile": {
    "name": "温庆京",
    "title": "Golang后端工程师",
    "age": 25,
    "gender": "男",
    "education": "本科",
    "experience": "0.9年",
    "location": "浙江杭州",
    "hometown": "江西赣州",
    "phone": "17620096266",
    "email": "3219431643@qq.com"
  },
  "skills": ["Go", "MySQL", "Redis", "Kafka", "Docker", "Gin", "gRPC", "K8s", "DDD", "TDD"],
  "jobIntention": {
    "position": "Golang后端工程师",
    "cities": ["深圳", "杭州", "厦门"],
    "availability": "周内到岗"
  },
  "education": [{
    "school": "江西财经大学现代经济管理学院",
    "major": "计算机科学与技术",
    "degree": "本科",
    "startDate": "2023.9",
    "endDate": "2025.7"
  }],
  "workExperience": [{
    "company": "杭州云缊科技有限公司",
    "role": "Golang后端开发",
    "startDate": "2025.5",
    "endDate": "2026.2",
    "description": ["重构项目架构...", "优化慢查询...", "..."]
  }],
  "projects": [{
    "name": "AI Gateway - LLM API统一网关",
    "startDate": "2025.12",
    "endDate": "2026.1",
    "description": ["OpenAI↔Anthropic双向协议转换...", "..."]
  }],
  "awards": [
    { "date": "2022.9", "title": "国家励志奖学金" },
    { "date": "2021.9", "title": "国家励志奖学金" }
  ]
}
```

### Visitor 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键自增 |
| ip | TEXT | 访客 IP |
| city | TEXT | 城市（IP 解析） |
| device | TEXT | 设备类型 |
| browser | TEXT | 浏览器 |
| os | TEXT | 操作系统 |
| visit_time | DATETIME | 访问时间 |
| duration | INTEGER | 停留时长（秒） |

## 前端页面设计

### 响应式策略

- `< 768px`：移动端布局
- `>= 768px`：PC 端布局

### 页面列表（共 4 个页面，各有移动端/PC端两种布局）

#### 1. 简历展示页（公开）

**移动端**：单栏垂直滚动
- 状态栏
- 个人信息区（头像、姓名、职位、基本信息、联系方式）
- 个人技能（标签式展示）
- 工作经历（卡片展示：公司名+时间+职位+描述）
- 项目经历（多张卡片）
- 教育经历
- 荣誉奖项
- 导出 PDF 按钮

**PC端**：顶部导航 + 两栏布局
- 顶部导航栏：Logo + 站点名 | 简历/编辑/访客导航 + 导出PDF按钮
- 左栏（360px）：个人信息卡片、技能卡片、教育卡片、荣誉卡片
- 右栏（自适应）：工作经历卡片、项目经历卡片（含3个项目）

#### 2. 简历编辑页（需认证）

**移动端**：顶部取消/保存导航 + 垂直表单
- 基本信息：姓名、职位、手机、邮箱输入框
- 技能编辑：标签 + X删除按钮 + 添加按钮
- 工作经历：卡片式表单（公司、职位、时间、工作内容文本框）+ 添加按钮
- 项目经历：卡片式表单（项目名、时间、描述文本框）+ 添加按钮
- 教育经历：学校、专业、学历、时间 + 添加按钮
- 荣誉奖项：日期 + 奖项名称 + 添加/删除

**PC端**：顶部导航（取消+保存按钮）+ 两栏表单
- 左栏（480px）：基本信息表单（双列排列）+ 技能编辑 + 教育经历 + 荣誉奖项
- 右栏（自适应）：工作经历表单 + 项目经历表单

#### 3. 访客统计页（需认证）

**移动端**：
- 3个统计卡片（总访问/今日/独立访客）
- 7日趋势折线图 + 筛选器
- 访客列表（卡片式：设备图标 + 地区/浏览器 + 时间/IP）

**PC端**：
- 4个统计卡片横排（总访问/今日/独立访客/平均停留）
- 趋势图全宽 + 时间筛选器（7天/30天/全部）
- 表格式访客记录（IP/地区/设备浏览器/访问时间/停留时长）

#### 4. Key 登录页

**移动端**：全屏居中
- 锁图标 + 标题 + 描述
- 密钥输入框 + 登录按钮

**PC端**：灰色背景 + 居中白色卡片（带阴影）
- 内容同移动端，卡片宽 420px

### 设计风格

- **风格**：简约专业
- **背景**：移动端白色，PC端浅灰(#F5F6F8)
- **主色**：#4A90D9（蓝色）
- **文字**：#1A1A1A（标题）、#666（正文）、#999（辅助）
- **字体**：Inter
- **圆角**：卡片 12px，按钮 8px，标签 14px
- **内容卡片**：白色背景 + #F0F0F0 边框

## 项目结构

```
jianli/
├── cmd/
│   └── server/
│       └── main.go          # 入口
├── internal/
│   ├── handler/              # HTTP 处理器
│   │   ├── resume.go
│   │   ├── visitor.go
│   │   └── auth.go
│   ├── middleware/            # Gin 中间件
│   │   ├── auth.go           # Key 认证
│   │   └── cors.go           # CORS（开发环境）
│   ├── config/               # 配置管理
│   │   └── config.go         # PORT, AUTH_KEY, DB_PATH 等
│   ├── model/                # 数据模型
│   │   ├── resume.go
│   │   └── visitor.go
│   ├── store/                # 数据存储层
│   │   ├── sqlite.go
│   │   ├── resume_store.go
│   │   └── visitor_store.go
│   └── pdf/                  # PDF 导出
│       └── export.go
├── web/                      # React 前端
│   ├── src/
│   │   ├── components/       # 共用组件
│   │   ├── pages/            # 页面
│   │   │   ├── ResumePage.tsx
│   │   │   ├── EditPage.tsx
│   │   │   ├── VisitorPage.tsx
│   │   │   └── LoginPage.tsx
│   │   ├── hooks/            # 自定义 hooks
│   │   ├── api/              # API 调用
│   │   └── App.tsx
│   └── package.json
├── go.mod
├── go.sum
└── Makefile
```

## 关键技术决策

| 决策 | 选择 | 理由 |
|------|------|------|
| 数据库 | SQLite | 单用户场景，零配置，单文件备份 |
| PDF导出 | chromedp | 高保真渲染，支持 CSS 样式 |
| 前端路由 | React Router | SPA 标准方案 |
| 内联编辑 | 表单编辑模式 | 独立编辑页面使用传统表单，可靠性优于 contentEditable |
| IP 解析 | ip2region 离线库 | 无外部依赖，响应快 |
| 部署 | go:embed 嵌入前端 | 单二进制分发 |

## 非功能需求

- **性能**：首屏加载 < 2s，API 响应 < 200ms
- **安全**：Key 认证保护编辑和统计功能，访客 IP 脱敏展示（显示前两段 + `***.***`，如 `112.17.***.***`）
- **兼容性**：Chrome/Safari/Firefox 最新版，iOS/Android 移动端浏览器

### 配置项

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `PORT` | `8080` | 服务监听端口 |
| `AUTH_KEY` | （必填） | 管理后台访问密钥 |
| `DB_PATH` | `./data/resume.db` | SQLite 数据库文件路径 |

### 数据库初始化

应用首次启动时自动建表（`AutoMigrate`），并从 Word 简历导入初始数据作为种子数据。
