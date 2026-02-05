package service

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrBatchOperationFailed    = errors.New("batch operation failed")
	ErrBatchLimitExceeded      = errors.New("batch operation limit exceeded")
	ErrInvalidBatchData        = errors.New("invalid batch data")
	ErrChapterNotInProject     = errors.New("chapter not in project")
)

const (
	MaxBatchSize = 100
)

// ChapterRepository 章节仓库接口
type ChapterRepository interface {
	Create(ctx context.Context, chapter *model.Chapter) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Chapter, error)
	Update(ctx context.Context, chapter *model.Chapter) error
	Delete(ctx context.Context, id uuid.UUID) error
	Transaction(ctx context.Context, fn func(context.Context) error) error
	UpdateChapterNumbers(ctx context.Context, updates []repository.ChapterNumberUpdate) error
}

// ProjectRepository 项目仓库接口
type ProjectRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.Project, error)
	UpdateStats(ctx context.Context, projectID uuid.UUID) error
}

// BatchOperationService 批量操作服务
type BatchOperationService struct {
	chapterRepo ChapterRepository
	projectRepo ProjectRepository
}

// NewBatchOperationService 创建批量操作服务
func NewBatchOperationService() *BatchOperationService {
	return &BatchOperationService{
		chapterRepo: repository.NewChapterRepository(),
		projectRepo: repository.NewProjectRepository(),
	}
}

// BatchUpdateRequest 批量更新请求
type BatchUpdateRequest struct {
	ChapterIDs []string               `json:"chapter_ids" binding:"required,min=1"`
	Updates    map[string]interface{} `json:"updates" binding:"required"`
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	ChapterIDs []string `json:"chapter_ids" binding:"required,min=1"`
}

// BatchStatusUpdateRequest 批量状态更新请求
type BatchStatusUpdateRequest struct {
	ChapterIDs []string `json:"chapter_ids" binding:"required,min=1"`
	Status     string   `json:"status" binding:"required,oneof=draft writing completed published"`
}

// ReorderRequest 重排序请求
type ReorderRequest struct {
	ChapterOrders []ChapterOrder `json:"chapter_orders" binding:"required,min=1"`
}

// ChapterOrder 章节顺序
type ChapterOrder struct {
	ChapterID     string `json:"chapter_id" binding:"required"`
	ChapterNumber int    `json:"chapter_number" binding:"required,min=1"`
}

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	Success      int      `json:"success"`
	Failed       int      `json:"failed"`
	Total        int      `json:"total"`
	FailedIDs    []string `json:"failed_ids,omitempty"`
	ErrorMessage string   `json:"error_message,omitempty"`
}

// BatchUpdate 批量更新章节
func (s *BatchOperationService) BatchUpdate(ctx context.Context, userID, projectID uuid.UUID, req *BatchUpdateRequest) (*BatchOperationResult, error) {
	// 验证批量大小
	if len(req.ChapterIDs) > MaxBatchSize {
		return nil, ErrBatchLimitExceeded
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	result := &BatchOperationResult{
		Total:     len(req.ChapterIDs),
		FailedIDs: []string{},
	}

	// 使用事务
	err = s.chapterRepo.Transaction(ctx, func(ctx context.Context) error {
		for _, chapterIDStr := range req.ChapterIDs {
			chapterID, err := uuid.Parse(chapterIDStr)
			if err != nil {
				result.Failed++
				result.FailedIDs = append(result.FailedIDs, chapterIDStr)
				continue
			}

			// 获取章节
			chapter, err := s.chapterRepo.FindByID(ctx, chapterID)
			if err != nil {
				result.Failed++
				result.FailedIDs = append(result.FailedIDs, chapterIDStr)
				continue
			}

			// 验证章节属于该项目
			if chapter.ProjectID != projectID {
				result.Failed++
				result.FailedIDs = append(result.FailedIDs, chapterIDStr)
				continue
			}

			// 应用更新
			if title, ok := req.Updates["title"].(string); ok {
				chapter.Title = title
			}
			if summary, ok := req.Updates["summary"].(string); ok {
				chapter.Summary = summary
			}
			if status, ok := req.Updates["status"].(string); ok {
				chapter.Status = status
			}

			// 更新章节
			if err := s.chapterRepo.Update(ctx, chapter); err != nil {
				result.Failed++
				result.FailedIDs = append(result.FailedIDs, chapterIDStr)
				continue
			}

			result.Success++
		}

		// 如果全部失败，回滚事务
		if result.Failed == result.Total {
			return ErrBatchOperationFailed
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// BatchDelete 批量删除章节
func (s *BatchOperationService) BatchDelete(ctx context.Context, userID, projectID uuid.UUID, req *BatchDeleteRequest) (*BatchOperationResult, error) {
	// 验证批量大小
	if len(req.ChapterIDs) > MaxBatchSize {
		return nil, ErrBatchLimitExceeded
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	result := &BatchOperationResult{
		Total:     len(req.ChapterIDs),
		FailedIDs: []string{},
	}

	// 使用事务
	err = s.chapterRepo.Transaction(ctx, func(ctx context.Context) error {
		for _, chapterIDStr := range req.ChapterIDs {
			chapterID, err := uuid.Parse(chapterIDStr)
			if err != nil {
				result.Failed++
				result.FailedIDs = append(result.FailedIDs, chapterIDStr)
				continue
			}

			// 获取章节
			chapter, err := s.chapterRepo.FindByID(ctx, chapterID)
			if err != nil {
				result.Failed++
				result.FailedIDs = append(result.FailedIDs, chapterIDStr)
				continue
			}

			// 验证章节属于该项目
			if chapter.ProjectID != projectID {
				result.Failed++
				result.FailedIDs = append(result.FailedIDs, chapterIDStr)
				continue
			}

			// 删除章节
			if err := s.chapterRepo.Delete(ctx, chapterID); err != nil {
				result.Failed++
				result.FailedIDs = append(result.FailedIDs, chapterIDStr)
				continue
			}

			result.Success++
		}

		// 如果全部失败，回滚事务
		if result.Failed == result.Total {
			return ErrBatchOperationFailed
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 更新项目统计
	s.projectRepo.UpdateStats(ctx, projectID)

	return result, nil
}

// BatchStatusUpdate 批量更新章节状态
func (s *BatchOperationService) BatchStatusUpdate(ctx context.Context, userID, projectID uuid.UUID, req *BatchStatusUpdateRequest) (*BatchOperationResult, error) {
	// 验证批量大小
	if len(req.ChapterIDs) > MaxBatchSize {
		return nil, ErrBatchLimitExceeded
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	result := &BatchOperationResult{
		Total:     len(req.ChapterIDs),
		FailedIDs: []string{},
	}

	// 使用事务
	err = s.chapterRepo.Transaction(ctx, func(ctx context.Context) error {
		for _, chapterIDStr := range req.ChapterIDs {
			chapterID, err := uuid.Parse(chapterIDStr)
			if err != nil {
				result.Failed++
				result.FailedIDs = append(result.FailedIDs, chapterIDStr)
				continue
			}

			// 获取章节
			chapter, err := s.chapterRepo.FindByID(ctx, chapterID)
			if err != nil {
				result.Failed++
				result.FailedIDs = append(result.FailedIDs, chapterIDStr)
				continue
			}

			// 验证章节属于该项目
			if chapter.ProjectID != projectID {
				result.Failed++
				result.FailedIDs = append(result.FailedIDs, chapterIDStr)
				continue
			}

			// 更新状态
			chapter.Status = req.Status

			// 更新章节
			if err := s.chapterRepo.Update(ctx, chapter); err != nil {
				result.Failed++
				result.FailedIDs = append(result.FailedIDs, chapterIDStr)
				continue
			}

			result.Success++
		}

		// 如果全部失败，回滚事务
		if result.Failed == result.Total {
			return ErrBatchOperationFailed
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Reorder 重排序章节
func (s *BatchOperationService) Reorder(ctx context.Context, userID, projectID uuid.UUID, req *ReorderRequest) error {
	// 验证批量大小
	if len(req.ChapterOrders) > MaxBatchSize {
		return ErrBatchLimitExceeded
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project.UserID != userID {
		return ErrProjectNotOwned
	}

	// 构建更新列表
	updates := make([]repository.ChapterNumberUpdate, 0, len(req.ChapterOrders))
	for _, order := range req.ChapterOrders {
		chapterID, err := uuid.Parse(order.ChapterID)
		if err != nil {
			return fmt.Errorf("invalid chapter id: %s", order.ChapterID)
		}

		// 验证章节属于该项目
		chapter, err := s.chapterRepo.FindByID(ctx, chapterID)
		if err != nil {
			return err
		}
		if chapter.ProjectID != projectID {
			return ErrChapterNotInProject
		}

		updates = append(updates, repository.ChapterNumberUpdate{
			ID:            chapterID,
			ChapterNumber: order.ChapterNumber,
		})
	}

	// 使用事务更新章节号
	return s.chapterRepo.UpdateChapterNumbers(ctx, updates)
}
