package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/service"
	"go-gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBatchOperationService 是 BatchOperationService 的 mock
type MockBatchOperationService struct {
	mock.Mock
}

func (m *MockBatchUpdate) BatchUpdate(ctx context.Context, userID, projectID uuid.UUID, req *service.BatchUpdateRequest) (*service.BatchOperationResult, error) {
	args := m.Called(ctx, userID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.BatchOperationResult), args.Error(1)
}

func (m *MockBatchDelete) BatchDelete(ctx context.Context, userID, projectID uuid.UUID, req *service.BatchDeleteRequest) (*service.BatchOperationResult, error) {
	args := m.Called(ctx, userID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.BatchOperationResult), args.Error(1)
}

func (m *MockBatchStatusUpdate) BatchStatusUpdate(ctx context.Context, userID, projectID uuid.UUID, req *service.BatchStatusUpdateRequest) (*service.BatchOperationResult, error) {
	args := m.Called(ctx, userID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.BatchOperationResult), args.Error(1)
}

func (m *MockReorder) Reorder(ctx context.Context, userID, projectID uuid.UUID, req *service.ReorderRequest) error {
	args := m.Called(ctx, userID, projectID, req)
	return args.Error(0)
}

// TestBatchOperationHandler_BatchUpdate 测试批量更新 API
func TestBatchOperationHandler_BatchUpdate(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()
	chapterID1 := uuid.New()
	chapterID2 := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockBatchOperationService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功批量更新",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"chapter_ids": []string{chapterID1.String(), chapterID2.String()},
				"updates": map[string]interface{}{
					"title":  "新标题",
					"status": "completed",
				},
			},
			setupAuth: true,
			mockSetup: func(m *MockBatchOperationService) {
				result := &service.BatchOperationResult{
					Success: 2,
					Failed:  0,
					Total:   2,
				}
				m.On("BatchUpdate", mock.Anything, userID, projectID, mock.AnythingOfType("*service.BatchUpdateRequest")).Return(result, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:      "未认证",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"chapter_ids": []string{chapterID1.String()},
				"updates":     map[string]interface{}{"title": "新标题"},
			},
			setupAuth:      false,
			mockSetup:      func(m *MockBatchOperationService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "无效的项目ID",
			projectID: "invalid-uuid",
			requestBody: map[string]interface{}{
				"chapter_ids": []string{chapterID1.String()},
				"updates":     map[string]interface{}{"title": "新标题"},
			},
			setupAuth:      true,
			mockSetup:      func(m *MockBatchOperationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "无效的请求体",
			projectID:      projectID.String(),
			requestBody:    "invalid json",
			setupAuth:      true,
			mockSetup:      func(m *MockBatchOperationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "批量大小超过限制",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"chapter_ids": []string{chapterID1.String()},
				"updates":     map[string]interface{}{"title": "新标题"},
			},
			setupAuth: true,
			mockSetup: func(m *MockBatchOperationService) {
				m.On("BatchUpdate", mock.Anything, userID, projectID, mock.AnythingOfType("*service.BatchUpdateRequest")).Return(nil, service.ErrBatchLimitExceeded)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockBatchOperationService)
			tt.mockSetup(mockService)

			handler := &BatchOperationHandler{
				batchService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.POST("/api/v1/projects/:id/chapters/batch-update", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.BatchUpdate(c)
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

			url := "/api/v1/projects/" + tt.projectID + "/chapters/batch-update"
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

// TestBatchOperationHandler_BatchDelete 测试批量删除 API
func TestBatchOperationHandler_BatchDelete(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()
	chapterID1 := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockBatchOperationService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功批量删除",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"chapter_ids": []string{chapterID1.String()},
			},
			setupAuth: true,
			mockSetup: func(m *MockBatchOperationService) {
				result := &service.BatchOperationResult{
					Success: 1,
					Failed:  0,
					Total:   1,
				}
				m.On("BatchDelete", mock.Anything, userID, projectID, mock.AnythingOfType("*service.BatchDeleteRequest")).Return(result, nil)
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
			mockSetup:      func(m *MockBatchOperationService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockBatchOperationService)
			tt.mockSetup(mockService)

			handler := &BatchOperationHandler{
				batchService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.POST("/api/v1/projects/:id/chapters/batch-delete", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.BatchDelete(c)
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

			url := "/api/v1/projects/" + tt.projectID + "/chapters/batch-delete"
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

// TestBatchOperationHandler_Reorder 测试重排序 API
func TestBatchOperationHandler_Reorder(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()
	chapterID1 := uuid.New()
	chapterID2 := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockBatchOperationService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功重排序",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"chapter_orders": []map[string]interface{}{
					{"chapter_id": chapterID1.String(), "chapter_number": 2},
					{"chapter_id": chapterID2.String(), "chapter_number": 1},
				},
			},
			setupAuth: true,
			mockSetup: func(m *MockBatchOperationService) {
				m.On("Reorder", mock.Anything, userID, projectID, mock.AnythingOfType("*service.ReorderRequest")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:      "未认证",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"chapter_orders": []map[string]interface{}{
					{"chapter_id": chapterID1.String(), "chapter_number": 1},
				},
			},
			setupAuth:      false,
			mockSetup:      func(m *MockBatchOperationService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "章节不属于项目",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"chapter_orders": []map[string]interface{}{
					{"chapter_id": chapterID1.String(), "chapter_number": 1},
				},
			},
			setupAuth: true,
			mockSetup: func(m *MockBatchOperationService) {
				m.On("Reorder", mock.Anything, userID, projectID, mock.AnythingOfType("*service.ReorderRequest")).Return(service.ErrChapterNotInProject)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockBatchOperationService)
			tt.mockSetup(mockService)

			handler := &BatchOperationHandler{
				batchService: mockService,
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
