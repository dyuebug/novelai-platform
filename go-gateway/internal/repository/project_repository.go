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
	ErrProjectNotFound = errors.New("project not found")
	ErrChapterNotFound = errors.New("chapter not found")
	ErrVersionNotFound = errors.New("version not found")
)

type ProjectRepository struct{}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{}
}

// Create 创建项目
func (r *ProjectRepository) Create(ctx context.Context, project *model.Project) error {
	return database.DB.WithContext(ctx).Create(project).Error
}

// FindByID 通过 ID 查找项目
func (r *ProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	var project model.Project
	result := database.DB.WithContext(ctx).
		Where("id = ? AND is_deleted = ?", id, false).
		First(&project)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, result.Error
	}
	return &project, nil
}

// FindByUserID 查找用户的所有项目
func (r *ProjectRepository) FindByUserID(ctx context.Context, userID uuid.UUID, opts *ListOptions) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	query := database.DB.WithContext(ctx).
		Model(&model.Project{}).
		Where("user_id = ? AND is_deleted = ?", userID, false)

	// 状态过滤
	if opts.Status != "" {
		query = query.Where("status = ?", opts.Status)
	}

	// 类型过滤
	if opts.Genre != "" {
		query = query.Where("genre = ?", opts.Genre)
	}

	// 计算总数
	query.Count(&total)

	// 排序
	if opts.Sort != "" {
		query = query.Order(opts.Sort)
	} else {
		query = query.Order("updated_at DESC")
	}

	// 分页
	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	result := query.Find(&projects)
	return projects, total, result.Error
}

// Update 更新项目
func (r *ProjectRepository) Update(ctx context.Context, project *model.Project) error {
	return database.DB.WithContext(ctx).Save(project).Error
}

// SoftDelete 软删除项目
func (r *ProjectRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Model(&model.Project{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": gorm.Expr("NOW()"),
		}).Error
}

// Restore 恢复项目
func (r *ProjectRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Model(&model.Project{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": false,
			"deleted_at": nil,
		}).Error
}

// FindDeletedByID 查找已删除的项目
func (r *ProjectRepository) FindDeletedByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	var project model.Project
	result := database.DB.WithContext(ctx).
		Where("id = ? AND is_deleted = ?", id, true).
		First(&project)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, result.Error
	}
	return &project, nil
}

// UpdateStats 更新项目统计
func (r *ProjectRepository) UpdateStats(ctx context.Context, projectID uuid.UUID) error {
	// 更新章节数和总字数
	return database.DB.WithContext(ctx).Exec(`
		UPDATE projects SET
			total_chapters = (SELECT COUNT(*) FROM chapters WHERE project_id = ?),
			total_words = (SELECT COALESCE(SUM(word_count), 0) FROM chapters WHERE project_id = ?)
		WHERE id = ?
	`, projectID, projectID, projectID).Error
}

// ListOptions 列表选项
type ListOptions struct {
	Page     int
	PageSize int
	Status   string
	Genre    string
	Sort     string
}
