package repository

import (
	"uptime-monitor/server/internal/model"

	"gorm.io/gorm"
)

// CreateMonitor inserts a new monitor record.
func CreateMonitor(m *model.Monitor) error {
	return DB.Create(m).Error
}

// GetMonitorByID returns a monitor by primary key.
func GetMonitorByID(id uint) (*model.Monitor, error) {
	var m model.Monitor
	err := DB.First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// UpdateMonitor saves changes to an existing monitor.
func UpdateMonitor(m *model.Monitor) error {
	return DB.Save(m).Error
}

// DeleteMonitor removes a monitor and its alert-channel associations.
func DeleteMonitor(id uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("monitor_id = ?", id).Delete(&model.MonitorAlertChannel{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Monitor{}, id).Error
	})
}

// ListMonitors returns a paginated, optionally keyword-filtered list of monitors.
func ListMonitors(page, pageSize int, keyword string) ([]model.Monitor, int64, error) {
	var monitors []model.Monitor
	var total int64

	q := DB.Model(&model.Monitor{})
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR url LIKE ?", like, like)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := q.Order("id DESC").Offset(offset).Limit(pageSize).Find(&monitors).Error; err != nil {
		return nil, 0, err
	}
	return monitors, total, nil
}

// GetMonitorsByAgent returns all enabled monitors whose agent_ids JSON array
// contains the given agent ID. Uses SQLite JSON1 LIKE matching.
func GetMonitorsByAgent(agentID uint) ([]model.Monitor, error) {
	var monitors []model.Monitor
	// agent_ids is stored as e.g. "[1,2,3]", so we match substrings.
	// A precise check would require JSON1 extension; LIKE is good enough for small ID spaces.
	err := DB.Where("enabled = ? AND (agent_ids LIKE ? OR agent_ids LIKE ? OR agent_ids LIKE ?)",
		true,
		`%[`+itoa(agentID)+`,%`,
		`%,`+itoa(agentID)+`,%`,
		`%,`+itoa(agentID)+`]%`,
	).Find(&monitors).Error
	return monitors, err
}

// ToggleMonitor flips the enabled flag of a monitor.
func ToggleMonitor(id uint, enabled bool) error {
	return DB.Model(&model.Monitor{}).Where("id = ?", id).Update("enabled", enabled).Error
}

// CountMonitors returns the total number of monitors.
func CountMonitors() (int64, error) {
	var count int64
	err := DB.Model(&model.Monitor{}).Count(&count).Error
	return count, err
}

// itoa converts a uint to its decimal string representation without importing strconv.
func itoa(n uint) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
