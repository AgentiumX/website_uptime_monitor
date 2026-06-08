package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"

	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
)

// AlertService handles alert evaluation logic.
type AlertService struct {
	failCounts sync.Map // key: "monitorID_agentID" -> *int
	notifySvc  *NotifyService
}

// NewAlertService creates a new AlertService.
func NewAlertService(notifySvc *NotifyService) *AlertService {
	return &AlertService{notifySvc: notifySvc}
}

// failKey returns the map key for a monitor+agent pair.
func failKey(monitorID, agentID uint) string {
	return fmt.Sprintf("%d_%d", monitorID, agentID)
}

// Evaluate checks a probe result against alert rules.
func (s *AlertService) Evaluate(agentID uint, result model.ProbeResult) {
	monitor, err := repository.GetMonitorByID(result.MonitorID)
	if err != nil {
		log.Printf("[AlertService] get monitor error: %v", err)
		return
	}
	if !monitor.Enabled {
		return
	}

	var failed bool
	var alertType, detail string

	// Check status code threshold
	if result.StatusCode >= monitor.StatusThreshold {
		failed = true
		alertType = "status_code"
		detail = fmt.Sprintf("状态码 %d >= 阈值 %d", result.StatusCode, monitor.StatusThreshold)
	}

	// Check latency threshold
	if result.DurationMs >= monitor.LatencyThreshold {
		failed = true
		if alertType == "" {
			alertType = "latency"
		}
		detail = fmt.Sprintf("响应时间 %dms >= 阈值 %dms", result.DurationMs, monitor.LatencyThreshold)
	}

	// Check content matching
	if monitor.MatchType == "contains" && !result.ContentMatched {
		failed = true
		if alertType == "" {
			alertType = "content_match"
		}
		detail = fmt.Sprintf("内容匹配失败：期望包含 [%s]", monitor.MatchContent)
	} else if monitor.MatchType == "not_contains" && result.ContentMatched {
		failed = true
		if alertType == "" {
			alertType = "content_match"
		}
		detail = fmt.Sprintf("内容匹配失败：期望不包含 [%s]", monitor.MatchContent)
	}

	key := failKey(monitor.ID, agentID)

	if failed {
		// Increment failure count
		val, _ := s.failCounts.LoadOrStore(key, new(int))
		countPtr := val.(*int)
		*countPtr++
		count := *countPtr

		// Check if threshold reached
		if count >= monitor.FailCount {
			// Check if already alerting
			activeAlert, err := repository.FindActiveAlert(monitor.ID, agentID)
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Printf("[AlertService] find active alert error: %v", err)
			}
			if activeAlert != nil {
				// Already alerting, skip
				return
			}

			// Create firing alert
			alert := &model.AlertHistory{
				MonitorID:   monitor.ID,
				AgentID:     agentID,
				AlertType:   alertType,
				Message:     detail,
				Status:      "firing",
				TriggeredAt: time.Now(),
			}
			if err := repository.CreateAlertHistory(alert); err != nil {
				log.Printf("[AlertService] create alert error: %v", err)
				return
			}

			// Send notification asynchronously
			go s.notifySvc.NotifyFiring(monitor, agentID, alertType, detail, count)
		}
	} else {
		// Success: resolve any active alert and reset counter
		val, loaded := s.failCounts.Load(key)
		if loaded {
			countPtr := val.(*int)
			if *countPtr > 0 {
				// Check for active alert to resolve
				activeAlert, err := repository.FindActiveAlert(monitor.ID, agentID)
				if err == nil && activeAlert != nil {
					if err := repository.ResolveAlert(activeAlert.ID); err != nil {
						log.Printf("[AlertService] resolve alert error: %v", err)
					} else {
						// Reload to get updated resolved_at
						resolved, err := repository.FindActiveAlert(monitor.ID, agentID)
						if err != nil {
							// Use the original alert with manually set resolved_at for notification
							now := time.Now()
							activeAlert.ResolvedAt = &now
							activeAlert.Status = "resolved"
							go s.notifySvc.NotifyResolved(monitor, activeAlert)
						} else {
							_ = resolved // should be nil after resolve
							now := time.Now()
							activeAlert.ResolvedAt = &now
							activeAlert.Status = "resolved"
							go s.notifySvc.NotifyResolved(monitor, activeAlert)
						}
					}
				}
			}
			*countPtr = 0
		}
	}
}
