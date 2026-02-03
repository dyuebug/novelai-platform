package service

import (
	"context"
	"errors"
	"unicode/utf8"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrChapterNotOwned = errors.New("chapter not owned by user")
)

type ChapterService struct {
	chapterRepo        *repository.ChapterRepository
	chapterVersionRepo *repository.ChapterVersionRepository
	projectRepo        *repository.ProjectRepository
}

func NewChapterService() *ChapterService {
	return &ChapterService{
		chapterRepo:        repository.NewChapterRepository(),
		chapterVersionRepo: repository.NewChapterVersionRepository(),
		projectRepo:        repository.NewProjectRepository(),
	}
}

// CreateChapterRequest 创建章节请求
type CreateChapterRequest struct {
	Title         string `json:"title" binding:"required,min=1,max=200"`
	Content       string `json:"content"`
	Summary       string `json:"summary"`
	ChapterNumber *int   `json:"chapter_number"`
	POVCharacter  string `json:"pov_character"`
	Location      string `json:"location"`
	TimeSetting   string `json:"time_setting"`
}

// UpdateChapterRequest 更新章节请求
type UpdateChapterRequest struct {
	Title        *string `json:"title" binding:"omitempty,min=1,max=200"`
	Content      *string `json:"content"`
	Summary      *string `json:"summary"`
	Status       *string `json:"status"`
	POVCharacter *string `json:"pov_character"`
	Location     *string `json:"location"`
	TimeSetting  *string `json:"time_setting"`
}

// ReorderChaptersRequest 重排序请求
type ReorderChaptersRequest struct {
	ChapterIDs []uuid.UUID `json:"chapter_ids" binding:"required,min=1"`
}

// ChapterListResponse 章节列表响应
type ChapterListResponse struct {
	Items      []model.Chapter `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// Create 创建章节
func (s *ChapterService) Create(ctx context.Context, userID, projectID uuid.UUID, req *CreateChapterRequest) (*model.Chapter, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	// 确定章节号
	chapterNumber := 1
	if req.ChapterNumber != nil {
		chapterNumber = *req.ChapterNumber
	} else {
		maxNum, err := s.chapterRepo.GetMaxChapterNumber(ctx, projectID)
		if err != nil {
			return nil, err
		}
		chapterNumber = maxNum + 1
	}

	// 计算字数
	wordCount := utf8.RuneCountInString(req.Content)

	chapter := &model.Chapter{
		ProjectID:     projectID,
		ChapterNumber: chapterNumber,
		Title:         req.Title,
		Content:       req.Content,
		Summary:       req.Summary,
		Status:        "draft",
		WordCount:     wordCount,
		POVCharacter:  req.POVCharacter,
		Location:      req.Location,
		TimeSetting:   req.TimeSetting,
	}

	if err := s.chapterRepo.Create(ctx, chapter); err != nil {
		return nil, err
	}

	// 创建初始版本
	if req.Content != "" {
		version := &model.ChapterVersion{
			ChapterID:     chapter.ID,
			VersionNumber: 1,
			Content:       req.Content,
			WordCount:     wordCount,
			Source:        "manual",
		}
		if err := s.chapterVersionRepo.Create(ctx, version); err != nil {
			// 版本创建失败不影响章节创建
		}
	}

	// 更新项目统计
	_ = s.projectRepo.UpdateStats(ctx, projectID)

	return chapter, nil
}

// Get 获取章节
func (s *ChapterService) Get(ctx context.Context, userID, chapterID uuid.UUID) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.FindByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, chapter.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrChapterNotOwned
	}

	return chapter, nil
}

// List 获取章节列表
func (s *ChapterService) List(ctx context.Context, userID, projectID uuid.UUID, page, pageSize int, status, sort string) (*ChapterListResponse, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	// 限制 pageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}
	if page <= 0 {
		page = 1
	}

	opts := &repository.ListOptions{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Sort:     sort,
	}

	chapters, total, err := s.chapterRepo.FindByProjectID(ctx, projectID, opts)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &ChapterListResponse{
		Items:      chapters,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update 更新章节
func (s *ChapterService) Update(ctx context.Context, userID, chapterID uuid.UUID, req *UpdateChapterRequest) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.FindByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, chapter.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrChapterNotOwned
	}

	// 记录旧内容用于版本控制
	oldContent := chapter.Content
	contentChanged := false

	// 更新字段
	if req.Title != nil {
		chapter.Title = *req.Title
	}
	if req.Content != nil {
		chapter.Content = *req.Content
		chapter.WordCount = utf8.RuneCountInString(*req.Content)
		contentChanged = oldContent != *req.Content
	}
	if req.Summary != nil {
		chapter.Summary = *req.Summary
	}
	if req.Status != nil {
		chapter.Status = *req.Status
	}
	if req.POVCharacter != nil {
		chapter.POVCharacter = *req.POVCharacter
	}
	if req.Location != nil {
		chapter.Location = *req.Location
	}
	if req.TimeSetting != nil {
		chapter.TimeSetting = *req.TimeSetting
	}

	if err := s.chapterRepo.Update(ctx, chapter); err != nil {
		return nil, err
	}

	// 内容变更时创建新版本
	if contentChanged && chapter.Content != "" {
		maxVersion, _ := s.chapterVersionRepo.GetMaxVersionNumber(ctx, chapterID)
		version := &model.ChapterVersion{
			ChapterID:     chapterID,
			VersionNumber: maxVersion + 1,
			Content:       chapter.Content,
			WordCount:     chapter.WordCount,
			Source:        "manual",
		}
		_ = s.chapterVersionRepo.Create(ctx, version)
	}

	// 更新项目统计
	_ = s.projectRepo.UpdateStats(ctx, chapter.ProjectID)

	return chapter, nil
}

// Delete 删除章节
func (s *ChapterService) Delete(ctx context.Context, userID, chapterID uuid.UUID) error {
	chapter, err := s.chapterRepo.FindByID(ctx, chapterID)
	if err != nil {
		return err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, chapter.ProjectID)
	if err != nil {
		return err
	}
	if project.UserID != userID {
		return ErrChapterNotOwned
	}

	projectID := chapter.ProjectID

	if err := s.chapterRepo.Delete(ctx, chapterID); err != nil {
		return err
	}

	// 更新项目统计
	_ = s.projectRepo.UpdateStats(ctx, projectID)

	return nil
}

// Reorder 重排序章节
func (s *ChapterService) Reorder(ctx context.Context, userID, projectID uuid.UUID, req *ReorderChaptersRequest) error {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project.UserID != userID {
		return ErrProjectNotOwned
	}

	// 构建更新列表
	updates := make([]repository.ChapterNumberUpdate, len(req.ChapterIDs))
	for i, id := range req.ChapterIDs {
		updates[i] = repository.ChapterNumberUpdate{
			ID:            id,
			ChapterNumber: i + 1,
		}
	}

	return s.chapterRepo.UpdateChapterNumbers(ctx, updates)
}
