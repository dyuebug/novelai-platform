package service

import (
	"context"
	"testing"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserLayoutRepository 模拟 UserLayoutRepository
type MockUserLayoutRepository struct {
	mock.Mock
}

func (m *MockUserLayoutRepository) Create(ctx context.Context, config *model.UserLayoutConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockUserLayoutRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.UserLayoutConfig, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserLayoutConfig), args.Error(1)
}

func (m *MockUserLayoutRepository) FindByUserIDAndType(ctx context.Context, userID uuid.UUID, layoutType string) (*model.UserLayoutConfig, error) {
	args := m.Called(ctx, userID, layoutType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserLayoutConfig), args.Error(1)
}

func (m *MockUserLayoutRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.UserLayoutConfig, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.UserLayoutConfig), args.Error(1)
}

func (m *MockUserLayoutRepository) Update(ctx context.Context, config *model.UserLayoutConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockUserLayoutRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserLayoutRepository) DeleteByUserIDAndType(ctx context.Context, userID uuid.UUID, layoutType string) error {
	args := m.Called(ctx, userID, layoutType)
	return args.Error(0)
}

func (m *MockUserLayoutRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// TestUserLayoutService_GetLayoutConfig 测试获取布局配置
func TestUserLayoutService_GetLayoutConfig(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	layoutType := "editor"

	t.Run("成功获取已存在的配置", func(t *testing.T) {
		mockRepo := new(MockUserLayoutRepository)
		service := &UserLayoutService{layoutRepo: mockRepo}

		expectedConfig := &model.UserLayoutConfig{
			ID:          uuid.New(),
			UserID:      userID,
			LayoutType:  layoutType,
			PanelStates: map[string]interface{}{"sidebar": "open"},
			PanelSizes:  map[string]interface{}{"sidebar": 220},
		}

		mockRepo.On("FindByUserIDAndType", ctx, userID, layoutType).Return(expectedConfig, nil)

		config, err := service.GetLayoutConfig(ctx, userID, layoutType)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, expectedConfig.ID, config.ID)
		assert.Equal(t, expectedConfig.UserID, config.UserID)
		assert.Equal(t, expectedConfig.LayoutType, config.LayoutType)
		mockRepo.AssertExpectations(t)
	})

	t.Run("配置不存在时返回默认配置", func(t *testing.T) {
		mockRepo := new(MockUserLayoutRepository)
		service := &UserLayoutService{layoutRepo: mockRepo}

		mockRepo.On("FindByUserIDAndType", ctx, userID, layoutType).Return((*model.UserLayoutConfig)(nil), repository.ErrLayoutConfigNotFound)

		config, err := service.GetLayoutConfig(ctx, userID, layoutType)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, userID, config.UserID)
		assert.Equal(t, layoutType, config.LayoutType)
		assert.NotNil(t, config.PanelStates)
		assert.NotNil(t, config.PanelSizes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("无效的布局类型", func(t *testing.T) {
		service := &UserLayoutService{}

		config, err := service.GetLayoutConfig(ctx, userID, "invalid-type")

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidLayoutType, err)
		assert.Nil(t, config)
	})
}

// TestUserLayoutService_SaveLayoutConfig 测试保存布局配置
func TestUserLayoutService_SaveLayoutConfig(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	layoutType := "editor"

	t.Run("创建新配置", func(t *testing.T) {
		mockRepo := new(MockUserLayoutRepository)
		service := &UserLayoutService{layoutRepo: mockRepo}

		req := &SaveLayoutConfigRequest{
			PanelStates: map[string]interface{}{"sidebar": "open"},
			PanelSizes:  map[string]interface{}{"sidebar": 220},
		}

		mockRepo.On("FindByUserIDAndType", ctx, userID, layoutType).Return((*model.UserLayoutConfig)(nil), repository.ErrLayoutConfigNotFound)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.UserLayoutConfig")).Return(nil)

		config, err := service.SaveLayoutConfig(ctx, userID, layoutType, req)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, userID, config.UserID)
		assert.Equal(t, layoutType, config.LayoutType)
		assert.Equal(t, req.PanelStates, config.PanelStates)
		assert.Equal(t, req.PanelSizes, config.PanelSizes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("更新已存在的配置", func(t *testing.T) {
		mockRepo := new(MockUserLayoutRepository)
		service := &UserLayoutService{layoutRepo: mockRepo}

		existingConfig := &model.UserLayoutConfig{
			ID:          uuid.New(),
			UserID:      userID,
			LayoutType:  layoutType,
			PanelStates: map[string]interface{}{"sidebar": "closed"},
			PanelSizes:  map[string]interface{}{"sidebar": 80},
		}

		req := &SaveLayoutConfigRequest{
			PanelStates: map[string]interface{}{"sidebar": "open"},
			PanelSizes:  map[string]interface{}{"sidebar": 220},
		}

		mockRepo.On("FindByUserIDAndType", ctx, userID, layoutType).Return(existingConfig, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*model.UserLayoutConfig")).Return(nil)

		config, err := service.SaveLayoutConfig(ctx, userID, layoutType, req)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, existingConfig.ID, config.ID)
		assert.Equal(t, req.PanelStates, config.PanelStates)
		assert.Equal(t, req.PanelSizes, config.PanelSizes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("无效的布局类型", func(t *testing.T) {
		service := &UserLayoutService{}

		req := &SaveLayoutConfigRequest{
			PanelStates: map[string]interface{}{},
			PanelSizes:  map[string]interface{}{},
		}

		config, err := service.SaveLayoutConfig(ctx, userID, "invalid-type", req)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidLayoutType, err)
		assert.Nil(t, config)
	})
}

// TestUserLayoutService_DeleteLayoutConfig 测试删除布局配置
func TestUserLayoutService_DeleteLayoutConfig(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	layoutType := "editor"

	t.Run("成功删除配置", func(t *testing.T) {
		mockRepo := new(MockUserLayoutRepository)
		service := &UserLayoutService{layoutRepo: mockRepo}

		existingConfig := &model.UserLayoutConfig{
			ID:         uuid.New(),
			UserID:     userID,
			LayoutType: layoutType,
		}

		mockRepo.On("FindByUserIDAndType", ctx, userID, layoutType).Return(existingConfig, nil)
		mockRepo.On("DeleteByUserIDAndType", ctx, userID, layoutType).Return(nil)

		err := service.DeleteLayoutConfig(ctx, userID, layoutType)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("配置不存在时不报错", func(t *testing.T) {
		mockRepo := new(MockUserLayoutRepository)
		service := &UserLayoutService{layoutRepo: mockRepo}

		mockRepo.On("FindByUserIDAndType", ctx, userID, layoutType).Return((*model.UserLayoutConfig)(nil), repository.ErrLayoutConfigNotFound)

		err := service.DeleteLayoutConfig(ctx, userID, layoutType)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("无效的布局类型", func(t *testing.T) {
		service := &UserLayoutService{}

		err := service.DeleteLayoutConfig(ctx, userID, "invalid-type")

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidLayoutType, err)
	})
}

// TestUserLayoutService_GetAllLayoutConfigs 测试获取所有布局配置
func TestUserLayoutService_GetAllLayoutConfigs(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("成功获取所有配置", func(t *testing.T) {
		mockRepo := new(MockUserLayoutRepository)
		service := &UserLayoutService{layoutRepo: mockRepo}

		expectedConfigs := []model.UserLayoutConfig{
			{
				ID:         uuid.New(),
				UserID:     userID,
				LayoutType: "editor",
			},
			{
				ID:         uuid.New(),
				UserID:     userID,
				LayoutType: "dashboard",
			},
		}

		mockRepo.On("FindByUserID", ctx, userID).Return(expectedConfigs, nil)

		configs, err := service.GetAllLayoutConfigs(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, configs)
		assert.Len(t, configs, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("没有配置时返回空数组", func(t *testing.T) {
		mockRepo := new(MockUserLayoutRepository)
		service := &UserLayoutService{layoutRepo: mockRepo}

		mockRepo.On("FindByUserID", ctx, userID).Return([]model.UserLayoutConfig(nil), nil)

		configs, err := service.GetAllLayoutConfigs(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, configs)
		assert.Len(t, configs, 0)
		mockRepo.AssertExpectations(t)
	})
}
