package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// Task mirrors the server's AgentTaskDTO — the monitoring task configuration
// that the agent receives from the server.
type Task struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	Method           string `json:"method"`
	Headers          string `json:"headers"`
	Cookie           string `json:"cookie"`
	BasicAuthUser    string `json:"basic_auth_user"`
	BasicAuthPass    string `json:"basic_auth_pass"`
	VerifySSL        bool   `json:"verify_ssl"`
	Frequency        int    `json:"frequency"`
	Proxy            string `json:"proxy"`
	MatchType        string `json:"match_type"`
	MatchContent     string `json:"match_content"`
	StatusThreshold  int    `json:"status_threshold"`
	LatencyThreshold int    `json:"latency_threshold"`
	FailCount        int    `json:"fail_count"`
	UpdatedAt        string `json:"updated_at"`
}

// taskRunner pairs a task with its cancellation function.
type taskRunner struct {
	task   Task
	cancel context.CancelFunc
}

// Scheduler manages task synchronisation and probe goroutines.
type Scheduler struct {
	cfg      *Config
	token    string
	results  *ResultQueue
	probeSem chan struct{}

	mu    sync.RWMutex
	tasks map[uint]*taskRunner
}

// NewScheduler creates a new Scheduler.
func NewScheduler(cfg *Config, token string, results *ResultQueue) *Scheduler {
	return &Scheduler{
		cfg:      cfg,
		token:    token,
		results:  results,
		probeSem: make(chan struct{}, cfg.Probe.MaxConcurrent),
		tasks:    make(map[uint]*taskRunner),
	}
}

// Start begins the task sync loop. It blocks until ctx is cancelled.
func (s *Scheduler) Start(ctx context.Context) {
	// Initial sync immediately.
	s.syncTasks(ctx)

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.syncTasks(ctx)
		case <-ctx.Done():
			s.stopAll()
			return
		}
	}
}

// tasksResponse mirrors the server's GET /api/v1/agent/tasks envelope.
type tasksResponse struct {
	Code int `json:"code"`
	Data struct {
		Tasks     []Task `json:"tasks"`
		UpdatedAt string `json:"updated_at"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// fetchTasks retrieves the current task list from the server.
func (s *Scheduler) fetchTasks() ([]Task, error) {
	req, err := http.NewRequest(http.MethodGet, s.cfg.Server.URL+"/api/v1/agent/tasks", nil)
	if err != nil {
		return nil, fmt.Errorf("构造请求失败: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var tr tasksResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if tr.Code != 0 {
		return nil, fmt.Errorf("获取任务失败 (code=%d): %s", tr.Code, tr.Msg)
	}

	return tr.Data.Tasks, nil
}

// syncTasks fetches the latest task list and reconciles with running tasks.
func (s *Scheduler) syncTasks(ctx context.Context) {
	newTasks, err := s.fetchTasks()
	if err != nil {
		log.Printf("同步任务失败: %v", err)
		return
	}

	// Build a lookup map for new tasks.
	newMap := make(map[uint]Task, len(newTasks))
	for _, t := range newTasks {
		newMap[t.ID] = t
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove tasks that are no longer in the list.
	for id, runner := range s.tasks {
		if _, exists := newMap[id]; !exists {
			log.Printf("停止已删除的任务: id=%d name=%s", id, runner.task.Name)
			runner.cancel()
			delete(s.tasks, id)
		}
	}

	// Add new tasks or restart changed tasks.
	for id, newTask := range newMap {
		existing, exists := s.tasks[id]
		if !exists {
			// New task — start it.
			s.startTask(ctx, newTask)
			log.Printf("启动新任务: id=%d name=%s", id, newTask.Name)
			continue
		}

		// Existing task — check if UpdatedAt changed.
		if existing.task.UpdatedAt != newTask.UpdatedAt {
			log.Printf("重启变更的任务: id=%d name=%s", id, newTask.Name)
			existing.cancel()
			delete(s.tasks, id)
			s.startTask(ctx, newTask)
		}
	}

	log.Printf("任务同步完成: 共 %d 个任务", len(s.tasks))
}

// startTask launches a probe goroutine for the given task.
// Must be called with s.mu held.
func (s *Scheduler) startTask(parentCtx context.Context, task Task) {
	ctx, cancel := context.WithCancel(parentCtx)
	s.tasks[task.ID] = &taskRunner{task: task, cancel: cancel}
	go s.runProbe(ctx, task)
}

// runProbe executes doProbe immediately, then on a ticker matching the task's
// frequency. It exits when ctx is cancelled.
func (s *Scheduler) runProbe(ctx context.Context, task Task) {
	// Execute once immediately.
	s.doProbe(task)

	freq := task.Frequency
	if freq <= 0 {
		freq = 60
	}
	ticker := time.NewTicker(time.Duration(freq) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.doProbe(task)
		case <-ctx.Done():
			return
		}
	}
}

// doProbe acquires the probe semaphore, executes the probe, and enqueues the result.
func (s *Scheduler) doProbe(task Task) {
	s.probeSem <- struct{}{}
	defer func() { <-s.probeSem }()

	result := ExecuteProbe(task, s.cfg.Probe.Timeout)
	s.results.Add(result)
}

// stopAll cancels all running task goroutines.
func (s *Scheduler) stopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, runner := range s.tasks {
		runner.cancel()
		delete(s.tasks, id)
	}
	log.Println("所有探测任务已停止")
}
