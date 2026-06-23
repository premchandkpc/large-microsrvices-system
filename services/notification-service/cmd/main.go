package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
	"github.com/premchandkpc/large-microsrvices-system/services/notification-service/internal/config"
	"github.com/premchandkpc/large-microsrvices-system/services/notification-service/internal/handler"
	"github.com/premchandkpc/large-microsrvices-system/services/notification-service/internal/provider"
	"github.com/premchandkpc/large-microsrvices-system/services/notification-service/internal/service"
	"go.uber.org/zap"
)

type WebSocketHub struct {
	clients map[string]*websocket.Conn
	mu      sync.RWMutex
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients: make(map[string]*websocket.Conn),
	}
}

func (h *WebSocketHub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = conn
}

func (h *WebSocketHub) Unregister(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, userID)
}

func (h *WebSocketHub) Send(userID string, message []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conn, ok := h.clients[userID]
	if !ok {
		return fmt.Errorf("user not connected")
	}
	return conn.WriteMessage(websocket.TextMessage, message)
}

func main() {
	cfg := config.Load()
	logger := initLogger(cfg)

	hub := NewWebSocketHub()

	emailProvider := provider.NewSMTPProvider(cfg, logger)

	svc := service.NewNotificationService(emailProvider, hub, logger)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.KafkaBrokers,
		Topic:       "notifications",
		GroupID:     "notification-service",
		MinBytes:    10,
		MaxBytes:    10e6,
		Compression: compress.Snappy,
		MaxWait:     1 * time.Second,
		StartOffset: kafka.LastOffset,
	})
	defer reader.Close()

	go func() {
		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				logger.Error("kafka read error", zap.Error(err))
				continue
			}
			if err := svc.HandleNotificationEvent(context.Background(), msg.Value); err != nil {
				logger.Error("handle notification error", zap.Error(err))
			}
		}
	}()

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "notification-service"})
	})

	v1 := router.Group("/api/v1/notifications")
	{
		v1.GET("", handler.GetNotifications(svc))
		v1.POST("/read", handler.MarkAsRead(svc))
		v1.POST("/send", handler.SendNotification(svc))
	}

	router.GET("/ws", func(c *gin.Context) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("ws upgrade failed", zap.Error(err))
			return
		}
		userID := c.Query("user_id")
		hub.Register(userID, conn)
		defer hub.Unregister(userID)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	go func() {
		logger.Info("notification-service starting", zap.Int("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func initLogger(cfg *config.Config) *zap.Logger {
	logger, _ := zap.NewProduction()
	if cfg.Environment == "development" {
		logger, _ = zap.NewDevelopment()
	}
	return logger
}
