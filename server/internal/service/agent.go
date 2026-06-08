package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"uptime-monitor/server/internal/config"
	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
)

// AgentService handles agent-related business logic.
type AgentService struct {
	cfg *config.AgentConfig
}

// NewAgentService creates a new AgentService.
func NewAgentService(cfg *config.AgentConfig) *AgentService {
	return &AgentService{cfg: cfg}
}

// Register validates the shared secret and creates a new agent with a random token.
func (s *AgentService) Register(req model.AgentRegisterRequest) (*model.AgentRegisterResponse, error) {
	if req.SharedSecret != s.cfg.SharedSecret {
		return nil, errors.New("密钥验证失败")
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, errors.New("生成令牌失败")
	}
	token := hex.EncodeToString(tokenBytes)

	agent := &model.Agent{
		Name:     req.Name,
		Location: req.Location,
		Token:    token,
		Status:   "offline",
	}
	if err := repository.CreateAgent(agent); err != nil {
		return nil, err
	}

	return &model.AgentRegisterResponse{
		AgentID: agent.ID,
		Token:   token,
	}, nil
}

// GetTasks returns the task list and the latest UpdatedAt for a given agent.
func (s *AgentService) GetTasks(agentID uint) ([]model.AgentTaskDTO, time.Time, error) {
	monitors, err := repository.GetMonitorsByAgent(agentID)
	if err != nil {
		return nil, time.Time{}, err
	}

	tasks := make([]model.AgentTaskDTO, len(monitors))
	var maxUpdated time.Time
	for i, m := range monitors {
		tasks[i] = model.AgentTaskDTO{
			ID:               m.ID,
			Name:             m.Name,
			URL:              m.URL,
			Method:           m.Method,
			Headers:          m.Headers,
			Cookie:           m.Cookie,
			BasicAuthUser:    m.BasicAuthUser,
			BasicAuthPass:    m.BasicAuthPass,
			VerifySSL:        m.VerifySSL,
			Frequency:        m.Frequency,
			Proxy:            m.Proxy,
			MatchType:        m.MatchType,
			MatchContent:     m.MatchContent,
			StatusThreshold:  m.StatusThreshold,
			LatencyThreshold: m.LatencyThreshold,
			FailCount:        m.FailCount,
			UpdatedAt:        m.UpdatedAt.Format(time.RFC3339),
		}
		if m.UpdatedAt.After(maxUpdated) {
			maxUpdated = m.UpdatedAt
		}
	}
	return tasks, maxUpdated, nil
}

// Heartbeat updates the agent's heartbeat timestamp.
func (s *AgentService) Heartbeat(agentID uint) error {
	return repository.UpdateHeartbeat(agentID)
}
