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
	ErrFileNotFound = errors.New("file not found")
)

// FileRepository 文件仓库
type FileRepository struct{}

// NewFileRepository 创建文件仓库
func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

// Create 创建文件记录
func (r *FileRepository) Create(ctx context.Context, file *model.File) error {
	return database.DB.WithContext(ctx).Create(file).Error
}

// FindByID 根据ID获取文件
func (r *FileRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.File, error) {
	var file model.File
	result := database.DB.WithContext(ctx).Where("id = ?", id).First(&file)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, result.Error
	}
	return &file, nil
}

// FindByHash 根据哈希查找文件（用于去重）
func (r *FileRepository) FindByHash(ctx context.Context, hash string) (*model.File, error) {
	var file model.File
	result := database.DB.WithContext(ctx).Where("hash = ?", hash).First(&file)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, result.Error
	}
	return &file, nil
}

// FindByUserID 获取用户的所有文件
func (r *FileRepository) FindByUserID(ctx context.Context, userID uuid.UUID, opts *ListOptions) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	query := database.DB.WithContext(ctx).
		Model(&model.File{}).
		Where("user_id = ?", userID)

	// 文件类型过滤
	if opts.Status != "" {
		query = query.Where("file_type = ?", opts.Status)
	}

	// 计算总数
	query.Count(&total)

	// 排序
	if opts.Sort != "" {
		query = query.Order(opts.Sort)
	} else {
		query = query.Order("created_at DESC")
	}

	// 分页
	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	result := query.Find(&files)
	return files, total, result.Error
}

// Update 更新文件记录
func (r *FileRepository) Update(ctx context.Context, file *model.File) error {
	return database.DB.WithContext(ctx).Save(file).Error
}

// Delete 删除文件记录（软删除）
func (r *FileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.File{}).Error
}

// HardDelete 硬删除文件记录
func (r *FileRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Unscoped().
		Where("id = ?", id).
		Delete(&model.File{}).Error
}

// CountByUserID 统计用户文件数
func (r *FileRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	result := database.DB.WithContext(ctx).
		Model(&model.File{}).
		Where("user_id = ?", userID).
		Count(&count)
	return count, result.Error
}

// GetTotalSizeByUserID 获取用户文件总大小
func (r *FileRepository) GetTotalSizeByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var totalSize int64
	result := database.DB.WithContext(ctx).
		Model(&model.File{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(file_size), 0)").
		Scan(&totalSize)
	return totalSize, result.Error
}
