package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// ---------------------------------------------------------------------------
// GORM Models
// ---------------------------------------------------------------------------

// User represents an admin user.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// Agent represents a remote probe node.
type Agent struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"not null" json:"name"`
	Token         string     `gorm:"uniqueIndex;not null" json:"token"`
	Status        string     `gorm:"default:offline" json:"status"`
	LastHeartbeat *time.Time `json:"last_heartbeat"`
	Location      string     `json:"location"`
	CreatedAt     time.Time  `json:"created_at"`
}

// Monitor represents a website monitoring task.
type Monitor struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"not null" json:"name"`
	URL              string    `gorm:"not null" json:"url"`
	Method           string    `gorm:"default:GET" json:"method"`
	Headers          string    `json:"headers"`
	Cookie           string    `json:"cookie"`
	BasicAuthUser    string    `json:"basic_auth_user"`
	BasicAuthPass    string    `json:"basic_auth_pass"`
	VerifySSL        bool      `gorm:"default:true" json:"verify_ssl"`
	Frequency        int       `gorm:"default:60" json:"frequency"`
	Proxy            string    `json:"proxy"`
	AgentIDs         JSONSlice `gorm:"type:text" json:"agent_ids"`
	MatchType        string    `gorm:"default:none" json:"match_type"`
	MatchContent     string    `json:"match_content"`
	StatusThreshold  int       `gorm:"default:400" json:"status_threshold"`
	LatencyThreshold int       `gorm:"default:3000" json:"latency_threshold"`
	FailCount        int       `gorm:"default:3" json:"fail_count"`
	Enabled          bool      `gorm:"default:true" json:"enabled"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// AlertChannel represents a notification channel (DingTalk, WeChat Work, etc.).
type AlertChannel struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	Type       string    `gorm:"not null" json:"type"`
	WebhookURL string    `gorm:"not null" json:"webhook_url"`
	Secret     string    `json:"secret"`
	Extra      string    `json:"extra"`
	Enabled    bool      `gorm:"default:true" json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
}

// MonitorAlertChannel is the join table linking monitors to alert channels.
type MonitorAlertChannel struct {
	MonitorID      uint `gorm:"primaryKey" json:"monitor_id"`
	AlertChannelID uint `gorm:"primaryKey" json:"alert_channel_id"`
}

// AlertHistory records individual alert events.
type AlertHistory struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	MonitorID   uint       `gorm:"index" json:"monitor_id"`
	AgentID     uint       `json:"agent_id"`
	AlertType   string     `json:"alert_type"`
	Message     string     `json:"message"`
	Status      string     `json:"status"`
	TriggeredAt time.Time  `json:"triggered_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`

	// Associations (preloaded on demand)
	Monitor *Monitor `gorm:"foreignKey:MonitorID" json:"monitor,omitempty"`
	Agent   *Agent   `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
}

// ---------------------------------------------------------------------------
// JSONSlice — custom type that serialises []uint to/from a JSON TEXT column
// ---------------------------------------------------------------------------

// JSONSlice is a []uint that implements GORM's Scanner/Valuer for SQLite TEXT.
type JSONSlice []uint

// Value implements driver.Valuer.
func (j JSONSlice) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	b, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan implements sql.Scanner.
func (j *JSONSlice) Scan(src interface{}) error {
	if src == nil {
		*j = JSONSlice{}
		return nil
	}

	var bytes []byte
	switch v := src.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return errors.New("unsupported type for JSONSlice")
	}

	var result []uint
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*j = result
	return nil
}

// ---------------------------------------------------------------------------
// API DTOs
// ---------------------------------------------------------------------------

// LoginRequest is the payload for admin login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AgentRegisterRequest is the payload for agent registration.
type AgentRegisterRequest struct {
	Name         string `json:"name" binding:"required"`
	Location     string `json:"location"`
	SharedSecret string `json:"shared_secret" binding:"required"`
}

// AgentRegisterResponse is returned after successful agent registration.
type AgentRegisterResponse struct {
	AgentID uint   `json:"agent_id"`
	Token   string `json:"token"`
}

// ProbeResult represents a single probe result from an agent.
type ProbeResult struct {
	MonitorID      uint    `json:"monitor_id"`
	StatusCode     int     `json:"status_code"`
	DurationMs     int     `json:"duration_ms"`
	ContentMatched bool    `json:"content_matched"`
	SSLExpiry      float64 `json:"ssl_expiry"`
	Success        bool    `json:"success"`
	ErrorMsg       string  `json:"error_msg"`
	Timestamp      int64   `json:"timestamp"`
}

// ReportRequest is the payload for batch probe result reporting.
type ReportRequest struct {
	Results []ProbeResult `json:"results" binding:"required"`
}

// AgentTaskDTO is the task information sent to an agent.
type AgentTaskDTO struct {
	ID               uint    `json:"id"`
	Name             string  `json:"name"`
	URL              string  `json:"url"`
	Method           string  `json:"method"`
	Headers          string  `json:"headers"`
	Cookie           string  `json:"cookie"`
	BasicAuthUser    string  `json:"basic_auth_user"`
	BasicAuthPass    string  `json:"basic_auth_pass"`
	VerifySSL        bool    `json:"verify_ssl"`
	Frequency        int     `json:"frequency"`
	Proxy            string  `json:"proxy"`
	MatchType        string  `json:"match_type"`
	MatchContent     string  `json:"match_content"`
	StatusThreshold  int     `json:"status_threshold"`
	LatencyThreshold int     `json:"latency_threshold"`
	FailCount        int     `json:"fail_count"`
}

// PageQuery is a common pagination parameter.
type PageQuery struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Keyword  string `form:"keyword" json:"keyword"`
}

// DashboardOverview holds dashboard summary statistics.
type DashboardOverview struct {
	TotalMonitors   int     `json:"total_monitors"`
	UptimeRate      float64 `json:"uptime_rate"`
	ActiveAlerts    int     `json:"active_alerts"`
	AvgResponseTime float64 `json:"avg_response_time"`
}
