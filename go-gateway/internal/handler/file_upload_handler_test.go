package handler

import (
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

// MockFileUploadService 是 FileUploadService 的 mock
type MockFileUploadService struct {
	mock.Mock
}

func (m *MockUploadFile) UploadFile(ctx context.Context, userID uuid.UUID, req *service.UploadFileRequest) (*model.File, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.File), args.Error(1)
}

func (m *MockGetFile) GetFile(ctx context.Context, userID, fileID uuid.UUID) (*model.File, error) {
	args := m.Called(ctx, userID, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.File), args.Error(1)
}

func (m *MockDeleteFile) DeleteFile(ctx context.Context, userID, fileID uuid.UUID) error {
	args := m.Called(ctx, userID, fileID)
	return args.Error(0)
}

func (m *MockListFiles) ListFiles(ctx context.Context, userID uuid.UUID, page, pageSize int, fileType, sort string) (*service.FileListResponse, error) {
	args := m.Called(ctx, userID, page, pageSize, fileType, sort)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.FileListResponse), args.Error(1)
}

func (m *MockGetUserStorageStats) GetUserStorageStats(ctx context.Context, userID uuid.UUID) (*service.StorageStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.StorageStats), args.Error(1)
}

// TestFileUploadHandler_GetFile 测试获取文件 API
func TestFileUploadHandler_GetFile(t *testing.T) {
	userID := uuid.New()
	fileID := uuid.New()

	tests := []struct {
		name           string
		fileID         string
		setupAuth      bool
		mockSetup      func(*MockFileUploadService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取文件",
			fileID:    fileID.String(),
			setupAuth: true,
			mockSetup: func(m *MockFileUploadService) {
				file := &model.File{
					ID:          fileID,
					UserID:      userID,
					FileName:    "test.jpg",
					StoragePath: "/uploads/test.jpg",
					FileSize:    1024,
					MimeType:    "image/jpeg",
					FileType:    "image",
					Hash:        "abc123",
				}
				m.On("GetFile", mock.Anything, userID, fileID).Return(file, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			fileID:         fileID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockFileUploadService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的文件ID",
			fileID:         "invalid-uuid",
			setupAuth:      true,
			mockSetup:      func(m *MockFileUploadService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "文件不存在",
			fileID:    fileID.String(),
			setupAuth: true,
			mockSetup: func(m *MockFileUploadService) {
				m.On("GetFile", mock.Anything, userID, fileID).Return(nil, repository.ErrFileNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:      "文件不属于用户",
			fileID:    fileID.String(),
			setupAuth: true,
			mockSetup: func(m *MockFileUploadService) {
				m.On("GetFile", mock.Anything, userID, fileID).Return(nil, service.ErrFileNotOwned)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockFileUploadService)
			tt.mockSetup(mockService)

			handler := &FileUploadHandler{
				fileService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.GET("/api/v1/files/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.GetFile(c)
			})

			// 创建请求
			url := "/api/v1/files/" + tt.fileID
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

// TestFileUploadHandler_DeleteFile 测试删除文件 API
func TestFileUploadHandler_DeleteFile(t *testing.T) {
	userID := uuid.New()
	fileID := uuid.New()

	tests := []struct {
		name           string
		fileID         string
		setupAuth      bool
		mockSetup      func(*MockFileUploadService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功删除文件",
			fileID:    fileID.String(),
			setupAuth: true,
			mockSetup: func(m *MockFileUploadService) {
				m.On("DeleteFile", mock.Anything, userID, fileID).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			fileID:         fileID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockFileUploadService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的文件ID",
			fileID:         "invalid-uuid",
			setupAuth:      true,
			mockSetup:      func(m *MockFileUploadService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "文件不存在",
			fileID:    fileID.String(),
			setupAuth: true,
			mockSetup: func(m *MockFileUploadService) {
				m.On("DeleteFile", mock.Anything, userID, fileID).Return(repository.ErrFileNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:      "文件不属于用户",
			fileID:    fileID.String(),
			setupAuth: true,
			mockSetup: func(m *MockFileUploadService) {
				m.On("DeleteFile", mock.Anything, userID, fileID).Return(service.ErrFileNotOwned)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockFileUploadService)
			tt.mockSetup(mockService)

			handler := &FileUploadHandler{
				fileService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.DELETE("/api/v1/files/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.DeleteFile(c)
			})

			// 创建请求
			url := "/api/v1/files/" + tt.fileID
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

// TestFileUploadHandler_ListFiles 测试获取文件列表 API
func TestFileUploadHandler_ListFiles(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		query          string
		setupAuth      bool
		mockSetup      func(*MockFileUploadService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取文件列表",
			query:     "page=1&page_size=20&file_type=image",
			setupAuth: true,
			mockSetup: func(m *MockFileUploadService) {
				resp := &service.FileListResponse{
					Items: []model.File{
						{
							ID:          uuid.New(),
							UserID:      userID,
							FileName:    "test.jpg",
							StoragePath: "/uploads/test.jpg",
							FileSize:    1024,
							MimeType:    "image/jpeg",
							FileType:    "image",
						},
					},
					Total:      1,
					Page:       1,
					PageSize:   20,
					TotalPages: 1,
				}
				m.On("ListFiles", mock.Anything, userID, 1, 20, "image", "").Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			query:          "page=1&page_size=20",
			setupAuth:      false,
			mockSetup:      func(m *MockFileUploadService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "默认分页参数",
			query:     "",
			setupAuth: true,
			mockSetup: func(m *MockFileUploadService) {
				resp := &service.FileListResponse{
					Items:      []model.File{},
					Total:      0,
					Page:       1,
					PageSize:   20,
					TotalPages: 0,
				}
				m.On("ListFiles", mock.Anything, userID, 1, 20, "", "").Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockFileUploadService)
			tt.mockSetup(mockService)

			handler := &FileUploadHandler{
				fileService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.GET("/api/v1/files", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.ListFiles(c)
			})

			// 创建请求
			url := "/api/v1/files"
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

// TestFileUploadHandler_GetStorageStats 测试获取存储统计 API
func TestFileUploadHandler_GetStorageStats(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		setupAuth      bool
		mockSetup      func(*MockFileUploadService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取存储统计",
			setupAuth: true,
			mockSetup: func(m *MockFileUploadService) {
				stats := &service.StorageStats{
					FileCount: 10,
					TotalSize: 1024000,
				}
				m.On("GetUserStorageStats", mock.Anything, userID).Return(stats, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			setupAuth:      false,
			mockSetup:      func(m *MockFileUploadService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置
			mockService := new(MockFileUploadService)
			tt.mockSetup(mockService)

			handler := &FileUploadHandler{
				fileService: mockService,
			}
			router := setupTestRouter()

			// 注册路由
			router.GET("/api/v1/files/storage-stats", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.GetStorageStats(c)
			})

			// 创建请求
			req := httptest.NewRequest(http.MethodGet, "/api/v1/files/storage-stats", nil)
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
