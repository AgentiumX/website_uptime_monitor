package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/service"
)

// MonitorHandler handles monitor CRUD API endpoints.
type MonitorHandler struct {
	monitorSvc *service.MonitorService
}

// NewMonitorHandler creates a new MonitorHandler.
func NewMonitorHandler(monitorSvc *service.MonitorService) *MonitorHandler {
	return &MonitorHandler{monitorSvc: monitorSvc}
}

// parseUintParam extracts a uint path parameter from the Gin context.
func parseUintParam(c *gin.Context, key string) (uint, error) {
	s := c.Param(key)
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}

// List handles GET /api/monitors — paginated list with optional keyword search.
func (h *MonitorHandler) List(c *gin.Context) {
	var q model.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		Error(c, http.StatusBadRequest, 40101, "参数错误")
		return
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}

	list, total, err := h.monitorSvc.List(q.Page, q.PageSize, q.Keyword)
	if err != nil {
		Error(c, http.StatusOK, 50001, "查询失败")
		return
	}

	Paginated(c, list, total, q.Page, q.PageSize)
}

// Get handles GET /api/monitors/:id.
func (h *MonitorHandler) Get(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		Error(c, http.StatusBadRequest, 40101, "无效的ID参数")
		return
	}

	m, err := h.monitorSvc.Get(id)
	if err != nil {
		Error(c, http.StatusOK, 50001, "监控任务不存在")
		return
	}

	Success(c, m)
}

// Create handles POST /api/monitors.
func (h *MonitorHandler) Create(c *gin.Context) {
	var m model.Monitor
	if err := c.ShouldBindJSON(&m); err != nil {
		Error(c, http.StatusBadRequest, 40101, "参数错误")
		return
	}

	if err := h.monitorSvc.Create(&m); err != nil {
		Error(c, http.StatusOK, 40101, err.Error())
		return
	}

	Success(c, m)
}

// Update handles PUT /api/monitors/:id.
func (h *MonitorHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		Error(c, http.StatusBadRequest, 40101, "无效的ID参数")
		return
	}

	var m model.Monitor
	if err := c.ShouldBindJSON(&m); err != nil {
		Error(c, http.StatusBadRequest, 40101, "参数错误")
		return
	}
	m.ID = id

	if err := h.monitorSvc.Update(&m); err != nil {
		Error(c, http.StatusOK, 50001, "更新失败")
		return
	}

	Success(c, m)
}

// Delete handles DELETE /api/monitors/:id.
func (h *MonitorHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		Error(c, http.StatusBadRequest, 40101, "无效的ID参数")
		return
	}

	if err := h.monitorSvc.Delete(id); err != nil {
		Error(c, http.StatusOK, 50001, "删除失败")
		return
	}

	Success(c, nil)
}

// ToggleEnabled handles PATCH /api/monitors/:id/enabled.
func (h *MonitorHandler) ToggleEnabled(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		Error(c, http.StatusBadRequest, 40101, "无效的ID参数")
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, 40101, "参数错误")
		return
	}

	if err := h.monitorSvc.Toggle(id, req.Enabled); err != nil {
		Error(c, http.StatusOK, 50001, "操作失败")
		return
	}

	Success(c, nil)
}
