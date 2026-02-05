package service

import (
	"context"
	"errors"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrLayoutConfigNotOwned = errors.New("layout config not owned by user")
	ErrInvalidLayoutType    = errors.New("invalid layout type")
)

// UserLayoutRepository 接口
type UserLayoutRepository interface {
	Create(ctx context.Context, config *model.UserLayoutConfig) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.UserLayoutConfig, error)
	FindByUserIDAndType(ctx context.Context, userID uuid.UUID, layoutType string) (*model.UserLayoutConfig, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.UserLayoutConfig, error)
	Update(ctx context.Context, config *model.UserLayoutConfig) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserIDAndType(ctx context.Context, userID uuid.UUID, layoutType string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

// 支持的布局类型
var validLayoutTypes = map[string]bool{
	"editor":        true,
	"dashboard":     true,
	"world-builder": true,
}

// UserLayoutService 用户布局配置服务
type UserLayoutService struct {
	layoutRepo UserLayoutRepository
}

// NewUserLayoutService 创建用户布局配置服务
func NewUserLayoutService() *UserLayoutService {
	return &UserLayoutService{
		layoutRepo: repository.NewUserLayoutRepository(),
	}
}

// SaveLayoutConfigRequest 保存布局配置请求
type SaveLayoutConfigRequest struct {
	PanelStates map[string]interface{} `json:"panel_states"`
	PanelSizes  map[string]interface{} `json:"panel_sizes"`
}

// LayoutConfigResponse 布局配置响应
type LayoutConfigResponse struct {
	ID          uuid.UUID              `json:"id"`
	UserID      uuid.UUID              `json:"user_id"`
	LayoutType  string                 `json:"layout_type"`
	PanelStates map[string]interface{} `json:"panel_states"`
	PanelSizes  map[string]interface{} `json:"panel_sizes"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// GetLayoutConfig 获取布局配置
func (s *UserLayoutService) GetLayoutConfig(ctx context.Context, userID uuid.UUID, layoutType string) (*model.UserLayoutConfig, error) {
	// 验证布局类型
	if !validLayoutTypes[layoutType] {
		return nil, ErrInvalidLayoutType
	}

	config, err := s.layoutRepo.FindByUserIDAndType(ctx, userID, layoutType)
	if err != nil {
		if errors.Is(err, repository.ErrLayoutConfigNotFound) {
			// 返回默认配置
			return &model.UserLayoutConfig{
				UserID:      userID,
				LayoutType:  layoutType,
				PanelStates: make(map[string]interface{}),
				PanelSizes:  make(map[string]interface{}),
			}, nil
		}
		return nil, err
	}

	return config, nil
}

// SaveLayoutConfig 保存布局配置（创建或更新）
func (s *UserLayoutService) SaveLayoutConfig(ctx context.Context, userID uuid.UUID, layoutType string, req *SaveLayoutConfigRequest) (*model.UserLayoutConfig, error) {
	// 验证布局类型
	if !validLayoutTypes[layoutType] {
		return nil, ErrInvalidLayoutType
	}

	// 初始化空 map（避免 nil）
	if req.PanelStates == nil {
		req.PanelStates = make(map[string]interface{})
	}
	if req.PanelSizes == nil {
		req.PanelSizes = make(map[string]interface{})
	}

	// 尝试查找现有配置
	existingConfig, err := s.layoutRepo.FindByUserIDAndType(ctx, userID, layoutType)
	if err != nil && !errors.Is(err, repository.ErrLayoutConfigNotFound) {
		return nil, err
	}

	if existingConfig != nil {
		// 更新现有配置
		existingConfig.PanelStates = req.PanelStates
		existingConfig.PanelSizes = req.PanelSizes

		if err := s.layoutRepo.Update(ctx, existingConfig); err != nil {
			return nil, err
		}

		return existingConfig, nil
	}

	// 创建新配置
	config := &model.UserLayoutConfig{
		UserID:      userID,
		LayoutType:  layoutType,
		PanelStates: req.PanelStates,
		PanelSizes:  req.PanelSizes,
	}

	if err := s.layoutRepo.Create(ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetAllLayoutConfigs 获取用户的所有布局配置
func (s *UserLayoutService) GetAllLayoutConfigs(ctx context.Context, userID uuid.UUID) ([]model.UserLayoutConfig, error) {
	configs, err := s.layoutRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 如果没有配置，返回空数组而不是 nil
	if configs == nil {
		configs = []model.UserLayoutConfig{}
	}

	return configs, nil
}

// DeleteLayoutConfig 删除布局配置
func (s *UserLayoutService) DeleteLayoutConfig(ctx context.Context, userID uuid.UUID, layoutType string) error {
	// 验证布局类型
	if !validLayoutTypes[layoutType] {
		return ErrInvalidLayoutType
	}

	// 验证配置所有权
	config, err := s.layoutRepo.FindByUserIDAndType(ctx, userID, layoutType)
	if err != nil {
		if errors.Is(err, repository.ErrLayoutConfigNotFound) {
			// 配置不存在，视为删除成功
			return nil
		}
		return err
	}

	if config.UserID != userID {
		return ErrLayoutConfigNotOwned
	}

	return s.layoutRepo.DeleteByUserIDAndType(ctx, userID, layoutType)
}

// ResetLayoutConfig 重置布局配置为默认值
func (s *UserLayoutService) ResetLayoutConfig(ctx context.Context, userID uuid.UUID, layoutType string) (*model.UserLayoutConfig, error) {
	// 验证布局类型
	if !validLayoutTypes[layoutType] {
		return nil, ErrInvalidLayoutType
	}

	// 删除现有配置
	if err := s.DeleteLayoutConfig(ctx, userID, layoutType); err != nil {
		return nil, err
	}

	// 返回默认配置
	return &model.UserLayoutConfig{
		UserID:      userID,
		LayoutType:  layoutType,
		PanelStates: make(map[string]interface{}),
		PanelSizes:  make(map[string]interface{}),
	}, nil
}

// DeleteAllLayoutConfigs 删除用户的所有布局配置
func (s *UserLayoutService) DeleteAllLayoutConfigs(ctx context.Context, userID uuid.UUID) error {
	return s.layoutRepo.DeleteByUserID(ctx, userID)
}
