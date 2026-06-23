package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"

	"github.com/google/uuid"
	"github.com/premchandkpc/large-microsrvices-system/services/notification-service/internal/model"
	"github.com/premchandkpc/large-microsrvices-system/services/notification-service/internal/provider"
	"go.uber.org/zap"
)

type WebSocketHub interface {
	Register(userID string, conn *websocket.Conn)
	Unregister(userID string)
	Send(userID string, message []byte) error
}

type NotificationService struct {
	emailProvider *provider.SMTPProvider
	wsHub         WebSocketHub
	logger        *zap.Logger
}

func NewNotificationService(emailProvider *provider.SMTPProvider, wsHub WebSocketHub, logger *zap.Logger) *NotificationService {
	return &NotificationService{
		emailProvider: emailProvider,
		wsHub:         wsHub,
		logger:        logger,
	}
}

func (s *NotificationService) Send(ctx context.Context, req *model.SendNotificationRequest) (*model.Notification, error) {
	notif := &model.Notification{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		TenantID:  req.TenantID,
		Channel:   req.Channel,
		Title:     req.Title,
		Body:      req.Body,
		Metadata:  req.Metadata,
		CreatedAt: time.Now(),
	}

	if err := s.dispatch(ctx, notif); err != nil {
		return nil, err
	}

	return notif, nil
}

func (s *NotificationService) dispatch(ctx context.Context, notif *model.Notification) error {
	switch notif.Channel {
	case model.ChannelEmail:
		email, ok := notif.Metadata["email"]
		if !ok {
			return fmt.Errorf("email address required for email channel")
		}
		go s.emailProvider.Send(email, notif.Title, notif.Body)

	case model.ChannelInApp:
		data, _ := json.Marshal(notif)
		go s.wsHub.Send(notif.UserID, data)

	case model.ChannelPush:
		// Firebase / APNs integration would go here

	case model.ChannelSlack:
		// Slack webhook integration would go here
	}

	return nil
}

func (s *NotificationService) HandleNotificationEvent(ctx context.Context, data []byte) error {
	var event model.NotificationEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("unmarshaling event: %w", err)
	}

	notif := &model.Notification{
		ID:        uuid.New().String(),
		UserID:    event.UserID,
		Title:     event.Title,
		Body:      event.Body,
		Channel:   event.Channel,
		CreatedAt: time.Now(),
	}

	return s.dispatch(ctx, notif)
}

func (s *NotificationService) GetNotifications(ctx context.Context, userID string) ([]model.Notification, error) {
	return nil, fmt.Errorf("not implemented: read from db")
}

func (s *NotificationService) MarkAsRead(ctx context.Context, userID string, ids []string) error {
	return fmt.Errorf("not implemented: update db")
}
