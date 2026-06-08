package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"uptime-monitor/server/internal/repository"
	"uptime-monitor/server/internal/service"
)

// DashboardHandler handles dashboard API endpoints.
type DashboardHandler struct {
	metricsSvc *service.MetricsService
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(metricsSvc *service.MetricsService) *DashboardHandler {
	return &DashboardHandler{metricsSvc: metricsSvc}
}

// Overview handles GET /api/dashboard/overview — returns dashboard summary statistics.
func (h *DashboardHandler) Overview(c *gin.Context) {
	totalMonitors, err := repository.CountMonitors()
	if err != nil {
		Error(c, http.StatusOK, 50001, "查询失败")
		return
	}

	activeAlerts, err := repository.CountActiveAlerts()
	if err != nil {
		Error(c, http.StatusOK, 50001, "查询失败")
		return
	}

	uptimeRate, avgLatencyMs := h.metricsSvc.GetOverallStats()

	// Get recent alerts
	recentAlerts, _, err := repository.ListAlertHistory(1, 10, 0, "", "", time.Time{}, time.Time{})
	if err != nil {
		Error(c, http.StatusOK, 50001, "查询失败")
		return
	}

	Success(c, gin.H{
		"total_monitors":  totalMonitors,
		"online_rate":     uptimeRate,
		"active_alerts":   activeAlerts,
		"avg_response_ms": avgLatencyMs,
		"recent_alerts":   recentAlerts,
	})
}
