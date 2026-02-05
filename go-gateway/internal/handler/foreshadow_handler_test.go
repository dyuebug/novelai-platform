package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/model"
	"go-gateway/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockForeshadowService 是 ForeshadowService 的 mock
type MockForeshadowService struct {
	mock.Mock
}

func (m *MockForeshadowService) Create(projectID string, req *service.CreateForeshadowRequest) (*model.Foreshadow, error) {
	args := m.Called(projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Foreshadow), args.Error(1)
}

func (m *MockForeshadowService) GetByID(id string) (*model.Foreshadow, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Foreshadow), args.Error(1)
}

func (m *MockForeshadowService) Update(id string, req *service.UpdateForeshadowRequest) (*model.Foreshadow, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Foreshadow), args.Error(1)
}

func (m *MockForeshadowService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockForeshadowService) List(projectID, status, priority string, page, pageSize int) ([]model.Foreshadow, int64, error) {
	args := m.Called(projectID, status, priority, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.Foreshadow), args.Get(1).(int64), args.Error(2)
}

func (m *MockForeshadowService) GetStats(projectID string, currentChapter int) (*model.ForeshadowStats, error) {
	args := m.Called(projectID, currentChapter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ForeshadowStats), args.Error(1)
}

func (m *MockForeshadowService) Resolve(id, chapterID string, chapterNum int, content string) (*model.Foreshadow, error) {
	args := m.Called(id, chapterID, chapterNum, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Foreshadow), args.Error(1)
}

func (m *MockForeshadowService) AddHint(foreshadowID, chapterID string, chapterNum int, content, hintType string) (*model.ForeshadowHint, error) {
	args := m.Called(foreshadowID, chapterID, chapterNum, content, hintType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ForeshadowHint), args.Error(1)
}

func (m *MockForeshadowService) GetHints(foreshadowID string) ([]model.ForeshadowHint, error) {
	args := m.Called(foreshadowID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ForeshadowHint), args.Error(1)
}

func (m *MockForeshadowService) DeleteHint(hintID string) error {
	args := m.Called(hintID)
	return args.Error(0)
}

func (m *MockForeshadowService) GetPendingReminders(projectID string, currentChapter int) ([]model.Foreshadow, error) {
	args := m.Called(projectID, currentChapter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Foreshadow), args.Error(1)
}

func (m *MockForeshadowService) GetOverdueForeshadows(projectID string, currentChapter int) ([]model.Foreshadow, error) {
	args := m.Called(projectID, currentChapter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Foreshadow), args.Error(1)
}

func (m *MockForeshadowService) CreateReminder(foreshadowID string, chapterNum int, message string) (*model.ForeshadowReminder, error) {
	args := m.Called(foreshadowID, chapterNum, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ForeshadowReminder), args.Error(1)
}

func (m *MockForeshadowService) GetUnreadReminders(projectID string) ([]model.ForeshadowReminder, error) {
	args := m.Called(projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ForeshadowReminder), args.Error(1)
}

func (m *MockForeshadowService) MarkReminderAsRead(reminderID string) error {
	args := m.Called(reminderID)
	return args.Error(0)
}

func (m *MockForeshadowService) MarkAllRemindersAsRead(projectID string) error {
	args := m.Called(projectID)
	return args.Error(0)
}

func (m *MockForeshadowService) CheckAndCreateReminders(projectID string, currentChapter int) ([]model.ForeshadowReminder, error) {
	args := m.Called(projectID, currentChapter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ForeshadowReminder), args.Error(1)
}

// TestForeshadowHandler_Create 测试创建伏笔 API
func TestForeshadowHandler_Create(t *testing.T) {
	projectID := uuid.New().String()
	foreshadowID := uuid.New().String()

	tests := []struct {
		name           string
		projectID      string
		requestBody    interface{}
		mockSetup      func(*MockForeshadowService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功创建伏笔",
			projectID: projectID,
			requestBody: map[string]interface{}{
				"title":       "神秘宝藏",
				"description": "主角在第一章发现的地图",
				"priority":    "high",
			},
			mockSetup: func(m *MockForeshadowService) {
				foreshadow := &model.Foreshadow{
					ID:          foreshadowID,
					Title:       "神秘宝藏",
					Description: "主角在第一章发现的地图",
					Priority:    "high",
					Status:      "pending",
				}
				m.On("Create", projectID, mock.AnythingOfType("*service.CreateForeshadowRequest")).Return(foreshadow, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "无效的请求体",
			projectID:      projectID,
			requestBody:    "invalid json",
			mockSetup:      func(m *MockForeshadowService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockForeshadowService)
			tt.mockSetup(mockService)

			handler := &ForeshadowHandler{
				svc: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/projects/:id/foreshadows", handler.Create)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			url := "/api/v1/projects/" + tt.projectID + "/foreshadows"
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
			} else {
				var resp model.Foreshadow
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "神秘宝藏", resp.Title)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestForeshadowHandler_Get 测试获取伏笔 API
func TestForeshadowHandler_Get(t *testing.T) {
	foreshadowID := uuid.New().String()

	tests := []struct {
		name           string
		foreshadowID   string
		mockSetup      func(*MockForeshadowService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:         "成功获取伏笔",
			foreshadowID: foreshadowID,
			mockSetup: func(m *MockForeshadowService) {
				foreshadow := &model.Foreshadow{
					ID:          foreshadowID,
					Title:       "神秘宝藏",
					Description: "主角在第一章发现的地图",
					Status:      "pending",
				}
				m.On("GetByID", foreshadowID).Return(foreshadow, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:         "伏笔不存在",
			foreshadowID: foreshadowID,
			mockSetup: func(m *MockForeshadowService) {
				m.On("GetByID", foreshadowID).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockForeshadowService)
			tt.mockSetup(mockService)

			handler := &ForeshadowHandler{
				svc: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/foreshadows/:id", handler.Get)

			url := "/api/v1/foreshadows/" + tt.foreshadowID
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotNil(t, resp["error"])
			} else {
				var resp model.Foreshadow
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "神秘宝藏", resp.Title)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestForeshadowHandler_List 测试获取伏笔列表 API
func TestForeshadowHandler_List(t *testing.T) {
	projectID := uuid.New().String()

	tests := []struct {
		name           string
		projectID      string
		query          string
		mockSetup      func(*MockForeshadowService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取伏笔列表",
			projectID: projectID,
			query:     "page=1&page_size=20&status=pending",
			mockSetup: func(m *MockForeshadowService) {
				foreshadows := []model.Foreshadow{
					{
						ID:     uuid.New().String(),
						Title:  "神秘宝藏",
						Status: "pending",
					},
				}
				m.On("List", projectID, "pending", "", 1, 20).Return(foreshadows, int64(1), nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:      "默认分页参数",
			projectID: projectID,
			query:     "",
			mockSetup: func(m *MockForeshadowService) {
				foreshadows := []model.Foreshadow{}
				m.On("List", projectID, "", "", 1, 20).Return(foreshadows, int64(0), nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockForeshadowService)
			tt.mockSetup(mockService)

			handler := &ForeshadowHandler{
				svc: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/projects/:id/foreshadows", handler.List)

			url := "/api/v1/projects/" + tt.projectID + "/foreshadows"
			if tt.query != "" {
				url += "?" + tt.query
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.NotNil(t, resp["error"])
			} else {
				assert.NotNil(t, resp["items"])
				assert.NotNil(t, resp["total"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestForeshadowHandler_Delete 测试删除伏笔 API
func TestForeshadowHandler_Delete(t *testing.T) {
	foreshadowID := uuid.New().String()

	tests := []struct {
		name           string
		foreshadowID   string
		mockSetup      func(*MockForeshadowService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:         "成功删除伏笔",
			foreshadowID: foreshadowID,
			mockSetup: func(m *MockForeshadowService) {
				m.On("Delete", foreshadowID).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  false,
		},
		{
			name:         "删除失败",
			foreshadowID: foreshadowID,
			mockSetup: func(m *MockForeshadowService) {
				m.On("Delete", foreshadowID).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockForeshadowService)
			tt.mockSetup(mockService)

			handler := &ForeshadowHandler{
				svc: mockService,
			}
			router := setupTestRouter()

			router.DELETE("/api/v1/foreshadows/:id", handler.Delete)

			url := "/api/v1/foreshadows/" + tt.foreshadowID
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestForeshadowHandler_GetStats 测试获取伏笔统计 API
func TestForeshadowHandler_GetStats(t *testing.T) {
	projectID := uuid.New().String()

	tests := []struct {
		name           string
		projectID      string
		query          string
		mockSetup      func(*MockForeshadowService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "成功获取统计",
			projectID: projectID,
			query:     "current_chapter=5",
			mockSetup: func(m *MockForeshadowService) {
				stats := &model.ForeshadowStats{
					Total:        10,
					Planted:      5,
					Resolved:     3,
					OverdueCount: 2,
				}
				m.On("GetStats", projectID, 5).Return(stats, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:      "默认章节参数",
			projectID: projectID,
			query:     "",
			mockSetup: func(m *MockForeshadowService) {
				stats := &model.ForeshadowStats{
					Total:        10,
					Planted:      5,
					Resolved:     5,
					OverdueCount: 0,
				}
				m.On("GetStats", projectID, 0).Return(stats, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockForeshadowService)
			tt.mockSetup(mockService)

			handler := &ForeshadowHandler{
				svc: mockService,
			}
			router := setupTestRouter()

			router.GET("/api/v1/projects/:id/foreshadows/stats", handler.GetStats)

			url := "/api/v1/projects/" + tt.projectID + "/foreshadows/stats"
			if tt.query != "" {
				url += "?" + tt.query
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectedError {
				var resp model.ForeshadowStats
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, 10, resp.Total)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestForeshadowHandler_Resolve 测试回收伏笔 API
func TestForeshadowHandler_Resolve(t *testing.T) {
	foreshadowID := uuid.New().String()
	chapterID := uuid.New().String()

	tests := []struct {
		name           string
		foreshadowID   string
		requestBody    interface{}
		mockSetup      func(*MockForeshadowService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:         "成功回收伏笔",
			foreshadowID: foreshadowID,
			requestBody: map[string]interface{}{
				"chapter_id":  chapterID,
				"chapter_num": 10,
				"content":     "宝藏被找到了",
			},
			mockSetup: func(m *MockForeshadowService) {
				foreshadow := &model.Foreshadow{
					ID:     foreshadowID,
					Title:  "神秘宝藏",
					Status: "resolved",
				}
				m.On("Resolve", foreshadowID, chapterID, 10, "宝藏被找到了").Return(foreshadow, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "无效的请求体",
			foreshadowID:   foreshadowID,
			requestBody:    "invalid json",
			mockSetup:      func(m *MockForeshadowService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockForeshadowService)
			tt.mockSetup(mockService)

			handler := &ForeshadowHandler{
				svc: mockService,
			}
			router := setupTestRouter()

			router.POST("/api/v1/foreshadows/:id/resolve", handler.Resolve)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			url := "/api/v1/foreshadows/" + tt.foreshadowID + "/resolve"
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
			} else {
				var resp model.Foreshadow
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, model.ForeshadowStatus("resolved"), resp.Status)
			}

			mockService.AssertExpectations(t)
		})
	}
}
