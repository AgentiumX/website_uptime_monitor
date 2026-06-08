package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
)

// NotifyService handles sending notifications to alert channels.
type NotifyService struct{}

// NewNotifyService creates a new NotifyService.
func NewNotifyService() *NotifyService {
	return &NotifyService{}
}

// NotifyFiring sends alert notifications to all channels associated with a monitor.
func (s *NotifyService) NotifyFiring(monitor *model.Monitor, agentID uint, alertType, detail string, count int) {
	channels, err := repository.GetAlertChannelsByMonitor(monitor.ID)
	if err != nil {
		log.Printf("[NotifyFiring] get channels error: %v", err)
		return
	}

	agent, err := repository.GetAgentByID(agentID)
	if err != nil {
		log.Printf("[NotifyFiring] get agent error: %v", err)
		return
	}

	title := fmt.Sprintf("🔴 告警：%s 异常", monitor.Name)
	text := fmt.Sprintf("**探测点**：%s (%s)\n"+
		"**监控地址**：%s\n"+
		"**告警类型**：%s\n"+
		"**详情**：%s\n"+
		"**连续失败**：%d 次\n"+
		"**时间**：%s",
		agent.Name, agent.Location,
		monitor.URL,
		alertType,
		detail,
		count,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	for _, ch := range channels {
		if !ch.Enabled {
			continue
		}
		if err := send(&ch, title, text); err != nil {
			log.Printf("[NotifyFiring] send to channel %s error: %v", ch.Name, err)
		}
	}
}

// NotifyResolved sends recovery notifications to all channels associated with a monitor.
func (s *NotifyService) NotifyResolved(monitor *model.Monitor, alert *model.AlertHistory) {
	channels, err := repository.GetAlertChannelsByMonitor(monitor.ID)
	if err != nil {
		log.Printf("[NotifyResolved] get channels error: %v", err)
		return
	}

	agent, err := repository.GetAgentByID(alert.AgentID)
	if err != nil {
		log.Printf("[NotifyResolved] get agent error: %v", err)
		return
	}

	duration := ""
	if alert.ResolvedAt != nil {
		d := alert.ResolvedAt.Sub(alert.TriggeredAt)
		duration = fmt.Sprintf("%v", d.Round(time.Second))
	}

	title := fmt.Sprintf("🟢 恢复：%s 已恢复正常", monitor.Name)
	text := fmt.Sprintf("**探测点**：%s\n"+
		"**恢复正常时间**：%s\n"+
		"**持续时长**：%s",
		agent.Name,
		time.Now().Format("2006-01-02 15:04:05"),
		duration,
	)

	for _, ch := range channels {
		if !ch.Enabled {
			continue
		}
		if err := send(&ch, title, text); err != nil {
			log.Printf("[NotifyResolved] send to channel %s error: %v", ch.Name, err)
		}
	}
}

// SendTest sends a test notification to the given alert channel.
func (s *NotifyService) SendTest(ch *model.AlertChannel) error {
	return send(ch, "🔔 测试通知", "这是一条测试消息，如果您收到此消息说明告警通道配置正确。")
}

// send dispatches a notification to the appropriate channel based on type.
func send(ch *model.AlertChannel, title, text string) error {
	switch ch.Type {
	case "dingtalk":
		return sendDingtalk(ch, title, text)
	case "wechat_work":
		return sendWechatWork(ch, title, text)
	case "feishu":
		return sendFeishu(ch, title, text)
	case "webhook":
		return sendWebhook(ch, title, text)
	default:
		return fmt.Errorf("unsupported channel type: %s", ch.Type)
	}
}

// sendDingtalk sends a DingTalk markdown notification with optional HMAC-SHA256 signing.
func sendDingtalk(ch *model.AlertChannel, title, text string) error {
	webhookURL := ch.WebhookURL
	if ch.Secret != "" {
		timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
		stringToSign := timestamp + "\n" + ch.Secret
		mac := hmac.New(sha256.New, []byte(ch.Secret))
		mac.Write([]byte(stringToSign))
		sign := url.QueryEscape(base64.StdEncoding.EncodeToString(mac.Sum(nil)))
		sep := "&"
		if !strings.Contains(webhookURL, "?") {
			sep = "?"
		}
		webhookURL = webhookURL + sep + "timestamp=" + timestamp + "&sign=" + sign
	}

	body := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
	}
	return postJSON(webhookURL, body)
}

// sendWechatWork sends a WeCom (企业微信) markdown notification.
func sendWechatWork(ch *model.AlertChannel, title, text string) error {
	body := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": title + "\n" + text,
		},
	}
	return postJSON(ch.WebhookURL, body)
}

// sendFeishu sends a Feishu (飞书) interactive card notification with optional signing.
func sendFeishu(ch *model.AlertChannel, title, text string) error {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	sign := ""
	if ch.Secret != "" {
		stringToSign := timestamp + "\n" + ch.Secret
		mac := hmac.New(sha256.New, []byte(ch.Secret))
		mac.Write([]byte(stringToSign))
		sign = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	}

	body := map[string]interface{}{
		"timestamp": timestamp,
		"sign":      sign,
		"msg_type":  "interactive",
		"card": map[string]interface{}{
			"header": map[string]interface{}{
				"title": map[string]string{
					"tag":     "plain_text",
					"content": title,
				},
			},
			"elements": []map[string]interface{}{
				{
					"tag": "div",
					"text": map[string]string{
						"tag":     "lark_md",
						"content": text,
					},
				},
			},
		},
	}
	return postJSON(ch.WebhookURL, body)
}

// sendWebhook sends a generic webhook notification.
func sendWebhook(ch *model.AlertChannel, title, text string) error {
	body := map[string]interface{}{
		"title":   title,
		"content": text,
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	}
	return postJSON(ch.WebhookURL, body)
}

// postJSON sends a JSON POST request to the given URL.
func postJSON(targetURL string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := http.Post(targetURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("notification POST returned status %d", resp.StatusCode)
	}
	return nil
}
