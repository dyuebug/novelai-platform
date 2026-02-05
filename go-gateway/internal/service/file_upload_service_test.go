package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFileRepository 是 FileRepository 的 mock
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Create(ctx context.Context, file *model.File) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockFileRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.File, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.File), args.Error(1)
}

func (m *MockFileRepository) FindByHash(ctx context.Context, hash string) (*model.File, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.File), args.Error(1)
}

func (m *MockFileRepository) FindByUserID(ctx context.Context, userID uuid.UUID, opts *repository.ListOptions) ([]model.File, int64, error) {
	args := m.Called(ctx, userID, opts)
	return args.Get(0).([]model.File), args.Get(1).(int64), args.Error(2)
}

func (m *MockFileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockFileRepository) GetTotalSizeByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// TestFileUploadService_GetFile 测试获取文件
func TestFileUploadService_GetFile(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	fileID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name      string
		userID    uuid.UUID
		fileID    uuid.UUID
		mockSetup func(*MockFileRepository)
		wantErr   error
	}{
		{
			name:   "成功获取文件",
			userID: userID,
			fileID: fileID,
			mockSetup: func(m *MockFileRepository) {
				file := &model.File{
					ID:          fileID,
					UserID:      userID,
					FileName:    "test.jpg",
					StoragePath: "/uploads/test.jpg",
					FileSize:    1024,
					MimeType:    "image/jpeg",
					FileType:    "image",
					Hash:        "abc123",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("FindByID", ctx, fileID).Return(file, nil)
			},
			wantErr: nil,
		},
		{
			name:   "文件不存在",
			userID: userID,
			fileID: fileID,
			mockSetup: func(m *MockFileRepository) {
				m.On("FindByID", ctx, fileID).Return(nil, errors.New("not found"))
			},
			wantErr: errors.New("not found"),
		},
		{
			name:   "文件不属于用户",
			userID: userID,
			fileID: fileID,
			mockSetup: func(m *MockFileRepository) {
				file := &model.File{
					ID:          fileID,
					UserID:      otherUserID, // 不同的用户
					FileName:    "test.jpg",
					StoragePath: "/uploads/test.jpg",
					FileSize:    1024,
					MimeType:    "image/jpeg",
					FileType:    "image",
					Hash:        "abc123",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("FindByID", ctx, fileID).Return(file, nil)
			},
			wantErr: ErrFileNotOwned,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockFileRepository)
			tt.mockSetup(mockRepo)

			service := &FileUploadService{
				fileRepo: mockRepo,
			}

			file, err := service.GetFile(ctx, tt.userID, tt.fileID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, file)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, tt.fileID, file.ID)
				assert.Equal(t, tt.userID, file.UserID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestFileUploadService_ListFiles 测试获取文件列表
func TestFileUploadService_ListFiles(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name         string
		page         int
		pageSize     int
		fileType     string
		sort         string
		mockSetup    func(*MockFileRepository)
		wantTotal    int64
		wantPageSize int
		wantErr      error
	}{
		{
			name:     "成功获取文件列表",
			page:     1,
			pageSize: 20,
			fileType: "image",
			sort:     "created_at DESC",
			mockSetup: func(m *MockFileRepository) {
				files := []model.File{
					{
						ID:          uuid.New(),
						UserID:      userID,
						FileName:    "test1.jpg",
						StoragePath: "/uploads/test1.jpg",
						FileSize:    1024,
						MimeType:    "image/jpeg",
						FileType:    "image",
						Hash:        "abc123",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          uuid.New(),
						UserID:      userID,
						FileName:    "test2.jpg",
						StoragePath: "/uploads/test2.jpg",
						FileSize:    2048,
						MimeType:    "image/jpeg",
						FileType:    "image",
						Hash:        "def456",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				opts := &repository.ListOptions{
					Page:     1,
					PageSize: 20,
					Status:   "image",
					Sort:     "created_at DESC",
				}
				m.On("FindByUserID", ctx, userID, opts).Return(files, int64(2), nil)
			},
			wantTotal:    2,
			wantPageSize: 20,
			wantErr:      nil,
		},
		{
			name:     "页码为0时默认为1",
			page:     0,
			pageSize: 20,
			fileType: "",
			sort:     "",
			mockSetup: func(m *MockFileRepository) {
				files := []model.File{}
				opts := &repository.ListOptions{
					Page:     1, // 应该被修正为1
					PageSize: 20,
					Status:   "",
					Sort:     "",
				}
				m.On("FindByUserID", ctx, userID, opts).Return(files, int64(0), nil)
			},
			wantTotal:    0,
			wantPageSize: 20,
			wantErr:      nil,
		},
		{
			name:     "页大小为0时默认为20",
			page:     1,
			pageSize: 0,
			fileType: "",
			sort:     "",
			mockSetup: func(m *MockFileRepository) {
				files := []model.File{}
				opts := &repository.ListOptions{
					Page:     1,
					PageSize: 20, // 应该被修正为20
					Status:   "",
					Sort:     "",
				}
				m.On("FindByUserID", ctx, userID, opts).Return(files, int64(0), nil)
			},
			wantTotal:    0,
			wantPageSize: 20,
			wantErr:      nil,
		},
		{
			name:     "页大小超过100时限制为100",
			page:     1,
			pageSize: 200,
			fileType: "",
			sort:     "",
			mockSetup: func(m *MockFileRepository) {
				files := []model.File{}
				opts := &repository.ListOptions{
					Page:     1,
					PageSize: 100, // 应该被限制为100
					Status:   "",
					Sort:     "",
				}
				m.On("FindByUserID", ctx, userID, opts).Return(files, int64(0), nil)
			},
			wantTotal:    0,
			wantPageSize: 100,
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockFileRepository)
			tt.mockSetup(mockRepo)

			service := &FileUploadService{
				fileRepo: mockRepo,
			}

			resp, err := service.ListFiles(ctx, userID, tt.page, tt.pageSize, tt.fileType, tt.sort)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantTotal, resp.Total)
				assert.Equal(t, tt.wantPageSize, resp.PageSize)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestFileUploadService_GetUserStorageStats 测试获取用户存储统计
func TestFileUploadService_GetUserStorageStats(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(*MockFileRepository)
		wantFileCount int64
		wantTotalSize int64
		wantErr       error
	}{
		{
			name: "成功获取存储统计",
			mockSetup: func(m *MockFileRepository) {
				m.On("CountByUserID", ctx, userID).Return(int64(10), nil)
				m.On("GetTotalSizeByUserID", ctx, userID).Return(int64(1024000), nil)
			},
			wantFileCount: 10,
			wantTotalSize: 1024000,
			wantErr:       nil,
		},
		{
			name: "获取文件数量失败",
			mockSetup: func(m *MockFileRepository) {
				m.On("CountByUserID", ctx, userID).Return(int64(0), errors.New("database error"))
			},
			wantErr: errors.New("database error"),
		},
		{
			name: "获取总大小失败",
			mockSetup: func(m *MockFileRepository) {
				m.On("CountByUserID", ctx, userID).Return(int64(10), nil)
				m.On("GetTotalSizeByUserID", ctx, userID).Return(int64(0), errors.New("database error"))
			},
			wantErr: errors.New("database error"),
		},
		{
			name: "用户没有文件",
			mockSetup: func(m *MockFileRepository) {
				m.On("CountByUserID", ctx, userID).Return(int64(0), nil)
				m.On("GetTotalSizeByUserID", ctx, userID).Return(int64(0), nil)
			},
			wantFileCount: 0,
			wantTotalSize: 0,
			wantErr:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockFileRepository)
			tt.mockSetup(mockRepo)

			service := &FileUploadService{
				fileRepo: mockRepo,
			}

			stats, err := service.GetUserStorageStats(ctx, userID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, stats)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, stats)
				assert.Equal(t, tt.wantFileCount, stats.FileCount)
				assert.Equal(t, tt.wantTotalSize, stats.TotalSize)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
