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

// MockChapterService 是 ChapterService 的 mock
type MockChapterService struct {
	mock.Mock
}

func (m *MockCreate) Create(ctx context.Context, userID, projectID uuid.UUID, req *service.CreateChapterRequest) (*model.Chapter, error) {
	args := m.Called(ctx, userID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chapter), args.Error(1)
}

func (m *MockList) List(ctx context.Context, userID, projectID uuid.UUID, page, pageSize int, status, sort string) (*service.ChapterListResponse, error) {
	args := m.Called(ctx, userID, projectID, page, pageSize, status, sort)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.ChapterListResponse), args.Error(1)
}

func (m *MockGet) Get(ctx context.Context, userID, chapterID uuid.UUID) (*model.Chapter, error) {
	args := m.Called(ctx, userID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chapter), args.Error(1)
}

func (m *MockUpdate) Update(ctx context.Context, userID, chapterID uuid.UUID, req *service.UpdateChapterRequest) (*model.Chapter, error) {
	args := m.Called(ctx, userID, chapterID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chapter), args.Error(1)
}

func (m *MockDelete) Delete(ctx context.Context, userID, chapterID uuid.UUID) error {
	args := m.Called(ctx, userID, chapterID)
	return args.Error(0)
}

func (m *MockReorder) Reorder(ctx context.Context, userID, projectID uuid.UUID, req *service.ReorderChaptersRequest) error {
	args := m.Called(ctx, userID, projectID, req)
	return args.Error(0)
}

// TestChapterHandler_Create 测试创建章节 API
func TestChapterHandler_Create(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()
	chapterID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockChapterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功创建章节",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"title":   "第一章",
				"content": "章节内容",
			},
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				chapter := &model.Chapter{
					ID:        chapterID,
					ProjectID: projectID,
					Title:     "第一章",
					Content:   "章节内容",
					Status:    "draft",
				}
				m.On("Create", mock.Anything, userID, projectID, mock.AnythingOfType("*service.CreateChapterRequest")).Return(chapter, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:      "未认证",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"title": "第一章",
			},
			setupAuth:      false,
			mockSetup:      func(m *MockChapterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的项目ID",
			projectID:      "invalid-uuid",
			requestBody:    map[string]interface{}{"title": "第一章"},
			setupAuth:      true,
			mockSetup:      func(m *MockChapterService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"title": "第一章",
			},
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				m.On("Create", mock.Anything, userID, projectID, mock.AnythingOfType("*service.CreateChapterRequest")).Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockChapterService)
			tt.mockSetup(mockService)

			handler := &ChapterHandler{
				chapterService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.POST("/api/v1/projects/:id/chapters", func(c *gin.Context) {
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

			url := "/api/v1/projects/" + tt.projectID + "/chapters"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
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

// TestChapterHandler_List 测试获取章节列表 API
func TestChapterHandler_List(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		query          string
		setupAuth      bool
		mockSetup      func(*MockChapterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取章节列表",
			projectID: projectID.String(),
			query:     "page=1&page_size=50&status=draft",
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				resp := &service.ChapterListResponse{
					Items: []model.Chapter{
						{
							ID:        uuid.New(),
							ProjectID: projectID,
							Title:     "第一章",
							Status:    "draft",
						},
					},
					Total:      1,
					Page:       1,
					PageSize:   50,
					TotalPages: 1,
				}
				m.On("List", mock.Anything, userID, projectID, 1, 50, "draft", "").Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			query:          "page=1&page_size=50",
			setupAuth:      false,
			mockSetup:      func(m *MockChapterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的项目ID",
			projectID:      "invalid-uuid",
			query:          "page=1&page_size=50",
			setupAuth:      true,
			mockSetup:      func(m *MockChapterService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			query:     "page=1&page_size=50",
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				m.On("List", mock.Anything, userID, projectID, 1, 50, "", "").Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockChapterService)
			tt.mockSetup(mockService)

			handler := &ChapterHandler{
				chapterService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.GET("/api/v1/projects/:id/chapters", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.List(c)
			})

			// 创建请求
			url := "/api/v1/projects/" + tt.projectID + "/chapters"
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

// TestChapterHandler_Get 测试获取单个章节 API
func TestChapterHandler_Get(t *testing.T) {
	userID := uuid.New()
	chapterID := uuid.New()

	tests := []struct {
		name           string
		chapterID      string
		setupAuth      bool
		mockSetup      func(*MockChapterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取章节",
			chapterID: chapterID.String(),
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				chapter := &model.Chapter{
					ID:      chapterID,
					Title:   "第一章",
					Content: "章节内容",
					Status:  "draft",
				}
				m.On("Get", mock.Anything, userID, chapterID).Return(chapter, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			chapterID:      chapterID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockChapterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的章节ID",
			chapterID:      "invalid-uuid",
			setupAuth:      true,
			mockSetup:      func(m *MockChapterService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "章节不存在",
			chapterID: chapterID.String(),
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				m.On("Get", mock.Anything, userID, chapterID).Return(nil, repository.ErrChapterNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:      "章节不属于用户",
			chapterID: chapterID.String(),
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				m.On("Get", mock.Anything, userID, chapterID).Return(nil, service.ErrChapterNotOwned)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockChapterService)
			tt.mockSetup(mockService)

			handler := &ChapterHandler{
				chapterService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.GET("/api/v1/chapters/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Get(c)
			})

			// 创建请求
			url := "/api/v1/chapters/" + tt.chapterID
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

// TestChapterHandler_Update 测试更新章节 API
func TestChapterHandler_Update(t *testing.T) {
	userID := uuid.New()
	chapterID := uuid.New()

	tests := []struct {
		name           string
		chapterID      string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockChapterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功更新章节",
			chapterID: chapterID.String(),
			requestBody: map[string]interface{}{
				"title":   "更新后的标题",
				"content": "更新后的内容",
			},
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				chapter := &model.Chapter{
					ID:      chapterID,
					Title:   "更新后的标题",
					Content: "更新后的内容",
				}
				m.On("Update", mock.Anything, userID, chapterID, mock.AnythingOfType("*service.UpdateChapterRequest")).Return(chapter, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:      "未认证",
			chapterID: chapterID.String(),
			requestBody: map[string]interface{}{
				"title": "更新后的标题",
			},
			setupAuth:      false,
			mockSetup:      func(m *MockChapterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的章节ID",
			chapterID:      "invalid-uuid",
			requestBody:    map[string]interface{}{"title": "更新后的标题"},
			setupAuth:      true,
			mockSetup:      func(m *MockChapterService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "章节不存在",
			chapterID: chapterID.String(),
			requestBody: map[string]interface{}{
				"title": "更新后的标题",
			},
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				m.On("Update", mock.Anything, userID, chapterID, mock.AnythingOfType("*service.UpdateChapterRequest")).Return(nil, repository.ErrChapterNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockChapterService)
			tt.mockSetup(mockService)

			handler := &ChapterHandler{
				chapterService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.PUT("/api/v1/chapters/:id", func(c *gin.Context) {
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

			url := "/api/v1/chapters/" + tt.chapterID
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

// TestChapterHandler_Delete 测试删除章节 API
func TestChapterHandler_Delete(t *testing.T) {
	userID := uuid.New()
	chapterID := uuid.New()

	tests := []struct {
		name           string
		chapterID      string
		setupAuth      bool
		mockSetup      func(*MockChapterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功删除章节",
			chapterID: chapterID.String(),
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				m.On("Delete", mock.Anything, userID, chapterID).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			chapterID:      chapterID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockChapterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的章节ID",
			chapterID:      "invalid-uuid",
			setupAuth:      true,
			mockSetup:      func(m *MockChapterService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "章节不存在",
			chapterID: chapterID.String(),
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				m.On("Delete", mock.Anything, userID, chapterID).Return(repository.ErrChapterNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockChapterService)
			tt.mockSetup(mockService)

			handler := &ChapterHandler{
				chapterService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.DELETE("/api/v1/chapters/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Delete(c)
			})

			// 创建请求
			url := "/api/v1/chapters/" + tt.chapterID
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

// TestChapterHandler_Reorder 测试重排序章节 API
func TestChapterHandler_Reorder(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()
	chapterID1 := uuid.New()
	chapterID2 := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockChapterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功重排序章节",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"chapter_ids": []string{chapterID2.String(), chapterID1.String()},
			},
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				m.On("Reorder", mock.Anything, userID, projectID, mock.AnythingOfType("*service.ReorderChaptersRequest")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:      "未认证",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"chapter_ids": []string{chapterID1.String()},
			},
			setupAuth:      false,
			mockSetup:      func(m *MockChapterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "无效的项目ID",
			projectID: "invalid-uuid",
			requestBody: map[string]interface{}{
				"chapter_ids": []string{chapterID1.String()},
			},
			setupAuth:      true,
			mockSetup:      func(m *MockChapterService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"chapter_ids": []string{chapterID1.String()},
			},
			setupAuth: true,
			mockSetup: func(m *MockChapterService) {
				m.On("Reorder", mock.Anything, userID, projectID, mock.AnythingOfType("*service.ReorderChaptersRequest")).Return(repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockChapterService)
			tt.mockSetup(mockService)

			handler := &ChapterHandler{
				chapterService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.POST("/api/v1/projects/:id/chapters/reorder", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Reorder(c)
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

			url := "/api/v1/projects/" + tt.projectID + "/chapters/reorder"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
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
