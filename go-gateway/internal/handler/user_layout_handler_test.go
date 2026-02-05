package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/model"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserLayoutService 是 UserLayoutService 的 mock
type MockUserLayoutService struct {
	mock.Mock
}

func (m *MockGetLayoutConfig) GetLayoutConfig(ctx context.Context, userID uuid.UUID, layoutType string) (*model.UserLayoutConfig, error) {
	args := m.Called(ctx, userID, layoutType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserLayoutConfig), args.Error(1)
}

func (m *MockSaveLayoutConfig) SaveLayoutConfig(ctx context.Context, userID uuid.UUID, layoutType string, req *service.SaveLayoutConfigRequest) (*model.UserLayoutConfig, error) {
	args := m.Called(ctx, userID, layoutType, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserLayoutConfig), args.Error(1)
}

func (m *MockDeleteLayoutConfig) DeleteLayoutConfig(ctx context.Context, userID uuid.UUID, layoutType string) error {
	args := m.Called(ctx, userID, layoutType)
	return args.Error(0)
}

func (m *MockGetAllLayoutConfigs) GetAllLayoutConfigs(ctx context.Context, userID uuid.UUID) ([]model.UserLayoutConfig, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.UserLayoutConfig), args.Error(1)
}

func (m *MockResetLayoutConfig) ResetLayoutConfig(ctx context.Context, userID uuid.UUID, layoutType string) (*model.UserLayoutConfig, error) {
	args := m.Called(ctx, userID, layoutType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserLayoutConfig), args.Error(1)
}

// setupTestRouter 设置测试路由
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// setUserIDContext 设置用户 ID 到上下文（模拟认证中间件）
func setUserIDContext(c *gin.Context, userID uuid.UUID) {
	claims := map[string]interface{}{
		"sub": userID.String(),
	}
	c.Set("claims", claims)
}

// TestUserLayoutHandler_GetLayoutConfig 测试获取布局配置 API
func TestUserLayoutHandler_GetLayoutConfig(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		layoutType     string
		setupAuth      bool
		mockSetup      func(*MockUserLayoutService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:       "成功获取布局配置",
			layoutType: "editor",
			setupAuth:  true,
			mockSetup: func(m *MockUserLayoutService) {
				config := &model.UserLayoutConfig{
					ID:          uuid.New(),
					UserID:      userID,
					LayoutType:  "editor",
					PanelStates: map[string]interface{}{"theme": "dark"},
				}
				m.On("GetLayoutConfig", mock.Anything, userID, "editor").Return(config, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			layoutType:     "editor",
			setupAuth:      false,
			mockSetup:      func(m *MockUserLayoutService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:       "无效的布局类型",
			layoutType: "invalid",
			setupAuth:  true,
			mockSetup: func(m *MockUserLayoutService) {
				m.On("GetLayoutConfig", mock.Anything, userID, "invalid").Return(nil, service.ErrInvalidLayoutType)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockUserLayoutService)
			tt.mockSetup(mockService)

			handler := NewUserLayoutHandler(mockService)
			router := setupTestRouter()

			// 注册路由
			router.GET("/api/v1/user/layout-config/:layoutType", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.GetLayoutConfig(c)
			})

			// 创建请求
			req := httptest.NewRequest(http.MethodGet, "/api/v1/user/layout-config/"+tt.layoutType, nil)
			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp response.Envelope
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.NotEqual(t, 0, resp.Code)
			} else {
				assert.Equal(t, 0, resp.Code)
				assert.NotNil(t, resp.Data)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestUserLayoutHandler_SaveLayoutConfig 测试保存布局配置 API
func TestUserLayoutHandler_SaveLayoutConfig(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		layoutType     string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockUserLayoutService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:       "成功保存布局配置",
			layoutType: "editor",
			requestBody: map[string]interface{}{
				"config": map[string]interface{}{"theme": "dark"},
			},
			setupAuth: true,
			mockSetup: func(m *MockUserLayoutService) {
				config := &model.UserLayoutConfig{
					ID:          uuid.New(),
					UserID:      userID,
					LayoutType:  "editor",
					PanelStates: map[string]interface{}{"theme": "dark"},
				}
				m.On("SaveLayoutConfig", mock.Anything, userID, "editor", mock.AnythingOfType("*service.SaveLayoutConfigRequest")).Return(config, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "无效的请求体",
			layoutType:     "editor",
			requestBody:    "invalid json",
			setupAuth:      true,
			mockSetup:      func(m *MockUserLayoutService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:       "未认证",
			layoutType: "editor",
			requestBody: map[string]interface{}{
				"config": map[string]interface{}{"theme": "dark"},
			},
			setupAuth:      false,
			mockSetup:      func(m *MockUserLayoutService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockUserLayoutService)
			tt.mockSetup(mockService)

			handler := NewUserLayoutHandler(mockService)
			router := setupTestRouter()

			// 注册路由
			router.POST("/api/v1/user/layout-config/:layoutType", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.SaveLayoutConfig(c)
			})

			// 创建请求
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/layout-config/"+tt.layoutType, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp response.Envelope
			err = json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.NotEqual(t, 0, resp.Code)
			} else {
				assert.Equal(t, 0, resp.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestUserLayoutHandler_DeleteLayoutConfig 测试删除布局配置 API
func TestUserLayoutHandler_DeleteLayoutConfig(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		layoutType     string
		setupAuth      bool
		mockSetup      func(*MockUserLayoutService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:       "成功删除布局配置",
			layoutType: "editor",
			setupAuth:  true,
			mockSetup: func(m *MockUserLayoutService) {
				m.On("DeleteLayoutConfig", mock.Anything, userID, "editor").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			layoutType:     "editor",
			setupAuth:      false,
			mockSetup:      func(m *MockUserLayoutService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:       "无效的布局类型",
			layoutType: "invalid",
			setupAuth:  true,
			mockSetup: func(m *MockUserLayoutService) {
				m.On("DeleteLayoutConfig", mock.Anything, userID, "invalid").Return(service.ErrInvalidLayoutType)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockUserLayoutService)
			tt.mockSetup(mockService)

			handler := NewUserLayoutHandler(mockService)
			router := setupTestRouter()

			// 注册路由
			router.DELETE("/api/v1/user/layout-config/:layoutType", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.DeleteLayoutConfig(c)
			})

			// 创建请求
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/user/layout-config/"+tt.layoutType, nil)
			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp response.Envelope
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.NotEqual(t, 0, resp.Code)
			} else {
				assert.Equal(t, 0, resp.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestUserLayoutHandler_GetAllLayoutConfigs 测试获取所有布局配置 API
func TestUserLayoutHandler_GetAllLayoutConfigs(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		setupAuth      bool
		mockSetup      func(*MockUserLayoutService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取所有布局配置",
			setupAuth: true,
			mockSetup: func(m *MockUserLayoutService) {
				configs := []model.UserLayoutConfig{
					{
						ID:          uuid.New(),
						UserID:      userID,
						LayoutType:  "editor",
						PanelStates: map[string]interface{}{"theme": "dark"},
					},
					{
						ID:          uuid.New(),
						UserID:      userID,
						LayoutType:  "dashboard",
						PanelStates: map[string]interface{}{"layout": "grid"},
					},
				}
				m.On("GetAllLayoutConfigs", mock.Anything, userID).Return(configs, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			setupAuth:      false,
			mockSetup:      func(m *MockUserLayoutService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockUserLayoutService)
			tt.mockSetup(mockService)

			handler := NewUserLayoutHandler(mockService)
			router := setupTestRouter()

			// 注册路由
			router.GET("/api/v1/user/layout-configs", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.GetAllLayoutConfigs(c)
			})

			// 创建请求
			req := httptest.NewRequest(http.MethodGet, "/api/v1/user/layout-configs", nil)
			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp response.Envelope
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.NotEqual(t, 0, resp.Code)
			} else {
				assert.Equal(t, 0, resp.Code)
				assert.NotNil(t, resp.Data)
			}

			mockService.AssertExpectations(t)
		})
	}
}
