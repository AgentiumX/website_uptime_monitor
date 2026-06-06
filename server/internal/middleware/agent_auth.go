package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"uptime-monitor/server/internal/handler"
	"uptime-monitor/server/internal/repository"
)

// AgentAuthMiddleware returns a Gin middleware that validates Agent Bearer tokens.
// On success it stores "agent_id" (uint) and "agent_name" (string) in the Gin context.
func AgentAuthMiddleware() gin.HandlerFunc {
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

		agent, err := repository.FindAgentByToken(parts[1])
		if err != nil {
			handler.Error(c, http.StatusUnauthorized, 40001, "认证令牌无效")
			c.Abort()
			return
		}

		c.Set("agent_id", agent.ID)
		c.Set("agent_name", agent.Name)
		c.Next()
	}
}
