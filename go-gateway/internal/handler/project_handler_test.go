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

// MockProjectService 是 ProjectService 的 mock
type MockProjectService struct {
	mock.Mock
}

func (m *MockCreate) Create(ctx context.Context, userID uuid.UUID, req *service.CreateProjectRequest) (*model.Project, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockList) List(ctx context.Context, userID uuid.UUID, page, pageSize int, status, genre, sort string) (*service.ProjectListResponse, error) {
	args := m.Called(ctx, userID, page, pageSize, status, genre, sort)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.ProjectListResponse), args.Error(1)
}

func (m *MockGet) Get(ctx context.Context, userID, projectID uuid.UUID) (*model.Project, error) {
	args := m.Called(ctx, userID, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockUpdate) Update(ctx context.Context, userID, projectID uuid.UUID, req *service.UpdateProjectRequest) (*model.Project, error) {
	args := m.Called(ctx, userID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockDelete) Delete(ctx context.Context, userID, projectID uuid.UUID) error {
	args := m.Called(ctx, userID, projectID)
	return args.Error(0)
}

func (m *MockRestore) Restore(ctx context.Context, userID, projectID uuid.UUID) (*model.Project, error) {
	args := m.Called(ctx, userID, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockUpdateMetadata) UpdateMetadata(ctx context.Context, userID, projectID uuid.UUID, req *service.UpdateMetadataRequest) (*model.Project, error) {
	args := m.Called(ctx, userID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockGetStatistics) GetStatistics(ctx context.Context, userID, projectID uuid.UUID) (*service.ProjectStatistics, error) {
	args := m.Called(ctx, userID, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.ProjectStatistics), args.Error(1)
}

// TestProjectHandler_Create 测试创建项目 API
func TestProjectHandler_Create(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockProjectService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "成功创建项目",
			requestBody: map[string]interface{}{
				"title":       "测试项目",
				"description": "这是一个测试项目",
				"genre":       "fantasy",
			},
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				project := &model.Project{
					ID:          projectID,
					UserID:      userID,
					Title:       "测试项目",
					Description: "这是一个测试项目",
					Genre:       "fantasy",
					Status:      "draft",
				}
				m.On("Create", mock.Anything, userID, mock.AnythingOfType("*service.CreateProjectRequest")).Return(project, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "未认证",
			requestBody: map[string]interface{}{
				"title": "测试项目",
			},
			setupAuth:      false,
			mockSetup:      func(m *MockProjectService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的请求体",
			requestBody:    "invalid json",
			setupAuth:      true,
			mockSetup:      func(m *MockProjectService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockProjectService)
			tt.mockSetup(mockService)

			handler := &ProjectHandler{
				projectService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.POST("/api/v1/projects", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Create(c)
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

			req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", bytes.NewBuffer(body))
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
				assert.NotNil(t, resp.Data)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestProjectHandler_List 测试获取项目列表 API
func TestProjectHandler_List(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		query          string
		setupAuth      bool
		mockSetup      func(*MockProjectService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取项目列表",
			query:     "page=1&page_size=20&status=draft",
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				resp := &service.ProjectListResponse{
					Items: []model.Project{
						{
							ID:     uuid.New(),
							UserID: userID,
							Title:  "测试项目",
							Status: "draft",
						},
					},
					Total:      1,
					Page:       1,
					PageSize:   20,
					TotalPages: 1,
				}
				m.On("List", mock.Anything, userID, 1, 20, "draft", "", "").Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			query:          "page=1&page_size=20",
			setupAuth:      false,
			mockSetup:      func(m *MockProjectService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "默认分页参数",
			query:     "",
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				resp := &service.ProjectListResponse{
					Items:      []model.Project{},
					Total:      0,
					Page:       1,
					PageSize:   20,
					TotalPages: 0,
				}
				m.On("List", mock.Anything, userID, 1, 20, "", "", "").Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockProjectService)
			tt.mockSetup(mockService)

			handler := &ProjectHandler{
				projectService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.GET("/api/v1/projects", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.List(c)
			})

			// 创建请求
			url := "/api/v1/projects"
			if tt.query != "" {
				url += "?" + tt.query
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
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

// TestProjectHandler_Get 测试获取单个项目 API
func TestProjectHandler_Get(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		setupAuth      bool
		mockSetup      func(*MockProjectService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取项目",
			projectID: projectID.String(),
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
					Status: "draft",
				}
				m.On("Get", mock.Anything, userID, projectID).Return(project, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockProjectService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的项目ID",
			projectID:      "invalid-uuid",
			setupAuth:      true,
			mockSetup:      func(m *MockProjectService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				m.On("Get", mock.Anything, userID, projectID).Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:      "项目不属于用户",
			projectID: projectID.String(),
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				m.On("Get", mock.Anything, userID, projectID).Return(nil, service.ErrProjectNotOwned)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockProjectService)
			tt.mockSetup(mockService)

			handler := &ProjectHandler{
				projectService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.GET("/api/v1/projects/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Get(c)
			})

			// 创建请求
			url := "/api/v1/projects/" + tt.projectID
			req := httptest.NewRequest(http.MethodGet, url, nil)
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

// TestProjectHandler_Update 测试更新项目 API
func TestProjectHandler_Update(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockProjectService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功更新项目",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"title":       "更新后的标题",
				"description": "更新后的描述",
			},
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				project := &model.Project{
					ID:          projectID,
					UserID:      userID,
					Title:       "更新后的标题",
					Description: "更新后的描述",
				}
				m.On("Update", mock.Anything, userID, projectID, mock.AnythingOfType("*service.UpdateProjectRequest")).Return(project, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:      "未认证",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"title": "更新后的标题",
			},
			setupAuth:      false,
			mockSetup:      func(m *MockProjectService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的项目ID",
			projectID:      "invalid-uuid",
			requestBody:    map[string]interface{}{"title": "更新后的标题"},
			setupAuth:      true,
			mockSetup:      func(m *MockProjectService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"title": "更新后的标题",
			},
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				m.On("Update", mock.Anything, userID, projectID, mock.AnythingOfType("*service.UpdateProjectRequest")).Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockProjectService)
			tt.mockSetup(mockService)

			handler := &ProjectHandler{
				projectService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.PUT("/api/v1/projects/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Update(c)
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

			url := "/api/v1/projects/" + tt.projectID
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
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
				assert.NotNil(t, resp.Data)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestProjectHandler_Delete 测试删除项目 API
func TestProjectHandler_Delete(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		setupAuth      bool
		mockSetup      func(*MockProjectService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功删除项目",
			projectID: projectID.String(),
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				m.On("Delete", mock.Anything, userID, projectID).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockProjectService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的项目ID",
			projectID:      "invalid-uuid",
			setupAuth:      true,
			mockSetup:      func(m *MockProjectService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				m.On("Delete", mock.Anything, userID, projectID).Return(repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockProjectService)
			tt.mockSetup(mockService)

			handler := &ProjectHandler{
				projectService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.DELETE("/api/v1/projects/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Delete(c)
			})

			// 创建请求
			url := "/api/v1/projects/" + tt.projectID
			req := httptest.NewRequest(http.MethodDelete, url, nil)
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

// TestProjectHandler_Restore 测试恢复项目 API
func TestProjectHandler_Restore(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		setupAuth      bool
		mockSetup      func(*MockProjectService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功恢复项目",
			projectID: projectID.String(),
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				m.On("Restore", mock.Anything, userID, projectID).Return(project, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockProjectService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "项目未被删除",
			projectID: projectID.String(),
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				m.On("Restore", mock.Anything, userID, projectID).Return(nil, service.ErrProjectNotDeleted)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockProjectService)
			tt.mockSetup(mockService)

			handler := &ProjectHandler{
				projectService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.POST("/api/v1/projects/:id/restore", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Restore(c)
			})

			// 创建请求
			url := "/api/v1/projects/" + tt.projectID + "/restore"
			req := httptest.NewRequest(http.MethodPost, url, nil)
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

// TestProjectHandler_GetStatistics 测试获取项目统计 API
func TestProjectHandler_GetStatistics(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		setupAuth      bool
		mockSetup      func(*MockProjectService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取统计信息",
			projectID: projectID.String(),
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				stats := &service.ProjectStatistics{
					TotalChapters:   10,
					TotalWords:      50000,
					CharacterCount:  5,
					LocationCount:   3,
					ForeshadowCount: 2,
					LastUpdated:     "2024-01-01T00:00:00Z",
				}
				m.On("GetStatistics", mock.Anything, userID, projectID).Return(stats, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockProjectService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			setupAuth: true,
			mockSetup: func(m *MockProjectService) {
				m.On("GetStatistics", mock.Anything, userID, projectID).Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockProjectService)
			tt.mockSetup(mockService)

			handler := &ProjectHandler{
				projectService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.GET("/api/v1/projects/:id/statistics", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.GetStatistics(c)
			})

			// 创建请求
			url := "/api/v1/projects/" + tt.projectID + "/statistics"
			req := httptest.NewRequest(http.MethodGet, url, nil)
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
