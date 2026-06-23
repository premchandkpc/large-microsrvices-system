package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/premchandkpc/large-microsrvices-system/services/api-gateway/internal/config"
	"github.com/premchandkpc/large-microsrvices-system/services/api-gateway/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	svc    *service.ServiceRegistry
	logger *zap.Logger
	cfg    *config.Config
	upgrader websocket.Upgrader
}

func NewHandler(svc *service.ServiceRegistry, logger *zap.Logger, cfg *config.Config) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
		cfg:    cfg,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if cfg.Environment == "development" || origin == "" {
					return true
				}
				for _, allowed := range cfg.CorsAllowedOrigins {
					if origin == allowed {
						return true
					}
				}
				return false
			},
		},
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "api-gateway"})
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Warn("login failed", zap.String("email", req.Email), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Name     string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		h.logger.Warn("registration failed", zap.String("email", req.Email), zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Logout(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	c.ShouldBindJSON(&req)
	h.svc.Logout(c.Request.Context(), userID.(string), req.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *Handler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	requesterID, _ := c.Get("user_id")
	roles, _ := c.Get("roles")

	user, err := h.svc.GetUser(c.Request.Context(), userID, requesterID.(string), roles.([]string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) UploadDocument(c *gin.Context) {
	userID, _ := c.Get("user_id")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	maxSize := int64(100 << 20) // 100MB
	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large. max 100MB"})
		return
	}

	doc, err := h.svc.UploadDocument(c.Request.Context(), file, header, userID.(string))
	if err != nil {
		h.logger.Error("upload failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}
	c.JSON(http.StatusCreated, doc)
}

func (h *Handler) ListDocuments(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	status := c.Query("status")

	docs, err := h.svc.ListDocuments(c.Request.Context(), userID.(string), page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docs)
}

func (h *Handler) GetDocument(c *gin.Context) {
	docID := c.Param("id")
	userID, _ := c.Get("user_id")

	doc, err := h.svc.GetDocument(c.Request.Context(), docID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	c.JSON(http.StatusOK, doc)
}

func (h *Handler) ProcessDocument(c *gin.Context) {
	docID := c.Param("id")
	var req struct {
		Pipeline string `json:"pipeline" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pipeline type required"})
		return
	}

	job, err := h.svc.ProcessDocument(c.Request.Context(), docID, req.Pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, job)
}

func (h *Handler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query required"})
		return
	}

	results, err := h.svc.Search(c.Request.Context(), query, c.QueryArray("filter"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *Handler) VectorSearch(c *gin.Context) {
	var req struct {
		Query  string   `json:"query" binding:"required"`
		Filter []string `json:"filter"`
		TopK   int      `json:"top_k"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.svc.VectorSearch(c.Request.Context(), req.Query, req.Filter, req.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *Handler) GetNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	notifs, err := h.svc.GetNotifications(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifs)
}

func (h *Handler) MarkAsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		NotificationIDs []string `json:"notification_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.svc.MarkNotificationsRead(c.Request.Context(), userID.(string), req.NotificationIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *Handler) WebSocketNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("websocket upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()

	h.svc.HandleWebSocket(c.Request.Context(), conn, userID.(string))
}
