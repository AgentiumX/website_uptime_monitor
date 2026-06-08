package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"uptime-monitor/server/internal/service"
)

// MetricsHandler handles metrics API endpoints.
type MetricsHandler struct {
	metricsSvc *service.MetricsService
}

// NewMetricsHandler creates a new MetricsHandler.
func NewMetricsHandler(metricsSvc *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{metricsSvc: metricsSvc}
}

// GetMonitorMetrics handles GET /api/monitors/:id/metrics — returns current metrics snapshot.
func (h *MetricsHandler) GetMonitorMetrics(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		Error(c, http.StatusBadRequest, 40101, "无效的ID参数")
		return
	}

	metrics := h.metricsSvc.GetMonitorMetrics(id)
	Success(c, metrics)
}
