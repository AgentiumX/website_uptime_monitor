package service

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
)

// MetricsService handles recording probe results to Prometheus gauges.
type MetricsService struct {
	probeSuccess    *prometheus.GaugeVec
	probeStatusCode *prometheus.GaugeVec
	probeDuration   *prometheus.GaugeVec
	probeSSLExpiry  *prometheus.GaugeVec
	probeMatched    *prometheus.GaugeVec
	registry        *prometheus.Registry
}

// NewMetricsService creates a new MetricsService and registers all gauges.
func NewMetricsService() *MetricsService {
	registry := prometheus.NewRegistry()

	probeSuccess := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "probe_success",
		Help: "Whether the probe was successful (1=success, 0=failure)",
	}, []string{"monitor_id", "agent_id", "url", "method"})

	probeStatusCode := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "probe_http_status_code",
		Help: "HTTP status code returned by the probe",
	}, []string{"monitor_id", "agent_id", "url"})

	probeDuration := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "probe_duration_seconds",
		Help: "Duration of the probe in seconds",
	}, []string{"monitor_id", "agent_id", "url"})

	probeSSLExpiry := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "probe_ssl_expiry_seconds",
		Help: "SSL certificate expiry time in seconds",
	}, []string{"monitor_id", "agent_id", "url"})

	probeMatched := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "probe_content_matched",
		Help: "Whether the content matched (1=matched, 0=not matched)",
	}, []string{"monitor_id", "agent_id", "url"})

	registry.MustRegister(probeSuccess, probeStatusCode, probeDuration, probeSSLExpiry, probeMatched)

	return &MetricsService{
		probeSuccess:    probeSuccess,
		probeStatusCode: probeStatusCode,
		probeDuration:   probeDuration,
		probeMatched:    probeMatched,
		probeSSLExpiry:  probeSSLExpiry,
		registry:        registry,
	}
}

// Record writes a probe result to the metrics store.
func (s *MetricsService) Record(result model.ProbeResult) {
	monitor, err := repository.GetMonitorByID(result.MonitorID)
	if err != nil {
		log.Printf("[MetricsService] get monitor error: %v", err)
		return
	}

	monitorID := strconv.FormatUint(uint64(result.MonitorID), 10)
	agentID := "0" // Agent ID not in result, use placeholder
	url := monitor.URL
	method := monitor.Method

	successVal := float64(0)
	if result.Success {
		successVal = 1
	}

	matchedVal := float64(0)
	if result.ContentMatched {
		matchedVal = 1
	}

	s.probeSuccess.WithLabelValues(monitorID, agentID, url, method).Set(successVal)
	s.probeStatusCode.WithLabelValues(monitorID, agentID, url).Set(float64(result.StatusCode))
	s.probeDuration.WithLabelValues(monitorID, agentID, url).Set(float64(result.DurationMs) / 1000.0)
	s.probeSSLExpiry.WithLabelValues(monitorID, agentID, url).Set(result.SSLExpiry)
	s.probeMatched.WithLabelValues(monitorID, agentID, url).Set(matchedVal)
}

// GetRegistry returns the Prometheus registry for the /metrics endpoint.
func (s *MetricsService) GetRegistry() *prometheus.Registry {
	return s.registry
}

// GetOverallStats returns average uptime rate and average latency from current gauge values.
func (s *MetricsService) GetOverallStats() (avgUptime float64, avgLatencyMs float64) {
	var successSum, successCount float64
	var durationSum, durationCount float64

	metrics, err := s.registry.Gather()
	if err != nil {
		log.Printf("[MetricsService] gather error: %v", err)
		return 0, 0
	}

	for _, mf := range metrics {
		name := mf.GetName()
		for _, m := range mf.GetMetric() {
			if m.GetGauge() == nil {
				continue
			}
			val := m.GetGauge().GetValue()
			if name == "probe_success" {
				successSum += val
				successCount++
			} else if name == "probe_duration_seconds" {
				durationSum += val
				durationCount++
			}
		}
	}

	if successCount > 0 {
		avgUptime = (successSum / successCount) * 100
	}
	if durationCount > 0 {
		avgLatencyMs = (durationSum / durationCount) * 1000
	}
	return avgUptime, avgLatencyMs
}

// GetMonitorMetrics returns current metrics snapshot for a specific monitor.
func (s *MetricsService) GetMonitorMetrics(monitorID uint) map[string]float64 {
	result := map[string]float64{
		"status_code": 0,
		"duration_s":  0,
		"ssl_expiry":  0,
		"success":     0,
		"matched":     0,
	}

	mid := strconv.FormatUint(uint64(monitorID), 10)
	metrics, err := s.registry.Gather()
	if err != nil {
		log.Printf("[MetricsService] gather error: %v", err)
		return result
	}

	for _, mf := range metrics {
		name := mf.GetName()
		for _, m := range mf.GetMetric() {
			// Check if this metric is for our monitor
			matchMonitor := false
			for _, label := range m.GetLabel() {
				if label.GetName() == "monitor_id" && label.GetValue() == mid {
					matchMonitor = true
					break
				}
			}
			if !matchMonitor || m.GetGauge() == nil {
				continue
			}

			val := m.GetGauge().GetValue()
			switch name {
			case "probe_success":
				result["success"] = val
			case "probe_http_status_code":
				result["status_code"] = val
			case "probe_duration_seconds":
				result["duration_s"] = val
			case "probe_ssl_expiry_seconds":
				result["ssl_expiry"] = val
			case "probe_content_matched":
				result["matched"] = val
			}
		}
	}
	return result
}
