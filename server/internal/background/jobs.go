package background

import (
	"context"
	"log"
	"time"

	"uptime-monitor/server/internal/config"
	"uptime-monitor/server/internal/repository"
)

// StartJobs launches background goroutines for periodic tasks.
// The goroutines run until ctx is cancelled.
func StartJobs(ctx context.Context, cfg *config.Config) {
	// Agent heartbeat detection: every 30s mark agents with stale heartbeat as offline.
	go func() {
		timeout := time.Duration(cfg.Agent.HeartbeatTimeout) * time.Second
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := repository.MarkOfflineAgents(timeout); err != nil {
					log.Printf("[background] MarkOfflineAgents error: %v", err)
				}
			}
		}
	}()

	// Alert history cleanup: every 24h delete old resolved alerts.
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := repository.CleanupAlertHistory(cfg.Alert.HistoryRetention); err != nil {
					log.Printf("[background] CleanupAlertHistory error: %v", err)
				}
			}
		}
	}()
}
