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

// MockLocationService 是 LocationService 的 mock
type MockLocationService struct {
	mock.Mock
}

func (m *MockCreate) Create(ctx context.Context, userID, projectID uuid.UUID, req *service.CreateLocationRequest) (*model.Location, error) {
	args := m.Called(ctx, userID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Location), args.Error(1)
}

func (m *MockList) List(ctx context.Context, userID, projectID uuid.UUID, page, pageSize int, locationType, sort string) (*service.LocationListResponse, error) {
	args := m.Called(ctx, userID, projectID, page, pageSize, locationType, sort)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.LocationListResponse), args.Error(1)
}

func (m *MockGetTree) GetTree(ctx context.Context, userID, projectID uuid.UUID) ([]model.Location, error) {
	args := m.Called(ctx, userID, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Location), args.Error(1)
}

func (m *MockGet) Get(ctx context.Context, userID, locationID uuid.UUID) (*model.Location, error) {
	args := m.Called(ctx, userID, locationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Location), args.Error(1)
}

func (m *MockUpdate) Update(ctx context.Context, userID, locationID uuid.UUID, req *service.UpdateLocationRequest) (*model.Location, error) {
	args := m.Called(ctx, userID, locationID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Location), args.Error(1)
}

func (m *MockDelete) Delete(ctx context.Context, userID, locationID uuid.UUID) error {
	args := m.Called(ctx, userID, locationID)
	return args.Error(0)
}

// MockOrganizationService 是 OrganizationService 的 mock
type MockOrganizationService struct {
	mock.Mock
}

func (m *MockCreate) Create(ctx context.Context, userID, projectID uuid.UUID, req *service.CreateOrganizationRequest) (*model.Organization, error) {
	args := m.Called(ctx, userID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Organization), args.Error(1)
}

func (m *MockList) List(ctx context.Context, userID, projectID uuid.UUID, page, pageSize int, orgType, sort string) (*service.OrganizationListResponse, error) {
	args := m.Called(ctx, userID, projectID, page, pageSize, orgType, sort)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.OrganizationListResponse), args.Error(1)
}

func (m *MockGet) Get(ctx context.Context, userID, orgID uuid.UUID) (*model.Organization, error) {
	args := m.Called(ctx, userID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Organization), args.Error(1)
}

func (m *MockUpdate) Update(ctx context.Context, userID, orgID uuid.UUID, req *service.UpdateOrganizationRequest) (*model.Organization, error) {
	args := m.Called(ctx, userID, orgID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Organization), args.Error(1)
}

func (m *MockDelete) Delete(ctx context.Context, userID, orgID uuid.UUID) error {
	args := m.Called(ctx, userID, orgID)
	return args.Error(0)
}

func (m *MockAddMember) AddMember(ctx context.Context, userID, orgID uuid.UUID, req *service.AddMemberRequest) (*model.OrganizationMember, error) {
	args := m.Called(ctx, userID, orgID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrganizationMember), args.Error(1)
}

func (m *MockGetMembers) GetMembers(ctx context.Context, userID, orgID uuid.UUID) ([]model.OrganizationMember, error) {
	args := m.Called(ctx, userID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.OrganizationMember), args.Error(1)
}

func (m *MockRemoveMember) RemoveMember(ctx context.Context, userID, memberID uuid.UUID) error {
	args := m.Called(ctx, userID, memberID)
	return args.Error(0)
}

// TestLocationHandler_Create 测试创建地点 API
func TestLocationHandler_Create(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()
	locationID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockLocationService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功创建地点",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"name":        "天都城",
				"type":        "city",
				"description": "主城描述",
			},
			setupAuth: true,
			mockSetup: func(m *MockLocationService) {
				location := &model.Location{
					ID:          locationID,
					ProjectID:   projectID,
					Name:        "天都城",
					Type:        "city",
					Description: "主城描述",
				}
				m.On("Create", mock.Anything, userID, projectID, mock.AnythingOfType("*service.CreateLocationRequest")).Return(location, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:      "未认证",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"name": "天都城",
			},
			setupAuth:      false,
			mockSetup:      func(m *MockLocationService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的项目ID",
			projectID:      "invalid-uuid",
			requestBody:    map[string]interface{}{"name": "天都城"},
			setupAuth:      true,
			mockSetup:      func(m *MockLocationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"name": "天都城",
			},
			setupAuth: true,
			mockSetup: func(m *MockLocationService) {
				m.On("Create", mock.Anything, userID, projectID, mock.AnythingOfType("*service.CreateLocationRequest")).Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockLocationService)
			tt.mockSetup(mockService)

			handler := &LocationHandler{
				locationService: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/projects/:id/locations", func(c *gin.Context) {
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

			url := "/api/v1/projects/" + tt.projectID + "/locations"
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

// TestLocationHandler_List 测试获取地点列表 API
func TestLocationHandler_List(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		query          string
		setupAuth      bool
		mockSetup      func(*MockLocationService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取地点列表",
			projectID: projectID.String(),
			query:     "page=1&page_size=50&type=city",
			setupAuth: true,
			mockSetup: func(m *MockLocationService) {
				resp := &service.LocationListResponse{
					Items: []model.Location{
						{
							ID:        uuid.New(),
							ProjectID: projectID,
							Name:      "天都城",
							Type:      "city",
						},
					},
					Total:      1,
					Page:       1,
					PageSize:   50,
					TotalPages: 1,
				}
				m.On("List", mock.Anything, userID, projectID, 1, 50, "city", "").Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			query:          "page=1&page_size=50",
			setupAuth:      false,
			mockSetup:      func(m *MockLocationService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			query:     "page=1&page_size=50",
			setupAuth: true,
			mockSetup: func(m *MockLocationService) {
				m.On("List", mock.Anything, userID, projectID, 1, 50, "", "").Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockLocationService)
			tt.mockSetup(mockService)

			handler := &LocationHandler{
				locationService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/projects/:id/locations", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.List(c)
			})

			url := "/api/v1/projects/" + tt.projectID + "/locations"
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

// TestLocationHandler_GetTree 测试获取地点树结构 API
func TestLocationHandler_GetTree(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		setupAuth      bool
		mockSetup      func(*MockLocationService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取地点树",
			projectID: projectID.String(),
			setupAuth: true,
			mockSetup: func(m *MockLocationService) {
				tree := []model.Location{
					{
						ID:   uuid.New(),
						Name: "天都城",
						Type: "city",
					},
				}
				m.On("GetTree", mock.Anything, userID, projectID).Return(tree, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			projectID:      projectID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockLocationService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			setupAuth: true,
			mockSetup: func(m *MockLocationService) {
				m.On("GetTree", mock.Anything, userID, projectID).Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockLocationService)
			tt.mockSetup(mockService)

			handler := &LocationHandler{
				locationService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/projects/:id/locations/tree", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.GetTree(c)
			})

			url := "/api/v1/projects/" + tt.projectID + "/locations/tree"
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

// TestLocationHandler_Get 测试获取单个地点 API
func TestLocationHandler_Get(t *testing.T) {
	userID := uuid.New()
	locationID := uuid.New()

	tests := []struct {
		name           string
		locationID     string
		setupAuth      bool
		mockSetup      func(*MockLocationService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:       "成功获取地点",
			locationID: locationID.String(),
			setupAuth:  true,
			mockSetup: func(m *MockLocationService) {
				location := &model.Location{
					ID:          locationID,
					Name:        "天都城",
					Type:        "city",
					Description: "主城描述",
				}
				m.On("Get", mock.Anything, userID, locationID).Return(location, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			locationID:     locationID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockLocationService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "无效的地点ID",
			locationID:     "invalid-uuid",
			setupAuth:      true,
			mockSetup:      func(m *MockLocationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:       "地点不存在",
			locationID: locationID.String(),
			setupAuth:  true,
			mockSetup: func(m *MockLocationService) {
				m.On("Get", mock.Anything, userID, locationID).Return(nil, repository.ErrLocationNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockLocationService)
			tt.mockSetup(mockService)

			handler := &LocationHandler{
				locationService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/locations/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Get(c)
			})

			url := "/api/v1/locations/" + tt.locationID
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

// TestLocationHandler_Delete 测试删除地点 API
func TestLocationHandler_Delete(t *testing.T) {
	userID := uuid.New()
	locationID := uuid.New()

	tests := []struct {
		name           string
		locationID     string
		setupAuth      bool
		mockSetup      func(*MockLocationService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:       "成功删除地点",
			locationID: locationID.String(),
			setupAuth:  true,
			mockSetup: func(m *MockLocationService) {
				m.On("Delete", mock.Anything, userID, locationID).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			locationID:     locationID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockLocationService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:       "地点不存在",
			locationID: locationID.String(),
			setupAuth:  true,
			mockSetup: func(m *MockLocationService) {
				m.On("Delete", mock.Anything, userID, locationID).Return(repository.ErrLocationNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockLocationService)
			tt.mockSetup(mockService)

			handler := &LocationHandler{
				locationService: mockService,
			}
			router := setupTestRouter()

			router.DELETE("/api/v1/locations/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Delete(c)
			})

			url := "/api/v1/locations/" + tt.locationID
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

// TestOrganizationHandler_Create 测试创建组织 API
func TestOrganizationHandler_Create(t *testing.T) {
	userID := uuid.New()
	projectID := uuid.New()
	orgID := uuid.New()

	tests := []struct {
		name           string
		projectID      string
		requestBody    interface{}
		setupAuth      bool
		mockSetup      func(*MockOrganizationService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功创建组织",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"name":        "天剑门",
				"type":        "sect",
				"description": "正道门派",
			},
			setupAuth: true,
			mockSetup: func(m *MockOrganizationService) {
				org := &model.Organization{
					ID:          orgID,
					ProjectID:   projectID,
					Name:        "天剑门",
					Type:        "sect",
					Description: "正道门派",
				}
				m.On("Create", mock.Anything, userID, projectID, mock.AnythingOfType("*service.CreateOrganizationRequest")).Return(org, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:      "未认证",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"name": "天剑门",
			},
			setupAuth:      false,
			mockSetup:      func(m *MockOrganizationService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "项目不存在",
			projectID: projectID.String(),
			requestBody: map[string]interface{}{
				"name": "天剑门",
			},
			setupAuth: true,
			mockSetup: func(m *MockOrganizationService) {
				m.On("Create", mock.Anything, userID, projectID, mock.AnythingOfType("*service.CreateOrganizationRequest")).Return(nil, repository.ErrProjectNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockOrganizationService)
			tt.mockSetup(mockService)

			handler := &OrganizationHandler{
				orgService: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/projects/:id/organizations", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Create(c)
			})

			var body []byte
			var err error
			body, err = json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			url := "/api/v1/projects/" + tt.projectID + "/organizations"
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

// TestOrganizationHandler_Get 测试获取单个组织 API
func TestOrganizationHandler_Get(t *testing.T) {
	userID := uuid.New()
	orgID := uuid.New()

	tests := []struct {
		name           string
		orgID          string
		setupAuth      bool
		mockSetup      func(*MockOrganizationService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取组织",
			orgID:     orgID.String(),
			setupAuth: true,
			mockSetup: func(m *MockOrganizationService) {
				org := &model.Organization{
					ID:          orgID,
					Name:        "天剑门",
					Type:        "sect",
					Description: "正道门派",
				}
				m.On("Get", mock.Anything, userID, orgID).Return(org, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "未认证",
			orgID:          orgID.String(),
			setupAuth:      false,
			mockSetup:      func(m *MockOrganizationService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:      "组织不存在",
			orgID:     orgID.String(),
			setupAuth: true,
			mockSetup: func(m *MockOrganizationService) {
				m.On("Get", mock.Anything, userID, orgID).Return(nil, repository.ErrOrganizationNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockOrganizationService)
			tt.mockSetup(mockService)

			handler := &OrganizationHandler{
				orgService: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/organizations/:id", func(c *gin.Context) {
				if tt.setupAuth {
					setUserIDContext(c, userID)
				}
				handler.Get(c)
			})

			url := "/api/v1/organizations/" + tt.orgID
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
