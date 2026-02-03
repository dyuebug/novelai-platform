package service

import (
	"context"
	"errors"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrProjectNotOwned = errors.New("project not owned by user")
	ErrProjectDeleted  = errors.New("project is deleted")
	ErrProjectNotDeleted = errors.New("project is not deleted")
)

type ProjectService struct {
	projectRepo *repository.ProjectRepository
}

func NewProjectService() *ProjectService {
	return &ProjectService{
		projectRepo: repository.NewProjectRepository(),
	}
}

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Title       string                 `json:"title" binding:"required,min=1,max=200"`
	Description string                 `json:"description" binding:"max=5000"`
	Genre       string                 `json:"genre"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=200"`
	Description *string `json:"description" binding:"omitempty,max=5000"`
	Genre       *string `json:"genre"`
	Status      *string `json:"status"`
	CoverURL    *string `json:"cover_url"`
}

// ProjectListResponse 项目列表响应
type ProjectListResponse struct {
	Items      []model.Project `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// Create 创建项目
func (s *ProjectService) Create(ctx context.Context, userID uuid.UUID, req *CreateProjectRequest) (*model.Project, error) {
	project := &model.Project{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Genre:       req.Genre,
		Status:      "draft",
		Metadata:    model.JSON(req.Metadata),
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// Get 获取项目
func (s *ProjectService) Get(ctx context.Context, userID, projectID uuid.UUID) (*model.Project, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 检查所有权
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	return project, nil
}

// List 获取项目列表
func (s *ProjectService) List(ctx context.Context, userID uuid.UUID, page, pageSize int, status, genre, sort string) (*ProjectListResponse, error) {
	// 限制 pageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page <= 0 {
		page = 1
	}

	opts := &repository.ListOptions{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Genre:    genre,
		Sort:     sort,
	}

	projects, total, err := s.projectRepo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &ProjectListResponse{
		Items:      projects,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update 更新项目
func (s *ProjectService) Update(ctx context.Context, userID, projectID uuid.UUID, req *UpdateProjectRequest) (*model.Project, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 检查所有权
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	// 更新字段
	if req.Title != nil {
		project.Title = *req.Title
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	if req.Genre != nil {
		project.Genre = *req.Genre
	}
	if req.Status != nil {
		project.Status = *req.Status
	}
	if req.CoverURL != nil {
		project.CoverURL = *req.CoverURL
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// Delete 删除项目 (软删除)
func (s *ProjectService) Delete(ctx context.Context, userID, projectID uuid.UUID) error {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}

	// 检查所有权
	if project.UserID != userID {
		return ErrProjectNotOwned
	}

	return s.projectRepo.SoftDelete(ctx, projectID)
}

// Restore 恢复项目
func (s *ProjectService) Restore(ctx context.Context, userID, projectID uuid.UUID) (*model.Project, error) {
	// 查找已删除的项目
	project, err := s.projectRepo.FindDeletedByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 检查所有权
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	if err := s.projectRepo.Restore(ctx, projectID); err != nil {
		return nil, err
	}

	// 重新获取项目
	return s.projectRepo.FindByID(ctx, projectID)
}
