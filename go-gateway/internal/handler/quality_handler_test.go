package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockQualityAIClient 是 QualityAIClient 的 mock
type MockQualityAIClient struct {
	mock.Mock
}

func (m *MockQualityAIClient) AnalyzeReadingPower(ctx context.Context, chapterID, content, provider, model string) (map[string]interface{}, error) {
	args := m.Called(ctx, chapterID, content, provider, model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockQualityAIClient) CheckConsistency(ctx context.Context, chapterID, content string, context map[string]interface{}, provider, model string) (map[string]interface{}, error) {
	args := m.Called(ctx, chapterID, content, context, provider, model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockQualityAIClient) MultiAgentReview(ctx context.Context, chapterID, content, provider, model string) (map[string]interface{}, error) {
	args := m.Called(ctx, chapterID, content, provider, model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockQualityAIClient) EvaluateQuality(ctx context.Context, chapterID, content string, context map[string]interface{}, provider, model string) (map[string]interface{}, error) {
	args := m.Called(ctx, chapterID, content, context, provider, model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// TestQualityHandler_AnalyzeReadingPower 测试追读力分析
func TestQualityHandler_AnalyzeReadingPower(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		chapterID      string
		requestBody    interface{}
		mockSetup      func(*MockQualityAIClient)
		expectedStatus int
	}{
		{
			name:      "成功分析追读力",
			chapterID: "chapter-123",
			requestBody: map[string]interface{}{
				"content": "这是一段测试内容",
			},
			mockSetup: func(m *MockQualityAIClient) {
				result := map[string]interface{}{
					"score":    85,
					"analysis": "追读力较强",
				}
				m.On("AnalyzeReadingPower", mock.Anything, "chapter-123", "这是一段测试内容", "", "").Return(result, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "缺少必需字段",
			chapterID: "chapter-123",
			requestBody: map[string]interface{}{
				"provider": "openai",
			},
			mockSetup:      func(m *MockQualityAIClient) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "AI 服务失败",
			chapterID: "chapter-123",
			requestBody: map[string]interface{}{
				"content": "这是一段测试内容",
			},
			mockSetup: func(m *MockQualityAIClient) {
				m.On("AnalyzeReadingPower", mock.Anything, "chapter-123", "这是一段测试内容", "", "").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockQualityAIClient)
			tt.mockSetup(mockClient)

			handler := &QualityHandler{
				aiClient: mockClient,
				logger:   logger,
			}
			router := setupTestRouter()

			router.POST("/api/v1/chapters/:id/analyze-reading-power", handler.AnalyzeReadingPower)

			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			url := "/api/v1/chapters/" + tt.chapterID + "/analyze-reading-power"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockClient.AssertExpectations(t)
		})
	}
}

// TestQualityHandler_CheckConsistency 测试一致性检查
func TestQualityHandler_CheckConsistency(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		chapterID      string
		requestBody    interface{}
		mockSetup      func(*MockQualityAIClient)
		expectedStatus int
	}{
		{
			name:      "成功检查一致性",
			chapterID: "chapter-123",
			requestBody: map[string]interface{}{
				"content": "这是一段测试内容",
				"characters": []map[string]interface{}{
					{"name": "张三", "role": "主角"},
				},
			},
			mockSetup: func(m *MockQualityAIClient) {
				result := map[string]interface{}{
					"consistent": true,
					"issues":     []interface{}{},
				}
				m.On("CheckConsistency", mock.Anything, "chapter-123", "这是一段测试内容", mock.AnythingOfType("map[string]interface {}"), "", "").Return(result, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "缺少必需字段",
			chapterID: "chapter-123",
			requestBody: map[string]interface{}{
				"characters": []map[string]interface{}{},
			},
			mockSetup:      func(m *MockQualityAIClient) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "AI 服务失败",
			chapterID: "chapter-123",
			requestBody: map[string]interface{}{
				"content": "这是一段测试内容",
			},
			mockSetup: func(m *MockQualityAIClient) {
				m.On("CheckConsistency", mock.Anything, "chapter-123", "这是一段测试内容", mock.AnythingOfType("map[string]interface {}"), "", "").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockQualityAIClient)
			tt.mockSetup(mockClient)

			handler := &QualityHandler{
				aiClient: mockClient,
				logger:   logger,
			}
			router := setupTestRouter()

			router.POST("/api/v1/chapters/:id/check-consistency", handler.CheckConsistency)

			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			url := "/api/v1/chapters/" + tt.chapterID + "/check-consistency"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockClient.AssertExpectations(t)
		})
	}
}

// TestQualityHandler_MultiAgentReview 测试多Agent审查
func TestQualityHandler_MultiAgentReview(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		chapterID      string
		requestBody    interface{}
		mockSetup      func(*MockQualityAIClient)
		expectedStatus int
	}{
		{
			name:      "成功多Agent审查",
			chapterID: "chapter-123",
			requestBody: map[string]interface{}{
				"content": "这是一段测试内容",
			},
			mockSetup: func(m *MockQualityAIClient) {
				result := map[string]interface{}{
					"reviews": []interface{}{
						map[string]interface{}{"agent": "agent1", "score": 85},
					},
				}
				m.On("MultiAgentReview", mock.Anything, "chapter-123", "这是一段测试内容", "", "").Return(result, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "缺少必需字段",
			chapterID: "chapter-123",
			requestBody: map[string]interface{}{
				"provider": "openai",
			},
			mockSetup:      func(m *MockQualityAIClient) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "AI 服务失败",
			chapterID: "chapter-123",
			requestBody: map[string]interface{}{
				"content": "这是一段测试内容",
			},
			mockSetup: func(m *MockQualityAIClient) {
				m.On("MultiAgentReview", mock.Anything, "chapter-123", "这是一段测试内容", "", "").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockQualityAIClient)
			tt.mockSetup(mockClient)

			handler := &QualityHandler{
				aiClient: mockClient,
				logger:   logger,
			}
			router := setupTestRouter()

			router.POST("/api/v1/chapters/:id/multi-agent-review", handler.MultiAgentReview)

			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			url := "/api/v1/chapters/" + tt.chapterID + "/multi-agent-review"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockClient.AssertExpectations(t)
		})
	}
}

// TestQualityHandler_EvaluateQuality 测试综合质量评估
func TestQualityHandler_EvaluateQuality(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		chapterID      string
		requestBody    interface{}
		mockSetup      func(*MockQualityAIClient)
		expectedStatus int
	}{
		{
			name:      "成功综合质量评估",
			chapterID: "chapter-123",
			requestBody: map[string]interface{}{
				"content": "这是一段测试内容",
				"characters": []map[string]interface{}{
					{"name": "张三", "role": "主角"},
				},
			},
			mockSetup: func(m *MockQualityAIClient) {
				result := map[string]interface{}{
					"overall_score": 85,
					"dimensions": map[string]interface{}{
						"plot":      80,
						"character": 90,
					},
				}
				m.On("EvaluateQuality", mock.Anything, "chapter-123", "这是一段测试内容", mock.AnythingOfType("map[string]interface {}"), "", "").Return(result, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "缺少必需字段",
			chapterID: "chapter-123",
			requestBody: map[string]interface{}{
				"characters": []map[string]interface{}{},
			},
			mockSetup:      func(m *MockQualityAIClient) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "AI 服务失败",
			chapterID: "chapter-123",
			requestBody: map[string]interface{}{
				"content": "这是一段测试内容",
			},
			mockSetup: func(m *MockQualityAIClient) {
				m.On("EvaluateQuality", mock.Anything, "chapter-123", "这是一段测试内容", mock.AnythingOfType("map[string]interface {}"), "", "").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockQualityAIClient)
			tt.mockSetup(mockClient)

			handler := &QualityHandler{
				aiClient: mockClient,
				logger:   logger,
			}
			router := setupTestRouter()

			router.POST("/api/v1/chapters/:id/evaluate-quality", handler.EvaluateQuality)

			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			url := "/api/v1/chapters/" + tt.chapterID + "/evaluate-quality"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockClient.AssertExpectations(t)
		})
	}
}
