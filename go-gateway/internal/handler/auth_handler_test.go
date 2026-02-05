package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService 是 AuthService 的 mock
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req *service.RegisterRequest) (*model.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req *service.LoginRequest) (*service.TokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.TokenResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*service.TokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.TokenResponse), args.Error(1)
}

func (m *MockAuthService) RequestPasswordReset(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockAuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(0)
}

// TestAuthHandler_Register 测试用户注册 API
func TestAuthHandler_Register(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "成功注册",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "Test1234",
			},
			mockSetup: func(m *MockAuthService) {
				user := &model.User{
					ID:        userID,
					Username:  "testuser",
					Email:     "test@example.com",
					Role:      "user",
					CreatedAt: time.Now(),
				}
				m.On("Register", mock.Anything, mock.AnythingOfType("*service.RegisterRequest")).Return(user, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "用户已存在",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "Test1234",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Register", mock.Anything, mock.AnythingOfType("*service.RegisterRequest")).Return(nil, repository.ErrUserAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
		},
		{
			name: "密码过弱",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "weakpass",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Register", mock.Anything, mock.AnythingOfType("*service.RegisterRequest")).Return(nil, service.ErrWeakPassword)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "无效的请求体",
			requestBody:    "invalid json",
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.mockSetup(mockService)

			handler := &AuthHandler{
				authService: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/auth/register", handler.Register)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
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

// TestAuthHandler_Login 测试用户登录 API
func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "成功登录",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "Test1234",
			},
			mockSetup: func(m *MockAuthService) {
				tokens := &service.TokenResponse{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
					TokenType:    "Bearer",
					ExpiresIn:    3600,
				}
				m.On("Login", mock.Anything, mock.AnythingOfType("*service.LoginRequest")).Return(tokens, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "无效的凭证",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Login", mock.Anything, mock.AnythingOfType("*service.LoginRequest")).Return(nil, service.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name: "账户被锁定",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "Test1234",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Login", mock.Anything, mock.AnythingOfType("*service.LoginRequest")).Return(nil, service.ErrAccountLocked)
			},
			expectedStatus: http.StatusLocked,
			expectedError:  true,
		},
		{
			name:           "无效的请求体",
			requestBody:    "invalid json",
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.mockSetup(mockService)

			handler := &AuthHandler{
				authService: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/auth/login", handler.Login)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestAuthHandler_Refresh 测试刷新令牌 API
func TestAuthHandler_Refresh(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "成功刷新令牌",
			requestBody: map[string]interface{}{
				"refresh_token": "valid_refresh_token",
			},
			mockSetup: func(m *MockAuthService) {
				tokens := &service.TokenResponse{
					AccessToken:  "new_access_token",
					RefreshToken: "new_refresh_token",
					TokenType:    "Bearer",
					ExpiresIn:    3600,
				}
				m.On("RefreshToken", mock.Anything, "valid_refresh_token").Return(tokens, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "无效的刷新令牌",
			requestBody: map[string]interface{}{
				"refresh_token": "invalid_token",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("RefreshToken", mock.Anything, "invalid_token").Return(nil, repository.ErrTokenNotFound)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name: "令牌已过期",
			requestBody: map[string]interface{}{
				"refresh_token": "expired_token",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("RefreshToken", mock.Anything, "expired_token").Return(nil, repository.ErrTokenExpired)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name: "令牌已使用",
			requestBody: map[string]interface{}{
				"refresh_token": "used_token",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("RefreshToken", mock.Anything, "used_token").Return(nil, repository.ErrTokenUsed)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.mockSetup(mockService)

			handler := &AuthHandler{
				authService: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/auth/refresh", handler.Refresh)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestAuthHandler_RequestPasswordReset 测试请求密码重置 API
func TestAuthHandler_RequestPasswordReset(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "成功请求密码重置",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("RequestPasswordReset", mock.Anything, "test@example.com").Return(nil)
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name: "请求过于频繁",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("RequestPasswordReset", mock.Anything, "test@example.com").Return(service.ErrResetCooldown)
			},
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "无效的请求体",
			requestBody:    "invalid json",
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.mockSetup(mockService)

			handler := &AuthHandler{
				authService: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/auth/password-reset/request", handler.RequestPasswordReset)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/password-reset/request", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestAuthHandler_ResetPassword 测试重置密码 API
func TestAuthHandler_ResetPassword(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "成功重置密码",
			requestBody: map[string]interface{}{
				"token":        "valid_reset_token",
				"new_password": "NewPass1234",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("ResetPassword", mock.Anything, "valid_reset_token", "NewPass1234").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "无效的重置令牌",
			requestBody: map[string]interface{}{
				"token":        "invalid_token",
				"new_password": "NewPass1234",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("ResetPassword", mock.Anything, "invalid_token", "NewPass1234").Return(repository.ErrTokenNotFound)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "令牌已过期",
			requestBody: map[string]interface{}{
				"token":        "expired_token",
				"new_password": "NewPass1234",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("ResetPassword", mock.Anything, "expired_token", "NewPass1234").Return(repository.ErrTokenExpired)
			},
			expectedStatus: http.StatusGone,
			expectedError:  true,
		},
		{
			name: "密码过弱",
			requestBody: map[string]interface{}{
				"token":        "valid_token",
				"new_password": "weakpass",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("ResetPassword", mock.Anything, "valid_token", "weakpass").Return(service.ErrWeakPassword)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.mockSetup(mockService)

			handler := &AuthHandler{
				authService: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/auth/password-reset/verify", handler.ResetPassword)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/password-reset/verify", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestAuthHandler_GetCurrentUser 测试获取当前用户 API
func TestAuthHandler_GetCurrentUser(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "成功获取当前用户",
			setupContext: func(c *gin.Context) {
				claims := map[string]interface{}{
					"sub":      "user-123",
					"username": "testuser",
					"email":    "test@example.com",
				}
				c.Set("claims", claims)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			setupContext:   func(c *gin.Context) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)

			handler := &AuthHandler{
				authService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/auth/me", func(c *gin.Context) {
				tt.setupContext(c)
				handler.GetCurrentUser(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
