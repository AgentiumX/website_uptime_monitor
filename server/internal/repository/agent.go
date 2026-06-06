package repository

import (
	"time"

	"uptime-monitor/server/internal/model"
)

// CreateAgent inserts a new agent record.
func CreateAgent(a *model.Agent) error {
	return DB.Create(a).Error
}

// GetAgentByID returns an agent by primary key.
func GetAgentByID(id uint) (*model.Agent, error) {
	var a model.Agent
	err := DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// DeleteAgent removes an agent by ID.
func DeleteAgent(id uint) error {
	return DB.Delete(&model.Agent{}, id).Error
}

// ListAgents returns all agents ordered by creation time descending.
func ListAgents() ([]model.Agent, error) {
	var agents []model.Agent
	err := DB.Order("id DESC").Find(&agents).Error
	return agents, err
}

// FindAgentByToken returns the agent with the given token.
func FindAgentByToken(token string) (*model.Agent, error) {
	var a model.Agent
	err := DB.Where("token = ?", token).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// UpdateHeartbeat sets last_heartbeat to now and status to online.
func UpdateHeartbeat(agentID uint) error {
	now := time.Now()
	return DB.Model(&model.Agent{}).Where("id = ?", agentID).
		Updates(map[string]interface{}{
			"last_heartbeat": now,
			"status":         "online",
		}).Error
}

// MarkOfflineAgents sets agents whose last heartbeat is older than timeout to offline.
func MarkOfflineAgents(timeout time.Duration) error {
	cutoff := time.Now().Add(-timeout)
	return DB.Model(&model.Agent{}).
		Where("status = ? AND last_heartbeat < ?", "online", cutoff).
		Update("status", "offline").Error
}
