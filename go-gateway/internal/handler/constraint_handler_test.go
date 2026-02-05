package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/grpcclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockAIClient 是 AIClient 的 mock
type MockAIClient struct {
	mock.Mock
}

func (m *MockCheckConstraints) CheckConstraints(ctx context.Context, req *grpcclient.CheckConstraintsRequest) (interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func (m *MockRequestExemption) RequestExemption(ctx context.Context, req *grpcclient.RequestExemptionRequest) (interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func (m *MockRevokeExemption) RevokeExemption(ctx context.Context, req *grpcclient.RevokeExemptionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// TestConstraintHandler_CheckConstraints 测试检查约束 API
func TestConstraintHandler_CheckConstraints(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		projectID      string
		chapterID      string
		requestBody    interface{}
		mockSetup      func(*MockAIClient)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功检查约束",
			projectID: "project-123",
			chapterID: "chapter-456",
			requestBody: map[string]interface{}{
				"content": "这是一段测试内容",
				"constraints": []map[string]interface{}{
					{"type": "length", "max": 1000},
				},
			},
			mockSetup: func(m *MockAIClient) {
				result := map[string]interface{}{
					"passed":     true,
					"violations": []interface{}{},
				}
				m.On("CheckConstraints", mock.Anything, mock.AnythingOfType("*grpcclient.CheckConstraintsRequest")).Return(result, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:      "无效的请求体",
			projectID: "project-123",
			chapterID: "chapter-456",
			requestBody: map[string]interface{}{
				"constraints": []map[string]interface{}{},
			},
			mockSetup:      func(m *MockAIClient) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "约束检查失败",
			projectID: "project-123",
			chapterID: "chapter-456",
			requestBody: map[string]interface{}{
				"content": "这是一段测试内容",
				"constraints": []map[string]interface{}{
					{"type": "length", "max": 1000},
				},
			},
			mockSetup: func(m *MockAIClient) {
				m.On("CheckConstraints", mock.Anything, mock.AnythingOfType("*grpcclient.CheckConstraintsRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockAIClient)
			tt.mockSetup(mockClient)

			handler := &ConstraintHandler{
				aiClient: mockClient,
				logger:   logger,
			}
			router := setupTestRouter()

			router.POST("/api/v1/projects/:projectId/chapters/:chapterId/check-constraints", handler.CheckConstraints)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			url := "/api/v1/projects/" + tt.projectID + "/chapters/" + tt.chapterID + "/check-constraints"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				var resp map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotNil(t, resp["error"])
			}

			mockClient.AssertExpectations(t)
		})
	}
}

// TestConstraintHandler_RequestExemption 测试请求豁免 API
func TestConstraintHandler_RequestExemption(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		projectID      string
		chapterID      string
		requestBody    interface{}
		mockSetup      func(*MockAIClient)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功请求豁免",
			projectID: "project-123",
			chapterID: "chapter-456",
			requestBody: map[string]interface{}{
				"constraint_id": "constraint-789",
				"reason":        "特殊情况需要豁免",
			},
			mockSetup: func(m *MockAIClient) {
				result := map[string]interface{}{
					"exemption_id": "exemption-001",
					"status":       "pending",
				}
				m.On("RequestExemption", mock.Anything, mock.AnythingOfType("*grpcclient.RequestExemptionRequest")).Return(result, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:      "无效的请求体",
			projectID: "project-123",
			chapterID: "chapter-456",
			requestBody: map[string]interface{}{
				"constraint_id": "constraint-789",
			},
			mockSetup:      func(m *MockAIClient) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "请求豁免失败",
			projectID: "project-123",
			chapterID: "chapter-456",
			requestBody: map[string]interface{}{
				"constraint_id": "constraint-789",
				"reason":        "特殊情况需要豁免",
			},
			mockSetup: func(m *MockAIClient) {
				m.On("RequestExemption", mock.Anything, mock.AnythingOfType("*grpcclient.RequestExemptionRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockAIClient)
			tt.mockSetup(mockClient)

			handler := &ConstraintHandler{
				aiClient: mockClient,
				logger:   logger,
			}
			router := setupTestRouter()

			router.POST("/api/v1/projects/:projectId/chapters/:chapterId/exemptions", handler.RequestExemption)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			url := "/api/v1/projects/" + tt.projectID + "/chapters/" + tt.chapterID + "/exemptions"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				var resp map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotNil(t, resp["error"])
			}

			mockClient.AssertExpectations(t)
		})
	}
}

// TestConstraintHandler_RevokeExemption 测试撤销豁免 API
func TestConstraintHandler_RevokeExemption(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		exemptionID    string
		mockSetup      func(*MockAIClient)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "成功撤销豁免",
			exemptionID: "exemption-001",
			mockSetup: func(m *MockAIClient) {
				m.On("RevokeExemption", mock.Anything, mock.AnythingOfType("*grpcclient.RevokeExemptionRequest")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "撤销豁免失败",
			exemptionID: "exemption-001",
			mockSetup: func(m *MockAIClient) {
				m.On("RevokeExemption", mock.Anything, mock.AnythingOfType("*grpcclient.RevokeExemptionRequest")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockAIClient)
			tt.mockSetup(mockClient)

			handler := &ConstraintHandler{
				aiClient: mockClient,
				logger:   logger,
			}
			router := setupTestRouter()

			router.DELETE("/api/v1/exemptions/:exemptionId", handler.RevokeExemption)

			url := "/api/v1/exemptions/" + tt.exemptionID
			req := httptest.NewRequest(http.MethodDelete, url, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotNil(t, resp["error"])
			}

			mockClient.AssertExpectations(t)
		})
	}
}
