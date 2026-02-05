package service

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockChapterRepository 是 ChapterRepository 的 mock
type MockChapterRepository struct {
	mock.Mock
}

func (m *MockChapterRepository) Create(ctx context.Context, chapter *model.Chapter) error {
	args := m.Called(ctx, chapter)
	return args.Error(0)
}

func (m *MockChapterRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Chapter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chapter), args.Error(1)
}

func (m *MockChapterRepository) Update(ctx context.Context, chapter *model.Chapter) error {
	args := m.Called(ctx, chapter)
	return args.Error(0)
}

func (m *MockChapterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChapterRepository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	// 直接执行传入的函数
	if err := fn(ctx); err != nil {
		return err
	}
	// 调用 mock 以验证调用
	args := m.Called(ctx, mock.Anything)
	return args.Error(0)
}

func (m *MockChapterRepository) UpdateChapterNumbers(ctx context.Context, updates []repository.ChapterNumberUpdate) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

// MockProjectRepository 是 ProjectRepository 的 mock
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockProjectRepository) UpdateStats(ctx context.Context, projectID uuid.UUID) error {
	args := m.Called(ctx, projectID)
	return args.Error(0)
}

// TestBatchOperationService_BatchUpdate 测试批量更新章节
func TestBatchOperationService_BatchUpdate(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	projectID := uuid.New()
	chapterID1 := uuid.New()
	chapterID2 := uuid.New()

	tests := []struct {
		name          string
		req           *BatchUpdateRequest
		mockSetup     func(*MockChapterRepository, *MockProjectRepository)
		wantSuccess   int
		wantFailed    int
		wantErr       error
	}{
		{
			name: "成功批量更新章节",
			req: &BatchUpdateRequest{
				ChapterIDs: []string{chapterID1.String(), chapterID2.String()},
				Updates: map[string]interface{}{
					"title":  "新标题",
					"status": "completed",
				},
			},
			mockSetup: func(chapterRepo *MockChapterRepository, projectRepo *MockProjectRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)

				chapter1 := &model.Chapter{
					ID:        chapterID1,
					ProjectID: projectID,
					Title:     "章节1",
					Status:    "draft",
				}
				chapter2 := &model.Chapter{
					ID:        chapterID2,
					ProjectID: projectID,
					Title:     "章节2",
					Status:    "draft",
				}

				chapterRepo.On("FindByID", ctx, chapterID1).Return(chapter1, nil)
				chapterRepo.On("FindByID", ctx, chapterID2).Return(chapter2, nil)
				chapterRepo.On("Update", ctx, mock.AnythingOfType("*model.Chapter")).Return(nil)
				chapterRepo.On("Transaction", ctx, mock.Anything).Return(nil)
			},
			wantSuccess: 2,
			wantFailed:  0,
			wantErr:     nil,
		},
		{
			name: "批量大小超过限制",
			req: &BatchUpdateRequest{
				ChapterIDs: make([]string, MaxBatchSize+1),
				Updates:    map[string]interface{}{"title": "新标题"},
			},
			mockSetup: func(chapterRepo *MockChapterRepository, projectRepo *MockProjectRepository) {
				// 不需要 mock，因为会在验证阶段失败
			},
			wantErr: ErrBatchLimitExceeded,
		},
		{
			name: "项目不属于用户",
			req: &BatchUpdateRequest{
				ChapterIDs: []string{chapterID1.String()},
				Updates:    map[string]interface{}{"title": "新标题"},
			},
			mockSetup: func(chapterRepo *MockChapterRepository, projectRepo *MockProjectRepository) {
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
			name: "部分章节更新失败",
			req: &BatchUpdateRequest{
				ChapterIDs: []string{chapterID1.String(), chapterID2.String()},
				Updates:    map[string]interface{}{"title": "新标题"},
			},
			mockSetup: func(chapterRepo *MockChapterRepository, projectRepo *MockProjectRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)

				chapter1 := &model.Chapter{
					ID:        chapterID1,
					ProjectID: projectID,
					Title:     "章节1",
				}

				chapterRepo.On("FindByID", ctx, chapterID1).Return(chapter1, nil)
				chapterRepo.On("FindByID", ctx, chapterID2).Return(nil, errors.New("not found"))
				chapterRepo.On("Update", ctx, mock.AnythingOfType("*model.Chapter")).Return(nil)
				chapterRepo.On("Transaction", ctx, mock.Anything).Return(nil)
			},
			wantSuccess: 1,
			wantFailed:  1,
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChapterRepo := new(MockChapterRepository)
			mockProjectRepo := new(MockProjectRepository)
			tt.mockSetup(mockChapterRepo, mockProjectRepo)

			service := &BatchOperationService{
				chapterRepo: mockChapterRepo,
				projectRepo: mockProjectRepo,
			}

			result, err := service.BatchUpdate(ctx, userID, projectID, tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.wantSuccess, result.Success)
				assert.Equal(t, tt.wantFailed, result.Failed)
			}

			mockChapterRepo.AssertExpectations(t)
			mockProjectRepo.AssertExpectations(t)
		})
	}
}

// TestBatchOperationService_BatchDelete 测试批量删除章节
func TestBatchOperationService_BatchDelete(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	projectID := uuid.New()
	chapterID1 := uuid.New()
	chapterID2 := uuid.New()

	tests := []struct {
		name        string
		req         *BatchDeleteRequest
		mockSetup   func(*MockChapterRepository, *MockProjectRepository)
		wantSuccess int
		wantFailed  int
		wantErr     error
	}{
		{
			name: "成功批量删除章节",
			req: &BatchDeleteRequest{
				ChapterIDs: []string{chapterID1.String(), chapterID2.String()},
			},
			mockSetup: func(chapterRepo *MockChapterRepository, projectRepo *MockProjectRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)
				projectRepo.On("UpdateStats", ctx, projectID).Return(nil)

				chapter1 := &model.Chapter{
					ID:        chapterID1,
					ProjectID: projectID,
					Title:     "章节1",
				}
				chapter2 := &model.Chapter{
					ID:        chapterID2,
					ProjectID: projectID,
					Title:     "章节2",
				}

				chapterRepo.On("FindByID", ctx, chapterID1).Return(chapter1, nil)
				chapterRepo.On("FindByID", ctx, chapterID2).Return(chapter2, nil)
				chapterRepo.On("Delete", ctx, chapterID1).Return(nil)
				chapterRepo.On("Delete", ctx, chapterID2).Return(nil)
				chapterRepo.On("Transaction", ctx, mock.Anything).Return(nil)
			},
			wantSuccess: 2,
			wantFailed:  0,
			wantErr:     nil,
		},
		{
			name: "批量大小超过限制",
			req: &BatchDeleteRequest{
				ChapterIDs: make([]string, MaxBatchSize+1),
			},
			mockSetup: func(chapterRepo *MockChapterRepository, projectRepo *MockProjectRepository) {
				// 不需要 mock
			},
			wantErr: ErrBatchLimitExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChapterRepo := new(MockChapterRepository)
			mockProjectRepo := new(MockProjectRepository)
			tt.mockSetup(mockChapterRepo, mockProjectRepo)

			service := &BatchOperationService{
				chapterRepo: mockChapterRepo,
				projectRepo: mockProjectRepo,
			}

			result, err := service.BatchDelete(ctx, userID, projectID, tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.wantSuccess, result.Success)
				assert.Equal(t, tt.wantFailed, result.Failed)
			}

			mockChapterRepo.AssertExpectations(t)
			mockProjectRepo.AssertExpectations(t)
		})
	}
}

// TestBatchOperationService_BatchStatusUpdate 测试批量更新状态
func TestBatchOperationService_BatchStatusUpdate(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	projectID := uuid.New()
	chapterID1 := uuid.New()

	tests := []struct {
		name        string
		req         *BatchStatusUpdateRequest
		mockSetup   func(*MockChapterRepository, *MockProjectRepository)
		wantSuccess int
		wantFailed  int
		wantErr     error
	}{
		{
			name: "成功批量更新状态",
			req: &BatchStatusUpdateRequest{
				ChapterIDs: []string{chapterID1.String()},
				Status:     "completed",
			},
			mockSetup: func(chapterRepo *MockChapterRepository, projectRepo *MockProjectRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)

				chapter := &model.Chapter{
					ID:        chapterID1,
					ProjectID: projectID,
					Title:     "章节1",
					Status:    "draft",
				}

				chapterRepo.On("FindByID", ctx, chapterID1).Return(chapter, nil)
				chapterRepo.On("Update", ctx, mock.AnythingOfType("*model.Chapter")).Return(nil)
				chapterRepo.On("Transaction", ctx, mock.Anything).Return(nil)
			},
			wantSuccess: 1,
			wantFailed:  0,
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChapterRepo := new(MockChapterRepository)
			mockProjectRepo := new(MockProjectRepository)
			tt.mockSetup(mockChapterRepo, mockProjectRepo)

			service := &BatchOperationService{
				chapterRepo: mockChapterRepo,
				projectRepo: mockProjectRepo,
			}

			result, err := service.BatchStatusUpdate(ctx, userID, projectID, tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.wantSuccess, result.Success)
				assert.Equal(t, tt.wantFailed, result.Failed)
			}

			mockChapterRepo.AssertExpectations(t)
			mockProjectRepo.AssertExpectations(t)
		})
	}
}

// TestBatchOperationService_Reorder 测试重排序章节
func TestBatchOperationService_Reorder(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	projectID := uuid.New()
	chapterID1 := uuid.New()
	chapterID2 := uuid.New()

	tests := []struct {
		name      string
		req       *ReorderRequest
		mockSetup func(*MockChapterRepository, *MockProjectRepository)
		wantErr   error
	}{
		{
			name: "成功重排序章节",
			req: &ReorderRequest{
				ChapterOrders: []ChapterOrder{
					{ChapterID: chapterID1.String(), ChapterNumber: 2},
					{ChapterID: chapterID2.String(), ChapterNumber: 1},
				},
			},
			mockSetup: func(chapterRepo *MockChapterRepository, projectRepo *MockProjectRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)

				chapter1 := &model.Chapter{
					ID:            chapterID1,
					ProjectID:     projectID,
					Title:         "章节1",
					ChapterNumber: 1,
				}
				chapter2 := &model.Chapter{
					ID:            chapterID2,
					ProjectID:     projectID,
					Title:         "章节2",
					ChapterNumber: 2,
				}

				chapterRepo.On("FindByID", ctx, chapterID1).Return(chapter1, nil)
				chapterRepo.On("FindByID", ctx, chapterID2).Return(chapter2, nil)
				chapterRepo.On("UpdateChapterNumbers", ctx, mock.AnythingOfType("[]repository.ChapterNumberUpdate")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "批量大小超过限制",
			req: &ReorderRequest{
				ChapterOrders: make([]ChapterOrder, MaxBatchSize+1),
			},
			mockSetup: func(chapterRepo *MockChapterRepository, projectRepo *MockProjectRepository) {
				// 不需要 mock
			},
			wantErr: ErrBatchLimitExceeded,
		},
		{
			name: "章节不属于项目",
			req: &ReorderRequest{
				ChapterOrders: []ChapterOrder{
					{ChapterID: chapterID1.String(), ChapterNumber: 1},
				},
			},
			mockSetup: func(chapterRepo *MockChapterRepository, projectRepo *MockProjectRepository) {
				project := &model.Project{
					ID:     projectID,
					UserID: userID,
					Title:  "测试项目",
				}
				projectRepo.On("FindByID", ctx, projectID).Return(project, nil)

				chapter := &model.Chapter{
					ID:            chapterID1,
					ProjectID:     uuid.New(), // 不同的项目
					Title:         "章节1",
					ChapterNumber: 1,
				}

				chapterRepo.On("FindByID", ctx, chapterID1).Return(chapter, nil)
			},
			wantErr: ErrChapterNotInProject,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChapterRepo := new(MockChapterRepository)
			mockProjectRepo := new(MockProjectRepository)
			tt.mockSetup(mockChapterRepo, mockProjectRepo)

			service := &BatchOperationService{
				chapterRepo: mockChapterRepo,
				projectRepo: mockProjectRepo,
			}

			err := service.Reorder(ctx, userID, projectID, tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}

			mockChapterRepo.AssertExpectations(t)
			mockProjectRepo.AssertExpectations(t)
		})
	}
}
