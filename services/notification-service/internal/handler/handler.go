package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/premchandkpc/large-microsrvices-system/services/notification-service/internal/model"
	"github.com/premchandkpc/large-microsrvices-system/services/notification-service/internal/service"
)

func GetNotifications(svc *service.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		notifs, err := svc.GetNotifications(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusOK, []interface{}{})
			return
		}
		c.JSON(http.StatusOK, notifs)
	}
}

func MarkAsRead(svc *service.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID string   `json:"user_id"`
			IDs    []string `json:"ids"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := svc.MarkAsRead(c.Request.Context(), req.UserID, req.IDs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func SendNotification(svc *service.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.SendNotificationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		notif, err := svc.Send(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, notif)
	}
}
