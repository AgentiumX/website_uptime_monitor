package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ProbeResult holds the outcome of a single HTTP probe.
type ProbeResult struct {
	MonitorID      uint   `json:"monitor_id"`
	StatusCode     int    `json:"status_code"`
	DurationMs     int64  `json:"duration_ms"`
	ContentMatched bool   `json:"content_matched"`
	SSLExpiry      int64  `json:"ssl_expiry"`
	Success        bool   `json:"success"`
	ErrorMsg       string `json:"error_msg"`
	Timestamp      int64  `json:"timestamp"`
}

// ExecuteProbe performs a single HTTP probe for the given task and returns the result.
func ExecuteProbe(task Task, timeoutSec int) ProbeResult {
	result := ProbeResult{
		MonitorID: task.ID,
		Timestamp: time.Now().Unix(),
	}

	// Build the request.
	var body io.Reader
	method := strings.ToUpper(task.Method)
	if method == "POST" || method == "PUT" || method == "PATCH" {
		body = strings.NewReader("")
	}

	req, err := http.NewRequest(method, task.URL, body)
	if err != nil {
		result.ErrorMsg = fmt.Sprintf("构造请求失败: %v", err)
		return result
	}

	// Set custom headers from JSON.
	if task.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(task.Headers), &headers); err == nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
	}

	// Set Cookie.
	if task.Cookie != "" {
		req.Header.Set("Cookie", task.Cookie)
	}

	// Set Basic Auth.
	if task.BasicAuthUser != "" {
		req.SetBasicAuth(task.BasicAuthUser, task.BasicAuthPass)
	}

	// Build transport with TLS and proxy settings.
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !task.VerifySSL,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	if task.Proxy != "" {
		proxyURL, err := url.Parse(task.Proxy)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	// Build client — do not follow redirects.
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeoutSec) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Execute the request.
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)
	result.DurationMs = duration.Milliseconds()

	if err != nil {
		result.ErrorMsg = fmt.Sprintf("请求失败: %v", err)
		result.Success = false
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	// Check SSL certificate expiry.
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		result.SSLExpiry = int64(time.Until(cert.NotAfter).Seconds())
	}

	// Read response body (limit to 1 MB).
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		result.ErrorMsg = fmt.Sprintf("读取响应体失败: %v", err)
		result.Success = false
		return result
	}

	// Content matching.
	bodyStr := string(respBody)
	switch task.MatchType {
	case "contains":
		result.ContentMatched = strings.Contains(bodyStr, task.MatchContent)
	case "not_contains":
		result.ContentMatched = !strings.Contains(bodyStr, task.MatchContent)
	default: // "none" or empty
		result.ContentMatched = true
	}

	// Determine overall success.
	threshold := task.StatusThreshold
	if threshold <= 0 {
		threshold = 400
	}
	result.Success = resp.StatusCode < threshold

	if !result.ContentMatched {
		result.Success = false
	}

	return result
}
