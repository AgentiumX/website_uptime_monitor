package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
	"uptime-monitor/server/internal/service"
)

// AgentHandler handles agent API endpoints.
type AgentHandler struct {
	agentSvc   *service.AgentService
	alertSvc   *service.AlertService
	metricsSvc *service.MetricsService
}

// NewAgentHandler creates a new AgentHandler.
func NewAgentHandler(agentSvc *service.AgentService, alertSvc *service.AlertService, metricsSvc *service.MetricsService) *AgentHandler {
	return &AgentHandler{
		agentSvc:   agentSvc,
		alertSvc:   alertSvc,
		metricsSvc: metricsSvc,
	}
}

// Register handles POST /api/v1/agent/register (public route).
func (h *AgentHandler) Register(c *gin.Context) {
	var req model.AgentRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, 40101, "参数错误")
		return
	}

	resp, err := h.agentSvc.Register(req)
	if err != nil {
		Error(c, http.StatusOK, 40001, err.Error())
		return
	}

	Success(c, resp)
}

// GetTasks handles GET /api/v1/agent/tasks.
func (h *AgentHandler) GetTasks(c *gin.Context) {
	val, exists := c.Get("agent_id")
	if !exists {
		Error(c, http.StatusOK, 40001, "未获取到Agent信息")
		return
	}
	agentID, ok := val.(uint)
	if !ok {
		Error(c, http.StatusOK, 40001, "Agent信息类型异常")
		return
	}

	tasks, updatedAt, err := h.agentSvc.GetTasks(agentID)
	if err != nil {
		Error(c, http.StatusOK, 50001, "获取任务列表失败")
		return
	}

	Success(c, gin.H{
		"tasks":      tasks,
		"updated_at": updatedAt,
	})
}

// Report handles POST /api/v1/agent/report.
func (h *AgentHandler) Report(c *gin.Context) {
	val, exists := c.Get("agent_id")
	if !exists {
		Error(c, http.StatusOK, 40001, "未获取到Agent信息")
		return
	}
	agentID, ok := val.(uint)
	if !ok {
		Error(c, http.StatusOK, 40001, "Agent信息类型异常")
		return
	}

	var req model.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, 40101, "参数错误")
		return
	}

	for _, result := range req.Results {
		h.metricsSvc.Record(result)
		h.alertSvc.Evaluate(agentID, result)
	}

	Success(c, nil)
}

// Heartbeat handles POST /api/v1/agent/heartbeat.
func (h *AgentHandler) Heartbeat(c *gin.Context) {
	val, exists := c.Get("agent_id")
	if !exists {
		Error(c, http.StatusOK, 40001, "未获取到Agent信息")
		return
	}
	agentID, ok := val.(uint)
	if !ok {
		Error(c, http.StatusOK, 40001, "Agent信息类型异常")
		return
	}

	if err := h.agentSvc.Heartbeat(agentID); err != nil {
		Error(c, http.StatusOK, 50001, "心跳上报失败")
		return
	}

	Success(c, nil)
}

// List handles GET /api/agents — returns all registered agents.
func (h *AgentHandler) List(c *gin.Context) {
	agents, err := repository.ListAgents()
	if err != nil {
		Error(c, http.StatusOK, 50001, "查询失败")
		return
	}
	Success(c, agents)
}

// Delete handles DELETE /api/agents/:id — removes an agent.
func (h *AgentHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		Error(c, http.StatusBadRequest, 40101, "无效的ID参数")
		return
	}

	if err := repository.DeleteAgent(id); err != nil {
		Error(c, http.StatusOK, 50001, "删除失败")
		return
	}

	Success(c, nil)
}
