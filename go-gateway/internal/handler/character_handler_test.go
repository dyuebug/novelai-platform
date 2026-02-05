package handler

import (
	"bytes"
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

// MockCharacterService 是 CharacterService 的 mock
type MockCharacterService struct {
	mock.Mock
}

func (m *MockCreate) Create(ctx context.Context, userID, projectID uuid.UUID, req *service.CreateCharacterRequest) (*model.Character, error) {
	args := m.Called(ctx, userID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Character), args.Error(1)
}

func (m *MockList) List(ctx context.Context, userID, projectID uuid.UUID, page, pageSize int, role, sort string) (*service.CharacterListResponse, error) {
	args := m.Called(ctx, userID, projectID, page, pageSize, role, sort)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.CharacterListResponse), args.Error(1)
}

func (m *MockGet) Get(ctx context.Context, userID, characterID uuid.UUID) (*model.Character, error) {
	args := m.Called(ctx, userID, characterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Character), args.Error(1)
}

func (m *MockUpdate) Update(ctx context.Context, userID, characterID uuid.UUID, req *service.UpdateCharacterRequest) (*model.Character, error) {
	args := m.Called(ctx, userID, characterID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Character), args.Error(1)
}

func (m *MockDelete) Delete(ctx context.Context, userID, characterID uuid.UUID) error {
	args := m.Called(ctx, userID, characterID)
	return args.Error(0)
}

func (m *MockCreateRelationship) CreateRelationship(ctx context.Context, userID, characterID uuid.UUID, req *service.CreateRelationshipRequest) (*model.CharacterRelationship, error) {
	args := m.Called(ctx, userID, characterID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CharacterRelationship), args.Error(1)
}

func (m *MockGetRelationships) GetRelationships(ctx context.Context, userID, characterID uuid.UUID) ([]model.CharacterRelationship, error) {
	args := m.Called(ctx, userID, characterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.CharacterRelationship), args.Error(1)
}

func (m *MockDeleteRelationship) DeleteRelationship(ctx context.Context, userID, relID uuid.UUID) error {
	args := m.Called(ctx, userID, relID)
	return args.Error(0)
}

func (m *MockCreateExperience) CreateExperience(ctx context.Context, userID, characterID uuid.UUID, req *service.CreateExperienceRequest) (*model.CharacterExperience, error) {
	args := m.Called(ctx, userID, characterID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CharacterExperience), args.Error(1)
}

func (m *MockGetExperiences) GetExperiences(ctx context.Context, userID, characterID uuid.UUID) ([]model.CharacterExperience, error) {
	args := m.Called(ctx, userID, characterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.CharacterExperience), args.Error(1)
}

func (m *MockDeleteExperience) DeleteExperience(ctx context.Context, userID, expID uuid.UUID) error {
	args := m.Called(ctx, userID, expID)
	return args.Error(0)
}

// TestCharacterHandler_Create 测试创建角色 API
func TestCharacterHandler_Create(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()
	characterID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockCharacterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功创建角色",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"name":        "张三",
				"role":        "protagonist",
				"description": "主角描述",
			},
			setupAuth: true,
			mockSetup: func(m *MockCharacterService) {
				character := &model.Character{
					ID:         characterID,
					ProjectID:  projectID,
					Name:       "张三",
					Role:       "protagonist",
					Background: "主角描述",
				}
				m.On("Create", mock.Anything, userID, projectID, mock.AnythingOfType("*service.CreateCharacterRequest")).Return(character, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:      "未认证",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"name": "张三",
			},
			setupAuth:      false,
			mockSetup:      func(m *MockCharacterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的项目ID",
			projectID:      "invalid-uuid",
			requestBody:    map[string]interface{}{"name": "张三"},
			setupAuth:      true,
			mockSetup:      func(m *MockCharacterService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"name": "张三",
			},
			setupAuth: true,
			mockSetup: func(m *MockCharacterService) {
				m.On("Create", mock.Anything, userID, projectID, mock.AnythingOfType("*service.CreateCharacterRequest")).Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCharacterService)
			tt.mockSetup(mockService)

			handler := &CharacterHandler{
				characterService: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/projects/:id/characters", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Create(c)
			})

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			url := "/api/v1/projects/" + tt.projectID + "/characters"
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

// TestCharacterHandler_List 测试获取角色列表 API
func TestCharacterHandler_List(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		query          string
		setupAuth      bool
		mockSetup      func(*MockCharacterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取角色列表",
			projectID: projectID.String(),
			query:     "page=1&page_size=50&role=protagonist",
			setupAuth: true,
			mockSetup: func(m *MockCharacterService) {
				resp := &service.CharacterListResponse{
					Items: []model.Character{
						{
							ID:        uuid.New(),
							ProjectID: projectID,
							Name:      "张三",
							Role:      "protagonist",
						},
					},
					Total:      1,
					Page:       1,
					PageSize:   50,
					TotalPages: 1,
				}
				m.On("List", mock.Anything, userID, projectID, 1, 50, "protagonist", "").Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			query:          "page=1&page_size=50",
			setupAuth:      false,
			mockSetup:      func(m *MockCharacterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			query:     "page=1&page_size=50",
			setupAuth: true,
			mockSetup: func(m *MockCharacterService) {
				m.On("List", mock.Anything, userID, projectID, 1, 50, "", "").Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCharacterService)
			tt.mockSetup(mockService)

			handler := &CharacterHandler{
				characterService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/projects/:id/characters", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.List(c)
			})

			url := "/api/v1/projects/" + tt.projectID + "/characters"
			if tt.query != "" {
				url += "?" + tt.query
			}
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

// TestCharacterHandler_Get 测试获取单个角色 API
func TestCharacterHandler_Get(t *testing.T) {
	userID := uuid.New()
	characterID := uuid.New()

	tests := []struct {
		name           string
		characterID    string
		setupAuth      bool
		mockSetup      func(*MockCharacterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "成功获取角色",
			characterID: characterID.String(),
			setupAuth:   true,
			mockSetup: func(m *MockCharacterService) {
				character := &model.Character{
					ID:         characterID,
					Name:       "张三",
					Role:       "protagonist",
					Background: "主角描述",
				}
				m.On("Get", mock.Anything, userID, characterID).Return(character, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			characterID:    characterID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockCharacterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的角色ID",
			characterID:    "invalid-uuid",
			setupAuth:      true,
			mockSetup:      func(m *MockCharacterService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:        "角色不存在",
			characterID: characterID.String(),
			setupAuth:   true,
			mockSetup: func(m *MockCharacterService) {
				m.On("Get", mock.Anything, userID, characterID).Return(nil, repository.ErrCharacterNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCharacterService)
			tt.mockSetup(mockService)

			handler := &CharacterHandler{
				characterService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/characters/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Get(c)
			})

			url := "/api/v1/characters/" + tt.characterID
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

// TestCharacterHandler_Update 测试更新角色 API
func TestCharacterHandler_Update(t *testing.T) {
	userID := uuid.New()
	characterID := uuid.New()

	tests := []struct {
		name           string
		characterID    string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockCharacterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "成功更新角色",
			characterID: characterID.String(),
			requestBody: map[string]interface{}{
				"name":        "李四",
				"description": "更新后的描述",
			},
			setupAuth: true,
			mockSetup: func(m *MockCharacterService) {
				character := &model.Character{
					ID:         characterID,
					Name:       "李四",
					Background: "更新后的描述",
				}
				m.On("Update", mock.Anything, userID, characterID, mock.AnythingOfType("*service.UpdateCharacterRequest")).Return(character, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "未认证",
			characterID: characterID.String(),
			requestBody: map[string]interface{}{
				"name": "李四",
			},
			setupAuth:      false,
			mockSetup:      func(m *MockCharacterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:        "角色不存在",
			characterID: characterID.String(),
			requestBody: map[string]interface{}{
				"name": "李四",
			},
			setupAuth: true,
			mockSetup: func(m *MockCharacterService) {
				m.On("Update", mock.Anything, userID, characterID, mock.AnythingOfType("*service.UpdateCharacterRequest")).Return(nil, repository.ErrCharacterNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCharacterService)
			tt.mockSetup(mockService)

			handler := &CharacterHandler{
				characterService: mockService,
			}
			router := setupTestRouter()

			router.PUT("/api/v1/characters/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Update(c)
			})

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			url := "/api/v1/characters/" + tt.characterID
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
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

// TestCharacterHandler_Delete 测试删除角色 API
func TestCharacterHandler_Delete(t *testing.T) {
	userID := uuid.New()
	characterID := uuid.New()

	tests := []struct {
		name           string
		characterID    string
		setupAuth      bool
		mockSetup      func(*MockCharacterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "成功删除角色",
			characterID: characterID.String(),
			setupAuth:   true,
			mockSetup: func(m *MockCharacterService) {
				m.On("Delete", mock.Anything, userID, characterID).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			characterID:    characterID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockCharacterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:        "角色不存在",
			characterID: characterID.String(),
			setupAuth:   true,
			mockSetup: func(m *MockCharacterService) {
				m.On("Delete", mock.Anything, userID, characterID).Return(repository.ErrCharacterNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCharacterService)
			tt.mockSetup(mockService)

			handler := &CharacterHandler{
				characterService: mockService,
			}
			router := setupTestRouter()

			router.DELETE("/api/v1/characters/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Delete(c)
			})

			url := "/api/v1/characters/" + tt.characterID
			req := httptest.NewRequest(http.MethodDelete, url, nil)
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
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestCharacterHandler_CreateRelationship 测试创建角色关系 API
func TestCharacterHandler_CreateRelationship(t *testing.T) {
	userID := uuid.New()
	characterID := uuid.New()
	relatedCharacterID := uuid.New()
	relationshipID := uuid.New()

	tests := []struct {
		name           string
		characterID    string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockCharacterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "成功创建角色关系",
			characterID: characterID.String(),
			requestBody: map[string]interface{}{
				"target_id":     relatedCharacterID.String(),
				"relation_type": "friend",
				"description":   "好友关系",
			},
			setupAuth: true,
			mockSetup: func(m *MockCharacterService) {
				rel := &model.CharacterRelationship{
					ID:           relationshipID,
					CharacterID:  characterID,
					TargetID:     relatedCharacterID,
					RelationType: "friend",
					Description:  "好友关系",
				}
				m.On("CreateRelationship", mock.Anything, userID, characterID, mock.AnythingOfType("*service.CreateRelationshipRequest")).Return(rel, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:        "未认证",
			characterID: characterID.String(),
			requestBody: map[string]interface{}{
				"related_character_id": relatedCharacterID.String(),
			},
			setupAuth:      false,
			mockSetup:      func(m *MockCharacterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCharacterService)
			tt.mockSetup(mockService)

			handler := &CharacterHandler{
				characterService: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/characters/:id/relationships", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.CreateRelationship(c)
			})

			var body []byte
			var err error
			body, err = json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			url := "/api/v1/characters/" + tt.characterID + "/relationships"
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

// TestCharacterHandler_GetRelationships 测试获取角色关系列表 API
func TestCharacterHandler_GetRelationships(t *testing.T) {
	userID := uuid.New()
	characterID := uuid.New()

	tests := []struct {
		name           string
		characterID    string
		setupAuth      bool
		mockSetup      func(*MockCharacterService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "成功获取角色关系列表",
			characterID: characterID.String(),
			setupAuth:   true,
			mockSetup: func(m *MockCharacterService) {
				rels := []model.CharacterRelationship{
					{
						ID:           uuid.New(),
						CharacterID:  characterID,
						RelationType: "friend",
					},
				}
				m.On("GetRelationships", mock.Anything, userID, characterID).Return(rels, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			characterID:    characterID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockCharacterService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCharacterService)
			tt.mockSetup(mockService)

			handler := &CharacterHandler{
				characterService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/characters/:id/relationships", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.GetRelationships(c)
			})

			url := "/api/v1/characters/" + tt.characterID + "/relationships"
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
