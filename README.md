<h1 align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Vue-3.5+-4FC08D?style=flat-square&logo=vuedotjs&logoColor=white" alt="Vue" />
  <img src="https://img.shields.io/badge/License-MIT-blue?style=flat-square" alt="License" />
</h1>

<h1 align="center">🔍 Uptime Monitor</h1>

<p align="center">
  <strong>分布式网站存活探测系统</strong><br>
  <em>实时监控 · 多渠道告警 · 苹果风格 Dashboard</em>
</p>

<p align="center">
  <a href="#-快速开始">快速开始</a> •
  <a href="#-功能特性">功能特性</a> •
  <a href="#-系统架构">架构</a> •
  <a href="#-api-文档">API</a> •
  <a href="#-配置说明">配置</a>
</p>

---

## ✨ 功能特性

### 📊 监控能力
- **HTTP 探测**：支持 GET / POST / PUT 请求方法
- **灵活配置**：自定义请求头、Cookie、Basic Auth、代理服务器
- **SSL 验证**：可选证书校验，自动检测证书过期时间
- **内容匹配**：支持包含/不包含关键词匹配验证
- **多探测点**：分布式 Agent 从不同地理位置发起探测

### 🔔 智能告警
- **多维度阈值**：状态码、响应时间、内容匹配三种告警条件
- **连续失败计数**：可配置 1-5 次连续失败后触发告警，避免误报
- **自动恢复**：服务恢复后自动解除告警并发送恢复通知
- **4 种通知渠道**：钉钉、企业微信、飞书、自定义 Webhook

### 📈 数据看板
- **实时概览**：监控总数、在线率、活跃告警、平均响应时间
- **趋势图表**：可用率趋势（24h/7d/30d）、健康度分布、延迟趋势
- **告警时间线**：历史告警记录，支持多维度筛选和导出

### 🚀 部署特性
- **单二进制**：Server 仅一个可执行文件（~39MB），内嵌 Web UI
- **零依赖**：无需安装数据库，SQLite 自动创建
- **跨平台**：Go 编译，支持 Linux / macOS / Windows

---

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────┐
│  Web Dashboard (Vue 3 SPA — 内嵌于 Server 二进制)    │
│  📊 仪表盘  ⚙️ 监控配置  🔔 告警管理  🤖 探测点管理  │
└───────────────────────┬─────────────────────────────┘
                        │ REST API (JSON)
┌───────────────────────▼─────────────────────────────┐
│  Server (Go / Gin — 单二进制部署)                     │
│                                                      │
│  Handler → Service → Repository                      │
│  ┌────────────┐  ┌────────────────────────────┐     │
│  │  SQLite    │  │  Prometheus Metrics        │     │
│  │  配置/告警  │  │  时序指标 (可用率/延迟/状态码) │     │
│  └────────────┘  └────────────────────────────┘     │
└───────────────────────▲─────────────────────────────┘
                        │ HTTP 轮询 (60s) + 结果回调
┌───────────────────────┴─────────────────────────────┐
│  Agent × N (Go — 分布式部署)                          │
│  🕐 任务调度器  🔍 HTTP 探测引擎  📤 结果上报器       │
└─────────────────────────────────────────────────────┘
```

---

## 🛠️ 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| **Server** | Go + Gin | HTTP API 服务，handler → service → repository 三层架构 |
| **数据库** | SQLite + GORM | 零配置嵌入式数据库，纯 Go 驱动（无需 CGO） |
| **指标** | Prometheus client_golang | 内存时序指标，支持 PromQL 查询 |
| **Web** | Vue 3 + TypeScript | 单页应用，Pinia 状态管理，Vue Router |
| **UI** | Element Plus + ECharts | 苹果风格定制主题，响应式数据可视化 |
| **Agent** | Go 标准库 | 无框架依赖，goroutine 驱动的轻量探测节点 |

---

## 🚀 快速开始

### 环境要求

- **Go** 1.22+
- **Node.js** 18+（仅开发时需要）

### 一键构建

```bash
# 克隆仓库
git clone https://github.com/AgentiumX/website_uptime_monitor.git
cd website_uptime_monitor

# 执行构建脚本（自动构建 Web → 嵌入 Server → 编译 Agent）
chmod +x build.sh
./build.sh
```

构建完成后产出：
```
server/uptime-server    # Server 二进制 (~39MB，内嵌 Web UI)
agent/uptime-agent      # Agent 二进制 (~9MB)
```

### 启动服务

**1. 启动 Server**

```bash
cd server
./uptime-server -c config.yaml
```

Server 启动在 `http://localhost:8080`，打开浏览器访问即可看到 Dashboard。

默认管理员账号：
- 用户名：`admin`
- 密码：`admin123`

**2. 启动 Agent**

```bash
cd agent
./uptime-agent -c agent.yaml
```

Agent 启动后自动向 Server 注册，开始接收并执行监控任务。

---

## 📖 使用指南

### 创建监控任务

1. 登录 Web Dashboard
2. 点击侧边栏「监控列表」→「新增监控」
3. 填写分步表单：
   - **基本信息**：任务名称、监控地址
   - **请求配置**：请求方法、Header、Cookie、认证、代理
   - **告警规则**：匹配类型、阈值、连续失败次数、检测频率
   - **通知配置**：选择告警通道和探测点

### 配置告警通道

1. 点击侧边栏「告警管理」→「新增通道」
2. 选择通道类型：
   - **钉钉**：填写 Webhook URL + 签名密钥（可选）
   - **企业微信**：填写 Webhook URL
   - **飞书**：填写 Webhook URL + 签名密钥（可选）
   - **自定义 Webhook**：填写 URL，接收 JSON POST
3. 点击「测试」验证通道是否可用

### 部署 Agent 探测点

在不同地理位置的服务器上部署 Agent：

```bash
# 修改 agent.yaml 中的 Server 地址和共享密钥
server:
  url: "http://your-server:8080"
  shared_secret: "your-secret"
agent:
  name: "北京电信"
  location: "北京"

# 启动
./uptime-agent -c agent.yaml
```

---

## 📡 API 文档

### Web API（需 JWT 认证）

```
POST   /api/auth/login              # 登录
GET    /api/monitors                # 监控列表
POST   /api/monitors                # 创建监控
PUT    /api/monitors/:id            # 更新监控
DELETE /api/monitors/:id            # 删除监控
GET    /api/monitors/:id/metrics    # 时序指标
GET    /api/dashboard/overview      # 仪表盘概览
GET    /api/alerts/channels         # 告警通道列表
GET    /api/alerts/history          # 告警历史
```

### Agent API（需 Agent Token）

```
POST   /api/v1/agent/register       # 注册（携带共享密钥）
GET    /api/v1/agent/tasks          # 获取监控任务
POST   /api/v1/agent/report         # 上报检测结果
POST   /api/v1/agent/heartbeat      # 心跳保活
```

完整 API 文档见 [设计规格](docs/superpowers/specs/2026-06-05-uptime-monitor-design.md)。

---

## ⚙️ 配置说明

### Server 配置 (`server/config.yaml`)

```yaml
server:
  port: 8080                          # 监听端口
  jwt_secret: "your-jwt-secret"       # JWT 签名密钥
  admin_username: "admin"             # 管理员用户名
  admin_password: "admin123"          # 管理员密码

agent:
  shared_secret: "agent-shared-secret" # Agent 注册密钥
  heartbeat_timeout: 120               # 心跳超时（秒）

database:
  path: "./data/uptime.db"             # SQLite 数据库路径

tsdb:
  path: "./data/tsdb"                  # 时序数据存储路径
  retention_raw: "360h"                # 原始数据保留 15 天
  retention_5m: "2160h"                # 5分钟聚合保留 90 天
  retention_1h: "8760h"               # 1小时聚合保留 1 年

alert:
  history_retention: 90                # 告警历史保留天数
```

### Agent 配置 (`agent/agent.yaml`)

```yaml
server:
  url: "http://localhost:8080"         # Server 地址
  shared_secret: "agent-shared-secret" # 注册密钥（需与 Server 一致）

agent:
  name: "本地探测点"                    # Agent 名称
  location: "本地"                     # 地理位置

probe:
  timeout: 30                          # 请求超时（秒）
  max_concurrent: 50                   # 最大并发探测数

report:
  interval: 30                         # 结果上报间隔（秒）
  batch_size: 100                      # 单次上报最大条数
  retry_max_age: 600                   # 失败结果最大保留时间（秒）
```

---

## 📁 项目结构

```
website_uptime_monitor/
├── server/                    # Server 端（Go/Gin）
│   ├── cmd/
│   │   ├── main.go           # 入口，路由注册，后台任务
│   │   └── web/dist/         # 嵌入的 Web 构建产物
│   ├── internal/
│   │   ├── config/           # 配置加载
│   │   ├── handler/          # HTTP Handler（认证/监控/告警/Agent）
│   │   ├── middleware/       # JWT + Agent Token 中间件
│   │   ├── service/          # 业务逻辑（告警评估/通知/指标）
│   │   ├── repository/       # 数据访问层（SQLite）
│   │   ├── model/            # 数据模型 + DTO
│   │   └── background/       # 后台任务（心跳检测/历史清理）
│   ├── config.yaml
│   └── go.mod
├── agent/                     # Agent 端（Go）
│   ├── main.go               # 入口，信号处理，优雅关闭
│   ├── config.go             # 配置加载
│   ├── register.go           # 注册 + Token 持久化
│   ├── scheduler.go          # 任务同步器（60s 轮询 + diff）
│   ├── probe.go              # HTTP 探测引擎
│   ├── reporter.go           # 结果上报器（30s 批量 + 重试）
│   ├── agent.yaml
│   └── go.mod
├── web/                       # Web 端（Vue 3）
│   ├── src/
│   │   ├── views/            # 页面（登录/仪表盘/监控/告警/Agent）
│   │   ├── components/       # 组件（布局/图表/通用）
│   │   ├── api/              # API 客户端
│   │   ├── stores/           # Pinia 状态管理
│   │   ├── router/           # 路由配置
│   │   └── styles/           # 苹果风格主题
│   └── package.json
├── docs/                      # 文档
│   └── superpowers/
│       ├── specs/            # 设计规格
│       └── plans/            # 实现计划
├── build.sh                   # 一键构建脚本
└── README.md
```

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

## 📄 License

[MIT License](LICENSE)
