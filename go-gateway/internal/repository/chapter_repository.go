package repository

import (
	"context"
	"errors"

	"go-gateway/internal/database"
	"go-gateway/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChapterRepository struct{}

func NewChapterRepository() *ChapterRepository {
	return &ChapterRepository{}
}

// Create 创建章节
func (r *ChapterRepository) Create(ctx context.Context, chapter *model.Chapter) error {
	return database.DB.WithContext(ctx).Create(chapter).Error
}

// FindByID 通过 ID 查找章节
func (r *ChapterRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Chapter, error) {
	var chapter model.Chapter
	result := database.DB.WithContext(ctx).
		Where("id = ?", id).
		First(&chapter)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrChapterNotFound
		}
		return nil, result.Error
	}
	return &chapter, nil
}

// FindByProjectID 查找项目的所有章节
func (r *ChapterRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts *ListOptions) ([]model.Chapter, int64, error) {
	var chapters []model.Chapter
	var total int64

	query := database.DB.WithContext(ctx).
		Model(&model.Chapter{}).
		Where("project_id = ?", projectID)

	// 状态过滤
	if opts.Status != "" {
		query = query.Where("status = ?", opts.Status)
	}

	// 计算总数
	query.Count(&total)

	// 排序
	if opts.Sort != "" {
		query = query.Order(opts.Sort)
	} else {
		query = query.Order("chapter_number ASC")
	}

	// 分页
	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	result := query.Find(&chapters)
	return chapters, total, result.Error
}

// Update 更新章节
func (r *ChapterRepository) Update(ctx context.Context, chapter *model.Chapter) error {
	return database.DB.WithContext(ctx).Save(chapter).Error
}

// Delete 删除章节
func (r *ChapterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Chapter{}).Error
}

// GetMaxChapterNumber 获取项目中最大的章节号
func (r *ChapterRepository) GetMaxChapterNumber(ctx context.Context, projectID uuid.UUID) (int, error) {
	var maxNum int
	result := database.DB.WithContext(ctx).
		Model(&model.Chapter{}).
		Where("project_id = ?", projectID).
		Select("COALESCE(MAX(chapter_number), 0)").
		Scan(&maxNum)
	return maxNum, result.Error
}

// UpdateChapterNumbers 批量更新章节号 (用于重排序)
func (r *ChapterRepository) UpdateChapterNumbers(ctx context.Context, updates []ChapterNumberUpdate) error {
	return database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, u := range updates {
			if err := tx.Model(&model.Chapter{}).
				Where("id = ?", u.ID).
				Update("chapter_number", u.ChapterNumber).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ChapterNumberUpdate 章节号更新
type ChapterNumberUpdate struct {
	ID            uuid.UUID
	ChapterNumber int
}

// CountByProjectID 统计项目章节数
func (r *ChapterRepository) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	var count int64
	result := database.DB.WithContext(ctx).
		Model(&model.Chapter{}).
		Where("project_id = ?", projectID).
		Count(&count)
	return count, result.Error
}
