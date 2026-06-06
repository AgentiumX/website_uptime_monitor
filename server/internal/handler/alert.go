package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
	"uptime-monitor/server/internal/service"
)

// AlertHandler handles alert channel and history API endpoints.
type AlertHandler struct {
	notifySvc *service.NotifyService
}

// NewAlertHandler creates a new AlertHandler.
func NewAlertHandler(notifySvc *service.NotifyService) *AlertHandler {
	return &AlertHandler{notifySvc: notifySvc}
}

// ListChannels handles GET /api/alerts/channels.
func (h *AlertHandler) ListChannels(c *gin.Context) {
	channels, err := repository.ListAlertChannels()
	if err != nil {
		Error(c, http.StatusOK, 50001, "查询失败")
		return
	}
	Success(c, channels)
}

// CreateChannel handles POST /api/alerts/channels.
func (h *AlertHandler) CreateChannel(c *gin.Context) {
	var ch model.AlertChannel
	if err := c.ShouldBindJSON(&ch); err != nil {
		Error(c, http.StatusBadRequest, 40101, "参数错误")
		return
	}

	if err := repository.CreateAlertChannel(&ch); err != nil {
		Error(c, http.StatusOK, 50001, "创建失败")
		return
	}

	Success(c, ch)
}

// UpdateChannel handles PUT /api/alerts/channels/:id.
func (h *AlertHandler) UpdateChannel(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		Error(c, http.StatusBadRequest, 40101, "无效的ID参数")
		return
	}

	var ch model.AlertChannel
	if err := c.ShouldBindJSON(&ch); err != nil {
		Error(c, http.StatusBadRequest, 40101, "参数错误")
		return
	}
	ch.ID = id

	if err := repository.UpdateAlertChannel(&ch); err != nil {
		Error(c, http.StatusOK, 50001, "更新失败")
		return
	}

	Success(c, ch)
}

// DeleteChannel handles DELETE /api/alerts/channels/:id.
func (h *AlertHandler) DeleteChannel(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		Error(c, http.StatusBadRequest, 40101, "无效的ID参数")
		return
	}

	if err := repository.DeleteAlertChannel(id); err != nil {
		Error(c, http.StatusOK, 50001, "删除失败")
		return
	}

	Success(c, nil)
}

// TestChannel handles POST /api/alerts/channels/:id/test.
func (h *AlertHandler) TestChannel(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		Error(c, http.StatusBadRequest, 40101, "无效的ID参数")
		return
	}

	ch, err := repository.GetAlertChannelByID(id)
	if err != nil {
		Error(c, http.StatusOK, 50001, "告警通道不存在")
		return
	}

	if err := h.notifySvc.SendTest(ch); err != nil {
		Error(c, http.StatusOK, 50001, "测试发送失败: "+err.Error())
		return
	}

	Success(c, nil)
}

// ListHistory handles GET /api/alerts/history — paginated with filters.
func (h *AlertHandler) ListHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	monitorID, _ := strconv.ParseUint(c.Query("monitor_id"), 10, 64)
	status := c.Query("status")
	alertType := c.Query("alert_type")

	var startTime, endTime time.Time
	if s := c.Query("start_time"); s != "" {
		startTime, _ = time.Parse("2006-01-02 15:04:05", s)
	}
	if s := c.Query("end_time"); s != "" {
		endTime, _ = time.Parse("2006-01-02 15:04:05", s)
	}

	records, total, err := repository.ListAlertHistory(
		page, pageSize,
		uint(monitorID), status, alertType,
		startTime, endTime,
	)
	if err != nil {
		Error(c, http.StatusOK, 50001, "查询失败")
		return
	}

	Paginated(c, records, total, page, pageSize)
}
