package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"uptime-monitor/server/internal/model"
	"uptime-monitor/server/internal/repository"
	"uptime-monitor/server/internal/service"
)

// AuthHandler handles authentication-related API endpoints.
type AuthHandler struct {
	authSvc *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Login handles POST /api/auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, 40101, "参数错误")
		return
	}

	token, err := h.authSvc.Login(req.Username, req.Password)
	if err != nil {
		Error(c, http.StatusOK, 40001, err.Error())
		return
	}

	Success(c, gin.H{"token": token})
}

// Logout handles POST /api/auth/logout (stateless — always succeeds).
func (h *AuthHandler) Logout(c *gin.Context) {
	Success(c, nil)
}

// Me handles GET /api/auth/me — returns the current user's information.
func (h *AuthHandler) Me(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		Error(c, http.StatusOK, 40001, "未获取到用户信息")
		return
	}
	userID, ok := val.(uint)
	if !ok {
		Error(c, http.StatusOK, 40001, "用户信息类型异常")
		return
	}

	user, err := repository.FindUserByID(userID)
	if err != nil {
		Error(c, http.StatusOK, 50001, "用户不存在")
		return
	}

	Success(c, user)
}
