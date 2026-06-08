package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// timestampedResult wraps a ProbeResult with its creation time for age-based expiry.
type timestampedResult struct {
	result    ProbeResult
	createdAt time.Time
}

// ResultQueue is a thread-safe queue for probe results awaiting upload.
type ResultQueue struct {
	mu      sync.Mutex
	results []timestampedResult
}

// NewResultQueue creates an empty ResultQueue.
func NewResultQueue() *ResultQueue {
	return &ResultQueue{}
}

// Add enqueues a probe result.
func (q *ResultQueue) Add(r ProbeResult) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.results = append(q.results, timestampedResult{
		result:    r,
		createdAt: time.Now(),
	})
}

// Drain removes expired results and returns up to maxBatch non-expired results.
// The returned results are removed from the queue.
func (q *ResultQueue) Drain(maxBatch int, maxAge time.Duration) []ProbeResult {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()

	// Remove expired results first.
	var active []timestampedResult
	for _, tr := range q.results {
		if now.Sub(tr.createdAt) <= maxAge {
			active = append(active, tr)
		}
	}
	q.results = active

	// Take up to maxBatch.
	n := maxBatch
	if n > len(q.results) {
		n = len(q.results)
	}

	batch := make([]ProbeResult, n)
	for i := 0; i < n; i++ {
		batch[i] = q.results[i].result
	}

	// Remove taken results from queue.
	q.results = q.results[n:]

	return batch
}

// Prepend puts results back at the front of the queue (used on upload failure).
func (q *ResultQueue) Prepend(results []ProbeResult) {
	q.mu.Lock()
	defer q.mu.Unlock()

	prepended := make([]timestampedResult, len(results))
	for i, r := range results {
		prepended[i] = timestampedResult{
			result:    r,
			createdAt: time.Now(),
		}
	}
	q.results = append(prepended, q.results...)
}

// Len returns the current queue length.
func (q *ResultQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.results)
}

// ---------------------------------------------------------------------------
// Reporter
// ---------------------------------------------------------------------------

// reportRequest is the JSON body for POST /api/v1/agent/report.
type reportRequest struct {
	Results []ProbeResult `json:"results"`
}

// reportResponse is the server's response envelope.
type reportResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Reporter periodically uploads probe results to the server and sends heartbeats.
type Reporter struct {
	cfg     *Config
	token   string
	queue   *ResultQueue
}

// NewReporter creates a new Reporter.
func NewReporter(cfg *Config, token string, queue *ResultQueue) *Reporter {
	return &Reporter{
		cfg:   cfg,
		token: token,
		queue: queue,
	}
}

// Start runs the report and heartbeat loops. It blocks until ctx is cancelled,
// then performs a final flush.
func (r *Reporter) Start(ctx context.Context) {
	reportTicker := time.NewTicker(time.Duration(r.cfg.Report.Interval) * time.Second)
	defer reportTicker.Stop()

	heartbeatTicker := time.NewTicker(60 * time.Second)
	defer heartbeatTicker.Stop()

	// Send initial heartbeat.
	r.heartbeat()

	for {
		select {
		case <-reportTicker.C:
			r.flush()
		case <-heartbeatTicker.C:
			r.heartbeat()
		case <-ctx.Done():
			// Final flush before shutdown.
			r.flush()
			return
		}
	}
}

// flush drains the queue and uploads results to the server.
func (r *Reporter) flush() {
	maxAge := time.Duration(r.cfg.Report.RetryMaxAge) * time.Second
	batch := r.queue.Drain(r.cfg.Report.BatchSize, maxAge)

	if len(batch) == 0 {
		return
	}

	payload, err := json.Marshal(reportRequest{Results: batch})
	if err != nil {
		log.Printf("序列化上报数据失败: %v", err)
		r.queue.Prepend(batch)
		return
	}

	url := r.cfg.Server.URL + "/api/v1/agent/report"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		log.Printf("构造上报请求失败: %v", err)
		r.queue.Prepend(batch)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("上报请求失败: %v", err)
		r.queue.Prepend(batch)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取上报响应失败: %v", err)
		r.queue.Prepend(batch)
		return
	}

	var rr reportResponse
	if err := json.Unmarshal(body, &rr); err != nil {
		log.Printf("解析上报响应失败: %v", err)
		r.queue.Prepend(batch)
		return
	}

	if rr.Code != 0 {
		log.Printf("上报失败 (code=%d): %s", rr.Code, rr.Msg)
		r.queue.Prepend(batch)
		return
	}

	log.Printf("成功上报 %d 条探测结果", len(batch))
}

// heartbeat sends a heartbeat signal to the server.
func (r *Reporter) heartbeat() {
	url := r.cfg.Server.URL + "/api/v1/agent/heartbeat"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Printf("构造心跳请求失败: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+r.token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("心跳请求失败: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var rr reportResponse
	if err := json.Unmarshal(body, &rr); err != nil {
		return
	}

	if rr.Code != 0 {
		log.Printf("心跳上报失败 (code=%d): %s", rr.Code, rr.Msg)
		return
	}
}
