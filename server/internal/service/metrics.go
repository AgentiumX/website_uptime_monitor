package service

import (
	"uptime-monitor/server/internal/model"
)

// MetricsService handles recording probe results to TSDB.
// Stub implementation — full logic will be added in Task 11.
type MetricsService struct{}

// NewMetricsService creates a new MetricsService.
func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

// Record writes a probe result to the metrics store.
// Stub — no-op until Task 11 implementation.
func (s *MetricsService) Record(result model.ProbeResult) {
	// TODO: implement metrics recording in Task 11
}
