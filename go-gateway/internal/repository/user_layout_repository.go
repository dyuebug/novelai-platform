package repository

import (
	"context"
	"errors"

	"go-gateway/internal/database"
	"go-gateway/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrLayoutConfigNotFound = errors.New("layout config not found")
)

// UserLayoutRepository 用户布局配置仓库
type UserLayoutRepository struct{}

// NewUserLayoutRepository 创建用户布局配置仓库
func NewUserLayoutRepository() *UserLayoutRepository {
	return &UserLayoutRepository{}
}

// Create 创建布局配置
func (r *UserLayoutRepository) Create(ctx context.Context, config *model.UserLayoutConfig) error {
	return database.DB.WithContext(ctx).Create(config).Error
}

// FindByID 根据ID获取布局配置
func (r *UserLayoutRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.UserLayoutConfig, error) {
	var config model.UserLayoutConfig
	result := database.DB.WithContext(ctx).Where("id = ?", id).First(&config)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrLayoutConfigNotFound
		}
		return nil, result.Error
	}
	return &config, nil
}

// FindByUserIDAndType 根据用户ID和布局类型获取配置
func (r *UserLayoutRepository) FindByUserIDAndType(ctx context.Context, userID uuid.UUID, layoutType string) (*model.UserLayoutConfig, error) {
	var config model.UserLayoutConfig
	result := database.DB.WithContext(ctx).
		Where("user_id = ? AND layout_type = ?", userID, layoutType).
		First(&config)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrLayoutConfigNotFound
		}
		return nil, result.Error
	}
	return &config, nil
}

// FindByUserID 获取用户的所有布局配置
func (r *UserLayoutRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.UserLayoutConfig, error) {
	var configs []model.UserLayoutConfig
	result := database.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("layout_type ASC").
		Find(&configs)
	return configs, result.Error
}

// Update 更新布局配置
func (r *UserLayoutRepository) Update(ctx context.Context, config *model.UserLayoutConfig) error {
	return database.DB.WithContext(ctx).Save(config).Error
}

// Upsert 创建或更新布局配置（根据 user_id + layout_type 唯一约束）
func (r *UserLayoutRepository) Upsert(ctx context.Context, config *model.UserLayoutConfig) error {
	// 使用原生 SQL 实现 ON CONFLICT DO UPDATE
	return database.DB.WithContext(ctx).Exec(`
		INSERT INTO user_layout_configs (id, user_id, layout_type, panel_states, panel_sizes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, layout_type)
		DO UPDATE SET
			panel_states = EXCLUDED.panel_states,
			panel_sizes = EXCLUDED.panel_sizes,
			updated_at = CURRENT_TIMESTAMP
	`, config.ID, config.UserID, config.LayoutType, config.PanelStates, config.PanelSizes).Error
}

// Delete 删除布局配置
func (r *UserLayoutRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.UserLayoutConfig{}).Error
}

// DeleteByUserIDAndType 根据用户ID和布局类型删除配置
func (r *UserLayoutRepository) DeleteByUserIDAndType(ctx context.Context, userID uuid.UUID, layoutType string) error {
	return database.DB.WithContext(ctx).
		Where("user_id = ? AND layout_type = ?", userID, layoutType).
		Delete(&model.UserLayoutConfig{}).Error
}

// DeleteByUserID 删除用户的所有布局配置
func (r *UserLayoutRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.UserLayoutConfig{}).Error
}
