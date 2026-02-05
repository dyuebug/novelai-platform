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

// ProjectStatistics 项目统计信息
type ProjectStatistics struct {
	CharacterCount  int64
	LocationCount   int64
	ForeshadowCount int64
}

// GetProjectStatistics 获取项目统计信息
func (r *ProjectRepository) GetProjectStatistics(ctx context.Context, projectID uuid.UUID) (*ProjectStatistics, error) {
	stats := &ProjectStatistics{}

	// 统计角色数量
	var characterCount int64
	if err := database.DB.WithContext(ctx).
		Model(&model.Character{}).
		Where("project_id = ?", projectID).
		Count(&characterCount).Error; err != nil {
		return nil, err
	}
	stats.CharacterCount = characterCount

	// 统计地点数量
	var locationCount int64
	if err := database.DB.WithContext(ctx).
		Model(&model.Location{}).
		Where("project_id = ?", projectID).
		Count(&locationCount).Error; err != nil {
		return nil, err
	}
	stats.LocationCount = locationCount

	// 统计伏笔数量
	var foreshadowCount int64
	if err := database.DB.WithContext(ctx).
		Model(&model.Foreshadow{}).
		Where("project_id = ?", projectID).
		Count(&foreshadowCount).Error; err != nil {
		return nil, err
	}
	stats.ForeshadowCount = foreshadowCount

	return stats, nil
}

// Search 搜索项目
func (r *ProjectRepository) Search(ctx context.Context, userID uuid.UUID, query string, page, pageSize int) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	// 使用 ILIKE 进行模糊搜索
	searchPattern := "%" + query + "%"

	queryDB := database.DB.WithContext(ctx).
		Model(&model.Project{}).
		Where("user_id = ? AND is_deleted = ?", userID, false).
		Where("title ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)

	// 计算总数
	queryDB.Count(&total)

	// 分页
	offset := (page - 1) * pageSize
	result := queryDB.
		Order("updated_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&projects)

	return projects, total, result.Error
}
