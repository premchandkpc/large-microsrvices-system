package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/premchandkpc/large-microsrvices-system/services/api-gateway/internal/config"
	"github.com/premchandkpc/large-microsrvices-system/services/api-gateway/internal/service"
	"go.uber.org/zap"
)

func newTestHandler() *Handler {
	cfg := &config.Config{Environment: "test"}
	logger, _ := zap.NewDevelopment()
	svc, _ := service.NewServiceRegistry(cfg, logger)
	return NewHandler(svc, logger, cfg)
}

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newTestHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

	h.Health(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestLoginValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newTestHandler()

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"missing email", `{"password":"testpass123"}`, http.StatusBadRequest},
		{"missing password", `{"email":"test@example.com"}`, http.StatusBadRequest},
		{"invalid email", `{"email":"bad","password":"testpass123"}`, http.StatusBadRequest},
		{"short password", `{"email":"test@example.com","password":"short"}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)

			h.Login(c)

			if w.Code != tt.wantStatus {
				t.Errorf("%s: expected %d, got %d", tt.name, w.Code, tt.wantStatus)
			}
		})
	}
}

func TestRegisterValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newTestHandler()

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"missing name", `{"email":"test@example.com","password":"testpass123"}`, http.StatusBadRequest},
		{"missing email", `{"password":"testpass123","name":"Test"}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", nil)

			h.Register(c)

			if w.Code != tt.wantStatus {
				t.Errorf("%s: expected %d, got %d", tt.name, w.Code, tt.wantStatus)
			}
		})
	}
}
