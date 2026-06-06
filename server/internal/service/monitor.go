package service

import (
	"errors"

	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
)

// MonitorService handles monitor-related business logic.
type MonitorService struct{}

// NewMonitorService creates a new MonitorService.
func NewMonitorService() *MonitorService {
	return &MonitorService{}
}

// List returns a paginated list of monitors, optionally filtered by keyword.
func (s *MonitorService) List(page, pageSize int, keyword string) ([]model.Monitor, int64, error) {
	return repository.ListMonitors(page, pageSize, keyword)
}

// Get returns a single monitor by ID.
func (s *MonitorService) Get(id uint) (*model.Monitor, error) {
	return repository.GetMonitorByID(id)
}

// Create validates and inserts a new monitor.
func (s *MonitorService) Create(m *model.Monitor) error {
	if m.URL == "" {
		return errors.New("监控地址不能为空")
	}
	if m.Frequency < 30 {
		return errors.New("检测频率不能低于30秒")
	}
	if m.FailCount < 1 || m.FailCount > 5 {
		return errors.New("连续失败次数必须在1-5之间")
	}
	return repository.CreateMonitor(m)
}

// Update saves changes to an existing monitor.
func (s *MonitorService) Update(m *model.Monitor) error {
	return repository.UpdateMonitor(m)
}

// Delete removes a monitor by ID.
func (s *MonitorService) Delete(id uint) error {
	return repository.DeleteMonitor(id)
}

// Toggle enables or disables a monitor.
func (s *MonitorService) Toggle(id uint, enabled bool) error {
	return repository.ToggleMonitor(id, enabled)
}
