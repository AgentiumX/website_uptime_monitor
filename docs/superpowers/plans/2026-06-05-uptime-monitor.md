# 网站存活探测系统 实现计划

> **面向 AI 代理的工作者：** 必需子技能：使用 superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 逐任务实现此计划。步骤使用复选框（`- [ ]`）语法来跟踪进度。

**目标：** 构建一个三端（Web/Server/Agent）分布式网站存活探测系统，支持 HTTP 可用性监控、多渠道告警、苹果风格 Dashboard。

**架构：** 经典分层架构。Server（Go/Gin）作为单体中心节点，内嵌 SQLite + Prometheus TSDB + Web 静态文件。Agent（Go）通过 HTTP 轮询获取任务并上报结果。Web 端（Vue 3 + Element Plus）构建后嵌入 Server。

**技术栈：** Go 1.22+ / Gin / GORM + SQLite (modernc.org/sqlite) / prometheus/client_golang / golang-jwt/jwt/v5 / Vue 3 / TypeScript / Element Plus / ECharts / Vite

---

## 文件结构总览

### Server (`server/`)
| 文件 | 职责 |
|---|---|
| `cmd/main.go` | 入口，初始化所有组件，启动后台任务，优雅关闭 |
| `internal/config/config.go` | 配置结构体定义 + YAML 加载 |
| `internal/model/model.go` | 所有数据模型 + API DTO 定义 |
| `internal/handler/response.go` | 统一响应格式 Success/Paginated/Error |
| `internal/handler/auth.go` | 登录/登出/当前用户 |
| `internal/handler/monitor.go` | 监控任务 CRUD + 启停 |
| `internal/handler/agent.go` | Agent 注册/任务获取/结果上报/心跳 |
| `internal/handler/alert.go` | 告警通道 CRUD + 告警历史查询 |
| `internal/handler/dashboard.go` | Dashboard 概览数据聚合 |
| `internal/handler/metrics.go` | 时序指标查询 |
| `internal/middleware/auth.go` | JWT 认证中间件 |
| `internal/middleware/agent_auth.go` | Agent Token 认证中间件 |
| `internal/service/auth.go` | JWT 生成/验证，密码校验 |
| `internal/service/monitor.go` | 监控业务逻辑 |
| `internal/service/agent.go` | Agent 注册/任务分配 |
| `internal/service/alert.go` | 告警评估（阈值检查、连续计数、创建/恢复告警） |
| `internal/service/metrics.go` | TSDB 写入 + 查询 |
| `internal/service/notify.go` | 通知发送（钉钉/企微/飞书/Webhook） |
| `internal/repository/db.go` | SQLite 初始化 + 自动迁移 + 种子数据 |
| `internal/repository/user.go` | 用户数据访问 |
| `internal/repository/monitor.go` | 监控数据访问 |
| `internal/repository/agent.go` | Agent 数据访问 |
| `internal/repository/alert.go` | 告警数据访问 |
| `internal/background/jobs.go` | 后台任务（心跳检测、历史清理） |
| `config.yaml` | 默认配置文件 |

### Agent (`agent/`)
| 文件 | 职责 |
|---|---|
| `main.go` | 入口，注册/加载 token，启动调度器/上报器，信号处理 |
| `config.go` | 配置结构体 + YAML 加载 |
| `register.go` | Agent 注册 + token 本地持久化 |
| `scheduler.go` | 任务同步器（60s 轮询，diff，启停 goroutine） |
| `probe.go` | HTTP 探测引擎（构造请求、执行、匹配、记录结果） |
| `reporter.go` | 结果上报器（30s 批量上报，失败重试） |
| `agent.yaml` | 默认配置文件 |

### Web (`web/`)
| 文件 | 职责 |
|---|---|
| `src/main.ts` | 应用入口 |
| `src/App.vue` | 根组件 |
| `src/router/index.ts` | 路由 + 导航守卫 |
| `src/stores/auth.ts` | 登录状态管理 |
| `src/api/request.ts` | Axios 封装 + JWT 拦截器 |
| `src/api/auth.ts` / `monitor.ts` / `alert.ts` / `agent.ts` / `dashboard.ts` | API 模块 |
| `src/styles/variables.scss` | 苹果风格 CSS 变量 |
| `src/styles/element-override.scss` | Element Plus 主题覆盖 |
| `src/views/login/LoginView.vue` | 登录页 |
| `src/views/dashboard/DashboardView.vue` | 仪表盘 |
| `src/views/monitor/MonitorList.vue` | 监控列表 |
| `src/views/monitor/MonitorForm.vue` | 监控创建/编辑（共用组件） |
| `src/views/monitor/MonitorDetail.vue` | 监控详情 |
| `src/views/alert/AlertChannelList.vue` | 告警通道列表 |
| `src/views/alert/AlertChannelForm.vue` | 告警通道表单（对话框） |
| `src/views/alert/AlertHistory.vue` | 告警历史 |
| `src/views/agent/AgentList.vue` | Agent 列表 |
| `src/components/layout/AppLayout.vue` | 应用布局 |
| `src/components/layout/Sidebar.vue` | 侧边导航 |
| `src/components/charts/UptimeChart.vue` | 可用率趋势图 |
| `src/components/charts/HealthDonut.vue` | 健康度环形图 |
| `src/components/charts/LatencyChart.vue` | 延迟趋势图 |
| `src/components/common/StatusDot.vue` | 状态指示灯 |
| `src/components/common/StatCard.vue` | 统计卡片 |

---

## 阶段一：Server 基础

### 任务 1：Server 项目初始化与配置

**文件：**
- 创建：
- 创建：
- 创建：

- [ ] **步骤 1：初始化 Go 模块并安装依赖**



- [ ] **步骤 2：编写配置文件 **



- [ ] **步骤 3：实现配置加载 **

yaml:"server"yaml:"agent"yaml:"database"yaml:"tsdb"yaml:"alert"yaml:"port"yaml:"jwt_secret"yaml:"admin_username"yaml:"admin_password"yaml:"shared_secret"yaml:"heartbeat_timeout"yaml:"path"yaml:"path"yaml:"history_retention"

- [ ] **步骤 4：验证编译  — 预期无错误**
- [ ] **步骤 5：提交 [main 1ea2fac] feat(server): init project, config loading
 2 files changed, 164 insertions(+)
 create mode 100644 server/go.mod
 create mode 100644 server/go.sum**

---

### 任务 2：数据模型定义

**文件：**
- 创建：

- [ ] **步骤 1：定义所有数据模型和 DTO**

完整代码见设计规格 3.1 节。关键要点：
- 6 个 GORM 模型：User, Agent, Monitor, AlertChannel, MonitorAlertChannel, AlertHistory
- JSONSlice 自定义类型（ 序列化到 SQLite TEXT 字段）
- API DTO：LoginRequest, AgentRegisterRequest/Response, ProbeResult, ReportRequest, AgentTaskDTO, PageQuery, DashboardOverview
- Monitor 包含所有配置字段：url, method, headers, cookie, basic_auth, verify_ssl, frequency, proxy, agent_ids, match_type, match_content, status_threshold, latency_threshold, fail_count, enabled

- [ ] **步骤 2：验证编译 **
- [ ] **步骤 3：提交 On branch main
Your branch is ahead of 'origin/main' by 3 commits.
  (use "git push" to publish your local commits)

Untracked files:
  (use "git add <file>..." to include in what will be committed)
	docs/superpowers/plans/

nothing added to commit but untracked files present (use "git add" to track)**

---

### 任务 3：SQLite Repository 层

**文件：**
- 创建：
- 创建：
- 创建：
- 创建：
- 创建：

- [ ] **步骤 1：实现数据库初始化 **



- [ ] **步骤 2：实现  — FindUserByUsername**
- [ ] **步骤 3：实现 **

关键方法：
-  — 分页 + 名称/URL 模糊搜索
-  /  /  / 
-  — PATCH 启用/禁用
-  — 用 JSON LIKE 匹配 agent_ids 中包含该 agent 的已启用任务

- [ ] **步骤 4：实现 **

关键方法：
-  /  /  / 
-  — 更新 status=online + last_heartbeat
-  — 将 last_heartbeat 超过 timeout 的 online agent 标记为 offline

- [ ] **步骤 5：实现 **

关键方法：
- AlertChannel CRUD： /  /  /  / 
-  — JOIN monitor_alert_channels 查询
-  — 事务内先删后建
-  /  / 
-  — 分页 + 预加载 Monitor/Agent
-  / 

- [ ] **步骤 6：验证编译并提交**



---

### 任务 4：统一响应格式

**文件：**
- 创建：

- [ ] **步骤 1：实现响应辅助函数**

json:"code"json:"data"json:"message"json:"list"json:"total"json:"page"json:"page_size"

- [ ] **步骤 2：验证编译并提交**



---

## 阶段二：Server API

### 任务 5：认证系统

**文件：**
- 创建：`server/internal/service/auth.go`
- 创建：`server/internal/middleware/auth.go`
- 创建：`server/internal/middleware/agent_auth.go`
- 创建：`server/internal/handler/auth.go`

- [ ] **步骤 1：实现 AuthService**

```go
// server/internal/service/auth.go
package service

import (
	"errors"
	"time"
	"uptime-monitor/server/internal/config"
	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg *config.ServerConfig
}

func NewAuthService(cfg *config.ServerConfig) *AuthService {
	return &AuthService{cfg: cfg}
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := repository.FindUserByUsername(username)
	if err != nil {
		return "", errors.New("用户名或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
```

- [ ] **步骤 2：实现 JWT 中间件 `middleware/auth.go`**

```go
package middleware

import (
	"strings"
	"uptime-monitor/server/internal/handler"
	"uptime-monitor/server/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handler.Error(c, 40001, "未提供认证令牌")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := authSvc.ValidateToken(tokenStr)
		if err != nil {
			handler.Error(c, 40001, "认证令牌无效")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
```

- [ ] **步骤 3：实现 Agent Token 中间件 `middleware/agent_auth.go`**

```go
package middleware

import (
	"strings"
	"uptime-monitor/server/internal/handler"
	"uptime-monitor/server/internal/repository"

	"github.com/gin-gonic/gin"
)

func AgentAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handler.Error(c, 40001, "未提供 Agent 令牌")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		agent, err := repository.FindAgentByToken(tokenStr)
		if err != nil {
			handler.Error(c, 40001, "Agent 令牌无效")
			c.Abort()
			return
		}
		c.Set("agent_id", agent.ID)
		c.Set("agent_name", agent.Name)
		c.Next()
	}
}
```

- [ ] **步骤 4：实现 Auth Handler `handler/auth.go`**

```go
package handler

import (
	"github.com/gin-gonic/gin"
	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 40101, "参数错误")
		return
	}
	token, err := h.authSvc.Login(req.Username, req.Password)
	if err != nil {
		Error(c, 40001, err.Error())
		return
	}
	Success(c, gin.H{"token": token})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	Success(c, nil) // JWT 无状态，客户端清除 token 即可
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	Success(c, gin.H{"id": userID})
}
```

- [ ] **步骤 5：验证编译并提交**

```bash
cd server && go build ./...
git add server/ && git commit -m "feat(server): implement auth system with JWT"
```

---

### 任务 6：监控任务 CRUD

**文件：**
- 创建：`server/internal/service/monitor.go`
- 创建：`server/internal/handler/monitor.go`

- [ ] **步骤 1：实现 MonitorService**

```go
// server/internal/service/monitor.go
package service

import (
	"errors"
	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
)

type MonitorService struct{}

func NewMonitorService() *MonitorService {
	return &MonitorService{}
}

func (s *MonitorService) List(page, pageSize int, keyword string) ([]model.Monitor, int64, error) {
	return repository.ListMonitors(page, pageSize, keyword)
}

func (s *MonitorService) Get(id uint) (*model.Monitor, error) {
	return repository.GetMonitor(id)
}

func (s *MonitorService) Create(m *model.Monitor) error {
	if m.URL == "" {
		return errors.New("监控地址不能为空")
	}
	if m.Frequency < 30 {
		m.Frequency = 30 // 最小 30 秒
	}
	if m.FailCount < 1 || m.FailCount > 5 {
		m.FailCount = 3
	}
	return repository.CreateMonitor(m)
}

func (s *MonitorService) Update(m *model.Monitor) error {
	return repository.UpdateMonitor(m)
}

func (s *MonitorService) Delete(id uint) error {
	return repository.DeleteMonitor(id)
}

func (s *MonitorService) Toggle(id uint, enabled bool) error {
	return repository.ToggleMonitor(id, enabled)
}
```

- [ ] **步骤 2：实现 Monitor Handler**

```go
// server/internal/handler/monitor.go
package handler

import (
	"strconv"
	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/service"

	"github.com/gin-gonic/gin"
)

type MonitorHandler struct {
	svc *service.MonitorService
}

func NewMonitorHandler(svc *service.MonitorService) *MonitorHandler {
	return &MonitorHandler{svc: svc}
}

func (h *MonitorHandler) List(c *gin.Context) {
	var q model.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		Error(c, 40101, "参数错误")
		return
	}
	list, total, err := h.svc.List(q.Page, q.PageSize, q.Keyword)
	if err != nil {
		Error(c, 50001, "查询失败")
		return
	}
	Paginated(c, list, total, q.Page, q.PageSize)
}

func (h *MonitorHandler) Get(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}
	m, err := h.svc.Get(id)
	if err != nil {
		Error(c, 50001, "查询失败")
		return
	}
	Success(c, m)
}

func (h *MonitorHandler) Create(c *gin.Context) {
	var m model.Monitor
	if err := c.ShouldBindJSON(&m); err != nil {
		Error(c, 40101, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.Create(&m); err != nil {
		Error(c, 40101, err.Error())
		return
	}
	Success(c, m)
}

func (h *MonitorHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}
	var m model.Monitor
	if err := c.ShouldBindJSON(&m); err != nil {
		Error(c, 40101, "参数错误")
		return
	}
	m.ID = id
	if err := h.svc.Update(&m); err != nil {
		Error(c, 50001, "更新失败")
		return
	}
	Success(c, m)
}

func (h *MonitorHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		Error(c, 50001, "删除失败")
		return
	}
	Success(c, nil)
}

func (h *MonitorHandler) ToggleEnabled(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 40101, "参数错误")
		return
	}
	if err := h.svc.Toggle(id, req.Enabled); err != nil {
		Error(c, 50001, "操作失败")
		return
	}
	Success(c, nil)
}

func parseUintParam(c *gin.Context, key string) (uint, error) {
	idStr := c.Param(key)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		Error(c, 40101, "无效的参数 ID")
		return 0, err
	}
	return uint(id), nil
}
```

- [ ] **步骤 3：验证编译并提交**

```bash
cd server && go build ./...
git add server/ && git commit -m "feat(server): implement monitor CRUD"
```

---

### 任务 7：Agent API

**文件：**
- 创建：`server/internal/service/agent.go`
- 创建：`server/internal/handler/agent.go`

- [ ] **步骤 1：实现 AgentService**

```go
// server/internal/service/agent.go
package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"uptime-monitor/server/internal/config"
	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
)

type AgentService struct {
	cfg *config.AgentConfig
}

func NewAgentService(cfg *config.AgentConfig) *AgentService {
	return &AgentService{cfg: cfg}
}

func (s *AgentService) Register(req model.AgentRegisterRequest) (*model.AgentRegisterResponse, error) {
	if req.SharedSecret != s.cfg.SharedSecret {
		return nil, errors.New("预共享密钥错误")
	}
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	agent := &model.Agent{
		Name:  req.Name,
		Token: token,
		Location: req.Location,
		Status: "online",
	}
	if err := repository.CreateAgent(agent); err != nil {
		return nil, err
	}
	return &model.AgentRegisterResponse{
		AgentID: agent.ID,
		Token:   token,
	}, nil
}

func (s *AgentService) GetTasks(agentID uint) ([]model.AgentTaskDTO, error) {
	monitors, err := repository.GetMonitorsByAgent(agentID)
	if err != nil {
		return nil, err
	}
	tasks := make([]model.AgentTaskDTO, len(monitors))
	for i, m := range monitors {
		tasks[i] = model.AgentTaskDTO{
			ID: m.ID, Name: m.Name, URL: m.URL, Method: m.Method,
			Headers: m.Headers, Cookie: m.Cookie,
			BasicAuthUser: m.BasicAuthUser, BasicAuthPass: m.BasicAuthPass,
			VerifySSL: m.VerifySSL, Frequency: m.Frequency, Proxy: m.Proxy,
			MatchType: m.MatchType, MatchContent: m.MatchContent,
			StatusThreshold: m.StatusThreshold, LatencyThreshold: m.LatencyThreshold,
			FailCount: m.FailCount, UpdatedAt: m.UpdatedAt,
		}
	}
	return tasks, nil
}

func (s *AgentService) Heartbeat(agentID uint) error {
	return repository.UpdateHeartbeat(agentID)
}
```

- [ ] **步骤 2：实现 Agent Handler**

```go
// server/internal/handler/agent.go
package handler

import (
	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/service"

	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	agentSvc   *service.AgentService
	alertSvc   *service.AlertService
	metricsSvc *service.MetricsService
}

func NewAgentHandler(agentSvc *service.AgentService, alertSvc *service.AlertService, metricsSvc *service.MetricsService) *AgentHandler {
	return &AgentHandler{agentSvc: agentSvc, alertSvc: alertSvc, metricsSvc: metricsSvc}
}

func (h *AgentHandler) Register(c *gin.Context) {
	var req model.AgentRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 40101, "参数错误")
		return
	}
	resp, err := h.agentSvc.Register(req)
	if err != nil {
		Error(c, 40001, err.Error())
		return
	}
	Success(c, resp)
}

func (h *AgentHandler) GetTasks(c *gin.Context) {
	agentID, _ := c.Get("agent_id")
	tasks, err := h.agentSvc.GetTasks(agentID.(uint))
	if err != nil {
		Error(c, 50001, "获取任务失败")
		return
	}
	Success(c, gin.H{"tasks": tasks})
}

func (h *AgentHandler) Report(c *gin.Context) {
	var req model.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 40101, "参数错误")
		return
	}
	agentID, _ := c.Get("agent_id")
	// 写入 TSDB + 评估告警
	for _, result := range req.Results {
		h.metricsSvc.Record(result)
		h.alertSvc.Evaluate(agentID.(uint), result)
	}
	Success(c, nil)
}

func (h *AgentHandler) Heartbeat(c *gin.Context) {
	agentID, _ := c.Get("agent_id")
	if err := h.agentSvc.Heartbeat(agentID.(uint)); err != nil {
		Error(c, 50001, "心跳上报失败")
		return
	}
	Success(c, nil)
}
```

- [ ] **步骤 3：验证编译并提交**

```bash
cd server && go build ./...
git add server/ && git commit -m "feat(server): implement agent API"
```

---

### 任务 8：告警通道 CRUD

**文件：**
- 创建：`server/internal/handler/alert.go`

- [ ] **步骤 1：实现 Alert Handler**

```go
// server/internal/handler/alert.go
package handler

import (
	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
	"uptime-monitor/server/internal/service"

	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	notifySvc *service.NotifyService
}

func NewAlertHandler(notifySvc *service.NotifyService) *AlertHandler {
	return &AlertHandler{notifySvc: notifySvc}
}

func (h *AlertHandler) ListChannels(c *gin.Context) {
	channels, err := repository.ListAlertChannels()
	if err != nil {
		Error(c, 50001, "查询失败")
		return
	}
	Success(c, channels)
}

func (h *AlertHandler) CreateChannel(c *gin.Context) {
	var ch model.AlertChannel
	if err := c.ShouldBindJSON(&ch); err != nil {
		Error(c, 40101, "参数错误")
		return
	}
	if err := repository.CreateAlertChannel(&ch); err != nil {
		Error(c, 50001, "创建失败")
		return
	}
	Success(c, ch)
}

func (h *AlertHandler) UpdateChannel(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}
	var ch model.AlertChannel
	if err := c.ShouldBindJSON(&ch); err != nil {
		Error(c, 40101, "参数错误")
		return
	}
	ch.ID = id
	if err := repository.UpdateAlertChannel(&ch); err != nil {
		Error(c, 50001, "更新失败")
		return
	}
	Success(c, ch)
}

func (h *AlertHandler) DeleteChannel(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}
	if err := repository.DeleteAlertChannel(id); err != nil {
		Error(c, 50001, "删除失败")
		return
	}
	Success(c, nil)
}

func (h *AlertHandler) TestChannel(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}
	ch, err := repository.GetAlertChannel(id)
	if err != nil {
		Error(c, 50001, "通道不存在")
		return
	}
	if err := h.notifySvc.SendTest(ch); err != nil {
		Error(c, 50001, "发送失败: "+err.Error())
		return
	}
	Success(c, nil)
}

func (h *AlertHandler) ListHistory(c *gin.Context) {
	var q struct {
		model.PageQuery
		MonitorID uint   `form:"monitor_id"`
		Status    string `form:"status"`
	}
	if err := c.ShouldBindQuery(&q); err != nil {
		Error(c, 40101, "参数错误")
		return
	}
	list, total, err := repository.ListAlertHistory(q.Page, q.PageSize, q.MonitorID, q.Status)
	if err != nil {
		Error(c, 50001, "查询失败")
		return
	}
	Paginated(c, list, total, q.Page, q.PageSize)
}
```

- [ ] **步骤 2：验证编译并提交**

---

### 任务 9：告警评估服务

**文件：**
- 创建：`server/internal/service/alert.go`

- [ ] **步骤 1：实现 AlertService**

```go
// server/internal/service/alert.go
package service

import (
	"fmt"
	"log"
	"sync"
	"time"
	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
)

type AlertService struct {
	failCounts sync.Map // key: "monitorID_agentID" -> *int
	notifySvc  *NotifyService
}

func NewAlertService(notifySvc *NotifyService) *AlertService {
	return &AlertService{notifySvc: notifySvc}
}

func failKey(monitorID, agentID uint) string {
	return fmt.Sprintf("%d_%d", monitorID, agentID)
}

func (s *AlertService) Evaluate(agentID uint, result model.ProbeResult) {
	monitor, err := repository.GetMonitor(result.MonitorID)
	if err != nil || !monitor.Enabled {
		return
	}

	// 检查是否触发告警条件
	failed := false
	alertTypes := []string{}

	if result.StatusCode >= monitor.StatusThreshold {
		failed = true
		alertTypes = append(alertTypes, "status_code")
	}
	if result.DurationMs >= int64(monitor.LatencyThreshold) {
		failed = true
		alertTypes = append(alertTypes, "latency")
	}
	if monitor.MatchType == "contains" && !result.ContentMatched {
		failed = true
		alertTypes = append(alertTypes, "content_match")
	}
	if monitor.MatchType == "not_contains" && result.ContentMatched {
		failed = true
		alertTypes = append(alertTypes, "content_match")
	}

	key := failKey(result.MonitorID, agentID)

	if failed {
		// 增加连续失败计数
		val, _ := s.failCounts.LoadOrStore(key, new(int))
		count := val.(*int)
		*count++

		if *count >= monitor.FailCount {
			// 检查是否已在告警中
			active, _ := repository.FindActiveAlert(result.MonitorID, agentID)
			if active != nil {
				return // 已在告警中，不重复
			}
			// 创建告警记录
			alertType := alertTypes[0]
			msg := fmt.Sprintf("状态码: %d, 延迟: %dms", result.StatusCode, result.DurationMs)
			history := &model.AlertHistory{
				MonitorID:   result.MonitorID,
				AgentID:     agentID,
				AlertType:   alertType,
				Message:     msg,
				Status:      "firing",
				TriggeredAt: time.Now(),
			}
			if err := repository.CreateAlertHistory(history); err != nil {
				log.Printf("创建告警记录失败: %v", err)
				return
			}
			// 发送通知
			go s.notifySvc.NotifyFiring(monitor, agentID, alertType, msg, *count)
		}
	} else {
		// 恢复正常
		val, loaded := s.failCounts.Load(key)
		if loaded {
			count := val.(*int)
			if *count > 0 {
				// 查找并解决告警
				active, _ := repository.FindActiveAlert(result.MonitorID, agentID)
				if active != nil {
					repository.ResolveAlert(active.ID)
					go s.notifySvc.NotifyResolved(monitor, active)
				}
			}
			*count = 0
		}
	}
}
```

- [ ] **步骤 2：验证编译并提交**

---

### 任务 10：通知服务

**文件：**
- 创建：`server/internal/service/notify.go`

- [ ] **步骤 1：实现 NotifyService**

```go
// server/internal/service/notify.go
package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
)

type NotifyService struct{}

func NewNotifyService() *NotifyService {
	return &NotifyService{}
}

// NotifyFiring 发送告警通知
func (s *NotifyService) NotifyFiring(monitor *model.Monitor, agentID uint, alertType, detail string, count int) {
	channels, err := repository.GetAlertChannelsByMonitor(monitor.ID)
	if err != nil || len(channels) == 0 {
		return
	}
	title := fmt.Sprintf("🔴 告警：%s 异常", monitor.Name)
	text := fmt.Sprintf(
		"**告警：%s 异常**\n\n- 监控地址：%s\n- 告警类型：%s\n- 详情：%s\n- 连续失败：%d 次\n- 时间：%s",
		monitor.Name, monitor.URL, alertType, detail, count, time.Now().Format("2006-01-02 15:04:05"),
	)
	for _, ch := range channels {
		if err := s.send(&ch, title, text); err != nil {
			log.Printf("发送告警通知失败 [%s]: %v", ch.Name, err)
		}
	}
}

// NotifyResolved 发送恢复通知
func (s *NotifyService) NotifyResolved(monitor *model.Monitor, alert *model.AlertHistory) {
	channels, err := repository.GetAlertChannelsByMonitor(monitor.ID)
	if err != nil || len(channels) == 0 {
		return
	}
	duration := "未知"
	if alert.ResolvedAt != nil {
		duration = alert.ResolvedAt.Sub(alert.TriggeredAt).Round(time.Second).String()
	}
	title := fmt.Sprintf("🟢 恢复：%s 已恢复正常", monitor.Name)
	text := fmt.Sprintf(
		"**恢复：%s 已恢复正常**\n\n- 监控地址：%s\n- 恢复时间：%s\n- 持续时长：%s",
		monitor.Name, monitor.URL, time.Now().Format("2006-01-02 15:04:05"), duration,
	)
	for _, ch := range channels {
		if err := s.send(&ch, title, text); err != nil {
			log.Printf("发送恢复通知失败 [%s]: %v", ch.Name, err)
		}
	}
}

// SendTest 发送测试消息
func (s *NotifyService) SendTest(ch *model.AlertChannel) error {
	return s.send(ch, "🔔 测试通知", "这是一条测试消息，如果您收到说明通道配置正确。")
}

func (s *NotifyService) send(ch *model.AlertChannel, title, text string) error {
	switch ch.Type {
	case "dingtalk":
		return s.sendDingtalk(ch, title, text)
	case "wechat_work":
		return s.sendWechatWork(ch, title, text)
	case "feishu":
		return s.sendFeishu(ch, title, text)
	case "webhook":
		return s.sendWebhook(ch, title, text)
	default:
		return fmt.Errorf("未知的通知类型: %s", ch.Type)
	}
}

func (s *NotifyService) sendDingtalk(ch *model.AlertChannel, title, text string) error {
	webhookURL := ch.WebhookURL
	if ch.Secret != "" {
		timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
		stringToSign := timestamp + "\n" + ch.Secret
		mac := hmac.New(sha256.New, []byte(ch.Secret))
		mac.Write([]byte(stringToSign))
		sign := url.QueryEscape(base64.StdEncoding.EncodeToString(mac.Sum(nil)))
		webhookURL += "&timestamp=" + timestamp + "&sign=" + sign
	}
	body := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
	}
	return postJSON(webhookURL, body)
}

func (s *NotifyService) sendWechatWork(ch *model.AlertChannel, title, text string) error {
	body := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": title + "\n" + text,
		},
	}
	return postJSON(ch.WebhookURL, body)
}

func (s *NotifyService) sendFeishu(ch *model.AlertChannel, title, text string) error {
	webhookURL := ch.WebhookURL
	var body map[string]interface{}
	if ch.Secret != "" {
		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		stringToSign := timestamp + "\n" + ch.Secret
		mac := hmac.New(sha256.New, []byte(stringToSign))
		sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		body = map[string]interface{}{
			"timestamp": timestamp,
			"sign":      sign,
			"msg_type":  "interactive",
			"card": map[string]interface{}{
				"header": map[string]interface{}{
					"title": map[string]string{"tag": "plain_text", "content": title},
				},
				"elements": []map[string]interface{}{
					{"tag": "div", "text": map[string]string{"tag": "lark_md", "content": text}},
				},
			},
		}
	} else {
		body = map[string]interface{}{
			"msg_type": "interactive",
			"card": map[string]interface{}{
				"header": map[string]interface{}{
					"title": map[string]string{"tag": "plain_text", "content": title},
				},
				"elements": []map[string]interface{}{
					{"tag": "div", "text": map[string]string{"tag": "lark_md", "content": text}},
				},
			},
		}
	}
	_ = webhookURL
	return postJSON(ch.WebhookURL, body)
}

func (s *NotifyService) sendWebhook(ch *model.AlertChannel, title, text string) error {
	// 自定义 webhook 直接发送 JSON
	body := map[string]interface{}{
		"title":   title,
		"content": text,
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	}
	return postJSON(ch.WebhookURL, body)
}

func postJSON(url string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}
```

注意：飞书签名部分 `hmac.New(sha256.New, []byte(stringToSign))` 应改为 `mac := hmac.New(sha256.New, []byte(""))` + `mac.Write([]byte(stringToSign))`，与标准飞书签名算法一致。实现时需参考飞书文档调整。

- [ ] **步骤 2：验证编译并提交**

---

### 任务 11：Metrics 与 Dashboard

**文件：**
- 创建：`server/internal/service/metrics.go`
- 创建：`server/internal/handler/metrics.go`
- 创建：`server/internal/handler/dashboard.go`

- [ ] **步骤 1：实现 MetricsService（使用 Prometheus client_golang 内存指标）**

使用 `prometheus.NewGaugeVec` 和 `prometheus.NewCounterVec` 注册指标，写入内存 Prometheus registry。查询时使用 `prometheus` 包的 API 读取。

关键设计：
- `Record(result ProbeResult)` — 更新 gauge 指标
- `QueryMetrics(monitorID, metricName, duration)` — 查询指定时间范围的时序数据
- 指标名：`probe_success`, `probe_http_status_code`, `probe_duration_seconds`, `probe_ssl_expiry_seconds`, `probe_content_matched`

```go
package service

import (
	"uptime-monitor/server/internal/model"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsService struct {
	probeSuccess    *prometheus.GaugeVec
	probeStatusCode *prometheus.GaugeVec
	probeDuration   *prometheus.GaugeVec
	probeSSLExpiry  *prometheus.GaugeVec
	probeMatched    *prometheus.GaugeVec
	registry        *prometheus.Registry
}

func NewMetricsService() *MetricsService {
	labels := []string{"monitor_id", "agent_id", "url"}
	s := &MetricsService{
		probeSuccess: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_success", Help: "Probe success (1=ok, 0=fail)",
		}, append(labels, "method")),
		probeStatusCode: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_http_status_code", Help: "HTTP status code",
		}, labels),
		probeDuration: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_duration_seconds", Help: "Response duration in seconds",
		}, labels),
		probeSSLExpiry: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_ssl_expiry_seconds", Help: "SSL cert expiry seconds",
		}, labels),
		probeMatched: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_content_matched", Help: "Content matched (1=yes, 0=no)",
		}, labels),
		registry: prometheus.NewRegistry(),
	}
	s.registry.MustRegister(s.probeSuccess, s.probeStatusCode, s.probeDuration, s.probeSSLExpiry, s.probeMatched)
	return s
}

func (s *MetricsService) Record(r model.ProbeResult) {
	monitorID := fmt.Sprintf("%d", r.MonitorID)
	agentID := fmt.Sprintf("%d", monitorID) // 从上报数据获取
	monitor, _ := repository.GetMonitor(r.MonitorID)
	urlLabel := ""
	method := "GET"
	if monitor != nil {
		urlLabel = monitor.URL
		method = monitor.Method
	}

	successVal := 0.0
	if r.Success {
		successVal = 1.0
	}
	matchedVal := 0.0
	if r.ContentMatched {
		matchedVal = 1.0
	}

	s.probeSuccess.WithLabelValues(monitorID, agentID, urlLabel, method).Set(successVal)
	s.probeStatusCode.WithLabelValues(monitorID, agentID, urlLabel).Set(float64(r.StatusCode))
	s.probeDuration.WithLabelValues(monitorID, agentID, urlLabel).Set(float64(r.DurationMs) / 1000.0)
	s.probeSSLExpiry.WithLabelValues(monitorID, agentID, urlLabel).Set(float64(r.SSLExpiry))
	s.probeMatched.WithLabelValues(monitorID, agentID, urlLabel).Set(matchedVal)
}
```

注意：上面代码为概念示意，实现时需修正 `agentID` 取值（应从方法参数传入）并添加 `"fmt"` 导入。

- [ ] **步骤 2：实现 Metrics Handler**

提供 `GET /api/monitors/:id/metrics` 接口，返回可用率、延迟趋势、状态码分布等聚合数据。内部读取 Prometheus registry 的当前值。

- [ ] **步骤 3：实现 Dashboard Handler**

提供 `GET /api/dashboard/overview` 接口：
- `total_monitors` — 查询 Monitor 总数
- `online_rate` — 从 probe_success gauge 计算
- `active_alerts` — 查询 firing 状态告警数
- `avg_response_ms` — 从 probe_duration gauge 平均值计算
- `recent_alerts` — 查询最近 10 条告警历史

- [ ] **步骤 4：验证编译并提交**

```bash
cd server && go build ./...
git add server/ && git commit -m "feat(server): implement metrics and dashboard"
```

---

## 阶段三：Server 组装

### 任务 12：Main 入口与后台任务

**文件：**
- 创建：`server/internal/background/jobs.go`
- 创建：`server/cmd/main.go`

- [ ] **步骤 1：实现后台任务 `background/jobs.go`**

```go
// server/internal/background/jobs.go
package background

import (
	"context"
	"log"
	"time"
	"uptime-monitor/server/internal/config"
	"uptime-monitor/server/internal/repository"
)

// StartJobs 启动所有后台任务，返回 cancel 函数
func StartJobs(ctx context.Context, cfg *config.Config) {
	// Agent 心跳检测：每 30 秒扫描
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				timeout := time.Duration(cfg.Agent.HeartbeatTimeout) * time.Second
				if err := repository.MarkOfflineAgents(timeout); err != nil {
					log.Printf("心跳检测失败: %v", err)
				}
			}
		}
	}()

	// 告警历史清理：每日凌晨
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := repository.CleanupAlertHistory(cfg.Alert.HistoryRetention); err != nil {
					log.Printf("告警历史清理失败: %v", err)
				}
			}
		}
	}()

	log.Println("后台任务已启动: 心跳检测(30s), 告警历史清理(24h)")
}
```

- [ ] **步骤 2：实现 main.go**

```go
// server/cmd/main.go
package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"uptime-monitor/server/internal/background"
	"uptime-monitor/server/internal/config"
	"uptime-monitor/server/internal/handler"
	"uptime-monitor/server/internal/middleware"
	"uptime-monitor/server/internal/repository"
	"uptime-monitor/server/internal/service"

	"github.com/gin-gonic/gin"
)

//go:embed all:web/dist
var webFS embed.FS

func main() {
	cfgPath := flag.String("c", "config.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 初始化数据库
	if err := repository.Init(&cfg.Database, cfg.Server.AdminUsername, cfg.Server.AdminPassword); err != nil {
		log.Fatal("初始化数据库失败:", err)
	}

	// 初始化服务层
	authSvc := service.NewAuthService(&cfg.Server)
	monitorSvc := service.NewMonitorService()
	agentSvc := service.NewAgentService(&cfg.Agent)
	notifySvc := service.NewNotifyService()
	alertSvc := service.NewAlertService(notifySvc)
	metricsSvc := service.NewMetricsService()

	// 初始化 Handler
	authHandler := handler.NewAuthHandler(authSvc)
	monitorHandler := handler.NewMonitorHandler(monitorSvc)
	agentHandler := handler.NewAgentHandler(agentSvc, alertSvc, metricsSvc)
	alertHandler := handler.NewAlertHandler(notifySvc)

	// 设置 Gin 路由
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 公开路由
	r.POST("/api/auth/login", authHandler.Login)

	// Agent 路由（注册不需要 token）
	r.POST("/api/v1/agent/register", agentHandler.Register)

	// Agent 路由（需要 Agent Token）
	agentGroup := r.Group("/api/v1/agent", middleware.AgentAuthMiddleware())
	agentGroup.GET("/tasks", agentHandler.GetTasks)
	agentGroup.POST("/report", agentHandler.Report)
	agentGroup.POST("/heartbeat", agentHandler.Heartbeat)

	// Web API 路由（需要 JWT）
	api := r.Group("/api", middleware.AuthMiddleware(authSvc))
	api.POST("/auth/logout", authHandler.Logout)
	api.GET("/auth/me", authHandler.Me)

	// 监控
	monitors := api.Group("/monitors")
	monitors.GET("", monitorHandler.List)
	monitors.POST("", monitorHandler.Create)
	monitors.GET("/:id", monitorHandler.Get)
	monitors.PUT("/:id", monitorHandler.Update)
	monitors.DELETE("/:id", monitorHandler.Delete)
	monitors.PATCH("/:id/enabled", monitorHandler.ToggleEnabled)

	// 告警
	alerts := api.Group("/alerts")
	alerts.GET("/channels", alertHandler.ListChannels)
	alerts.POST("/channels", alertHandler.CreateChannel)
	alerts.PUT("/channels/:id", alertHandler.UpdateChannel)
	alerts.DELETE("/channels/:id", alertHandler.DeleteChannel)
	alerts.POST("/channels/:id/test", alertHandler.TestChannel)
	alerts.GET("/history", alertHandler.ListHistory)

	// 嵌入 Web 静态文件
	distFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Fatal("加载 Web 静态文件失败:", err)
	}
	r.NoRoute(gin.WrapH(http.FileServer(http.FS(distFS))))

	// 启动后台任务
	ctx, cancel := context.WithCancel(context.Background())
	background.StartJobs(ctx, cfg)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	go func() {
		log.Printf("Server 启动在 http://localhost%s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatal("Server 启动失败:", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭 Server...")
	cancel()
}
```

- [ ] **步骤 3：创建 `server/web/dist/` 占位目录**

```bash
mkdir -p server/web/dist
echo '<html><body>Web UI placeholder</body></html>' > server/web/dist/index.html
```

- [ ] **步骤 4：验证编译并提交**

```bash
cd server && go build ./cmd/main.go
git add server/ && git commit -m "feat(server): implement main entry, routing, background jobs"
```

---

## 阶段四：Agent 端

### 任务 13：Agent 项目初始化

**文件：**
- 创建：`agent/go.mod`
- 创建：`agent/config.go`
- 创建：`agent/agent.yaml`

- [ ] **步骤 1：初始化 Go 模块**

```bash
mkdir -p agent && cd agent
go mod init uptime-monitor/agent
go get gopkg.in/yaml.v3
```

- [ ] **步骤 2：实现配置加载 `config.go`**

```go
// agent/config.go
package main

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Agent  AgentConfig  `yaml:"agent"`
	Probe  ProbeConfig  `yaml:"probe"`
	Report ReportConfig `yaml:"report"`
}

type ServerConfig struct {
	URL          string `yaml:"url"`
	SharedSecret string `yaml:"shared_secret"`
}

type AgentConfig struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
}

type ProbeConfig struct {
	Timeout       int `yaml:"timeout"`
	MaxConcurrent int `yaml:"max_concurrent"`
}

type ReportConfig struct {
	Interval    int `yaml:"interval"`
	BatchSize   int `yaml:"batch_size"`
	RetryMaxAge int `yaml:"retry_max_age"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// 默认值
	if cfg.Probe.Timeout == 0 {
		cfg.Probe.Timeout = 30
	}
	if cfg.Probe.MaxConcurrent == 0 {
		cfg.Probe.MaxConcurrent = 50
	}
	if cfg.Report.Interval == 0 {
		cfg.Report.Interval = 30
	}
	if cfg.Report.BatchSize == 0 {
		cfg.Report.BatchSize = 100
	}
	if cfg.Report.RetryMaxAge == 0 {
		cfg.Report.RetryMaxAge = 600
	}
	return &cfg, nil
}
```

- [ ] **步骤 3：编写 `agent.yaml`**

```yaml
server:
  url: "http://localhost:8080"
  shared_secret: "agent-shared-secret"
agent:
  name: "本地探测点"
  location: "本地"
probe:
  timeout: 30
  max_concurrent: 50
report:
  interval: 30
  batch_size: 100
  retry_max_age: 600
```

- [ ] **步骤 4：验证编译并提交**

---

### 任务 14：Agent 注册

**文件：**
- 创建：`agent/register.go`

- [ ] **步骤 1：实现注册与 token 持久化**

```go
// agent/register.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type registerReq struct {
	Name         string `json:"name"`
	Location     string `json:"location"`
	SharedSecret string `json:"shared_secret"`
}

type registerResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AgentID uint   `json:"agent_id"`
		Token   string `json:"token"`
	} `json:"data"`
}

const tokenFile = ".agent_token"

func Register(cfg *Config) (string, error) {
	// 先尝试加载已保存的 token
	token, err := loadToken()
	if err == nil && token != "" {
		fmt.Println("已加载保存的 Agent token")
		return token, nil
	}

	// 向 Server 注册
	body := registerReq{
		Name:         cfg.Agent.Name,
		Location:     cfg.Agent.Location,
		SharedSecret: cfg.Server.SharedSecret,
	}
	data, _ := json.Marshal(body)
	resp, err := http.Post(cfg.Server.URL+"/api/v1/agent/register", "application/json", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("注册请求失败: %w", err)
	}
	defer resp.Body.Close()

	respData, _ := io.ReadAll(resp.Body)
	var result registerResp
	if err := json.Unmarshal(respData, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}
	if result.Code != 0 {
		return "", fmt.Errorf("注册失败: %s", result.Message)
	}

	// 保存 token 到本地文件
	if err := saveToken(result.Data.Token); err != nil {
		fmt.Printf("警告: 保存 token 失败: %v\n", err)
	}
	fmt.Printf("注册成功，Agent ID: %d\n", result.Data.AgentID)
	return result.Data.Token, nil
}

func loadToken() (string, error) {
	data, err := os.ReadFile(tokenFilePath())
	return string(data), err
}

func saveToken(token string) error {
	return os.WriteFile(tokenFilePath(), []byte(token), 0600)
}

func tokenFilePath() string {
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), tokenFile)
}
```

- [ ] **步骤 2：验证编译并提交**

---

### 任务 15：任务同步器

**文件：**
- 创建：`agent/scheduler.go`

- [ ] **步骤 1：实现任务同步器**

```go
// agent/scheduler.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Task struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	Method           string `json:"method"`
	Headers          string `json:"headers"`
	Cookie           string `json:"cookie"`
	BasicAuthUser    string `json:"basic_auth_user"`
	BasicAuthPass    string `json:"basic_auth_pass"`
	VerifySSL        bool   `json:"verify_ssl"`
	Frequency        int    `json:"frequency"`
	Proxy            string `json:"proxy"`
	MatchType        string `json:"match_type"`
	MatchContent     string `json:"match_content"`
	StatusThreshold  int    `json:"status_threshold"`
	LatencyThreshold int    `json:"latency_threshold"`
	FailCount        int    `json:"fail_count"`
	UpdatedAt        string `json:"updated_at"`
}

type Scheduler struct {
	cfg       *Config
	token     string
	tasks     map[uint]*taskRunner
	mu        sync.RWMutex
	results   *ResultQueue
	probeSem  chan struct{} // 并发控制
}

type taskRunner struct {
	task   Task
	cancel context.CancelFunc
}

func NewScheduler(cfg *Config, token string, results *ResultQueue) *Scheduler {
	return &Scheduler{
		cfg:      cfg,
		token:    token,
		tasks:    make(map[uint]*taskRunner),
		results:  results,
		probeSem: make(chan struct{}, cfg.Probe.MaxConcurrent),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	// 立即同步一次
	s.syncTasks(ctx)

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			s.stopAll()
			return
		case <-ticker.C:
			s.syncTasks(ctx)
		}
	}
}

func (s *Scheduler) syncTasks(ctx context.Context) {
	newTasks, err := s.fetchTasks()
	if err != nil {
		log.Printf("获取任务失败: %v", err)
		return
	}

	newMap := make(map[uint]Task)
	for _, t := range newTasks {
		newMap[t.ID] = t
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 删除不再需要的任务
	for id, runner := range s.tasks {
		if _, exists := newMap[id]; !exists {
			log.Printf("停止任务: %s (ID: %d)", runner.task.Name, id)
			runner.cancel()
			delete(s.tasks, id)
		}
	}

	// 新增或更新任务
	for id, task := range newMap {
		runner, exists := s.tasks[id]
		if exists && runner.task.UpdatedAt == task.UpdatedAt {
			continue // 未变更
		}
		if exists {
			runner.cancel() // 配置变更，先停旧的
		}
		taskCtx, cancel := context.WithCancel(ctx)
		s.tasks[id] = &taskRunner{task: task, cancel: cancel}
		log.Printf("启动任务: %s (ID: %d, 频率: %ds)", task.Name, id, task.Frequency)
		go s.runProbe(taskCtx, task)
	}

	log.Printf("任务同步完成: 当前运行 %d 个任务", len(s.tasks))
}

func (s *Scheduler) fetchTasks() ([]Task, error) {
	req, _ := http.NewRequest("GET", s.cfg.Server.URL+"/api/v1/agent/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+s.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Code int    `json:"code"`
		Data struct {
			Tasks []Task `json:"tasks"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("API 错误: code=%d", result.Code)
	}
	return result.Data.Tasks, nil
}

func (s *Scheduler) stopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, runner := range s.tasks {
		runner.cancel()
		delete(s.tasks, id)
	}
	log.Println("所有任务已停止")
}

func (s *Scheduler) runProbe(ctx context.Context, task Task) {
	ticker := time.NewTicker(time.Duration(task.Frequency) * time.Second)
	defer ticker.Stop()

	// 立即执行一次
	s.doProbe(task)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.doProbe(task)
		}
	}
}

func (s *Scheduler) doProbe(task Task) {
	s.probeSem <- struct{}{}
	defer func() { <-s.probeSem }()

	result := ExecuteProbe(task, s.cfg.Probe.Timeout)
	s.results.Add(result)
}
```

- [ ] **步骤 2：验证编译并提交**

---

### 任务 16：HTTP 探测引擎

**文件：**
- 创建：`agent/probe.go`

- [ ] **步骤 1：实现探测引擎**

```go
// agent/probe.go
package main

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"encoding/json"
	"fmt"
)

type ProbeResult struct {
	MonitorID      uint   `json:"monitor_id"`
	StatusCode     int    `json:"status_code"`
	DurationMs     int64  `json:"duration_ms"`
	ContentMatched bool   `json:"content_matched"`
	SSLExpiry      int64  `json:"ssl_expiry"`
	Success        bool   `json:"success"`
	ErrorMsg       string `json:"error_msg"`
	Timestamp      int64  `json:"timestamp"`
}

func ExecuteProbe(task Task, timeoutSec int) ProbeResult {
	result := ProbeResult{
		MonitorID: task.ID,
		Timestamp: time.Now().Unix(),
	}

	// 构建请求
	var bodyReader io.Reader
	if task.Method == "POST" || task.Method == "PUT" {
		bodyReader = strings.NewReader("")
	}
	req, err := http.NewRequest(task.Method, task.URL, bodyReader)
	if err != nil {
		result.ErrorMsg = fmt.Sprintf("构建请求失败: %v", err)
		return result
	}

	// 设置自定义请求头
	if task.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(task.Headers), &headers); err == nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
	}

	// 设置 Cookie
	if task.Cookie != "" {
		req.Header.Set("Cookie", task.Cookie)
	}

	// 设置 Basic Auth
	if task.BasicAuthUser != "" {
		req.SetBasicAuth(task.BasicAuthUser, task.BasicAuthPass)
	}

	// 配置 HTTP Client
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !task.VerifySSL,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	if task.Proxy != "" {
		proxyURL, err := url.Parse(task.Proxy)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	client := &http.Client{
		Timeout:   time.Duration(timeoutSec) * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 不跟随重定向
		},
	}

	// 执行请求
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)
	result.DurationMs = duration.Milliseconds()

	if err != nil {
		result.ErrorMsg = fmt.Sprintf("请求失败: %v", err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	// 检查 SSL 证书过期时间
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		result.SSLExpiry = time.Until(cert.NotAfter).Seconds()
	}

	// 读取响应体进行内容匹配
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 最多读 1MB
	bodyStr := string(bodyBytes)

	switch task.MatchType {
	case "contains":
		result.ContentMatched = strings.Contains(bodyStr, task.MatchContent)
	case "not_contains":
		result.ContentMatched = !strings.Contains(bodyStr, task.MatchContent)
	default:
		result.ContentMatched = true // 无匹配规则视为匹配
	}

	// 判断是否成功：状态码 < 400
	result.Success = resp.StatusCode < 400

	return result
}
```

- [ ] **步骤 2：验证编译并提交**

---

### 任务 17：结果上报器与主入口

**文件：**
- 创建：`agent/reporter.go`
- 创建：`agent/main.go`

- [ ] **步骤 1：实现结果队列和上报器 `reporter.go`**

```go
// agent/reporter.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type timestampedResult struct {
	result    ProbeResult
	createdAt time.Time
}

type ResultQueue struct {
	mu      sync.Mutex
	results []timestampedResult
}

func NewResultQueue() *ResultQueue {
	return &ResultQueue{}
}

func (q *ResultQueue) Add(r ProbeResult) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.results = append(q.results, timestampedResult{result: r, createdAt: time.Now()})
}

func (q *ResultQueue) Drain(maxBatch int, maxAge time.Duration) []ProbeResult {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	// 清理过期结果
	alive := make([]timestampedResult, 0, len(q.results))
	for _, tr := range q.results {
		if now.Sub(tr.createdAt) <= maxAge {
			alive = append(alive, tr)
		}
	}
	q.results = alive

	if len(q.results) == 0 {
		return nil
	}

	batch := make([]ProbeResult, 0, maxBatch)
	for i := 0; i < len(q.results) && i < maxBatch; i++ {
		batch = append(batch, q.results[i].result)
	}
	q.results = q.results[len(batch):]
	return batch
}

type Reporter struct {
	cfg   *Config
	token string
	queue *ResultQueue
}

func NewReporter(cfg *Config, token string, queue *ResultQueue) *Reporter {
	return &Reporter{cfg: cfg, token: token, queue: queue}
}

func (r *Reporter) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(r.cfg.Report.Interval) * time.Second)
	defer ticker.Stop()

	// 心跳 ticker
	hbTicker := time.NewTicker(60 * time.Second)
	defer hbTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			// 关闭前尝试最后一次上报
			r.flush()
			return
		case <-ticker.C:
			r.flush()
		case <-hbTicker.C:
			r.heartbeat()
		}
	}
}

// 注意：上面的 ctx context.Context 需要导入 "context"

func (r *Reporter) flush() {
	maxAge := time.Duration(r.cfg.Report.RetryMaxAge) * time.Second
	batch := r.queue.Drain(r.cfg.Report.BatchSize, maxAge)
	if len(batch) == 0 {
		return
	}

	body := map[string]interface{}{"results": batch}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", r.cfg.Server.URL+"/api/v1/agent/report", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("上报失败: %v，%d 条结果将在下次重试", len(batch))
		// 放回队列
		r.queue.mu.Lock()
		for _, result := range batch {
			r.queue.results = append([]timestampedResult{{result: result, createdAt: time.Now()}}, r.queue.results...)
		}
		r.queue.mu.Unlock()
		return
	}
	defer resp.Body.Close()
	log.Printf("上报成功: %d 条结果", len(batch))
}

func (r *Reporter) heartbeat() {
	req, _ := http.NewRequest("POST", r.cfg.Server.URL+"/api/v1/agent/heartbeat", nil)
	req.Header.Set("Authorization", "Bearer "+r.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("心跳失败: %v", err)
		return
	}
	resp.Body.Close()
}
```

- [ ] **步骤 2：实现 main.go**

```go
// agent/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfgPath := flag.String("c", "agent.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := LoadConfig(*cfgPath)
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 注册或加载 token
	token, err := Register(cfg)
	if err != nil {
		log.Fatal("Agent 注册失败:", err)
	}

	// 创建共享结果队列
	results := NewResultQueue()

	// 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())

	// 启动调度器
	scheduler := NewScheduler(cfg, token, results)
	go scheduler.Start(ctx)

	// 启动上报器
	reporter := NewReporter(cfg, token, results)
	go reporter.Start(ctx)

	fmt.Printf("Agent [%s] 已启动，等待任务...\n", cfg.Agent.Name)

	// 等待信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("正在关闭 Agent...")
	cancel()

	// 等待上报完成（最多 10 秒）
	time.Sleep(10 * time.Second)
	fmt.Println("Agent 已关闭")
}
```

- [ ] **步骤 3：验证编译并提交**

```bash
cd agent && go build -o uptime-agent .
git add agent/ && git commit -m "feat(agent): implement complete agent with probe, scheduler, reporter"
```

---

## 阶段五：Web 前端基础

### 任务 18：Web 项目初始化

- [ ] **步骤 1：使用 Vite 创建 Vue 3 + TypeScript 项目**

```bash
npm create vite@latest web -- --template vue-ts
cd web
npm install
npm install vue-router@4 pinia axios element-plus @element-plus/icons-vue
npm install echarts vue-echarts
npm install -D sass
```

- [ ] **步骤 2：配置 `vite.config.ts`**

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  },
  resolve: {
    alias: {
      '@': '/src'
    }
  }
})
```

- [ ] **步骤 3：提交**

```bash
git add web/ && git commit -m "feat(web): init Vue 3 + TypeScript project"
```

---

### 任务 19：路由与认证状态

**文件：**
- 创建：`web/src/router/index.ts`
- 创建：`web/src/stores/auth.ts`

- [ ] **步骤 1：实现 Pinia auth store**

```typescript
// web/src/stores/auth.ts
import { defineStore } from 'pinia'
import { login, logout } from '@/api/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userId: null as number | null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
  },
  actions: {
    async login(username: string, password: string) {
      const res = await login(username, password)
      this.token = res.data.token
      localStorage.setItem('token', this.token)
    },
    async logout() {
      try { await logout() } catch {}
      this.token = ''
      this.userId = null
      localStorage.removeItem('token')
    },
  },
})
```

- [ ] **步骤 2：实现路由配置**

```typescript
// web/src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginView.vue'),
    meta: { guest: true },
  },
  {
    path: '/',
    component: () => import('@/components/layout/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/dashboard' },
      { path: 'dashboard', name: 'Dashboard', component: () => import('@/views/dashboard/DashboardView.vue') },
      { path: 'monitors', name: 'MonitorList', component: () => import('@/views/monitor/MonitorList.vue') },
      { path: 'monitors/create', name: 'MonitorCreate', component: () => import('@/views/monitor/MonitorForm.vue') },
      { path: 'monitors/:id/edit', name: 'MonitorEdit', component: () => import('@/views/monitor/MonitorForm.vue') },
      { path: 'monitors/:id', name: 'MonitorDetail', component: () => import('@/views/monitor/MonitorDetail.vue') },
      { path: 'alerts/channels', name: 'AlertChannels', component: () => import('@/views/alert/AlertChannelList.vue') },
      { path: 'alerts/history', name: 'AlertHistory', component: () => import('@/views/alert/AlertHistory.vue') },
      { path: 'agents', name: 'AgentList', component: () => import('@/views/agent/AgentList.vue') },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    next('/login')
  } else if (to.meta.guest && auth.isLoggedIn) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
```

- [ ] **步骤 3：更新 `main.ts`**

```typescript
// web/src/main.ts
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import './styles/variables.scss'
import './styles/element-override.scss'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(ElementPlus)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.mount('#app')
```

- [ ] **步骤 4：提交**

---

### 任务 20：布局与苹果风格主题

**文件：**
- 创建：`web/src/styles/variables.scss`
- 创建：`web/src/styles/element-override.scss`
- 创建：`web/src/components/layout/AppLayout.vue`
- 创建：`web/src/components/layout/Sidebar.vue`

- [ ] **步骤 1：定义苹果风格 CSS 变量 `variables.scss`**

```scss
// web/src/styles/variables.scss
:root {
  --color-primary: #0071E3;
  --color-primary-light: #64ACFF;
  --color-success: #34C759;
  --color-warning: #FF9500;
  --color-danger: #FF3B30;
  --color-text-primary: #1D1D1F;
  --color-text-secondary: #86868B;
  --color-bg: #FBFBFD;
  --color-bg-card: #FFFFFF;
  --color-border: #E5E5E7;
  --radius-card: 16px;
  --radius-btn: 12px;
  --shadow-card: 0 1px 3px rgba(0, 0, 0, 0.04);
  --shadow-card-hover: 0 4px 12px rgba(0, 0, 0, 0.08);
  --font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text', 'Helvetica Neue', Arial, sans-serif;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: var(--font-family);
  color: var(--color-text-primary);
  background: var(--color-bg);
  -webkit-font-smoothing: antialiased;
}
```

- [ ] **步骤 2：Element Plus 主题覆盖 `element-override.scss`**

```scss
// web/src/styles/element-override.scss
// 覆盖 Element Plus 默认样式以匹配苹果风格
.el-button--primary {
  background-color: var(--color-primary) !important;
  border-color: var(--color-primary) !important;
  border-radius: var(--radius-btn) !important;
  font-weight: 500 !important;
}
.el-card {
  border-radius: var(--radius-card) !important;
  border: none !important;
  box-shadow: var(--shadow-card) !important;
  &:hover {
    box-shadow: var(--shadow-card-hover) !important;
  }
}
.el-table {
  border-radius: var(--radius-card) !important;
  th.el-table__cell {
    background-color: #F5F5F7 !important;
    color: var(--color-text-secondary) !important;
    font-weight: 600 !important;
    font-size: 12px !important;
    text-transform: uppercase !important;
    letter-spacing: 0.5px !important;
  }
}
.el-input__wrapper {
  border-radius: 10px !important;
}
.el-menu {
  border-right: none !important;
}
```

- [ ] **步骤 3：实现 AppLayout.vue**

```vue
<!-- web/src/components/layout/AppLayout.vue -->
<template>
  <div class="app-layout">
    <Sidebar />
    <div class="app-main">
      <header class="app-header">
        <div class="header-right">
          <el-badge :value="0" :hidden="true">
            <el-icon><Bell /></el-icon>
          </el-badge>
          <el-dropdown>
            <span class="user-info">管理员 <el-icon><ArrowDown /></el-icon></span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>
      <main class="app-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import Sidebar from './Sidebar.vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

const handleLogout = async () => {
  await auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
}
.app-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-left: 240px;
}
.app-header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 32px;
  background: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
}
.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}
.user-info {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  font-size: 14px;
  color: var(--color-text-secondary);
}
.app-content {
  flex: 1;
  padding: 32px;
}
</style>
```

- [ ] **步骤 4：实现 Sidebar.vue**

```vue
<!-- web/src/components/layout/Sidebar.vue -->
<template>
  <aside class="sidebar">
    <div class="sidebar-logo">
      <h2>Uptime</h2>
      <span>Monitor</span>
    </div>
    <el-menu :default-active="activeMenu" router>
      <el-menu-item index="/dashboard">
        <el-icon><DataAnalysis /></el-icon>
        <span>仪表盘</span>
      </el-menu-item>
      <el-menu-item index="/monitors">
        <el-icon><Monitor /></el-icon>
        <span>监控列表</span>
      </el-menu-item>
      <el-menu-item index="/alerts/channels">
        <el-icon><Bell /></el-icon>
        <span>告警管理</span>
      </el-menu-item>
      <el-menu-item index="/agents">
        <el-icon><Connection /></el-icon>
        <span>探测点</span>
      </el-menu-item>
    </el-menu>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activeMenu = computed(() => {
  const path = route.path
  if (path.startsWith('/monitors')) return '/monitors'
  if (path.startsWith('/alerts')) return '/alerts/channels'
  return path
})
</script>

<style scoped>
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 240px;
  background: var(--color-bg-card);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  z-index: 100;
}
.sidebar-logo {
  padding: 24px 20px;
  border-bottom: 1px solid var(--color-border);
}
.sidebar-logo h2 {
  font-size: 20px;
  font-weight: 700;
  color: var(--color-text-primary);
  display: inline;
}
.sidebar-logo span {
  font-size: 20px;
  font-weight: 300;
  color: var(--color-primary);
}
.el-menu {
  flex: 1;
  padding-top: 12px;
}
</style>
```

- [ ] **步骤 5：提交**

---

### 任务 21：API 客户端层

**文件：**
- 创建：`web/src/api/request.ts`
- 创建：`web/src/api/auth.ts`
- 创建：`web/src/api/monitor.ts`
- 创建：`web/src/api/alert.ts`
- 创建：`web/src/api/agent.ts`
- 创建：`web/src/api/dashboard.ts`

- [ ] **步骤 1：实现 Axios 封装 `request.ts`**

```typescript
// web/src/api/request.ts
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import router from '@/router'

const request = axios.create({
  baseURL: '/api',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  return config
})

request.interceptors.response.use(
  (response) => {
    const res = response.data
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      if (res.code === 40001) {
        const auth = useAuthStore()
        auth.logout()
        router.push('/login')
      }
      return Promise.reject(new Error(res.message))
    }
    return res
  },
  (error) => {
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

export default request
```

- [ ] **步骤 2：实现各 API 模块**

```typescript
// web/src/api/auth.ts
import request from './request'

export const login = (username: string, password: string) =>
  request.post('/auth/login', { username, password })

export const logout = () => request.post('/auth/logout')

export const getMe = () => request.get('/auth/me')
```

```typescript
// web/src/api/monitor.ts
import request from './request'

export const listMonitors = (params: { page: number; page_size: number; keyword?: string }) =>
  request.get('/monitors', { params })

export const getMonitor = (id: number) => request.get(`/monitors/${id}`)

export const createMonitor = (data: any) => request.post('/monitors', data)

export const updateMonitor = (id: number, data: any) => request.put(`/monitors/${id}`, data)

export const deleteMonitor = (id: number) => request.delete(`/monitors/${id}`)

export const toggleMonitor = (id: number, enabled: boolean) =>
  request.patch(`/monitors/${id}/enabled`, { enabled })

export const getMonitorMetrics = (id: number) => request.get(`/monitors/${id}/metrics`)
```

```typescript
// web/src/api/alert.ts
import request from './request'

export const listAlertChannels = () => request.get('/alerts/channels')

export const createAlertChannel = (data: any) => request.post('/alerts/channels', data)

export const updateAlertChannel = (id: number, data: any) => request.put(`/alerts/channels/${id}`, data)

export const deleteAlertChannel = (id: number) => request.delete(`/alerts/channels/${id}`)

export const testAlertChannel = (id: number) => request.post(`/alerts/channels/${id}/test`)

export const listAlertHistory = (params: { page: number; page_size: number; monitor_id?: number; status?: string }) =>
  request.get('/alerts/history', { params })
```

```typescript
// web/src/api/agent.ts
import request from './request'

export const listAgents = () => request.get('/agents')

export const deleteAgent = (id: number) => request.delete(`/agents/${id}`)
```

```typescript
// web/src/api/dashboard.ts
import request from './request'

export const getDashboardOverview = () => request.get('/dashboard/overview')
```

- [ ] **步骤 3：提交**

---

### 任务 22：通用组件

**文件：**
- 创建：`web/src/components/common/StatusDot.vue`
- 创建：`web/src/components/common/StatCard.vue`

- [ ] **步骤 1：实现 StatusDot.vue**

```vue
<!-- web/src/components/common/StatusDot.vue -->
<template>
  <span class="status-dot" :class="statusClass" />
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ status: 'online' | 'offline' | 'warning' }>()

const statusClass = computed(() => ({
  'dot-online': props.status === 'online',
  'dot-offline': props.status === 'offline',
  'dot-warning': props.status === 'warning',
}))
</script>

<style scoped>
.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
}
.dot-online { background: var(--color-success); }
.dot-offline { background: var(--color-danger); }
.dot-warning { background: var(--color-warning); }
</style>
```

- [ ] **步骤 2：实现 StatCard.vue**

```vue
<!-- web/src/components/common/StatCard.vue -->
<template>
  <div class="stat-card">
    <div class="stat-label">{{ label }}</div>
    <div class="stat-value" :style="{ color: valueColor }">
      {{ value }}<span v-if="suffix" class="stat-suffix">{{ suffix }}</span>
    </div>
    <div class="stat-trend" :class="trendClass" v-if="trend">
      {{ trend }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  label: string
  value: string | number
  suffix?: string
  valueColor?: string
  trend?: string
  trendUp?: boolean
}>()

const trendClass = computed(() => props.trendUp ? 'trend-up' : 'trend-down')
</script>

<style scoped>
.stat-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-card);
  padding: 20px;
  box-shadow: var(--shadow-card);
}
.stat-label {
  font-size: 12px;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 1px;
}
.stat-value {
  font-size: 36px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-top: 4px;
}
.stat-suffix {
  font-size: 18px;
  color: var(--color-text-secondary);
}
.stat-trend {
  font-size: 12px;
  margin-top: 4px;
}
.trend-up { color: var(--color-success); }
.trend-down { color: var(--color-danger); }
</style>
```

- [ ] **步骤 3：提交**

---

## 阶段六：Web 页面

### 任务 23：登录页

**文件：**
- 创建：`web/src/views/login/LoginView.vue`

- [ ] **步骤 1：实现登录页**

```vue
<!-- web/src/views/login/LoginView.vue -->
<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <h1>Uptime<span>Monitor</span></h1>
        <p>网站存活探测系统</p>
      </div>
      <el-form :model="form" :rules="rules" ref="formRef" @submit.prevent="handleLogin">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="用户名" prefix-icon="User" size="large" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码" prefix-icon="Lock" size="large" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleLogin" size="large" style="width: 100%;">
            登 录
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const auth = useAuthStore()
const formRef = ref()
const loading = ref(false)

const form = reactive({ username: '', password: '' })
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

const handleLogin = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    await auth.login(form.username, form.password)
    ElMessage.success('登录成功')
    router.push('/dashboard')
  } catch (e: any) {
    ElMessage.error(e.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-card {
  background: var(--color-bg-card);
  border-radius: 24px;
  padding: 48px 40px;
  width: 400px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}
.login-header {
  text-align: center;
  margin-bottom: 36px;
}
.login-header h1 {
  font-size: 28px;
  font-weight: 700;
}
.login-header h1 span {
  color: var(--color-primary);
  font-weight: 300;
}
.login-header p {
  color: var(--color-text-secondary);
  margin-top: 8px;
  font-size: 14px;
}
</style>
```

- [ ] **步骤 2：提交**

---

### 任务 24：Dashboard 看板

**文件：**
- 创建：`web/src/views/dashboard/DashboardView.vue`
- 创建：`web/src/components/charts/UptimeChart.vue`
- 创建：`web/src/components/charts/HealthDonut.vue`

- [ ] **步骤 1：实现 DashboardView.vue**

```vue
<!-- web/src/views/dashboard/DashboardView.vue -->
<template>
  <div class="dashboard">
    <h1 class="page-title">仪表盘</h1>

    <!-- 统计卡片 -->
    <div class="stats-row">
      <StatCard label="监控总数" :value="overview.total_monitors" :trend="'本周新增 ' + 0" :trendUp="true" />
      <StatCard label="在线率" :value="(overview.online_rate * 100).toFixed(1) + '%'" :valueColor="overview.online_rate > 0.99 ? '#34C759' : '#FF9500'" />
      <StatCard label="活跃告警" :value="overview.active_alerts" :valueColor="overview.active_alerts > 0 ? '#FF3B30' : '#34C759'" />
      <StatCard label="平均响应" :value="overview.avg_response_ms.toFixed(0)" suffix="ms" />
    </div>

    <!-- 图表行 -->
    <div class="charts-row">
      <div class="chart-wide">
        <el-card shadow="never">
          <template #header>
            <div class="chart-header">
              <span>可用率趋势</span>
              <div class="time-switch">
                <el-button :type="range === '24h' ? 'primary' : 'default'" size="small" @click="range = '24h'">24h</el-button>
                <el-button :type="range === '7d' ? 'primary' : 'default'" size="small" @click="range = '7d'">7d</el-button>
                <el-button :type="range === '30d' ? 'primary' : 'default'" size="small" @click="range = '30d'">30d</el-button>
              </div>
            </div>
          </template>
          <UptimeChart :range="range" />
        </el-card>
      </div>
      <div class="chart-narrow">
        <el-card shadow="never">
          <template #header><span>健康度分布</span></template>
          <HealthDonut :healthy="healthyCount" :warning="warningCount" :danger="dangerCount" />
        </el-card>
      </div>
    </div>

    <!-- 告警 + TOP5 -->
    <div class="bottom-row">
      <el-card shadow="never">
        <template #header>
          <div class="chart-header">
            <span>最近告警</span>
            <el-button text type="primary" @click="$router.push('/alerts/history')">查看全部 →</el-button>
          </div>
        </template>
        <div v-for="alert in overview.recent_alerts" :key="alert.id" class="alert-item" :class="'alert-' + alert.status">
          <StatusDot :status="alert.status === 'firing' ? 'offline' : 'online'" />
          <div class="alert-info">
            <div class="alert-name">{{ alert.monitor?.name || '未知' }}</div>
            <div class="alert-detail">{{ alert.alert_type }} · {{ alert.message }}</div>
          </div>
          <div class="alert-time">{{ formatTime(alert.triggered_at) }}</div>
        </div>
        <el-empty v-if="!overview.recent_alerts?.length" description="暂无告警" />
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { getDashboardOverview } from '@/api/dashboard'
import StatCard from '@/components/common/StatCard.vue'
import StatusDot from '@/components/common/StatusDot.vue'
import UptimeChart from '@/components/charts/UptimeChart.vue'
import HealthDonut from '@/components/charts/HealthDonut.vue'

const range = ref('24h')
const overview = ref<any>({
  total_monitors: 0,
  online_rate: 0,
  active_alerts: 0,
  avg_response_ms: 0,
  recent_alerts: [],
})

const healthyCount = computed(() => Math.round(overview.value.total_monitors * overview.value.online_rate))
const warningCount = computed(() => 0)
const dangerCount = computed(() => overview.value.total_monitors - healthyCount.value)

const formatTime = (t: string) => {
  if (!t) return ''
  const d = new Date(t)
  const diff = (Date.now() - d.getTime()) / 1000
  if (diff < 60) return '刚刚'
  if (diff < 3600) return Math.floor(diff / 60) + ' 分钟前'
  if (diff < 86400) return Math.floor(diff / 3600) + ' 小时前'
  return Math.floor(diff / 86400) + ' 天前'
}

onMounted(async () => {
  const res = await getDashboardOverview()
  overview.value = res.data
})
</script>

<style scoped>
.page-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 24px;
}
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}
.charts-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;
  margin-bottom: 24px;
}
.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}
.time-switch {
  display: flex;
  gap: 4px;
}
.bottom-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
.alert-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  margin-bottom: 8px;
}
.alert-firing { background: #FFF5F5; border-left: 3px solid #FF3B30; }
.alert-resolved { background: #F0FFF4; border-left: 3px solid #34C759; }
.alert-info { flex: 1; }
.alert-name { font-weight: 500; font-size: 14px; }
.alert-detail { font-size: 12px; color: var(--color-text-secondary); margin-top: 2px; }
.alert-time { font-size: 11px; color: var(--color-text-secondary); }
</style>
```

- [ ] **步骤 2：实现 UptimeChart.vue**

```vue
<!-- web/src/components/charts/UptimeChart.vue -->
<template>
  <v-chart :option="chartOption" style="height: 300px;" autoresize />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([LineChart, GridComponent, TooltipComponent, CanvasRenderer])

const props = defineProps<{ range: string }>()

// TODO: 从 API 获取真实数据，这里使用示例
const chartOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: 50, right: 20, top: 20, bottom: 30 },
  xAxis: { type: 'category', data: generateTimeLabels(props.range) },
  yAxis: { type: 'value', min: 98, max: 100, axisLabel: { formatter: '{value}%' } },
  series: [{
    type: 'line',
    data: generateMockData(),
    smooth: true,
    lineStyle: { color: '#0071E3', width: 2.5 },
    areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [
      { offset: 0, color: 'rgba(0, 113, 227, 0.15)' },
      { offset: 1, color: 'rgba(0, 113, 227, 0)' }
    ]}},
    itemStyle: { color: '#0071E3' },
  }],
}))

function generateTimeLabels(range: string) {
  const labels = []
  const count = range === '24h' ? 24 : range === '7d' ? 7 : 30
  for (let i = count; i >= 0; i--) {
    labels.push(i === 0 ? '现在' : `${i}${range === '24h' ? 'h' : 'd'}前`)
  }
  return labels
}

function generateMockData() {
  return Array.from({ length: 25 }, () => 99 + Math.random())
}
</script>
```

- [ ] **步骤 3：实现 HealthDonut.vue**

```vue
<!-- web/src/components/charts/HealthDonut.vue -->
<template>
  <v-chart :option="chartOption" style="height: 250px;" autoresize />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([PieChart, LegendComponent, CanvasRenderer])

const props = defineProps<{ healthy: number; warning: number; danger: number }>()

const chartOption = computed(() => ({
  series: [{
    type: 'pie',
    radius: ['55%', '75%'],
    label: { show: true, position: 'center', formatter: () => `${props.healthy + props.warning + props.danger}`, fontSize: 28, fontWeight: 600 },
    data: [
      { value: props.healthy, name: '健康', itemStyle: { color: '#34C759' } },
      { value: props.warning, name: '警告', itemStyle: { color: '#FF9500' } },
      { value: props.danger, name: '异常', itemStyle: { color: '#FF3B30' } },
    ].filter(d => d.value > 0),
  }],
}))
</script>
```

- [ ] **步骤 4：提交**

---

### 任务 25-29：其他页面

由于篇幅限制，以下页面提供实现要点，具体代码参考任务 23-24 的模式：

**任务 25：监控列表页 `MonitorList.vue`**
- 使用 `el-table` 展示所有监控任务
- 列：名称、URL、状态（StatusDot）、可用率、延迟、操作（编辑/删除/启停）
- 顶部搜索框 + 新增按钮
- 分页组件

**任务 26：监控表单 `MonitorForm.vue`**
- 分步表单（el-steps）
- Step 1：基本信息（名称、URL）
- Step 2：请求配置（方法、Headers JSON 编辑器、Cookie、Basic Auth、代理、SSL 验证）
- Step 3：告警规则（匹配类型、匹配内容、状态码阈值、延迟阈值、连续次数、频率）
- Step 4：通知通道（多选告警通道、选择 Agent、启用开关）
- 编辑模式通过 route.params.id 判断，加载已有数据

**任务 27：监控详情页 `MonitorDetail.vue`**
- 使用 UptimeChart 和 LatencyChart 展示趋势
- 状态码分布饼图
- 各探测点对比表
- 告警时间线（el-timeline）
- 历史记录表 + CSV 导出

**任务 28：告警通道管理**
- `AlertChannelList.vue`：表格列出所有通道，操作列含编辑/删除/测试按钮
- `AlertChannelForm.vue`：el-dialog 弹窗表单，字段：名称、类型（下拉：钉钉/企微/飞书/Webhook）、Webhook URL、签名密钥
- `AlertHistory.vue`：表格 + 筛选（监控/状态/时间范围）+ 分页

**任务 29：Agent 列表页 `AgentList.vue`**
- 表格展示所有 Agent
- 列：名称、位置、状态（StatusDot）、最后心跳时间、操作（删除）

- [ ] **步骤 1-5：逐个实现上述页面并提交**

---

## 阶段七：集成与部署

### 任务 30：嵌入 Web 到 Server

- [ ] **步骤 1：更新 `server/cmd/main.go` 中的 embed 指令**

确保 `server/web/dist/` 目录包含 Vue 构建产物。在构建 Server 前，先构建 Web：

```bash
# 构建 Web
cd web && npm run build
# 复制构建产物到 server/web/dist/
cp -r dist/* ../server/web/dist/
# 构建 Server
cd ../server && go build -o uptime-server cmd/main.go
```

- [ ] **步骤 2：添加构建脚本 `build.sh`**

```bash
#!/bin/bash
set -e

echo "构建 Web..."
cd web
npm install
npm run build
cd ..

echo "复制 Web 构建产物..."
rm -rf server/web/dist
cp -r web/dist server/web/dist

echo "构建 Server..."
cd server
go build -o uptime-server cmd/main.go
cd ..

echo "构建完成！"
echo "Server 二进制: server/uptime-server"
echo "运行: ./server/uptime-server -c server/config.yaml"
```

- [ ] **步骤 3：提交构建脚本**

```bash
git add build.sh && git commit -m "feat: add build script for web+server"
```

---

### 任务 31：端到端集成测试

- [ ] **步骤 1：启动 Server**

```bash
cd server && ./uptime-server -c config.yaml
```

- [ ] **步骤 2：测试 Web 登录**

打开浏览器访问 `http://localhost:8080`，使用 admin/admin123 登录，验证 Dashboard 加载。

- [ ] **步骤 3：注册 Agent**

```bash
cd agent && ./uptime-agent -c agent.yaml
```

验证 Agent 注册成功并获取到 token。

- [ ] **步骤 4：创建监控任务**

在 Web 界面创建一个监控任务（例如监控 `https://www.baidu.com`），分配给刚注册的 Agent。

- [ ] **步骤 5：验证检测结果**

等待 1-2 分钟，在 Dashboard 查看：
- 监控状态变为在线
- 可用率和延迟数据开始更新
- 历史记录中有检测数据

- [ ] **步骤 6：测试告警**

创建一个会失败的监控任务（例如监控不存在的域名），验证：
- 连续失败计数增加
- 达到阈值后触发告警
- 告警记录出现在告警历史中
- 通知发送到配置的通道（如果配置了）

- [ ] **步骤 7：提交最终集成**

```bash
git add -A && git commit -m "feat: complete uptime monitor system"
```

---

## 总结

本计划共 31 个任务，分 7 个阶段：
1. **Server 基础**（任务 1-4）：项目初始化、数据模型、Repository、响应格式
2. **Server API**（任务 5-11）：认证、监控 CRUD、Agent API、告警、通知、Metrics
3. **Server 组装**（任务 12）：Main 入口、后台任务、路由
4. **Agent**（任务 13-17）：配置、注册、任务同步、探测引擎、结果上报
5. **Web 基础**（任务 18-22）：项目初始化、路由、布局主题、API 客户端
6. **Web 页面**（任务 23-29）：登录、Dashboard、监控管理、告警管理、Agent 管理
7. **集成部署**（任务 30-31）：嵌入构建、端到端测试

**预计工作量：** 2-3 周（全职开发）

**下一步：** 使用 `superpowers:subagent-driven-development` 或 `superpowers:executing-plans` 逐任务实现。
