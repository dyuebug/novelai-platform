package service

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSearchProjectRepository 是 ProjectRepository 的 mock（用于搜索）
type MockSearchProjectRepository struct {
	mock.Mock
}

func (m *MockSearchProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockSearchProjectRepository) Search(ctx context.Context, userID uuid.UUID, query string, page, pageSize int) ([]model.Project, int64, error) {
	args := m.Called(ctx, userID, query, page, pageSize)
	return args.Get(0).([]model.Project), args.Get(1).(int64), args.Error(2)
}

// MockSearchChapterRepository 是 ChapterRepository 的 mock（用于搜索）
type MockSearchChapterRepository struct {
	mock.Mock
}

func (m *MockSearchChapterRepository) Search(ctx context.Context, userID, projectID uuid.UUID, query string, page, pageSize int) ([]model.Chapter, int64, error) {
	args := m.Called(ctx, userID, projectID, query, page, pageSize)
	return args.Get(0).([]model.Chapter), args.Get(1).(int64), args.Error(2)
}

func (m *MockSearchChapterRepository) AdvancedFilter(ctx context.Context, projectID uuid.UUID, req interface{}) ([]model.Chapter, int64, error) {
	args := m.Called(ctx, projectID, req)
	return args.Get(0).([]model.Chapter), args.Get(1).(int64), args.Error(2)
}

// MockCharacterRepository 是 CharacterRepository 的 mock
type MockCharacterRepository struct {
	mock.Mock
}

func (m *MockCharacterRepository) Search(ctx context.Context, projectID uuid.UUID, query string, page, pageSize int) ([]model.Character, int64, error) {
	args := m.Called(ctx, projectID, query, page, pageSize)
	return args.Get(0).([]model.Character), args.Get(1).(int64), args.Error(2)
}

// MockLocationRepository 是 LocationRepository 的 mock
type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) Search(ctx context.Context, projectID uuid.UUID, query string, page, pageSize int) ([]model.Location, int64, error) {
	args := m.Called(ctx, projectID, query, page, pageSize)
	return args.Get(0).([]model.Location), args.Get(1).(int64), args.Error(2)
}

// TestSearchService_GlobalSearch 测试全局搜索
func TestSearchService_GlobalSearch(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name      string
		req       *GlobalSearchRequest
		mockSetup func(*MockSearchProjectRepository, *MockSearchChapterRepository)
		wantTotal int64
		wantErr   error
	}{
		{
			name: "成功全局搜索",
			req: &GlobalSearchRequest{
				Query:    "测试",
				Page:     1,
				PageSize: 20,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository) {
				projects := []model.Project{
					{
						ID:          uuid.New(),
						UserID:      userID,
						Title:       "测试项目",
						Description: "这是一个测试项目",
						Genre:       "fantasy",
						Status:      "in_progress",
					},
				}
				projectRepo.On("Search", ctx, userID, "测试", 1, 20).Return(projects, int64(1), nil)

				chapters := []model.Chapter{
					{
						ID:            uuid.New(),
						ProjectID:     uuid.New(),
						Title:         "测试章节",
						Summary:       "这是一个测试章节",
						ChapterNumber: 1,
						Status:        "draft",
						WordCount:     1000,
					},
				}
				chapterRepo.On("Search", ctx, userID, uuid.Nil, "测试", 1, 20).Return(chapters, int64(1), nil)
			},
			wantTotal: 2,
			wantErr:   nil,
		},
		{
			name: "页码为0时默认为1",
			req: &GlobalSearchRequest{
				Query:    "测试",
				Page:     0,
				PageSize: 20,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository) {
				projectRepo.On("Search", ctx, userID, "测试", 1, 20).Return([]model.Project{}, int64(0), nil)
				chapterRepo.On("Search", ctx, userID, uuid.Nil, "测试", 1, 20).Return([]model.Chapter{}, int64(0), nil)
			},
			wantTotal: 0,
			wantErr:   nil,
		},
		{
			name: "页大小为0时默认为20",
			req: &GlobalSearchRequest{
				Query:    "测试",
				Page:     1,
				PageSize: 0,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository) {
				projectRepo.On("Search", ctx, userID, "测试", 1, 20).Return([]model.Project{}, int64(0), nil)
				chapterRepo.On("Search", ctx, userID, uuid.Nil, "测试", 1, 20).Return([]model.Chapter{}, int64(0), nil)
			},
			wantTotal: 0,
			wantErr:   nil,
		},
		{
			name: "页大小超过100时限制为100",
			req: &GlobalSearchRequest{
				Query:    "测试",
				Page:     1,
				PageSize: 200,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository) {
				projectRepo.On("Search", ctx, userID, "测试", 1, 100).Return([]model.Project{}, int64(0), nil)
				chapterRepo.On("Search", ctx, userID, uuid.Nil, "测试", 1, 100).Return([]model.Chapter{}, int64(0), nil)
			},
			wantTotal: 0,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProjectRepo := new(MockSearchProjectRepository)
			mockChapterRepo := new(MockSearchChapterRepository)
			tt.mockSetup(mockProjectRepo, mockChapterRepo)

			service := &SearchService{
				projectRepo: mockProjectRepo,
				chapterRepo: mockChapterRepo,
			}

			resp, err := service.GlobalSearch(ctx, userID, tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantTotal, resp.Total)
			}

			mockProjectRepo.AssertExpectations(t)
			mockChapterRepo.AssertExpectations(t)
		})
	}
}

// TestSearchService_ProjectSearch 测试项目内搜索
func TestSearchService_ProjectSearch(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	projectID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name      string
		req       *ProjectSearchRequest
		mockSetup func(*MockSearchProjectRepository, *MockSearchChapterRepository, *MockCharacterRepository, *MockLocationRepository)
		wantTotal int64
		wantErr   error
	}{
		{
			name: "成功搜索章节",
			req: &ProjectSearchRequest{
				Query:    "测试",
				Type:     "chapter",
				Page:     1,
				PageSize: 20,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository, characterRepo *MockCharacterRepository, locationRepo *MockLocationRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)

				chapters := []model.Chapter{
					{
						ID:            uuid.New(),
						ProjectID:     projectID,
						Title:         "测试章节",
						Summary:       "这是一个测试章节",
						ChapterNumber: 1,
						Status:        "draft",
						WordCount:     1000,
					},
				}
				chapterRepo.On("Search", ctx, userID, projectID, "测试", 1, 20).Return(chapters, int64(1), nil)
			},
			wantTotal: 1,
			wantErr:   nil,
		},
		{
			name: "项目不属于用户",
			req: &ProjectSearchRequest{
				Query:    "测试",
				Type:     "chapter",
				Page:     1,
				PageSize: 20,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository, characterRepo *MockCharacterRepository, locationRepo *MockLocationRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: otherUserID, // 不同的用户
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)
			},
			wantErr: ErrProjectNotOwned,
		},
		{
			name: "项目不存在",
			req: &ProjectSearchRequest{
				Query:    "测试",
				Type:     "chapter",
				Page:     1,
				PageSize: 20,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository, characterRepo *MockCharacterRepository, locationRepo *MockLocationRepository) {
				projectRepo.On("FindByID", ctx, projectID).Return(nil, errors.New("not found"))
			},
			wantErr: errors.New("not found"),
		},
		{
			name: "成功搜索角色",
			req: &ProjectSearchRequest{
				Query:    "主角",
				Type:     "character",
				Page:     1,
				PageSize: 20,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository, characterRepo *MockCharacterRepository, locationRepo *MockLocationRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)

				characters := []model.Character{
					{
						ID:         uuid.New(),
						ProjectID:  projectID,
						Name:       "主角",
						Background: "这是主角的背景",
						Role:       "protagonist",
						Status:     "active",
					},
				}
				characterRepo.On("Search", ctx, projectID, "主角", 1, 20).Return(characters, int64(1), nil)
			},
			wantTotal: 1,
			wantErr:   nil,
		},
		{
			name: "成功搜索地点",
			req: &ProjectSearchRequest{
				Query:    "城市",
				Type:     "location",
				Page:     1,
				PageSize: 20,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository, characterRepo *MockCharacterRepository, locationRepo *MockLocationRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)

				locations := []model.Location{
					{
						ID:          uuid.New(),
						ProjectID:   projectID,
						Name:        "测试城市",
						Description: "这是一个测试城市",
						Type:        "city",
					},
				}
				locationRepo.On("Search", ctx, projectID, "城市", 1, 20).Return(locations, int64(1), nil)
			},
			wantTotal: 1,
			wantErr:   nil,
		},
		{
			name: "搜索所有类型",
			req: &ProjectSearchRequest{
				Query:    "测试",
				Type:     "all",
				Page:     1,
				PageSize: 20,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository, characterRepo *MockCharacterRepository, locationRepo *MockLocationRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)

				chapterRepo.On("Search", ctx, userID, projectID, "测试", 1, 10).Return([]model.Chapter{}, int64(0), nil)
				characterRepo.On("Search", ctx, projectID, "测试", 1, 10).Return([]model.Character{}, int64(0), nil)
				locationRepo.On("Search", ctx, projectID, "测试", 1, 10).Return([]model.Location{}, int64(0), nil)
			},
			wantTotal: 0,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProjectRepo := new(MockSearchProjectRepository)
			mockChapterRepo := new(MockSearchChapterRepository)
			mockCharacterRepo := new(MockCharacterRepository)
			mockLocationRepo := new(MockLocationRepository)
			tt.mockSetup(mockProjectRepo, mockChapterRepo, mockCharacterRepo, mockLocationRepo)

			service := &SearchService{
				projectRepo:   mockProjectRepo,
				chapterRepo:   mockChapterRepo,
				characterRepo: mockCharacterRepo,
				locationRepo:  mockLocationRepo,
			}

			resp, err := service.ProjectSearch(ctx, userID, projectID, tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantTotal, resp.Total)
			}

			mockProjectRepo.AssertExpectations(t)
			mockChapterRepo.AssertExpectations(t)
			mockCharacterRepo.AssertExpectations(t)
			mockLocationRepo.AssertExpectations(t)
		})
	}
}

// TestSearchService_AdvancedFilter 测试高级筛选
func TestSearchService_AdvancedFilter(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name      string
		req       *AdvancedFilterRequest
		mockSetup func(*MockSearchProjectRepository, *MockSearchChapterRepository)
		wantTotal int64
		wantErr   error
	}{
		{
			name: "成功高级筛选",
			req: &AdvancedFilterRequest{
				Tags:         []string{"tag1", "tag2"},
				Status:       []string{"draft"},
				WordCountMin: 1000,
				WordCountMax: 5000,
				Page:         1,
				PageSize:     20,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)

				chapters := []model.Chapter{
					{
						ID:            uuid.New(),
						ProjectID:     projectID,
						Title:         "测试章节",
						Summary:       "这是一个测试章节",
						ChapterNumber: 1,
						Status:        "draft",
						WordCount:     2000,
					},
				}
				chapterRepo.On("AdvancedFilter", ctx, projectID, mock.AnythingOfType("*service.AdvancedFilterRequest")).Return(chapters, int64(1), nil)
			},
			wantTotal: 1,
			wantErr:   nil,
		},
		{
			name: "项目不属于用户",
			req: &AdvancedFilterRequest{
				Page:     1,
				PageSize: 20,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: uuid.New(), // 不同的用户
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)
			},
			wantErr: ErrProjectNotOwned,
		},
		{
			name: "页码和页大小默认值",
			req: &AdvancedFilterRequest{
				Page:     0,
				PageSize: 0,
			},
			mockSetup: func(projectRepo *MockSearchProjectRepository, chapterRepo *MockSearchChapterRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)
				chapterRepo.On("AdvancedFilter", ctx, projectID, mock.AnythingOfType("*service.AdvancedFilterRequest")).Return([]model.Chapter{}, int64(0), nil)
			},
			wantTotal: 0,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProjectRepo := new(MockSearchProjectRepository)
			mockChapterRepo := new(MockSearchChapterRepository)
			tt.mockSetup(mockProjectRepo, mockChapterRepo)

			service := &SearchService{
				projectRepo: mockProjectRepo,
				chapterRepo: mockChapterRepo,
			}

			resp, err := service.AdvancedFilter(ctx, userID, projectID, tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantTotal, resp.Total)
			}

			mockProjectRepo.AssertExpectations(t)
			mockChapterRepo.AssertExpectations(t)
		})
	}
}

// TestSearchService_HelperFunctions 测试辅助函数
func TestSearchService_HelperFunctions(t *testing.T) {
	t.Run("truncate函数", func(t *testing.T) {
		// 测试短文本
		result := truncate("短文本", 100)
		assert.Equal(t, "短文本", result)

		// 测试长文本
		longText := "这是一个很长的文本，需要被截断。这是一个很长的文本，需要被截断。"
		result = truncate(longText, 20)
		assert.Equal(t, 20+3, len(result)) // 20个字符 + "..."
		assert.True(t, len(result) <= 23)
	})

	t.Run("highlightText函数", func(t *testing.T) {
		// 测试空查询
		result := highlightText("测试文本", "")
		assert.Equal(t, "测试文本", result)

		// 测试匹配
		result = highlightText("这是测试文本", "测试")
		assert.Contains(t, result, "<mark>")
		assert.Contains(t, result, "</mark>")

		// 测试不匹配
		result = highlightText("这是文本", "测试")
		assert.Equal(t, "这是文本", result)
	})
}
