package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/service"
	"go-gateway/pkg/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOAuthService 是 OAuthService 的 mock
type MockOAuthService struct {
	mock.Mock
}

func (m *MockOAuthService) GetEnabledProviders() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockOAuthService) GetAuthorizationURL(provider, redirectURI string) (string, error) {
	args := m.Called(provider, redirectURI)
	return args.String(0), args.Error(1)
}

func (m *MockOAuthService) HandleCallback(ctx context.Context, provider, code, state string) (*service.TokenResponse, string, error) {
	args := m.Called(ctx, provider, code, state)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}
	return args.Get(0).(*service.TokenResponse), args.String(1), args.Error(2)
}

// TestOAuthHandler_GetProviders 测试获取已启用的 OAuth 提供商
func TestOAuthHandler_GetProviders(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockOAuthService)
		expectedStatus int
	}{
		{
			name: "成功获取提供商列表",
			mockSetup: func(m *MockOAuthService) {
				m.On("GetEnabledProviders").Return([]string{"github", "google"})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "空提供商列表",
			mockSetup: func(m *MockOAuthService) {
				m.On("GetEnabledProviders").Return([]string{})
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockOAuthService)
			tt.mockSetup(mockService)

			handler := &OAuthHandler{
				oauthService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/oauth/providers", handler.GetProviders)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/oauth/providers", nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp response.Envelope
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, 0, resp.Code)
			assert.NotNil(t, resp.Data)

			mockService.AssertExpectations(t)
		})
	}
}

// TestOAuthHandler_Authorize 测试获取授权 URL
func TestOAuthHandler_Authorize(t *testing.T) {
	tests := []struct {
		name           string
		provider       string
		redirectURI    string
		mockSetup      func(*MockOAuthService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "成功获取授权 URL",
			provider:    "github",
			redirectURI: "http://localhost:3000/callback",
			mockSetup: func(m *MockOAuthService) {
				m.On("GetAuthorizationURL", "github", "http://localhost:3000/callback").Return("https://github.com/login/oauth/authorize?client_id=xxx", nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "提供商未启用",
			provider:    "facebook",
			redirectURI: "http://localhost:3000/callback",
			mockSetup: func(m *MockOAuthService) {
				m.On("GetAuthorizationURL", "facebook", "http://localhost:3000/callback").Return("", service.ErrOAuthProviderNotEnabled)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:        "内部错误",
			provider:    "github",
			redirectURI: "http://localhost:3000/callback",
			mockSetup: func(m *MockOAuthService) {
				m.On("GetAuthorizationURL", "github", "http://localhost:3000/callback").Return("", assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockOAuthService)
			tt.mockSetup(mockService)

			handler := &OAuthHandler{
				oauthService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/oauth/:provider/authorize", handler.Authorize)

			url := "/api/v1/oauth/" + tt.provider + "/authorize?redirect_uri=" + tt.redirectURI
			req := httptest.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

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

// TestOAuthHandler_Callback 测试 OAuth 回调
func TestOAuthHandler_Callback(t *testing.T) {
	tests := []struct {
		name           string
		provider       string
		code           string
		state          string
		errorParam     string
		mockSetup      func(*MockOAuthService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:     "成功回调（返回 JSON）",
			provider: "github",
			code:     "valid_code",
			state:    "valid_state",
			mockSetup: func(m *MockOAuthService) {
				tokens := &service.TokenResponse{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
					TokenType:    "Bearer",
					ExpiresIn:    3600,
				}
				m.On("HandleCallback", mock.Anything, "github", "valid_code", "valid_state").Return(tokens, "", nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "OAuth 错误参数",
			provider:       "github",
			errorParam:     "access_denied",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup:      func(m *MockOAuthService) {},
		},
		{
			name:           "缺少授权码",
			provider:       "github",
			state:          "valid_state",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup:      func(m *MockOAuthService) {},
		},
		{
			name:           "缺少 state 参数",
			provider:       "github",
			code:           "valid_code",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup:      func(m *MockOAuthService) {},
		},
		{
			name:     "无效的 state",
			provider: "github",
			code:     "valid_code",
			state:    "invalid_state",
			mockSetup: func(m *MockOAuthService) {
				m.On("HandleCallback", mock.Anything, "github", "valid_code", "invalid_state").Return(nil, "", service.ErrOAuthInvalidState)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:     "无效的授权码",
			provider: "github",
			code:     "invalid_code",
			state:    "valid_state",
			mockSetup: func(m *MockOAuthService) {
				m.On("HandleCallback", mock.Anything, "github", "invalid_code", "valid_state").Return(nil, "", service.ErrOAuthInvalidCode)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:     "获取用户信息失败",
			provider: "github",
			code:     "valid_code",
			state:    "valid_state",
			mockSetup: func(m *MockOAuthService) {
				m.On("HandleCallback", mock.Anything, "github", "valid_code", "valid_state").Return(nil, "", service.ErrOAuthUserInfoFailed)
			},
			expectedStatus: http.StatusBadGateway,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockOAuthService)
			tt.mockSetup(mockService)

			handler := &OAuthHandler{
				oauthService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/oauth/:provider/callback", handler.Callback)

			url := "/api/v1/oauth/" + tt.provider + "/callback?"
			if tt.code != "" {
				url += "code=" + tt.code + "&"
			}
			if tt.state != "" {
				url += "state=" + tt.state + "&"
			}
			if tt.errorParam != "" {
				url += "error=" + tt.errorParam + "&"
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				var resp response.Envelope
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotEqual(t, 0, resp.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestOAuthHandler_CallbackPost 测试 POST 方式的 OAuth 回调
func TestOAuthHandler_CallbackPost(t *testing.T) {
	tests := []struct {
		name           string
		provider       string
		requestBody    interface{}
		mockSetup      func(*MockOAuthService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:     "成功回调",
			provider: "github",
			requestBody: map[string]interface{}{
				"code":  "valid_code",
				"state": "valid_state",
			},
			mockSetup: func(m *MockOAuthService) {
				tokens := &service.TokenResponse{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
					TokenType:    "Bearer",
					ExpiresIn:    3600,
				}
				m.On("HandleCallback", mock.Anything, "github", "valid_code", "valid_state").Return(tokens, "", nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:     "无效的请求体",
			provider: "github",
			requestBody: map[string]interface{}{
				"code": "valid_code",
			},
			mockSetup:      func(m *MockOAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:     "无效的 state",
			provider: "github",
			requestBody: map[string]interface{}{
				"code":  "valid_code",
				"state": "invalid_state",
			},
			mockSetup: func(m *MockOAuthService) {
				m.On("HandleCallback", mock.Anything, "github", "valid_code", "invalid_state").Return(nil, "", service.ErrOAuthInvalidState)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockOAuthService)
			tt.mockSetup(mockService)

			handler := &OAuthHandler{
				oauthService: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/oauth/:provider/callback", handler.CallbackPost)

			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			url := "/api/v1/oauth/" + tt.provider + "/callback"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

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
