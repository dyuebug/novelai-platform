package service

import (
	"context"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
)

// ========== 世界设定服务 ==========

type WorldSettingService struct {
	settingRepo *repository.WorldSettingRepository
	projectRepo *repository.ProjectRepository
}

func NewWorldSettingService() *WorldSettingService {
	return &WorldSettingService{
		settingRepo: repository.NewWorldSettingRepository(),
		projectRepo: repository.NewProjectRepository(),
	}
}

// CreateWorldSettingRequest 创建世界设定请求
type CreateWorldSettingRequest struct {
	ParentID *uuid.UUID             `json:"parent_id"`
	Category string                 `json:"category" binding:"required"`
	Title    string                 `json:"title" binding:"required,min=1,max=200"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
}

// UpdateWorldSettingRequest 更新世界设定请求
type UpdateWorldSettingRequest struct {
	ParentID  *uuid.UUID `json:"parent_id"`
	Category  *string    `json:"category"`
	Title     *string    `json:"title" binding:"omitempty,min=1,max=200"`
	Content   *string    `json:"content"`
	SortOrder *int       `json:"sort_order"`
}

// WorldSettingListResponse 世界设定列表响应
type WorldSettingListResponse struct {
	Items      []model.WorldSetting `json:"items"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// Create 创建世界设定
func (s *WorldSettingService) Create(ctx context.Context, userID, projectID uuid.UUID, req *CreateWorldSettingRequest) (*model.WorldSetting, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	setting := &model.WorldSetting{
		ProjectID: projectID,
		ParentID:  req.ParentID,
		Category:  req.Category,
		Title:     req.Title,
		Content:   req.Content,
		Metadata:  model.JSON(req.Metadata),
	}

	if err := s.settingRepo.Create(ctx, setting); err != nil {
		return nil, err
	}

	return setting, nil
}

// Get 获取世界设定
func (s *WorldSettingService) Get(ctx context.Context, userID, settingID uuid.UUID) (*model.WorldSetting, error) {
	setting, err := s.settingRepo.FindByID(ctx, settingID)
	if err != nil {
		return nil, err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, setting.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrWorldSettingNotOwned
	}

	return setting, nil
}

// List 获取世界设定列表
func (s *WorldSettingService) List(ctx context.Context, userID, projectID uuid.UUID, page, pageSize int, category, sort string) (*WorldSettingListResponse, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}
	if page <= 0 {
		page = 1
	}

	opts := &repository.ListOptions{
		Page:     page,
		PageSize: pageSize,
		Status:   category,
		Sort:     sort,
	}

	settings, total, err := s.settingRepo.FindByProjectID(ctx, projectID, opts)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &WorldSettingListResponse{
		Items:      settings,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetTree 获取世界设定树结构
func (s *WorldSettingService) GetTree(ctx context.Context, userID, projectID uuid.UUID, category string) ([]model.WorldSetting, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	return s.settingRepo.FindRootsByProjectID(ctx, projectID, category)
}

// GetByCategory 按分类获取设定
func (s *WorldSettingService) GetByCategory(ctx context.Context, userID, projectID uuid.UUID, category string) ([]model.WorldSetting, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	return s.settingRepo.FindByCategory(ctx, projectID, category)
}

// Update 更新世界设定
func (s *WorldSettingService) Update(ctx context.Context, userID, settingID uuid.UUID, req *UpdateWorldSettingRequest) (*model.WorldSetting, error) {
	setting, err := s.settingRepo.FindByID(ctx, settingID)
	if err != nil {
		return nil, err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, setting.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrWorldSettingNotOwned
	}

	// 更新字段
	if req.ParentID != nil {
		setting.ParentID = req.ParentID
	}
	if req.Category != nil {
		setting.Category = *req.Category
	}
	if req.Title != nil {
		setting.Title = *req.Title
	}
	if req.Content != nil {
		setting.Content = *req.Content
	}
	if req.SortOrder != nil {
		setting.SortOrder = *req.SortOrder
	}

	if err := s.settingRepo.Update(ctx, setting); err != nil {
		return nil, err
	}

	return setting, nil
}

// Delete 删除世界设定
func (s *WorldSettingService) Delete(ctx context.Context, userID, settingID uuid.UUID) error {
	setting, err := s.settingRepo.FindByID(ctx, settingID)
	if err != nil {
		return err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, setting.ProjectID)
	if err != nil {
		return err
	}
	if project.UserID != userID {
		return ErrWorldSettingNotOwned
	}

	return s.settingRepo.Delete(ctx, settingID)
}
