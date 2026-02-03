package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"
)

// OAuthHandler OAuth 处理器
type OAuthHandler struct {
	oauthService *service.OAuthService
}

// NewOAuthHandler 创建 OAuth 处理器
func NewOAuthHandler(oauthService *service.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
	}
}

// GetProviders 获取已启用的 OAuth 提供商
// GET /api/v1/oauth/providers
func (h *OAuthHandler) GetProviders(c *gin.Context) {
	providers := h.oauthService.GetEnabledProviders()
	response.Success(c, http.StatusOK, "enabled oauth providers", gin.H{
		"providers": providers,
	})
}

// Authorize 获取授权 URL
// GET /api/v1/oauth/:provider/authorize
func (h *OAuthHandler) Authorize(c *gin.Context) {
	provider := c.Param("provider")
	redirectURI := c.Query("redirect_uri") // 前端回调地址

	if redirectURI == "" {
		redirectURI = c.Query("redirect") // 兼容
	}

	authURL, err := h.oauthService.GetAuthorizationURL(provider, redirectURI)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOAuthProviderNotEnabled):
			response.Error(c, http.StatusBadRequest, "oauth provider not enabled")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to generate authorization url")
		}
		return
	}

	response.Success(c, http.StatusOK, "authorization url generated", gin.H{
		"authorization_url": authURL,
	})
}

// Callback 处理 OAuth 回调
// GET /api/v1/oauth/:provider/callback
func (h *OAuthHandler) Callback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	// 检查是否有错误
	if errorParam != "" {
		errorDesc := c.Query("error_description")
		response.Error(c, http.StatusBadRequest, "oauth error: "+errorParam+": "+errorDesc)
		return
	}

	if code == "" {
		response.Error(c, http.StatusBadRequest, "missing authorization code")
		return
	}

	if state == "" {
		response.Error(c, http.StatusBadRequest, "missing state parameter")
		return
	}

	tokens, frontendRedirectURI, err := h.oauthService.HandleCallback(c.Request.Context(), provider, code, state)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOAuthInvalidState):
			response.Error(c, http.StatusBadRequest, "invalid or expired state")
		case errors.Is(err, service.ErrOAuthInvalidCode):
			response.Error(c, http.StatusBadRequest, "invalid authorization code")
		case errors.Is(err, service.ErrOAuthUserInfoFailed):
			response.Error(c, http.StatusBadGateway, "failed to get user info from provider")
		case errors.Is(err, service.ErrOAuthProviderNotEnabled):
			response.Error(c, http.StatusBadRequest, "oauth provider not enabled")
		default:
			response.Error(c, http.StatusInternalServerError, "oauth callback failed: "+err.Error())
		}
		return
	}

	// 如果有前端回调地址，重定向到前端并带上 token
	if frontendRedirectURI != "" {
		// 构建重定向 URL，将 token 作为 query 参数
		redirectURL := frontendRedirectURI +
			"?access_token=" + tokens.AccessToken +
			"&refresh_token=" + tokens.RefreshToken +
			"&token_type=Bearer" +
			"&expires_in=" + strconv.Itoa(tokens.ExpiresIn)

		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// 否则直接返回 token
	response.Success(c, http.StatusOK, "oauth login successful", tokens)
}

// CallbackPost 处理 OAuth 回调 (POST 方式，用于前端中转)
// POST /api/v1/oauth/:provider/callback
func (h *OAuthHandler) CallbackPost(c *gin.Context) {
	provider := c.Param("provider")

	var req struct {
		Code  string `json:"code" binding:"required"`
		State string `json:"state" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	tokens, _, err := h.oauthService.HandleCallback(c.Request.Context(), provider, req.Code, req.State)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOAuthInvalidState):
			response.Error(c, http.StatusBadRequest, "invalid or expired state")
		case errors.Is(err, service.ErrOAuthInvalidCode):
			response.Error(c, http.StatusBadRequest, "invalid authorization code")
		case errors.Is(err, service.ErrOAuthUserInfoFailed):
			response.Error(c, http.StatusBadGateway, "failed to get user info from provider")
		case errors.Is(err, service.ErrOAuthProviderNotEnabled):
			response.Error(c, http.StatusBadRequest, "oauth provider not enabled")
		default:
			response.Error(c, http.StatusInternalServerError, "oauth callback failed: "+err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "oauth login successful", tokens)
}
