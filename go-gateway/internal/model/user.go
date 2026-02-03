package model

import (
	"time"

	"github.com/google/uuid"
)

// User 用户模型
type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username     string     `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string     `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"size:255" json:"-"`
	AvatarURL    string     `gorm:"size:500" json:"avatar_url,omitempty"`
	Role         string     `gorm:"size:20;default:user" json:"role"`
	Status       string     `gorm:"size:20;default:active" json:"status"`
	OAuthProvider string    `gorm:"size:50" json:"oauth_provider,omitempty"`
	OAuthID      string     `gorm:"size:255" json:"oauth_id,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`

	// 账户锁定
	FailedAttempts int        `gorm:"default:0" json:"-"`
	LockedUntil    *time.Time `json:"-"`
}

// UserSettings 用户设置模型
type UserSettings struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	DefaultProvider  string    `gorm:"size:50;default:openai" json:"default_provider"`
	DefaultModel     string    `gorm:"size:100" json:"default_model,omitempty"`
	OpenAIAPIKey     string    `gorm:"size:255" json:"-"` // 加密存储
	OpenAIBaseURL    string    `gorm:"size:500" json:"openai_base_url,omitempty"`
	AnthropicAPIKey  string    `gorm:"size:255" json:"-"`
	GeminiAPIKey     string    `gorm:"size:255" json:"-"`
	DefaultTemp      float64   `gorm:"type:decimal(3,2);default:0.7" json:"default_temperature"`
	DefaultMaxTokens int       `gorm:"default:4096" json:"default_max_tokens"`
	Theme            string    `gorm:"size:20;default:light" json:"theme"`
	Language         string    `gorm:"size:10;default:zh-CN" json:"language"`
	EditorFontSize   int       `gorm:"default:16" json:"editor_font_size"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// RefreshToken 刷新令牌模型
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	Token     string    `gorm:"size:255;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	RevokedAt *time.Time

	User User `gorm:"foreignKey:UserID"`
}

// PasswordResetToken 密码重置令牌
type PasswordResetToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	Token     string    `gorm:"size:255;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	UsedAt    *time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID"`
}

// IsLocked 检查账户是否被锁定
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (UserSettings) TableName() string {
	return "user_settings"
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
