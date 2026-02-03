package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"go-gateway/internal/config"
	"go-gateway/internal/repository"
)

var (
	ErrOAuthProviderNotEnabled = errors.New("oauth provider not enabled")
	ErrOAuthInvalidState       = errors.New("invalid oauth state")
	ErrOAuthInvalidCode        = errors.New("invalid oauth code")
	ErrOAuthUserInfoFailed     = errors.New("failed to get user info")
	ErrOAuthPKCEVerifyFailed   = errors.New("PKCE verification failed")
)

// OAuthService OAuth 服务
type OAuthService struct {
	cfg         *config.OAuthConfig
	authService *AuthService
	httpClient  *http.Client

	// state 存储 (生产环境应使用 Redis)
	stateMu    sync.RWMutex
	stateStore map[string]*OAuthState
}

// OAuthState OAuth 状态
type OAuthState struct {
	Provider      string
	CodeVerifier  string // PKCE code_verifier
	RedirectURI   string // 前端回调地址
	CreatedAt     time.Time
	ExpiresAt     time.Time
}

// OAuthUserInfo OAuth 用户信息
type OAuthUserInfo struct {
	Provider    string
	ProviderID  string
	Email       string
	Username    string
	DisplayName string
	AvatarURL   string
}

// NewOAuthService 创建 OAuth 服务
func NewOAuthService(cfg *config.OAuthConfig, authService *AuthService) *OAuthService {
	svc := &OAuthService{
		cfg:         cfg,
		authService: authService,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		stateStore: make(map[string]*OAuthState),
	}

	// 启动清理过期 state 的 goroutine
	go svc.cleanupExpiredStates()

	return svc
}

// GetAuthorizationURL 获取授权 URL
func (s *OAuthService) GetAuthorizationURL(provider string, frontendRedirectURI string) (string, error) {
	providerCfg, err := s.getProviderConfig(provider)
	if err != nil {
		return "", err
	}

	// 生成 state
	state, err := s.generateState()
	if err != nil {
		return "", err
	}

	// 生成 PKCE code_verifier 和 code_challenge
	codeVerifier, err := s.generateCodeVerifier()
	if err != nil {
		return "", err
	}
	codeChallenge := s.generateCodeChallenge(codeVerifier)

	// 存储 state
	s.stateMu.Lock()
	s.stateStore[state] = &OAuthState{
		Provider:     provider,
		CodeVerifier: codeVerifier,
		RedirectURI:  frontendRedirectURI,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(10 * time.Minute),
	}
	s.stateMu.Unlock()

	// 构建授权 URL
	params := url.Values{}
	params.Set("client_id", providerCfg.ClientID)
	params.Set("redirect_uri", providerCfg.RedirectURL)
	params.Set("response_type", "code")
	params.Set("scope", strings.Join(providerCfg.Scopes, " "))
	params.Set("state", state)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")

	return fmt.Sprintf("%s?%s", providerCfg.AuthURL, params.Encode()), nil
}

// HandleCallback 处理 OAuth 回调
func (s *OAuthService) HandleCallback(ctx context.Context, provider, code, state string) (*TokenResponse, string, error) {
	// 验证 state
	s.stateMu.RLock()
	oauthState, exists := s.stateStore[state]
	s.stateMu.RUnlock()

	if !exists || oauthState.Provider != provider || time.Now().After(oauthState.ExpiresAt) {
		return nil, "", ErrOAuthInvalidState
	}

	// 删除已使用的 state
	s.stateMu.Lock()
	delete(s.stateStore, state)
	s.stateMu.Unlock()

	providerCfg, err := s.getProviderConfig(provider)
	if err != nil {
		return nil, "", err
	}

	// 交换 code 获取 access_token
	accessToken, err := s.exchangeCode(providerCfg, code, oauthState.CodeVerifier)
	if err != nil {
		return nil, "", err
	}

	// 获取用户信息
	userInfo, err := s.getUserInfo(provider, providerCfg, accessToken)
	if err != nil {
		return nil, "", err
	}

	// 查找或创建用户
	tokens, err := s.findOrCreateUser(ctx, userInfo)
	if err != nil {
		return nil, "", err
	}

	return tokens, oauthState.RedirectURI, nil
}

// getProviderConfig 获取提供商配置
func (s *OAuthService) getProviderConfig(provider string) (*config.OAuthProviderConfig, error) {
	var cfg *config.OAuthProviderConfig

	switch strings.ToLower(provider) {
	case "github":
		cfg = &s.cfg.GitHub
	case "google":
		cfg = &s.cfg.Google
	case "linuxdo":
		cfg = &s.cfg.LinuxDo
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}

	if !cfg.Enabled {
		return nil, ErrOAuthProviderNotEnabled
	}

	return cfg, nil
}

// generateState 生成随机 state
func (s *OAuthService) generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// generateCodeVerifier 生成 PKCE code_verifier
func (s *OAuthService) generateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generateCodeChallenge 生成 PKCE code_challenge (S256)
func (s *OAuthService) generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// exchangeCode 交换授权码获取 access_token
func (s *OAuthService) exchangeCode(cfg *config.OAuthProviderConfig, code, codeVerifier string) (string, error) {
	data := url.Values{}
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", cfg.RedirectURL)
	data.Set("grant_type", "authorization_code")
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequest("POST", cfg.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Error       string `json:"error"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.Error != "" {
		return "", fmt.Errorf("token error: %s", tokenResp.Error)
	}

	return tokenResp.AccessToken, nil
}

// getUserInfo 获取用户信息
func (s *OAuthService) getUserInfo(provider string, cfg *config.OAuthProviderConfig, accessToken string) (*OAuthUserInfo, error) {
	req, err := http.NewRequest("GET", cfg.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrOAuthUserInfoFailed
	}

	// 根据不同提供商解析用户信息
	switch strings.ToLower(provider) {
	case "github":
		return s.parseGitHubUserInfo(body)
	case "google":
		return s.parseGoogleUserInfo(body)
	case "linuxdo":
		return s.parseLinuxDoUserInfo(body)
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}

// parseGitHubUserInfo 解析 GitHub 用户信息
func (s *OAuthService) parseGitHubUserInfo(body []byte) (*OAuthUserInfo, error) {
	var data struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	displayName := data.Name
	if displayName == "" {
		displayName = data.Login
	}

	return &OAuthUserInfo{
		Provider:    "github",
		ProviderID:  fmt.Sprintf("%d", data.ID),
		Email:       data.Email,
		Username:    data.Login,
		DisplayName: displayName,
		AvatarURL:   data.AvatarURL,
	}, nil
}

// parseGoogleUserInfo 解析 Google 用户信息
func (s *OAuthService) parseGoogleUserInfo(body []byte) (*OAuthUserInfo, error) {
	var data struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		VerifiedEmail bool   `json:"verified_email"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	// 从邮箱生成用户名
	username := strings.Split(data.Email, "@")[0]

	return &OAuthUserInfo{
		Provider:    "google",
		ProviderID:  data.ID,
		Email:       data.Email,
		Username:    username,
		DisplayName: data.Name,
		AvatarURL:   data.Picture,
	}, nil
}

// parseLinuxDoUserInfo 解析 LinuxDo 用户信息
func (s *OAuthService) parseLinuxDoUserInfo(body []byte) (*OAuthUserInfo, error) {
	var data struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar_url"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	displayName := data.Name
	if displayName == "" {
		displayName = data.Username
	}

	return &OAuthUserInfo{
		Provider:    "linuxdo",
		ProviderID:  fmt.Sprintf("%d", data.ID),
		Email:       data.Email,
		Username:    data.Username,
		DisplayName: displayName,
		AvatarURL:   data.Avatar,
	}, nil
}

// findOrCreateUser 查找或创建用户
func (s *OAuthService) findOrCreateUser(ctx context.Context, info *OAuthUserInfo) (*TokenResponse, error) {
	// 先尝试通过 OAuth 绑定查找用户
	user, err := repository.GetUserByOAuthProvider(ctx, info.Provider, info.ProviderID)
	if err == nil && user != nil {
		// 用户已存在，直接生成 token
		return s.authService.GenerateTokens(ctx, user)
	}

	// 如果有邮箱，尝试通过邮箱查找
	if info.Email != "" {
		user, err = repository.GetUserByEmail(ctx, info.Email)
		if err == nil && user != nil {
			// 绑定 OAuth 到现有用户
			if err := repository.BindOAuthToUser(ctx, user.ID, info.Provider, info.ProviderID); err != nil {
				return nil, err
			}
			return s.authService.GenerateTokens(ctx, user)
		}
	}

	// 创建新用户
	// 确保用户名唯一
	username := info.Username
	if username == "" {
		username = fmt.Sprintf("%s_%s", info.Provider, info.ProviderID)
	}

	// 检查用户名是否已存在，如果存在则添加随机后缀
	for i := 0; i < 5; i++ {
		_, err := repository.GetUserByUsername(ctx, username)
		if err != nil {
			break // 用户名可用
		}
		// 添加随机后缀
		suffix := make([]byte, 4)
		rand.Read(suffix)
		username = fmt.Sprintf("%s_%s", info.Username, base64.RawURLEncoding.EncodeToString(suffix)[:6])
	}

	// 创建用户 (OAuth 用户不需要密码)
	user, err = repository.CreateOAuthUser(ctx, &repository.CreateOAuthUserParams{
		Username:    username,
		Email:       info.Email,
		DisplayName: info.DisplayName,
		AvatarURL:   info.AvatarURL,
		Provider:    info.Provider,
		ProviderID:  info.ProviderID,
	})
	if err != nil {
		return nil, err
	}

	return s.authService.GenerateTokens(ctx, user)
}

// cleanupExpiredStates 清理过期的 state
func (s *OAuthService) cleanupExpiredStates() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.stateMu.Lock()
		now := time.Now()
		for state, data := range s.stateStore {
			if now.After(data.ExpiresAt) {
				delete(s.stateStore, state)
			}
		}
		s.stateMu.Unlock()
	}
}

// GetEnabledProviders 获取已启用的提供商列表
func (s *OAuthService) GetEnabledProviders() []string {
	var providers []string
	if s.cfg.GitHub.Enabled {
		providers = append(providers, "github")
	}
	if s.cfg.Google.Enabled {
		providers = append(providers, "google")
	}
	if s.cfg.LinuxDo.Enabled {
		providers = append(providers, "linuxdo")
	}
	return providers
}
