package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWorldSettingService 是 WorldSettingService 的 mock
type MockWorldSettingService struct {
	mock.Mock
}

func (m *MockCreate) Create(ctx context.Context, userID, projectID uuid.UUID, req *service.CreateWorldSettingRequest) (*model.WorldSetting, error) {
	args := m.Called(ctx, userID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WorldSetting), args.Error(1)
}

func (m *MockList) List(ctx context.Context, userID, projectID uuid.UUID, page, pageSize int, category, sort string) (*service.WorldSettingListResponse, error) {
	args := m.Called(ctx, userID, projectID, page, pageSize, category, sort)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.WorldSettingListResponse), args.Error(1)
}

func (m *MockGetTree) GetTree(ctx context.Context, userID, projectID uuid.UUID, category string) ([]model.WorldSetting, error) {
	args := m.Called(ctx, userID, projectID, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.WorldSetting), args.Error(1)
}

func (m *MockGetByCategory) GetByCategory(ctx context.Context, userID, projectID uuid.UUID, category string) ([]model.WorldSetting, error) {
	args := m.Called(ctx, userID, projectID, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.WorldSetting), args.Error(1)
}

func (m *MockGet) Get(ctx context.Context, userID, settingID uuid.UUID) (*model.WorldSetting, error) {
	args := m.Called(ctx, userID, settingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WorldSetting), args.Error(1)
}

func (m *MockUpdate) Update(ctx context.Context, userID, settingID uuid.UUID, req *service.UpdateWorldSettingRequest) (*model.WorldSetting, error) {
	args := m.Called(ctx, userID, settingID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WorldSetting), args.Error(1)
}

func (m *MockDelete) Delete(ctx context.Context, userID, settingID uuid.UUID) error {
	args := m.Called(ctx, userID, settingID)
	return args.Error(0)
}

// TestWorldSettingHandler_Create 测试创建世界设定 API
func TestWorldSettingHandler_Create(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()
	settingID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		projectID      string
		requestBody    interface{}
		mockSetup      func(*MockWorldSettingService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功创建世界设定",
			userID:    userID,
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"title":    "魔法体系",
				"category": "magic",
				"content":  "世界的魔法体系设定",
			},
			mockSetup: func(m *MockWorldSettingService) {
				setting := &model.WorldSetting{
					ID:        settingID,
					ProjectID: projectID,
					Title:     "魔法体系",
					Category:  "magic",
					Content:   "世界的魔法体系设定",
				}
				m.On("Create", mock.Anything, userID, projectID, mock.AnythingOfType("*service.CreateWorldSettingRequest")).Return(setting, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			requestBody:    map[string]interface{}{"name": "test"},
			mockSetup:      func(m *MockWorldSettingService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的项目ID",
			userID:         userID,
			projectID:      "invalid-uuid",
			requestBody:    map[string]interface{}{"name": "test"},
			mockSetup:      func(m *MockWorldSettingService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			userID:    userID,
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"title":    "test",
				"category": "magic",
			},
			mockSetup: func(m *MockWorldSettingService) {
				m.On("Create", mock.Anything, userID, projectID, mock.AnythingOfType("*service.CreateWorldSettingRequest")).Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWorldSettingService)
			tt.mockSetup(mockService)

			handler := &WorldSettingHandler{
				settingService: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/projects/:id/world-settings", func(c *gin.Context) {
				if tt.userID != uuid.Nil {
					setUserIDContext(c, tt.userID)
				}
				handler.Create(c)
			})

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			url := "/api/v1/projects/" + tt.projectID + "/world-settings"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp response.Envelope
			err = json.Unmarshal(w.Body.Bytes(), &resp)
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

// TestWorldSettingHandler_List 测试获取世界设定列表 API
func TestWorldSettingHandler_List(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		projectID      string
		query          string
		mockSetup      func(*MockWorldSettingService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取世界设定列表",
			userID:    userID,
			projectID: projectID.String(),
			query:     "page=1&page_size=20&category=magic",
			mockSetup: func(m *MockWorldSettingService) {
				result := &service.WorldSettingListResponse{
					Items: []model.WorldSetting{
						{
							ID:       uuid.New(),
							Title:    "魔法体系",
							Category: "magic",
						},
					},
					Total:    1,
					Page:     1,
					PageSize: 20,
				}
				m.On("List", mock.Anything, userID, projectID, 1, 20, "magic", "").Return(result, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			mockSetup:      func(m *MockWorldSettingService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			userID:    userID,
			projectID: projectID.String(),
			mockSetup: func(m *MockWorldSettingService) {
				m.On("List", mock.Anything, userID, projectID, 1, 50, "", "").Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWorldSettingService)
			tt.mockSetup(mockService)

			handler := &WorldSettingHandler{
				settingService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/projects/:id/world-settings", func(c *gin.Context) {
				if tt.userID != uuid.Nil {
					setUserIDContext(c, tt.userID)
				}
				handler.List(c)
			})

			url := "/api/v1/projects/" + tt.projectID + "/world-settings"
			if tt.query != "" {
				url += "?" + tt.query
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestWorldSettingHandler_GetTree 测试获取世界设定树 API
func TestWorldSettingHandler_GetTree(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		projectID      string
		query          string
		mockSetup      func(*MockWorldSettingService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取世界设定树",
			userID:    userID,
			projectID: projectID.String(),
			query:     "category=magic",
			mockSetup: func(m *MockWorldSettingService) {
				tree := []model.WorldSetting{
					{
						ID:       uuid.New(),
						Title:    "魔法体系",
						Category: "magic",
					},
				}
				m.On("GetTree", mock.Anything, userID, projectID, "magic").Return(tree, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			mockSetup:      func(m *MockWorldSettingService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			userID:    userID,
			projectID: projectID.String(),
			mockSetup: func(m *MockWorldSettingService) {
				m.On("GetTree", mock.Anything, userID, projectID, "").Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWorldSettingService)
			tt.mockSetup(mockService)

			handler := &WorldSettingHandler{
				settingService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/projects/:id/world-settings/tree", func(c *gin.Context) {
				if tt.userID != uuid.Nil {
					setUserIDContext(c, tt.userID)
				}
				handler.GetTree(c)
			})

			url := "/api/v1/projects/" + tt.projectID + "/world-settings/tree"
			if tt.query != "" {
				url += "?" + tt.query
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestWorldSettingHandler_GetByCategory 测试按分类获取设定 API
func TestWorldSettingHandler_GetByCategory(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		projectID      string
		category       string
		mockSetup      func(*MockWorldSettingService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功按分类获取设定",
			userID:    userID,
			projectID: projectID.String(),
			category:  "magic",
			mockSetup: func(m *MockWorldSettingService) {
				settings := []model.WorldSetting{
					{
						ID:       uuid.New(),
						Title:    "魔法体系",
						Category: "magic",
					},
				}
				m.On("GetByCategory", mock.Anything, userID, projectID, "magic").Return(settings, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			category:       "magic",
			mockSetup:      func(m *MockWorldSettingService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			userID:    userID,
			projectID: projectID.String(),
			category:  "magic",
			mockSetup: func(m *MockWorldSettingService) {
				m.On("GetByCategory", mock.Anything, userID, projectID, "magic").Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWorldSettingService)
			tt.mockSetup(mockService)

			handler := &WorldSettingHandler{
				settingService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/projects/:id/world-settings/category/:category", func(c *gin.Context) {
				if tt.userID != uuid.Nil {
					setUserIDContext(c, tt.userID)
				}
				handler.GetByCategory(c)
			})

			url := "/api/v1/projects/" + tt.projectID + "/world-settings/category/" + tt.category
			req := httptest.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestWorldSettingHandler_Get 测试获取单个世界设定 API
func TestWorldSettingHandler_Get(t *testing.T) {
	userID := uuid.New()
	settingID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		settingID      string
		mockSetup      func(*MockWorldSettingService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取世界设定",
			userID:    userID,
			settingID: settingID.String(),
			mockSetup: func(m *MockWorldSettingService) {
				setting := &model.WorldSetting{
					ID:       settingID,
					Title:    "魔法体系",
					Category: "magic",
				}
				m.On("Get", mock.Anything, userID, settingID).Return(setting, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			settingID:      settingID.String(),
			mockSetup:      func(m *MockWorldSettingService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的设定ID",
			userID:         userID,
			settingID:      "invalid-uuid",
			mockSetup:      func(m *MockWorldSettingService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "世界设定不存在",
			userID:    userID,
			settingID: settingID.String(),
			mockSetup: func(m *MockWorldSettingService) {
				m.On("Get", mock.Anything, userID, settingID).Return(nil, repository.ErrWorldSettingNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWorldSettingService)
			tt.mockSetup(mockService)

			handler := &WorldSettingHandler{
				settingService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/world-settings/:id", func(c *gin.Context) {
				if tt.userID != uuid.Nil {
					setUserIDContext(c, tt.userID)
				}
				handler.Get(c)
			})

			url := "/api/v1/world-settings/" + tt.settingID
			req := httptest.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestWorldSettingHandler_Update 测试更新世界设定 API
func TestWorldSettingHandler_Update(t *testing.T) {
	userID := uuid.New()
	settingID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		settingID      string
		requestBody    interface{}
		mockSetup      func(*MockWorldSettingService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功更新世界设定",
			userID:    userID,
			settingID: settingID.String(),
			requestBody: map[string]interface{}{
				"title":   "更新的魔法体系",
				"content": "更新后的描述",
			},
			mockSetup: func(m *MockWorldSettingService) {
				setting := &model.WorldSetting{
					ID:      settingID,
					Title:   "更新的魔法体系",
					Content: "更新后的描述",
				}
				m.On("Update", mock.Anything, userID, settingID, mock.AnythingOfType("*service.UpdateWorldSettingRequest")).Return(setting, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			settingID:      settingID.String(),
			requestBody:    map[string]interface{}{"name": "test"},
			mockSetup:      func(m *MockWorldSettingService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "世界设定不存在",
			userID:    userID,
			settingID: settingID.String(),
			requestBody: map[string]interface{}{
				"name": "test",
			},
			mockSetup: func(m *MockWorldSettingService) {
				m.On("Update", mock.Anything, userID, settingID, mock.AnythingOfType("*service.UpdateWorldSettingRequest")).Return(nil, repository.ErrWorldSettingNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWorldSettingService)
			tt.mockSetup(mockService)

			handler := &WorldSettingHandler{
				settingService: mockService,
			}
			router := setupTestRouter()

			router.PUT("/api/v1/world-settings/:id", func(c *gin.Context) {
				if tt.userID != uuid.Nil {
					setUserIDContext(c, tt.userID)
				}
				handler.Update(c)
			})

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			url := "/api/v1/world-settings/" + tt.settingID
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestWorldSettingHandler_Delete 测试删除世界设定 API
func TestWorldSettingHandler_Delete(t *testing.T) {
	userID := uuid.New()
	settingID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		settingID      string
		mockSetup      func(*MockWorldSettingService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功删除世界设定",
			userID:    userID,
			settingID: settingID.String(),
			mockSetup: func(m *MockWorldSettingService) {
				m.On("Delete", mock.Anything, userID, settingID).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			settingID:      settingID.String(),
			mockSetup:      func(m *MockWorldSettingService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "世界设定不存在",
			userID:    userID,
			settingID: settingID.String(),
			mockSetup: func(m *MockWorldSettingService) {
				m.On("Delete", mock.Anything, userID, settingID).Return(repository.ErrWorldSettingNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWorldSettingService)
			tt.mockSetup(mockService)

			handler := &WorldSettingHandler{
				settingService: mockService,
			}
			router := setupTestRouter()

			router.DELETE("/api/v1/world-settings/:id", func(c *gin.Context) {
				if tt.userID != uuid.Nil {
					setUserIDContext(c, tt.userID)
				}
				handler.Delete(c)
			})

			url := "/api/v1/world-settings/" + tt.settingID
			req := httptest.NewRequest(http.MethodDelete, url, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}
