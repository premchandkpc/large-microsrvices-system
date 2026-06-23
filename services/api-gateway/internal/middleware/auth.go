package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/premchandkpc/large-microsrvices-system/services/api-gateway/internal/config"
	"github.com/premchandkpc/large-microsrvices-system/services/api-gateway/internal/service"
	"go.uber.org/zap"
)

func Authentication(cfg *config.Config, svc *service.ServiceRegistry, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		userID, roles, err := svc.ValidateToken(c.Request.Context(), token)
		if err != nil {
			logger.Warn("token validation failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("user_id", userID)
		c.Set("roles", roles)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("roles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no roles found"})
			return
		}

		roleList := userRoles.([]string)
		for _, required := range roles {
			for _, userRole := range roleList {
				if userRole == required {
					c.Next()
					return
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}
