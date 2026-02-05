package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserLayoutConfig 用户布局配置模型
type UserLayoutConfig struct {
	ID         uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID              `gorm:"type:uuid;uniqueIndex:idx_user_layout;not null;column:user_id" json:"user_id"`
	LayoutType string                 `gorm:"size:50;uniqueIndex:idx_user_layout;not null;column:layout_type" json:"layout_type"` // editor, dashboard, world-builder
	PanelStates map[string]interface{} `gorm:"type:jsonb;column:panel_states" json:"panel_states"` // 面板状态
	PanelSizes  map[string]interface{} `gorm:"type:jsonb;column:panel_sizes" json:"panel_sizes"`   // 面板尺寸
	CreatedAt  time.Time              `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time              `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (UserLayoutConfig) TableName() string {
	return "user_layout_configs"
}

// BeforeCreate GORM钩子：创建前设置UUID
func (u *UserLayoutConfig) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
