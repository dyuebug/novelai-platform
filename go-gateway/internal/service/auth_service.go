package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"go-gateway/internal/config"
	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptCost           = 12
	AccessTokenTTL       = 15 * time.Minute
	RefreshTokenTTL      = 7 * 24 * time.Hour
	MaxFailedAttempts    = 5
	AccountLockDuration  = 30 * time.Minute
	PasswordResetTTL     = 1 * time.Hour
	PasswordResetCooldown = 5 * time.Minute
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account is locked")
	ErrWeakPassword       = errors.New("password does not meet requirements")
	ErrResetCooldown      = errors.New("password reset request too frequent")
)

type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	resetTokenRepo   *repository.PasswordResetTokenRepository
	emailService     *EmailService
	config           config.AuthConfig
}

func NewAuthService(cfg config.AuthConfig, emailService *EmailService) (*AuthService, error) {
	// 验证 JWT_SECRET 长度
	if len(cfg.JWTSecret) < 32 {
		log.Printf("WARNING: JWT_SECRET is shorter than recommended 32 characters (current: %d). Consider using a stronger secret.", len(cfg.JWTSecret))
	}

	return &AuthService{
		userRepo:         repository.NewUserRepository(),
		refreshTokenRepo: repository.NewRefreshTokenRepository(),
		resetTokenRepo:   repository.NewPasswordResetTokenRepository(),
		emailService:     emailService,
		config:           cfg,
	}, nil
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*model.User, error) {
	// 验证密码强度
	if !isStrongPassword(req.Password) {
		return nil, ErrWeakPassword
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), BcryptCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
		Status:       "active",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 检查账户是否被锁定
	if user.IsLocked() {
		return nil, ErrAccountLocked
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// 增加失败次数
		attempts, _ := s.userRepo.IncrementFailedAttempts(ctx, user.ID)
		if attempts >= MaxFailedAttempts {
			s.userRepo.LockAccount(ctx, user.ID, AccountLockDuration)
		}
		return nil, ErrInvalidCredentials
	}

	// 更新最后登录时间
	s.userRepo.UpdateLastLogin(ctx, user.ID)

	// 生成令牌
	return s.GenerateTokens(ctx, user)
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// 查找刷新令牌
	rt, err := s.refreshTokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if time.Now().After(rt.ExpiresAt) {
		return nil, repository.ErrTokenExpired
	}

	// 检查是否已撤销
	if rt.RevokedAt != nil {
		return nil, repository.ErrTokenUsed
	}

	// 撤销旧令牌
	s.refreshTokenRepo.Revoke(ctx, refreshToken)

	// 获取用户
	user, err := s.userRepo.FindByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	// 生成新令牌
	return s.GenerateTokens(ctx, user)
}

// GenerateTokens 生成访问令牌和刷新令牌
func (s *AuthService) GenerateTokens(ctx context.Context, user *model.User) (*TokenResponse, error) {
	now := time.Now()

	// 生成 Access Token (HMAC-SHA256)
	accessClaims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"iss":   s.config.Issuer,
		"aud":   s.config.Audience,
		"iat":   now.Unix(),
		"exp":   now.Add(AccessTokenTTL).Unix(),
		"roles": []string{user.Role},
		"email": user.Email,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	// 生成 Refresh Token (随机字符串)
	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, err
	}
	refreshTokenString := base64.URLEncoding.EncodeToString(refreshTokenBytes)

	// 保存刷新令牌
	rt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: now.Add(RefreshTokenTTL),
	}
	if err := s.refreshTokenRepo.Create(ctx, rt); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int(AccessTokenTTL.Seconds()),
	}, nil
}

// ValidateToken 验证访问令牌
func (s *AuthService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

// RequestPasswordReset 请求密码重置
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// 不泄露用户是否存在
		return nil
	}

	// 检查冷却时间
	_, err = s.resetTokenRepo.FindRecentByUserID(ctx, user.ID, PasswordResetCooldown)
	if err == nil {
		return ErrResetCooldown
	}

	// 生成重置令牌
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return err
	}
	tokenString := base64.URLEncoding.EncodeToString(tokenBytes)

	resetToken := &model.PasswordResetToken{
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(PasswordResetTTL),
	}
	if err := s.resetTokenRepo.Create(ctx, resetToken); err != nil {
		return err
	}

	// 发送密码重置邮件
	if s.emailService != nil {
		if err := s.emailService.SendPasswordResetEmail(user.Email, tokenString, ""); err != nil {
			// 记录错误但不阻止流程
			// log.Printf("Failed to send password reset email: %v", err)
		}
	}

	return nil
}

// ResetPassword 重置密码
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// 验证密码强度
	if !isStrongPassword(newPassword) {
		return ErrWeakPassword
	}

	// 查找令牌
	resetToken, err := s.resetTokenRepo.FindByToken(ctx, token)
	if err != nil {
		return err
	}

	// 检查是否过期
	if time.Now().After(resetToken.ExpiresAt) {
		return repository.ErrTokenExpired
	}

	// 检查是否已使用
	if resetToken.UsedAt != nil {
		return repository.ErrTokenUsed
	}

	// 哈希新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), BcryptCost)
	if err != nil {
		return err
	}

	// 更新密码
	user, err := s.userRepo.FindByID(ctx, resetToken.UserID)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// 标记令牌已使用
	s.resetTokenRepo.MarkUsed(ctx, token)

	// 撤销所有刷新令牌
	s.refreshTokenRepo.RevokeAllForUser(ctx, user.ID)

	return nil
}

// isStrongPassword 检查密码强度
// 要求: 8-128 字符，至少 1 大写，1 小写，1 数字
func isStrongPassword(password string) bool {
	if len(password) < 8 || len(password) > 128 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}
