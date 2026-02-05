package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/grpcclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAIService 是 AIService 的 mock
type MockAIService struct {
	mock.Mock
}

// MockStream 是 Stream 的 mock
type MockStream struct {
	mock.Mock
	responses []*grpcclient.GenerateResponse
	index     int
}

func (m *MockStream) Recv() (*grpcclient.GenerateResponse, error) {
	if m.index >= len(m.responses) {
		return nil, io.EOF
	}
	resp := m.responses[m.index]
	m.index++
	return resp, nil
}

func (m *MockStream) CloseSend() error {
	return nil
}

func (m *MockAIService) GenerateStream(ctx context.Context, req *grpcclient.GenerateRequest) (StreamInterface, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(StreamInterface), args.Error(1)
}

// TestStreamHandler_Generate 测试流式生成 API
func TestStreamHandler_Generate(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockSetup      func(*MockAIService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:  "缺少 prompt 参数",
			query: "",
			mockSetup: func(m *MockAIService) {
				// 不需要 mock，因为会在参数验证阶段失败
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:  "AI 服务不可用",
			query: "prompt=测试提示词",
			mockSetup: func(m *MockAIService) {
				m.On("GenerateStream", mock.Anything, mock.AnythingOfType("*grpcclient.GenerateRequest")).Return(nil, errors.New("service unavailable"))
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAIService)
			tt.mockSetup(mockService)

			handler := &StreamHandler{
				ai: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/stream/generate", handler.Generate)

			url := "/api/v1/stream/generate"
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
