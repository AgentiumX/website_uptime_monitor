package service

import (
	"uptime-monitor/server/internal/model"
)

// NotifyService handles sending notifications to alert channels.
// Stub implementation — full logic will be added in Task 10.
type NotifyService struct{}

// NewNotifyService creates a new NotifyService.
func NewNotifyService() *NotifyService {
	return &NotifyService{}
}

// SendTest sends a test notification to the given alert channel.
// Stub — returns nil until Task 10 implementation.
func (s *NotifyService) SendTest(ch *model.AlertChannel) error {
	// TODO: implement notification sending in Task 10
	return nil
}
