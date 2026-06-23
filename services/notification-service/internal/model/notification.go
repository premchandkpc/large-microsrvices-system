package model

import "time"

type Notification struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	TenantID  string            `json:"tenant_id"`
	Channel   string            `json:"channel"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	IsRead    bool              `json:"is_read"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
}

type SendNotificationRequest struct {
	UserID   string            `json:"user_id"`
	TenantID string            `json:"tenant_id"`
	Channel  string            `json:"channel"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	Metadata map[string]string `json:"metadata"`
	Priority string            `json:"priority"`
}

type NotificationEvent struct {
	Type      string `json:"type"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Channel   string `json:"channel"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	ChannelEmail     = "email"
	ChannelInApp     = "in_app"
	ChannelPush      = "push"
	ChannelSlack     = "slack"
	ChannelWebhook   = "webhook"

	PriorityLow    = "low"
	PriorityNormal = "normal"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"
)
