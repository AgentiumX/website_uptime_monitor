package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

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

// GetTasks returns the task list for a given agent.
func (s *AgentService) GetTasks(agentID uint) ([]model.AgentTaskDTO, error) {
	monitors, err := repository.GetMonitorsByAgent(agentID)
	if err != nil {
		return nil, err
	}

	tasks := make([]model.AgentTaskDTO, len(monitors))
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
		}
	}
	return tasks, nil
}

// Heartbeat updates the agent's heartbeat timestamp.
func (s *AgentService) Heartbeat(agentID uint) error {
	return repository.UpdateHeartbeat(agentID)
}
