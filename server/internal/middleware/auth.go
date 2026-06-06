package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"uptime-monitor/server/internal/handler"
	"uptime-monitor/server/internal/service"
)

// AuthMiddleware returns a Gin middleware that validates JWT Bearer tokens.
// On success it stores "user_id" (uint) in the Gin context.
func AuthMiddleware(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handler.Error(c, http.StatusUnauthorized, 40001, "认证令牌无效")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			handler.Error(c, http.StatusUnauthorized, 40001, "认证令牌无效")
			c.Abort()
			return
		}

		claims, err := authSvc.ValidateToken(parts[1])
		if err != nil {
			handler.Error(c, http.StatusUnauthorized, 40001, "认证令牌无效")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
