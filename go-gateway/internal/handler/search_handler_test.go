package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/repository"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSearchService 是 SearchService 的 mock
type MockSearchService struct {
	mock.Mock
}

func (m *MockGlobalSearch) GlobalSearch(ctx context.Context, userID uuid.UUID, req *service.GlobalSearchRequest) (*service.SearchResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.SearchResponse), args.Error(1)
}

func (m *MockProjectSearch) ProjectSearch(ctx context.Context, userID, projectID uuid.UUID, req *service.ProjectSearchRequest) (*service.SearchResponse, error) {
	args := m.Called(ctx, userID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.SearchResponse), args.Error(1)
}

func (m *MockAdvancedFilter) AdvancedFilter(ctx context.Context, userID, projectID uuid.UUID, req *service.AdvancedFilterRequest) (*service.SearchResponse, error) {
	args := m.Called(ctx, userID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.SearchResponse), args.Error(1)
}

// TestSearchHandler_GlobalSearch 测试全局搜索 API
func TestSearchHandler_GlobalSearch(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		query          string
		setupAuth      bool
		mockSetup      func(*MockSearchService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功全局搜索",
			query:     "q=测试&page=1&page_size=20",
			setupAuth: true,
			mockSetup: func(m *MockSearchService) {
				resp := &service.SearchResponse{
					Results: []service.SearchResult{
						{
							Type:      "project",
							ID:        uuid.New().String(),
							Title:     "测试项目",
							Snippet:   "这是一个测试项目",
							Highlight: "测试项目",
						},
					},
					Total:      1,
					Page:       1,
					PageSize:   20,
					TotalPages: 1,
				}
				m.On("GlobalSearch", mock.Anything, userID, mock.AnythingOfType("*service.GlobalSearchRequest")).Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			query:          "q=测试",
			setupAuth:      false,
			mockSetup:      func(m *MockSearchService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "缺少查询参数",
			query:          "",
			setupAuth:      true,
			mockSetup:      func(m *MockSearchService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockSearchService)
			tt.mockSetup(mockService)

			handler := &SearchHandler{
				searchService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.GET("/api/v1/search", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.GlobalSearch(c)
			})

			// 创建请求
			url := "/api/v1/search"
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

// TestSearchHandler_ProjectSearch 测试项目内搜索 API
func TestSearchHandler_ProjectSearch(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		query          string
		setupAuth      bool
		mockSetup      func(*MockSearchService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功搜索章节",
			projectID: projectID.String(),
			query:     "q=测试&type=chapter&page=1&page_size=20",
			setupAuth: true,
			mockSetup: func(m *MockSearchService) {
				resp := &service.SearchResponse{
					Results: []service.SearchResult{
						{
							Type:      "chapter",
							ID:        uuid.New().String(),
							Title:     "测试章节",
							Snippet:   "这是一个测试章节",
							Highlight: "测试章节",
						},
					},
					Total:      1,
					Page:       1,
					PageSize:   20,
					TotalPages: 1,
				}
				m.On("ProjectSearch", mock.Anything, userID, projectID, mock.AnythingOfType("*service.ProjectSearchRequest")).Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			query:          "q=测试",
			setupAuth:      false,
			mockSetup:      func(m *MockSearchService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的项目ID",
			projectID:      "invalid-uuid",
			query:          "q=测试",
			setupAuth:      true,
			mockSetup:      func(m *MockSearchService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			query:     "q=测试",
			setupAuth: true,
			mockSetup: func(m *MockSearchService) {
				m.On("ProjectSearch", mock.Anything, userID, projectID, mock.AnythingOfType("*service.ProjectSearchRequest")).Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:      "项目不属于用户",
			projectID: projectID.String(),
			query:     "q=测试",
			setupAuth: true,
			mockSetup: func(m *MockSearchService) {
				m.On("ProjectSearch", mock.Anything, userID, projectID, mock.AnythingOfType("*service.ProjectSearchRequest")).Return(nil, service.ErrProjectNotOwned)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockSearchService)
			tt.mockSetup(mockService)

			handler := &SearchHandler{
				searchService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.GET("/api/v1/projects/:id/search", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.ProjectSearch(c)
			})

			// 创建请求
			url := "/api/v1/projects/" + tt.projectID + "/search"
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

// TestSearchHandler_AdvancedFilter 测试高级筛选 API
func TestSearchHandler_AdvancedFilter(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockSearchService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功高级筛选",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"tags":           []string{"tag1", "tag2"},
				"status":         []string{"draft"},
				"word_count_min": 1000,
				"word_count_max": 5000,
				"page":           1,
				"page_size":      20,
			},
			setupAuth: true,
			mockSetup: func(m *MockSearchService) {
				resp := &service.SearchResponse{
					Results: []service.SearchResult{
						{
							Type:    "chapter",
							ID:      uuid.New().String(),
							Title:   "测试章节",
							Snippet: "这是一个测试章节",
						},
					},
					Total:      1,
					Page:       1,
					PageSize:   20,
					TotalPages: 1,
				}
				m.On("AdvancedFilter", mock.Anything, userID, projectID, mock.AnythingOfType("*service.AdvancedFilterRequest")).Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:      "未认证",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"tags": []string{"tag1"},
			},
			setupAuth:      false,
			mockSetup:      func(m *MockSearchService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "无效的项目ID",
			projectID: "invalid-uuid",
			requestBody: map[string]interface{}{
				"tags": []string{"tag1"},
			},
			setupAuth:      true,
			mockSetup:      func(m *MockSearchService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "无效的请求体",
			projectID:      projectID.String(),
			requestBody:    "invalid json",
			setupAuth:      true,
			mockSetup:      func(m *MockSearchService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "项目不属于用户",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"tags": []string{"tag1"},
			},
			setupAuth: true,
			mockSetup: func(m *MockSearchService) {
				m.On("AdvancedFilter", mock.Anything, userID, projectID, mock.AnythingOfType("*service.AdvancedFilterRequest")).Return(nil, service.ErrProjectNotOwned)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockSearchService)
			tt.mockSetup(mockService)

			handler := &SearchHandler{
				searchService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.POST("/api/v1/projects/:id/advanced-filter", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.AdvancedFilter(c)
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

			url := "/api/v1/projects/" + tt.projectID + "/advanced-filter"
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
