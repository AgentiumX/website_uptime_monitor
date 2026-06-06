package repository

import (
	"time"

	"gorm.io/gorm"

	"uptime-monitor/server/internal/model"
)

// ---------------------------------------------------------------------------
// AlertChannel CRUD
// ---------------------------------------------------------------------------

// CreateAlertChannel inserts a new alert channel.
func CreateAlertChannel(ch *model.AlertChannel) error {
	return DB.Create(ch).Error
}

// GetAlertChannelByID returns an alert channel by primary key.
func GetAlertChannelByID(id uint) (*model.AlertChannel, error) {
	var ch model.AlertChannel
	err := DB.First(&ch, id).Error
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

// UpdateAlertChannel saves changes to an existing alert channel.
func UpdateAlertChannel(ch *model.AlertChannel) error {
	return DB.Save(ch).Error
}

// DeleteAlertChannel removes an alert channel and its monitor associations.
func DeleteAlertChannel(id uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("alert_channel_id = ?", id).Delete(&model.MonitorAlertChannel{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.AlertChannel{}, id).Error
	})
}

// ListAlertChannels returns all alert channels.
func ListAlertChannels() ([]model.AlertChannel, error) {
	var channels []model.AlertChannel
	err := DB.Order("id DESC").Find(&channels).Error
	return channels, err
}

// GetAlertChannelsByMonitor returns alert channels associated with a monitor via JOIN.
func GetAlertChannelsByMonitor(monitorID uint) ([]model.AlertChannel, error) {
	var channels []model.AlertChannel
	err := DB.Joins("JOIN monitor_alert_channels ON monitor_alert_channels.alert_channel_id = alert_channels.id").
		Where("monitor_alert_channels.monitor_id = ?", monitorID).
		Find(&channels).Error
	return channels, err
}

// SetMonitorAlertChannels replaces the alert channel associations for a monitor in a transaction.
func SetMonitorAlertChannels(monitorID uint, channelIDs []uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("monitor_id = ?", monitorID).Delete(&model.MonitorAlertChannel{}).Error; err != nil {
			return err
		}
		for _, cid := range channelIDs {
			mac := model.MonitorAlertChannel{
				MonitorID:      monitorID,
				AlertChannelID: cid,
			}
			if err := tx.Create(&mac).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ---------------------------------------------------------------------------
// AlertHistory CRUD
// ---------------------------------------------------------------------------

// CreateAlertHistory inserts a new alert history record.
func CreateAlertHistory(h *model.AlertHistory) error {
	return DB.Create(h).Error
}

// FindActiveAlert returns the currently firing alert for a given monitor+agent pair.
func FindActiveAlert(monitorID, agentID uint) (*model.AlertHistory, error) {
	var h model.AlertHistory
	err := DB.Where("monitor_id = ? AND agent_id = ? AND status = ?", monitorID, agentID, "firing").
		First(&h).Error
	if err != nil {
		return nil, err
	}
	return &h, nil
}

// ResolveAlert marks a firing alert as resolved with the current timestamp.
func ResolveAlert(id uint) error {
	now := time.Now()
	return DB.Model(&model.AlertHistory{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      "resolved",
			"resolved_at": now,
		}).Error
}

// ListAlertHistory returns a paginated list of alert history records with preloaded associations.
func ListAlertHistory(page, pageSize int, monitorID uint, status string) ([]model.AlertHistory, int64, error) {
	var records []model.AlertHistory
	var total int64

	q := DB.Model(&model.AlertHistory{})
	if monitorID > 0 {
		q = q.Where("monitor_id = ?", monitorID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := q.Preload("Monitor").Preload("Agent").
		Order("id DESC").
		Offset(offset).Limit(pageSize).
		Find(&records).Error
	return records, total, err
}

// CountActiveAlerts returns the number of currently firing alerts.
func CountActiveAlerts() (int64, error) {
	var count int64
	err := DB.Model(&model.AlertHistory{}).Where("status = ?", "firing").Count(&count).Error
	return count, err
}

// CleanupAlertHistory deletes resolved alerts older than the given retention period.
func CleanupAlertHistory(retentionDays int) error {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	return DB.Where("status = ? AND resolved_at < ?", "resolved", cutoff).
		Delete(&model.AlertHistory{}).Error
}
