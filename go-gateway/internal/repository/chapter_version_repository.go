package repository

import (
	"context"
	"errors"

	"go-gateway/internal/database"
	"go-gateway/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChapterVersionRepository struct{}

func NewChapterVersionRepository() *ChapterVersionRepository {
	return &ChapterVersionRepository{}
}

// Create 创建版本
func (r *ChapterVersionRepository) Create(ctx context.Context, version *model.ChapterVersion) error {
	return database.DB.WithContext(ctx).Create(version).Error
}

// FindByChapterID 查找章节的所有版本
func (r *ChapterVersionRepository) FindByChapterID(ctx context.Context, chapterID uuid.UUID, opts *ListOptions) ([]model.ChapterVersion, int64, error) {
	var versions []model.ChapterVersion
	var total int64

	query := database.DB.WithContext(ctx).
		Model(&model.ChapterVersion{}).
		Where("chapter_id = ?", chapterID)

	// 计算总数
	query.Count(&total)

	// 排序 (默认按版本号降序)
	if opts.Sort != "" {
		query = query.Order(opts.Sort)
	} else {
		query = query.Order("version_number DESC")
	}

	// 分页
	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	result := query.Find(&versions)
	return versions, total, result.Error
}

// FindByVersionNumber 通过版本号查找
func (r *ChapterVersionRepository) FindByVersionNumber(ctx context.Context, chapterID uuid.UUID, versionNumber int) (*model.ChapterVersion, error) {
	var version model.ChapterVersion
	result := database.DB.WithContext(ctx).
		Where("chapter_id = ? AND version_number = ?", chapterID, versionNumber).
		First(&version)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrVersionNotFound
		}
		return nil, result.Error
	}
	return &version, nil
}

// GetMaxVersionNumber 获取章节的最大版本号
func (r *ChapterVersionRepository) GetMaxVersionNumber(ctx context.Context, chapterID uuid.UUID) (int, error) {
	var maxNum int
	result := database.DB.WithContext(ctx).
		Model(&model.ChapterVersion{}).
		Where("chapter_id = ?", chapterID).
		Select("COALESCE(MAX(version_number), 0)").
		Scan(&maxNum)
	return maxNum, result.Error
}

// GetTwoVersions 获取两个版本用于对比
func (r *ChapterVersionRepository) GetTwoVersions(ctx context.Context, chapterID uuid.UUID, v1, v2 int) (*model.ChapterVersion, *model.ChapterVersion, error) {
	version1, err := r.FindByVersionNumber(ctx, chapterID, v1)
	if err != nil {
		return nil, nil, err
	}

	version2, err := r.FindByVersionNumber(ctx, chapterID, v2)
	if err != nil {
		return nil, nil, err
	}

	return version1, version2, nil
}
