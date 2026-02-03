package repository

import (
	"context"
	"errors"
	"time"

	"go-gateway/internal/database"
	"go-gateway/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrTokenNotFound     = errors.New("token not found")
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenUsed         = errors.New("token already used")
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	result := database.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrUserAlreadyExists
		}
		return result.Error
	}
	return nil
}

// FindByEmail 通过邮箱查找用户
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	result := database.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByID 通过 ID 查找用户
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	result := database.DB.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByUsername 通过用户名查找用户
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := database.DB.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return database.DB.WithContext(ctx).Save(user).Error
}

// UpdateLastLogin 更新最后登录时间
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return database.DB.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at":   now,
			"failed_attempts": 0,
			"locked_until":    nil,
		}).Error
}

// IncrementFailedAttempts 增加失败尝试次数
func (r *UserRepository) IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) (int, error) {
	var user model.User
	err := database.DB.WithContext(ctx).Model(&user).
		Where("id = ?", userID).
		Update("failed_attempts", gorm.Expr("failed_attempts + 1")).Error
	if err != nil {
		return 0, err
	}

	database.DB.WithContext(ctx).Where("id = ?", userID).First(&user)
	return user.FailedAttempts, nil
}

// LockAccount 锁定账户
func (r *UserRepository) LockAccount(ctx context.Context, userID uuid.UUID, duration time.Duration) error {
	lockedUntil := time.Now().Add(duration)
	return database.DB.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Update("locked_until", lockedUntil).Error
}

// RefreshToken 相关操作

type RefreshTokenRepository struct{}

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{}
}

// Create 创建刷新令牌
func (r *RefreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	return database.DB.WithContext(ctx).Create(token).Error
}

// FindByToken 通过令牌查找
func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	result := database.DB.WithContext(ctx).Where("token = ?", token).First(&rt)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, result.Error
	}
	return &rt, nil
}

// Revoke 撤销令牌
func (r *RefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	now := time.Now()
	return database.DB.WithContext(ctx).Model(&model.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked_at", now).Error
}

// RevokeAllForUser 撤销用户所有令牌
func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return database.DB.WithContext(ctx).Model(&model.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

// DeleteExpired 删除过期令牌
func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return database.DB.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&model.RefreshToken{}).Error
}

// PasswordResetToken 相关操作

type PasswordResetTokenRepository struct{}

func NewPasswordResetTokenRepository() *PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{}
}

// Create 创建密码重置令牌
func (r *PasswordResetTokenRepository) Create(ctx context.Context, token *model.PasswordResetToken) error {
	return database.DB.WithContext(ctx).Create(token).Error
}

// FindByToken 通过令牌查找
func (r *PasswordResetTokenRepository) FindByToken(ctx context.Context, token string) (*model.PasswordResetToken, error) {
	var prt model.PasswordResetToken
	result := database.DB.WithContext(ctx).Where("token = ?", token).First(&prt)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, result.Error
	}
	return &prt, nil
}

// FindRecentByUserID 查找用户最近的重置请求
func (r *PasswordResetTokenRepository) FindRecentByUserID(ctx context.Context, userID uuid.UUID, within time.Duration) (*model.PasswordResetToken, error) {
	var prt model.PasswordResetToken
	since := time.Now().Add(-within)
	result := database.DB.WithContext(ctx).
		Where("user_id = ? AND created_at > ? AND used_at IS NULL", userID, since).
		Order("created_at DESC").
		First(&prt)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, result.Error
	}
	return &prt, nil
}

// MarkUsed 标记令牌已使用
func (r *PasswordResetTokenRepository) MarkUsed(ctx context.Context, token string) error {
	now := time.Now()
	return database.DB.WithContext(ctx).Model(&model.PasswordResetToken{}).
		Where("token = ?", token).
		Update("used_at", now).Error
}

// OAuth 相关函数 (包级别函数，供 OAuthService 使用)

// GetUserByOAuthProvider 通过 OAuth 提供商查找用户
func GetUserByOAuthProvider(ctx context.Context, provider, providerID string) (*model.User, error) {
	var user model.User
	result := database.DB.WithContext(ctx).
		Where("oauth_provider = ? AND oauth_id = ?", provider, providerID).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByEmail 通过邮箱查找用户
func GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	result := database.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByUsername 通过用户名查找用户
func GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := database.DB.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// BindOAuthToUser 绑定 OAuth 到现有用户
func BindOAuthToUser(ctx context.Context, userID uuid.UUID, provider, providerID string) error {
	return database.DB.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"oauth_provider": provider,
			"oauth_id":       providerID,
		}).Error
}

// CreateOAuthUserParams OAuth 用户创建参数
type CreateOAuthUserParams struct {
	Username    string
	Email       string
	DisplayName string
	AvatarURL   string
	Provider    string
	ProviderID  string
}

// CreateOAuthUser 创建 OAuth 用户
func CreateOAuthUser(ctx context.Context, params *CreateOAuthUserParams) (*model.User, error) {
	user := &model.User{
		ID:            uuid.New(),
		Username:      params.Username,
		Email:         params.Email,
		AvatarURL:     params.AvatarURL,
		Role:          "user",
		Status:        "active",
		OAuthProvider: params.Provider,
		OAuthID:       params.ProviderID,
	}

	result := database.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, ErrUserAlreadyExists
		}
		return nil, result.Error
	}

	return user, nil
}
