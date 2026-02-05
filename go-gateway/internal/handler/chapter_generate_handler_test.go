package handler

import (
	"bytes"
	"context"
	"encoding/json"
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

// MockChapterGenerateService 是 ChapterGenerateService 的 mock
type MockChapterGenerateService struct {
	mock.Mock
}

func (m *MockChapterGenerateService) Get(ctx context.Context, userID, chapterID uuid.UUID) (*model.Chapter, error) {
	args := m.Called(ctx, userID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chapter), args.Error(1)
}

// MockProjectGenerateService 是 ProjectGenerateService 的 mock
type MockProjectGenerateService struct {
	mock.Mock
}

func (m *MockProjectGenerateService) Get(ctx context.Context, userID, projectID uuid.UUID) (*model.Project, error) {
	args := m.Called(ctx, userID, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

// MockAIGenerateService 是 AIGenerateService 的 mock
type MockAIGenerateService struct {
	mock.Mock
}

func (m *MockAIGenerateService) GenerateStream(ctx context.Context, req *service.GenerateRequest) (service.Stream, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(service.Stream), args.Error(1)
}

// TestChapterGenerateHandler_GenerateStream 测试流式生成章节内容
func TestChapterGenerateHandler_GenerateStream(t *testing.T) {
	userID := uuid.New()
	chapterID := uuid.New()

	tests := []struct {
		name           string
		chapterID      string
		requestBody    interface{}
		setAuth        bool
		mockSetup      func(*MockChapterGenerateService, *MockProjectGenerateService, *MockAIGenerateService)
		expectedStatus int
	}{
		{
			name:      "未授权",
			chapterID: chapterID.String(),
			requestBody: map[string]interface{}{
				"instruction": "生成章节内容",
			},
			setAuth:        false,
			mockSetup:      func(cs *MockChapterGenerateService, ps *MockProjectGenerateService, ai *MockAIGenerateService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:      "无效的章节 ID",
			chapterID: "invalid-uuid",
			requestBody: map[string]interface{}{
				"instruction": "生成章节内容",
			},
			setAuth:        true,
			mockSetup:      func(cs *MockChapterGenerateService, ps *MockProjectGenerateService, ai *MockAIGenerateService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "章节不存在",
			chapterID: chapterID.String(),
			requestBody: map[string]interface{}{
				"instruction": "生成章节内容",
			},
			setAuth: true,
			mockSetup: func(cs *MockChapterGenerateService, ps *MockProjectGenerateService, ai *MockAIGenerateService) {
				cs.On("Get", mock.Anything, userID, chapterID).Return(nil, repository.ErrChapterNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChapterService := new(MockChapterGenerateService)
			mockProjectService := new(MockProjectGenerateService)
			mockAIService := new(MockAIGenerateService)
			tt.mockSetup(mockChapterService, mockProjectService, mockAIService)

			handler := &ChapterGenerateHandler{
				ai:             mockAIService,
				chapterService: mockChapterService,
				projectService: mockProjectService,
			}
			router := setupTestRouter()

			if tt.setAuth {
				router.POST("/api/v1/chapters/:id/generate-stream", func(c *gin.Context) {
					setUserIDContext(c, userID)
					handler.GenerateStream(c)
				})
			} else {
				router.POST("/api/v1/chapters/:id/generate-stream", handler.GenerateStream)
			}

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			url := "/api/v1/chapters/" + tt.chapterID + "/generate-stream"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockChapterService.AssertExpectations(t)
			mockProjectService.AssertExpectations(t)
			mockAIService.AssertExpectations(t)
		})
	}
}

// TestChapterGenerateHandler_PartialRegenerateStream 测试局部重写流式生成
func TestChapterGenerateHandler_PartialRegenerateStream(t *testing.T) {
	userID := uuid.New()
	chapterID := uuid.New()

	tests := []struct {
		name           string
		chapterID      string
		requestBody    interface{}
		setAuth        bool
		mockSetup      func(*MockChapterGenerateService, *MockAIGenerateService)
		expectedStatus int
	}{
		{
			name:      "未授权",
			chapterID: chapterID.String(),
			requestBody: map[string]interface{}{
				"selection":   "选中的文本",
				"instruction": "重写指令",
			},
			setAuth:        false,
			mockSetup:      func(cs *MockChapterGenerateService, ai *MockAIGenerateService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:      "缺少必需字段",
			chapterID: chapterID.String(),
			requestBody: map[string]interface{}{
				"instruction": "重写指令",
			},
			setAuth:        true,
			mockSetup:      func(cs *MockChapterGenerateService, ai *MockAIGenerateService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "章节不存在",
			chapterID: chapterID.String(),
			requestBody: map[string]interface{}{
				"selection":   "选中的文本",
				"instruction": "重写指令",
			},
			setAuth: true,
			mockSetup: func(cs *MockChapterGenerateService, ai *MockAIGenerateService) {
				cs.On("Get", mock.Anything, userID, chapterID).Return(nil, repository.ErrChapterNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChapterService := new(MockChapterGenerateService)
			mockAIService := new(MockAIGenerateService)
			tt.mockSetup(mockChapterService, mockAIService)

			handler := &ChapterGenerateHandler{
				ai:             mockAIService,
				chapterService: mockChapterService,
			}
			router := setupTestRouter()

			if tt.setAuth {
				router.POST("/api/v1/chapters/:id/partial-regenerate-stream", func(c *gin.Context) {
					setUserIDContext(c, userID)
					handler.PartialRegenerateStream(c)
				})
			} else {
				router.POST("/api/v1/chapters/:id/partial-regenerate-stream", handler.PartialRegenerateStream)
			}

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			url := "/api/v1/chapters/" + tt.chapterID + "/partial-regenerate-stream"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockChapterService.AssertExpectations(t)
			mockAIService.AssertExpectations(t)
		})
	}
}

// TestChapterGenerateHandler_PolishStream 测试润色流式生成
func TestChapterGenerateHandler_PolishStream(t *testing.T) {
	userID := uuid.New()
	chapterID := uuid.New()

	tests := []struct {
		name           string
		chapterID      string
		requestBody    interface{}
		setAuth        bool
		mockSetup      func(*MockChapterGenerateService, *MockAIGenerateService)
		expectedStatus int
	}{
		{
			name:      "未授权",
			chapterID: chapterID.String(),
			requestBody: map[string]interface{}{
				"instruction": "润色指令",
			},
			setAuth:        false,
			mockSetup:      func(cs *MockChapterGenerateService, ai *MockAIGenerateService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "无效的请求体",
			chapterID:      chapterID.String(),
			requestBody:    "invalid json",
			setAuth:        true,
			mockSetup:      func(cs *MockChapterGenerateService, ai *MockAIGenerateService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "章节不存在",
			chapterID: chapterID.String(),
			requestBody: map[string]interface{}{
				"instruction": "润色指令",
			},
			setAuth: true,
			mockSetup: func(cs *MockChapterGenerateService, ai *MockAIGenerateService) {
				cs.On("Get", mock.Anything, userID, chapterID).Return(nil, repository.ErrChapterNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChapterService := new(MockChapterGenerateService)
			mockAIService := new(MockAIGenerateService)
			tt.mockSetup(mockChapterService, mockAIService)

			handler := &ChapterGenerateHandler{
				ai:             mockAIService,
				chapterService: mockChapterService,
			}
			router := setupTestRouter()

			if tt.setAuth {
				router.POST("/api/v1/chapters/:id/polish-stream", func(c *gin.Context) {
					setUserIDContext(c, userID)
					handler.PolishStream(c)
				})
			} else {
				router.POST("/api/v1/chapters/:id/polish-stream", handler.PolishStream)
			}

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			url := "/api/v1/chapters/" + tt.chapterID + "/polish-stream"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockChapterService.AssertExpectations(t)
			mockAIService.AssertExpectations(t)
		})
	}
}
