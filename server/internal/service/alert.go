package service

import (
	"uptime-monitor/server/internal/model"
)

// AlertService handles alert evaluation logic.
// Stub implementation — full logic will be added in Task 9.
type AlertService struct{}

// NewAlertService creates a new AlertService.
func NewAlertService() *AlertService {
	return &AlertService{}
}

// Evaluate checks a probe result against alert rules.
// Stub — no-op until Task 9 implementation.
func (s *AlertService) Evaluate(agentID uint, result model.ProbeResult) {
	// TODO: implement alert evaluation in Task 9
}
