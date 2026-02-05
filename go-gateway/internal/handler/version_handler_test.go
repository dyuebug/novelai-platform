package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"
	"go-gateway/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVersionService 是 VersionService 的 mock
type MockVersionService struct {
	mock.Mock
}

func (m *MockListVersions) ListVersions(ctx context.Context, userID, chapterID uuid.UUID, page, pageSize int) (*service.VersionListResponse, error) {
	args := m.Called(ctx, userID, chapterID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.VersionListResponse), args.Error(1)
}

func (m *MockGetVersion) GetVersion(ctx context.Context, userID, chapterID uuid.UUID, versionNumber int) (*model.ChapterVersion, error) {
	args := m.Called(ctx, userID, chapterID, versionNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChapterVersion), args.Error(1)
}

func (m *MockRestoreVersion) RestoreVersion(ctx context.Context, userID, chapterID uuid.UUID, versionNumber int) (*model.Chapter, error) {
	args := m.Called(ctx, userID, chapterID, versionNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chapter), args.Error(1)
}

func (m *MockDiffVersions) DiffVersions(ctx context.Context, userID, chapterID uuid.UUID, v1, v2 int) (*service.DiffResponse, error) {
	args := m.Called(ctx, userID, chapterID, v1, v2)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.DiffResponse), args.Error(1)
}

// TestVersionHandler_ListVersions 测试获取版本列表 API
func TestVersionHandler_ListVersions(t *testing.T) {
	userID := uuid.New()
	chapterID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		chapterID      string
		query          string
		mockSetup      func(*MockVersionService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取版本列表",
			userID:    userID,
			chapterID: chapterID.String(),
			query:     "page=1&page_size=20",
			mockSetup: func(m *MockVersionService) {
				result := &service.VersionListResponse{
					Items: []model.ChapterVersion{
						{
							ID:            uuid.New(),
							ChapterID:     chapterID,
							VersionNumber: 1,
							Content:       "版本1内容",
						},
					},
					Total:    1,
					Page:     1,
					PageSize: 20,
				}
				m.On("ListVersions", mock.Anything, userID, chapterID, 1, 20).Return(result, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			chapterID:      chapterID.String(),
			mockSetup:      func(m *MockVersionService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的章节ID",
			userID:         userID,
			chapterID:      "invalid-uuid",
			mockSetup:      func(m *MockVersionService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "章节不存在",
			userID:    userID,
			chapterID: chapterID.String(),
			mockSetup: func(m *MockVersionService) {
				m.On("ListVersions", mock.Anything, userID, chapterID, 1, 20).Return(nil, repository.ErrChapterNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockVersionService)
			tt.mockSetup(mockService)

			handler := &VersionHandler{
				versionService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/chapters/:id/versions", func(c *gin.Context) {
				if tt.userID != uuid.Nil {
					setUserIDContext(c, tt.userID)
				}
				handler.ListVersions(c)
			})

			url := "/api/v1/chapters/" + tt.chapterID + "/versions"
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

// TestVersionHandler_GetVersion 测试获取指定版本 API
func TestVersionHandler_GetVersion(t *testing.T) {
	userID := uuid.New()
	chapterID := uuid.New()
	versionID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		chapterID      string
		versionNumber  string
		mockSetup      func(*MockVersionService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:          "成功获取版本",
			userID:        userID,
			chapterID:     chapterID.String(),
			versionNumber: "1",
			mockSetup: func(m *MockVersionService) {
				version := &model.ChapterVersion{
					ID:            versionID,
					ChapterID:     chapterID,
					VersionNumber: 1,
					Content:       "版本1内容",
				}
				m.On("GetVersion", mock.Anything, userID, chapterID, 1).Return(version, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			chapterID:      chapterID.String(),
			versionNumber:  "1",
			mockSetup:      func(m *MockVersionService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的版本号",
			userID:         userID,
			chapterID:      chapterID.String(),
			versionNumber:  "invalid",
			mockSetup:      func(m *MockVersionService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:          "版本不存在",
			userID:        userID,
			chapterID:     chapterID.String(),
			versionNumber: "999",
			mockSetup: func(m *MockVersionService) {
				m.On("GetVersion", mock.Anything, userID, chapterID, 999).Return(nil, repository.ErrVersionNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockVersionService)
			tt.mockSetup(mockService)

			handler := &VersionHandler{
				versionService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/chapters/:id/versions/:n", func(c *gin.Context) {
				if tt.userID != uuid.Nil {
					setUserIDContext(c, tt.userID)
				}
				handler.GetVersion(c)
			})

			url := "/api/v1/chapters/" + tt.chapterID + "/versions/" + tt.versionNumber
			req := httptest.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestVersionHandler_RestoreVersion 测试恢复版本 API
func TestVersionHandler_RestoreVersion(t *testing.T) {
	userID := uuid.New()
	chapterID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		chapterID      string
		versionNumber  string
		mockSetup      func(*MockVersionService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:          "成功恢复版本",
			userID:        userID,
			chapterID:     chapterID.String(),
			versionNumber: "1",
			mockSetup: func(m *MockVersionService) {
				chapter := &model.Chapter{
					ID:      chapterID,
					Title:   "第一章",
					Content: "恢复后的内容",
				}
				m.On("RestoreVersion", mock.Anything, userID, chapterID, 1).Return(chapter, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			chapterID:      chapterID.String(),
			versionNumber:  "1",
			mockSetup:      func(m *MockVersionService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:          "版本不存在",
			userID:        userID,
			chapterID:     chapterID.String(),
			versionNumber: "999",
			mockSetup: func(m *MockVersionService) {
				m.On("RestoreVersion", mock.Anything, userID, chapterID, 999).Return(nil, repository.ErrVersionNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockVersionService)
			tt.mockSetup(mockService)

			handler := &VersionHandler{
				versionService: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/chapters/:id/versions/:n/restore", func(c *gin.Context) {
				if tt.userID != uuid.Nil {
					setUserIDContext(c, tt.userID)
				}
				handler.RestoreVersion(c)
			})

			url := "/api/v1/chapters/" + tt.chapterID + "/versions/" + tt.versionNumber + "/restore"
			req := httptest.NewRequest(http.MethodPost, url, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestVersionHandler_DiffVersions 测试对比版本 API
func TestVersionHandler_DiffVersions(t *testing.T) {
	userID := uuid.New()
	chapterID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		chapterID      string
		query          string
		mockSetup      func(*MockVersionService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功对比版本",
			userID:    userID,
			chapterID: chapterID.String(),
			query:     "v1=1&v2=2",
			mockSetup: func(m *MockVersionService) {
				diff := &service.DiffResponse{
					Version1:  1,
					Version2:  2,
					Content1:  "版本1内容",
					Content2:  "版本2内容",
					Additions: 1,
					Deletions: 0,
				}
				m.On("DiffVersions", mock.Anything, userID, chapterID, 1, 2).Return(diff, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			chapterID:      chapterID.String(),
			query:          "v1=1&v2=2",
			mockSetup:      func(m *MockVersionService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的v1参数",
			userID:         userID,
			chapterID:      chapterID.String(),
			query:          "v1=invalid&v2=2",
			mockSetup:      func(m *MockVersionService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "无效的v2参数",
			userID:         userID,
			chapterID:      chapterID.String(),
			query:          "v1=1&v2=invalid",
			mockSetup:      func(m *MockVersionService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "版本不存在",
			userID:    userID,
			chapterID: chapterID.String(),
			query:     "v1=1&v2=999",
			mockSetup: func(m *MockVersionService) {
				m.On("DiffVersions", mock.Anything, userID, chapterID, 1, 999).Return(nil, repository.ErrVersionNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockVersionService)
			tt.mockSetup(mockService)

			handler := &VersionHandler{
				versionService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/chapters/:id/versions/diff", func(c *gin.Context) {
				if tt.userID != uuid.Nil {
					setUserIDContext(c, tt.userID)
				}
				handler.DiffVersions(c)
			})

			url := "/api/v1/chapters/" + tt.chapterID + "/versions/diff"
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
