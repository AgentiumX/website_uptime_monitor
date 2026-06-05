# 网站存活探测系统 — 设计规格说明

> 日期：2026-06-05
> 状态：已审批

## 1. 概述

一个分布式网站存活探测系统，由三个独立端组成：**Web 端**（前端）、**Server 端**（Go/Gin 后端）、**Agent 端**（Go 探测节点）。用于监控 50–500 个网站的可用性、响应时间、状态码和内容匹配，并在异常时通过多种渠道告警。

### 1.1 技术选型

| 端 | 技术 | 说明 |
|---|---|---|
| Web | Vue 3 + TypeScript + Element Plus | 苹果风格高级感 UI，ECharts 图表 |
| Server | Go + Gin + SQLite + 内嵌 Prometheus TSDB | 单二进制部署，Web 静态文件嵌入 |
| Agent | Go（标准库，无 Web 框架） | 不起端口，goroutine 驱动 |

### 1.2 架构选型

**经典分层架构**：Server 作为单体中心节点，handler → service → repository 三层分层。Agent 通过 HTTP 轮询与 Server 通信。Web SPA 构建后嵌入 Server 二进制静态服务。

## 2. 系统架构

### 2.1 组件关系

```
┌──────────────────────────────────────────────────┐
│  Web 端 (Vue 3 SPA — 构建后嵌入 Server)            │
│  Dashboard · 任务配置 · 告警管理 · Agent 管理       │
└─────────────────────┬────────────────────────────┘
                      │ REST API (JSON)
┌─────────────────────▼────────────────────────────┐
│  Server 端 (Go/Gin — 单二进制)                     │
│  Handler → Service → Repository                   │
│  ┌──────────┐  ┌──────────────────────────┐      │
│  │ SQLite   │  │ 内嵌 Prometheus TSDB      │      │
│  │ 配置/告警 │  │ 时序指标                   │      │
│  └──────────┘  └──────────────────────────┘      │
└─────────────────────▲────────────────────────────┘
                      │ HTTP 轮询 (60s) / 回调结果
┌─────────────────────┴────────────────────────────┐
│  Agent 端 × N (Go — 分布在多地理位置)              │
│  任务调度器 · HTTP 探测器 · 结果上报器              │
└──────────────────────────────────────────────────┘
```

### 2.2 数据流

**正常检测流**：Agent 轮询 → Server 返回任务列表 → Agent 按频率执行检测 → Agent 回调结果 → Server 写入 TSDB + 评估告警规则 → 触发通知（如需要）

**告警流**：检测结果入库 → AlertService 评估规则（状态码/响应时间/内容匹配）→ 检查连续次数阈值 → 达到阈值 → NotifyService 发送通知（钉钉/企微/飞书/Webhook）

**数据查询流**：Web Dashboard → Server MetricsService → 查询 Prometheus TSDB（时序数据）+ SQLite（配置/告警历史）→ 聚合计算 → 返回前端渲染

## 3. 数据模型

### 3.1 SQLite 表结构

#### users（管理员）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | INTEGER PRIMARY KEY | 自增主键 |
| username | TEXT NOT NULL UNIQUE | 用户名 |
| password_hash | TEXT NOT NULL | bcrypt 哈希 |
| created_at | DATETIME | 创建时间 |

#### agents（探测点）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | INTEGER PRIMARY KEY | 自增主键 |
| name | TEXT NOT NULL | Agent 名称，如 "北京电信" |
| token | TEXT NOT NULL UNIQUE | 注册后生成的通信 token |
| status | TEXT DEFAULT 'offline' | online / offline |
| last_heartbeat | DATETIME | 最后心跳时间 |
| location | TEXT | 地理位置描述 |
| created_at | DATETIME | 创建时间 |

#### monitors（监控任务）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | INTEGER PRIMARY KEY | 自增主键 |
| name | TEXT NOT NULL | 任务名称 |
| url | TEXT NOT NULL | 监控地址 |
| method | TEXT DEFAULT 'GET' | GET / POST / PUT |
| headers | TEXT | JSON 格式自定义请求头 |
| cookie | TEXT | Cookie 值 |
| basic_auth_user | TEXT | HTTP 验证用户名 |
| basic_auth_pass | TEXT | HTTP 验证密码 |
| verify_ssl | BOOLEAN DEFAULT true | 是否验证 SSL 证书 |
| frequency | INTEGER DEFAULT 60 | 检测频率（秒） |
| proxy | TEXT | 代理服务器地址 |
| agent_ids | TEXT | JSON 数组，指定探测点 ID |
| match_type | TEXT DEFAULT 'none' | none / contains / not_contains |
| match_content | TEXT | 匹配内容 |
| status_threshold | INTEGER DEFAULT 400 | 状态码 >= 此值则告警 |
| latency_threshold | INTEGER DEFAULT 3000 | 响应时间(ms) >= 此值则告警 |
| fail_count | INTEGER DEFAULT 3 | 连续几次超过阈值后告警（1~5） |
| enabled | BOOLEAN DEFAULT true | 是否启用 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

#### alert_channels（告警通道）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | INTEGER PRIMARY KEY | 自增主键 |
| name | TEXT NOT NULL | 通道名称 |
| type | TEXT NOT NULL | dingtalk / wechat_work / feishu / webhook |
| webhook_url | TEXT NOT NULL | Webhook 地址 |
| secret | TEXT | 签名密钥（可选） |
| extra | TEXT | JSON 额外配置 |
| enabled | BOOLEAN DEFAULT true | 是否启用 |
| created_at | DATETIME | 创建时间 |

#### monitor_alert_channels（监控-告警通道关联）

| 字段 | 类型 | 说明 |
|---|---|---|
| monitor_id | INTEGER REFERENCES monitors(id) | 监控任务 ID |
| alert_channel_id | INTEGER REFERENCES alert_channels(id) | 告警通道 ID |

联合主键 (monitor_id, alert_channel_id)

#### alert_history（告警历史）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | INTEGER PRIMARY KEY | 自增主键 |
| monitor_id | INTEGER REFERENCES monitors(id) | 监控任务 ID |
| agent_id | INTEGER REFERENCES agents(id) | 触发告警的 Agent ID |
| alert_type | TEXT | status_code / latency / content_match |
| message | TEXT | 告警详情 |
| status | TEXT | firing / resolved |
| triggered_at | DATETIME | 触发时间 |
| resolved_at | DATETIME | 恢复时间 |

### 3.2 Prometheus TSDB 指标

| 指标名 | 标签 | 说明 |
|---|---|---|
| probe_success | monitor_id, agent_id, url, method | 1=成功 0=失败 |
| probe_http_status_code | monitor_id, agent_id, url | HTTP 状态码 |
| probe_duration_seconds | monitor_id, agent_id, url | 响应耗时（秒） |
| probe_ssl_expiry_seconds | monitor_id, agent_id, url | SSL 证书到期时间（秒） |
| probe_content_matched | monitor_id, agent_id, url | 1=匹配 0=不匹配 |

**PromQL 查询示例**：
- 24h 可用率：`avg_over_time(probe_success[24h])`
- P99 延迟：`quantile_over_time(0.99, probe_duration_seconds[1h])`

### 3.3 数据保留策略

- **SQLite**：告警历史保留 90 天，定期清理
- **TSDB**：原始数据 15 天，5 分钟聚合保留 90 天，1 小时聚合保留 1 年

## 4. Server API 接口

### 4.1 认证机制

- 管理员登录：`POST /api/auth/login` → 返回 JWT token
- JWT 中间件保护所有 `/api/*` 路由（Agent 接口除外）
- Agent 通信使用注册时获取的独立 token（Bearer）

### 4.2 Web 端接口（需 JWT）

```
POST   /api/auth/login              # 登录
POST   /api/auth/logout             # 登出
GET    /api/auth/me                 # 当前用户信息

GET    /api/monitors                # 监控列表（分页+搜索）
POST   /api/monitors                # 创建监控
GET    /api/monitors/:id            # 监控详情
PUT    /api/monitors/:id            # 更新监控
DELETE /api/monitors/:id            # 删除监控
PATCH  /api/monitors/:id/enabled    # 启用/禁用监控

GET    /api/monitors/:id/results    # 检测结果（分页，支持时间范围筛选）
GET    /api/monitors/:id/metrics    # 时序指标（可用率、延迟趋势、状态码分布）

GET    /api/alerts/channels         # 告警通道列表
POST   /api/alerts/channels         # 创建告警通道
PUT    /api/alerts/channels/:id     # 更新告警通道
DELETE /api/alerts/channels/:id     # 删除告警通道
POST   /api/alerts/channels/:id/test # 测试告警通道（发送测试消息）

GET    /api/alerts/history          # 告警历史（分页，按监控/时间/类型筛选）

GET    /api/agents                  # Agent 列表
DELETE /api/agents/:id              # 移除 Agent

GET    /api/dashboard/overview      # Dashboard 概览数据
```

### 4.3 Agent 端接口（需 Agent Token）

```
POST   /api/v1/agent/register       # Agent 注册（携带预共享密钥）
                                    # 请求：{ name, location, shared_secret }
                                    # 响应：{ agent_id, token }

GET    /api/v1/agent/tasks          # 获取分配给当前 Agent 的任务列表
                                    # 请求头：Authorization: Bearer <agent_token>
                                    # 响应：{ tasks: [...], updated_at }

POST   /api/v1/agent/report         # 批量上报检测结果
                                    # 请求：{ results: [{ monitor_id, status_code,
                                    #   duration_ms, content_matched, ssl_expiry,
                                    #   success, error_msg, timestamp }] }

POST   /api/v1/agent/heartbeat      # 心跳上报（每60秒）
```

### 4.4 统一响应格式

```json
// 成功
{ "code": 0, "data": { ... }, "message": "ok" }

// 分页
{ "code": 0, "data": { "list": [...], "total": 100, "page": 1, "page_size": 20 } }

// 错误
{ "code": 40001, "data": null, "message": "用户名或密码错误" }
```

### 4.5 错误码规范

| 范围 | 说明 |
|---|---|
| 0 | 成功 |
| 40001–40099 | 认证错误 |
| 40101–40199 | 参数校验错误 |
| 50001–50099 | 服务端内部错误 |

## 5. Agent 端设计

### 5.1 进程模型

单 Go 进程，不起端口，3 个核心 goroutine 协作：

- **任务同步器**：每 60 秒 HTTP 拉取 Server 任务列表
- **探测引擎**：每个任务一个 goroutine，按 frequency 定时执行 HTTP 探测
- **结果上报器**：每 30 秒批量上报 pendingResults

### 5.2 任务同步器

1. HTTP GET `/api/v1/agent/tasks`
2. 拿到最新任务列表 `newTasks[]`
3. 与本地 `currentTasks` 做 diff：
   - 新增的任务：启动新的探测 goroutine
   - 删除的任务：通过 `context cancel` 停止对应 goroutine
   - 变更的任务（配置修改）：停止旧的，启动新的
4. 更新 `currentTasks`

### 5.3 探测引擎

每个监控任务对应一个独立 goroutine，内部用 `time.Ticker` 按 `frequency` 定时执行：

1. 构造 `http.Request`（方法、URL、Header、Cookie、BasicAuth）
2. 配置 `http.Client`（代理、超时、是否验证 SSL、禁止跟随重定向）
3. 执行请求，记录状态码、响应时间、响应体
4. 执行内容匹配检查（contains / not_contains）
5. 将结果写入 `pendingResults` 队列

### 5.4 结果上报器

1. 每 30 秒从 `pendingResults` 取出一批结果（最多 100 条）
2. HTTP POST `/api/v1/agent/report`
3. 成功：清除已上报的结果
4. 失败：保留结果，下次重试（最多保留 10 分钟，超过丢弃）
5. 无结果时不上报（节省流量）

### 5.5 配置文件 (agent.yaml)

```yaml
server:
  url: "http://server-host:8080"
  shared_secret: "your-secret-key"

agent:
  name: "北京电信"
  location: "北京"

probe:
  timeout: 30            # 默认请求超时(秒)
  max_concurrent: 50     # 最大并发探测数

report:
  interval: 30           # 上报间隔(秒)
  batch_size: 100        # 单次上报最大条数
  retry_max_age: 600     # 失败结果最大保留时间(秒)
```

### 5.6 优雅关闭

监听 `SIGINT` / `SIGTERM`：停止任务同步器 → cancel 所有探测 goroutine → 等待 pendingResults 上报完毕（最多等 10 秒）→ 退出

### 5.7 Agent 注册流程

使用预共享密钥注册：Agent 首次启动时带预共享密钥向 Server 注册，Server 验证密钥后返回唯一 token，Agent 保存 token 用于后续通信。

## 6. 告警系统

### 6.1 告警评估流程

每次 Agent 上报检测结果时，AlertService 执行：

1. 检查监控是否启用（未启用则跳过）
2. 逐项检查告警条件（满足任一即计数）：
   - 状态码 >= `status_threshold`
   - 响应时间 >= `latency_threshold` ms
   - 内容匹配告警（contains / not_contains）
3. 记录连续失败次数（进程内存 `sync.Map`，key: `{monitor_id}_{agent_id}`）
4. 连续失败 >= `fail_count`？ 否 → 结束
5. 检查是否已处于告警中 → 是 → 不重复告警
6. 创建告警记录（alert_history，状态: firing）
7. 查询关联的告警通道（monitor_alert_channels）
8. 并发发送通知到所有通道

### 6.2 告警恢复

当检测结果恢复正常时：
1. 重置连续失败计数为 0
2. 如果当前有 firing 状态的告警 → 更新为 resolved，记录 resolved_at
3. 向关联通道发送恢复通知

### 6.3 告警消息模板

```
🔴 告警：{monitor_name} 异常
├── 探测点：{agent_name} ({agent_location})
├── 监控地址：{url}
├── 告警类型：{status_code / latency / content_match}
├── 详情：{具体数值/阈值对比}
├── 连续失败：{current_count} 次
└── 时间：{timestamp}

🟢 恢复：{monitor_name} 已恢复正常
├── 探测点：{agent_name}
├── 恢复正常时间：{timestamp}
└── 持续时长：{resolved_at - triggered_at}
```

### 6.4 通知渠道

| 渠道 | 协议 | 关键点 |
|---|---|---|
| 钉钉 | POST JSON | 支持签名验证（secret），Markdown 消息格式 |
| 企业微信 | POST JSON | 支持签名验证，Markdown 消息格式 |
| 飞书 | POST JSON | 支持签名验证，富文本消息卡片 |
| 自定义 Webhook | POST JSON | 用户自定义模板变量：`{{monitor_name}}`、`{{status}}`、`{{url}}` 等 |

### 6.5 失败计数存储

使用 Server 进程内存中的 `sync.Map` 存储连续失败计数。Server 重启时计数归零，对监控系统可接受（最多多告警一次，不会漏告警）。

## 7. Web 端设计

### 7.1 设计风格

苹果官网风格，基于 Element Plus 组件库 + 自定义主题定制：

- **色彩体系**：主文字 #1D1D1F、次文字 #86868B、主色 #0071E3、健康 #34C759、告警 #FF3B30、背景 #FBFBFD
- **设计原则**：大量留白、圆角卡片(16px) + 细微投影、SF Pro 系统字体栈、字号层级清晰
- **图表**：ECharts，配色与主色调统一
- **暗色模式**：通过 CSS 变量切换

### 7.2 页面结构

左侧边栏导航，顶部显示当前用户名和通知铃铛。登录后默认进入仪表盘。

导航项：仪表盘、监控列表、告警管理、探测点、设置

### 7.3 Dashboard 看板

- **顶部统计卡片（4列）**：监控总数、在线率、活跃告警数、平均响应时间 — 每张卡片含数值和趋势箭头
- **可用率趋势图（宽）**：24h / 7d / 30d 切换，折线图 + 面积填充，异常点红点标注
- **健康度分布（窄）**：环形图，健康/警告/异常比例
- **最近告警列表**：firing(红) / warning(橙) / resolved(绿)，含时间和详情
- **TOP5 监控站点**：按可用率排序，含延迟数据，状态指示灯

### 7.4 监控列表页

表格展示所有监控任务，支持搜索/筛选/排序。每行显示：名称、URL、状态指示灯、可用率、平均延迟、最近检测时间、操作按钮（编辑/删除/启停）。

### 7.5 任务配置表单

分步表单：
1. **基本信息**：任务名称、监控地址
2. **请求配置**：请求方法、HTTP 请求头、Cookie、HTTP 验证用户名/密码、代理服务器、是否验证证书
3. **告警规则**：匹配类型（contains/not_contains）、匹配内容、状态码阈值、响应时间阈值、连续失败次数(1~5)、监控频率
4. **通知通道**：选择关联的告警通道、选择探测点（agent_ids）、是否启用

带表单校验，提交后可立即启用。

### 7.6 监控详情页

单站点深度视图：可用率趋势图、响应时间分布图、状态码统计饼图、各探测点对比表、告警时间线、历史记录表（可导出 CSV）。

### 7.7 登录页

居中卡片式登录表单，品牌 Logo + 标题，用户名/密码输入框，登录按钮。全屏渐变背景，简洁大气。

## 8. Server 端设计

### 8.1 项目结构

```
server/
├── cmd/
│   └── main.go              # 入口，初始化所有组件
├── internal/
│   ├── config/
│   │   └── config.go        # 配置加载（YAML）
│   ├── middleware/
│   │   ├── auth.go          # JWT 认证中间件
│   │   └── cors.go          # CORS 中间件
│   ├── handler/
│   │   ├── auth.go          # 登录/登出
│   │   ├── monitor.go       # 监控任务 CRUD
│   │   ├── alert.go         # 告警通道 + 历史
│   │   ├── agent.go         # Agent 管理
│   │   ├── dashboard.go     # Dashboard 概览
│   │   └── metrics.go       # 时序指标查询
│   ├── service/
│   │   ├── auth.go          # 认证逻辑
│   │   ├── monitor.go       # 监控业务逻辑
│   │   ├── alert.go         # 告警评估逻辑
│   │   ├── agent.go         # Agent 管理逻辑
│   │   ├── metrics.go       # TSDB 读写
│   │   └── notify.go        # 通知发送（钉钉/企微/飞书/Webhook）
│   ├── repository/
│   │   ├── sqlite.go        # SQLite 初始化 + 通用方法
│   │   ├── user.go          # 用户数据访问
│   │   ├── monitor.go       # 监控数据访问
│   │   ├── alert.go         # 告警数据访问
│   │   └── agent.go         # Agent 数据访问
│   ├── model/
│   │   └── model.go         # 所有数据模型定义
│   └── prom/
│       └── tsdb.go          # 内嵌 Prometheus TSDB 初始化
├── web/
│   └── dist/                # Vue 3 构建产物（go:embed 嵌入）
├── config.yaml              # 配置文件
├── go.mod
└── go.sum
```

### 8.2 配置文件 (config.yaml)

```yaml
server:
  port: 8080
  jwt_secret: "your-jwt-secret"
  admin_username: "admin"
  admin_password: "admin123"     # 首次启动时 bcrypt 加密存储

agent:
  shared_secret: "your-shared-secret"
  heartbeat_timeout: 120         # 心跳超时(秒)，超过则标记 offline

database:
  path: "./data/uptime.db"       # SQLite 文件路径

tsdb:
  path: "./data/tsdb"            # Prometheus TSDB 存储路径
  retention_raw: 360h            # 原始数据保留 15 天
  retention_5m: 2160h            # 5分钟聚合保留 90 天
  retention_1h: 8760h            # 1小时聚合保留 1 年

alert:
  history_retention: 90          # 告警历史保留天数
```

### 8.3 单二进制部署

- Web 前端构建产物通过 `go:embed` 嵌入 Server 二进制
- Server 的 Gin 静态中间件服务 `/*` 路径的 SPA 文件
- SQLite 文件自动创建在配置指定路径
- Prometheus TSDB 文件自动创建在配置指定路径
- 部署只需：一个二进制文件 + 一个 config.yaml

## 9. Web 端项目结构

```
web/
├── src/
│   ├── main.ts
│   ├── App.vue
│   ├── router/
│   │   └── index.ts           # 路由配置 + 导航守卫
│   ├── stores/
│   │   ├── auth.ts            # 登录状态
│   │   └── monitor.ts         # 监控数据
│   ├── api/
│   │   ├── request.ts         # Axios 封装（拦截器、JWT 注入）
│   │   ├── auth.ts            # 登录 API
│   │   ├── monitor.ts         # 监控 API
│   │   ├── alert.ts           # 告警 API
│   │   ├── agent.ts           # Agent API
│   │   └── dashboard.ts       # Dashboard API
│   ├── views/
│   │   ├── login/
│   │   │   └── LoginView.vue
│   │   ├── dashboard/
│   │   │   └── DashboardView.vue
│   │   ├── monitor/
│   │   │   ├── MonitorList.vue
│   │   │   ├── MonitorCreate.vue
│   │   │   ├── MonitorEdit.vue
│   │   │   └── MonitorDetail.vue
│   │   ├── alert/
│   │   │   ├── AlertChannelList.vue
│   │   │   ├── AlertChannelForm.vue
│   │   │   └── AlertHistory.vue
│   │   └── agent/
│   │       └── AgentList.vue
│   ├── components/
│   │   ├── layout/
│   │   │   ├── AppLayout.vue
│   │   │   └── Sidebar.vue
│   │   ├── charts/
│   │   │   ├── UptimeChart.vue
│   │   │   ├── LatencyChart.vue
│   │   │   └── HealthDonut.vue
│   │   └── common/
│   │       ├── StatusDot.vue
│   │       └── StatCard.vue
│   └── styles/
│       ├── variables.scss      # 苹果风格 CSS 变量
│       └── element-override.scss # Element Plus 主题覆盖
├── index.html
├── vite.config.ts
├── tsconfig.json
└── package.json
```

## 10. 部署方式

### 10.1 单二进制部署

```bash
# 构建 Web
cd web && npm install && npm run build

# 构建 Server（自动嵌入 web/dist/）
cd ../server && go build -o uptime-server cmd/main.go

# 部署
./uptime-server -c config.yaml
```

### 10.2 Agent 部署

```bash
cd agent && go build -o uptime-agent
./uptime-agent -c agent.yaml
```

### 10.3 目录结构

```
/opt/uptime-monitor/
├── uptime-server          # Server 二进制
├── config.yaml            # Server 配置
├── data/                  # 自动创建
│   ├── uptime.db          # SQLite
│   └── tsdb/              # Prometheus TSDB
```
