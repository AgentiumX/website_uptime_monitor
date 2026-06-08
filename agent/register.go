package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// registerRequest is the JSON body sent to the server for registration.
type registerRequest struct {
	Name         string `json:"name"`
	Location     string `json:"location"`
	SharedSecret string `json:"shared_secret"`
}

// registerResponse mirrors the server's response envelope.
type registerResponse struct {
	Code int `json:"code"`
	Data struct {
		AgentID uint   `json:"agent_id"`
		Token   string `json:"token"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// Register obtains an agent token — either from local cache or by registering
// with the server using the pre-shared secret.
func Register(cfg *Config) (string, error) {
	// Try to load a previously saved token.
	if token, err := loadToken(); err == nil && token != "" {
		log.Println("已加载保存的 Agent token")
		return token, nil
	}

	// No cached token — register with server.
	body, err := json.Marshal(registerRequest{
		Name:         cfg.Agent.Name,
		Location:     cfg.Agent.Location,
		SharedSecret: cfg.Server.SharedSecret,
	})
	if err != nil {
		return "", fmt.Errorf("序列化注册请求失败: %w", err)
	}

	url := cfg.Server.URL + "/api/v1/agent/register"
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("注册请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取注册响应失败: %w", err)
	}

	var regResp registerResponse
	if err := json.Unmarshal(respBody, &regResp); err != nil {
		return "", fmt.Errorf("解析注册响应失败: %w", err)
	}

	if regResp.Code != 0 {
		return "", fmt.Errorf("注册失败 (code=%d): %s", regResp.Code, regResp.Msg)
	}

	if regResp.Data.Token == "" {
		return "", fmt.Errorf("注册响应中 token 为空")
	}

	if err := saveToken(regResp.Data.Token); err != nil {
		log.Printf("警告: 保存 token 失败: %v", err)
	}

	log.Printf("Agent 注册成功, agent_id=%d", regResp.Data.AgentID)
	return regResp.Data.Token, nil
}

// tokenFilePath returns the path to the .agent_token file next to the executable.
func tokenFilePath() string {
	exe, err := os.Executable()
	if err != nil {
		// Fallback to current directory.
		return ".agent_token"
	}
	return filepath.Join(filepath.Dir(exe), ".agent_token")
}

// loadToken reads the persisted agent token from disk.
func loadToken() (string, error) {
	data, err := os.ReadFile(tokenFilePath())
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(data)), nil
}

// saveToken persists the agent token to disk with restricted permissions.
func saveToken(token string) error {
	return os.WriteFile(tokenFilePath(), []byte(token), 0600)
}
