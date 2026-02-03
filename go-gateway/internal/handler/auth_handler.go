package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/repository"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 用户注册
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserAlreadyExists):
			response.Error(c, http.StatusConflict, "user already exists")
		case errors.Is(err, service.ErrWeakPassword):
			response.Error(c, http.StatusBadRequest, "password must be 8-128 chars with at least 1 uppercase, 1 lowercase, and 1 digit")
		default:
			response.Error(c, http.StatusInternalServerError, "registration failed")
		}
		return
	}

	response.Success(c, http.StatusCreated, "user registered successfully", gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}

// Login 用户登录
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	tokens, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			response.Error(c, http.StatusUnauthorized, "invalid email or password")
		case errors.Is(err, service.ErrAccountLocked):
			response.Error(c, http.StatusLocked, "account is locked, please try again later")
		default:
			response.Error(c, http.StatusInternalServerError, "login failed")
		}
		return
	}

	response.Success(c, http.StatusOK, "login successful", tokens)
}

// Refresh 刷新令牌
// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrTokenNotFound):
			response.Error(c, http.StatusUnauthorized, "invalid refresh token")
		case errors.Is(err, repository.ErrTokenExpired):
			response.Error(c, http.StatusUnauthorized, "refresh token expired")
		case errors.Is(err, repository.ErrTokenUsed):
			response.Error(c, http.StatusUnauthorized, "refresh token already used")
		default:
			response.Error(c, http.StatusInternalServerError, "token refresh failed")
		}
		return
	}

	response.Success(c, http.StatusOK, "token refreshed successfully", tokens)
}

// RequestPasswordReset 请求密码重置
// POST /api/v1/auth/password-reset/request
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	err := h.authService.RequestPasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, service.ErrResetCooldown) {
			response.Error(c, http.StatusTooManyRequests, "please wait before requesting another reset")
			return
		}
		// 不泄露其他错误
	}

	// 无论成功与否都返回 202
	response.Success(c, http.StatusAccepted, "if the email exists, a reset link will be sent", nil)
}

// ResetPassword 重置密码
// POST /api/v1/auth/password-reset/verify
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8,max=128"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrTokenNotFound):
			response.Error(c, http.StatusBadRequest, "invalid reset token")
		case errors.Is(err, repository.ErrTokenExpired):
			response.Error(c, http.StatusGone, "reset token expired")
		case errors.Is(err, repository.ErrTokenUsed):
			response.Error(c, http.StatusBadRequest, "reset token already used")
		case errors.Is(err, service.ErrWeakPassword):
			response.Error(c, http.StatusBadRequest, "password must be 8-128 chars with at least 1 uppercase, 1 lowercase, and 1 digit")
		default:
			response.Error(c, http.StatusInternalServerError, "password reset failed")
		}
		return
	}

	response.Success(c, http.StatusOK, "password reset successfully", nil)
}

// GetCurrentUser 获取当前用户信息
// GET /api/v1/auth/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// 从 JWT 中间件获取用户信息
	claims, exists := c.Get("claims")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	response.Success(c, http.StatusOK, "current user", claims)
}
