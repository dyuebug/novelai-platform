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
	ErrWorldSettingNotFound = errors.New("world setting not found")
)

type WorldSettingRepository struct{}

func NewWorldSettingRepository() *WorldSettingRepository {
	return &WorldSettingRepository{}
}

// Create 创建世界设定
func (r *WorldSettingRepository) Create(ctx context.Context, setting *model.WorldSetting) error {
	return database.DB.WithContext(ctx).Create(setting).Error
}

// FindByID 通过 ID 查找世界设定
func (r *WorldSettingRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.WorldSetting, error) {
	var setting model.WorldSetting
	result := database.DB.WithContext(ctx).
		Where("id = ?", id).
		First(&setting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrWorldSettingNotFound
		}
		return nil, result.Error
	}
	return &setting, nil
}

// FindByProjectID 查找项目的所有世界设定
func (r *WorldSettingRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts *ListOptions) ([]model.WorldSetting, int64, error) {
	var settings []model.WorldSetting
	var total int64

	query := database.DB.WithContext(ctx).
		Model(&model.WorldSetting{}).
		Where("project_id = ?", projectID)

	// 分类过滤
	if opts.Status != "" {
		query = query.Where("category = ?", opts.Status)
	}

	// 计算总数
	query.Count(&total)

	// 排序
	if opts.Sort != "" {
		query = query.Order(opts.Sort)
	} else {
		query = query.Order("category ASC, sort_order ASC, created_at DESC")
	}

	// 分页
	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	result := query.Find(&settings)
	return settings, total, result.Error
}

// FindRootsByProjectID 查找项目的根设定（无父级）
func (r *WorldSettingRepository) FindRootsByProjectID(ctx context.Context, projectID uuid.UUID, category string) ([]model.WorldSetting, error) {
	var settings []model.WorldSetting
	query := database.DB.WithContext(ctx).
		Where("project_id = ? AND parent_id IS NULL", projectID)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	result := query.Order("sort_order ASC, created_at DESC").Find(&settings)
	return settings, result.Error
}

// FindChildren 查找子设定
func (r *WorldSettingRepository) FindChildren(ctx context.Context, parentID uuid.UUID) ([]model.WorldSetting, error) {
	var settings []model.WorldSetting
	result := database.DB.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("sort_order ASC, created_at DESC").
		Find(&settings)
	return settings, result.Error
}

// FindWithChildren 查找设定及其子设定
func (r *WorldSettingRepository) FindWithChildren(ctx context.Context, id uuid.UUID) (*model.WorldSetting, error) {
	var setting model.WorldSetting
	result := database.DB.WithContext(ctx).
		Preload("Children").
		Where("id = ?", id).
		First(&setting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrWorldSettingNotFound
		}
		return nil, result.Error
	}
	return &setting, nil
}

// FindByCategory 按分类查找
func (r *WorldSettingRepository) FindByCategory(ctx context.Context, projectID uuid.UUID, category string) ([]model.WorldSetting, error) {
	var settings []model.WorldSetting
	result := database.DB.WithContext(ctx).
		Where("project_id = ? AND category = ?", projectID, category).
		Order("sort_order ASC, created_at DESC").
		Find(&settings)
	return settings, result.Error
}

// Update 更新世界设定
func (r *WorldSettingRepository) Update(ctx context.Context, setting *model.WorldSetting) error {
	return database.DB.WithContext(ctx).Save(setting).Error
}

// Delete 删除世界设定
func (r *WorldSettingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.WorldSetting{}).Error
}
