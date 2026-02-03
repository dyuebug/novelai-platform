package service

import (
	"errors"
	"time"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"
)

// ForeshadowService 伏笔服务
type ForeshadowService struct {
	repo         *repository.ForeshadowRepository
	hintRepo     *repository.ForeshadowHintRepository
	reminderRepo *repository.ForeshadowReminderRepository
}

// NewForeshadowService 创建伏笔服务
func NewForeshadowService() *ForeshadowService {
	return &ForeshadowService{
		repo:         repository.NewForeshadowRepository(),
		hintRepo:     repository.NewForeshadowHintRepository(),
		reminderRepo: repository.NewForeshadowReminderRepository(),
	}
}

// CreateForeshadowRequest 创建伏笔请求
type CreateForeshadowRequest struct {
	Title               string   `json:"title" binding:"required"`
	Description         string   `json:"description"`
	Content             string   `json:"content"`
	Priority            string   `json:"priority"`
	PlantChapterID      string   `json:"plant_chapter_id"`
	PlantChapterNum     int      `json:"plant_chapter_num"`
	PlantPosition       string   `json:"plant_position"`
	RemindChapterNum    int      `json:"remind_chapter_num"`
	RelatedCharacterIDs []string `json:"related_character_ids"`
	Tags                []string `json:"tags"`
}

// UpdateForeshadowRequest 更新伏笔请求
type UpdateForeshadowRequest struct {
	Title               *string  `json:"title"`
	Description         *string  `json:"description"`
	Content             *string  `json:"content"`
	Priority            *string  `json:"priority"`
	Status              *string  `json:"status"`
	PlantChapterID      *string  `json:"plant_chapter_id"`
	PlantChapterNum     *int     `json:"plant_chapter_num"`
	PlantPosition       *string  `json:"plant_position"`
	ResolveChapterID    *string  `json:"resolve_chapter_id"`
	ResolveChapterNum   *int     `json:"resolve_chapter_num"`
	ResolveContent      *string  `json:"resolve_content"`
	RemindChapterNum    *int     `json:"remind_chapter_num"`
	RemindEnabled       *bool    `json:"remind_enabled"`
	RelatedCharacterIDs []string `json:"related_character_ids"`
	Tags                []string `json:"tags"`
}

// Create 创建伏笔
func (s *ForeshadowService) Create(projectID string, req *CreateForeshadowRequest) (*model.Foreshadow, error) {
	priority := model.ForeshadowPriorityMedium
	if req.Priority != "" {
		priority = model.ForeshadowPriority(req.Priority)
	}

	foreshadow := &model.Foreshadow{
		ProjectID:           projectID,
		Title:               req.Title,
		Description:         req.Description,
		Content:             req.Content,
		Priority:            priority,
		Status:              model.ForeshadowStatusPlanted,
		PlantChapterNum:     req.PlantChapterNum,
		PlantPosition:       req.PlantPosition,
		RemindChapterNum:    req.RemindChapterNum,
		RemindEnabled:       true,
		RelatedCharacterIDs: req.RelatedCharacterIDs,
		Tags:                req.Tags,
	}

	if req.PlantChapterID != "" {
		foreshadow.PlantChapterID = &req.PlantChapterID
	}

	if err := s.repo.Create(foreshadow); err != nil {
		return nil, err
	}

	return foreshadow, nil
}

// GetByID 获取伏笔
func (s *ForeshadowService) GetByID(id string) (*model.Foreshadow, error) {
	return s.repo.GetByID(id)
}

// Update 更新伏笔
func (s *ForeshadowService) Update(id string, req *UpdateForeshadowRequest) (*model.Foreshadow, error) {
	foreshadow, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		foreshadow.Title = *req.Title
	}
	if req.Description != nil {
		foreshadow.Description = *req.Description
	}
	if req.Content != nil {
		foreshadow.Content = *req.Content
	}
	if req.Priority != nil {
		foreshadow.Priority = model.ForeshadowPriority(*req.Priority)
	}
	if req.Status != nil {
		foreshadow.Status = model.ForeshadowStatus(*req.Status)
		// 如果状态变为已回收，记录时间
		if foreshadow.Status == model.ForeshadowStatusResolved {
			now := time.Now()
			foreshadow.ResolvedAt = &now
		}
	}
	if req.PlantChapterID != nil {
		foreshadow.PlantChapterID = req.PlantChapterID
	}
	if req.PlantChapterNum != nil {
		foreshadow.PlantChapterNum = *req.PlantChapterNum
	}
	if req.PlantPosition != nil {
		foreshadow.PlantPosition = *req.PlantPosition
	}
	if req.ResolveChapterID != nil {
		foreshadow.ResolveChapterID = req.ResolveChapterID
	}
	if req.ResolveChapterNum != nil {
		foreshadow.ResolveChapterNum = *req.ResolveChapterNum
	}
	if req.ResolveContent != nil {
		foreshadow.ResolveContent = *req.ResolveContent
	}
	if req.RemindChapterNum != nil {
		foreshadow.RemindChapterNum = *req.RemindChapterNum
	}
	if req.RemindEnabled != nil {
		foreshadow.RemindEnabled = *req.RemindEnabled
	}
	if req.RelatedCharacterIDs != nil {
		foreshadow.RelatedCharacterIDs = req.RelatedCharacterIDs
	}
	if req.Tags != nil {
		foreshadow.Tags = req.Tags
	}

	if err := s.repo.Update(foreshadow); err != nil {
		return nil, err
	}

	return foreshadow, nil
}

// Delete 删除伏笔
func (s *ForeshadowService) Delete(id string) error {
	return s.repo.Delete(id)
}

// List 获取伏笔列表
func (s *ForeshadowService) List(projectID string, status string, priority string, page, pageSize int) ([]model.Foreshadow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.ListByProject(projectID, status, priority, page, pageSize)
}

// GetStats 获取伏笔统计
func (s *ForeshadowService) GetStats(projectID string, currentChapter int) (*model.ForeshadowStats, error) {
	return s.repo.GetStats(projectID, currentChapter)
}

// Resolve 回收伏笔
func (s *ForeshadowService) Resolve(id string, chapterID string, chapterNum int, content string) (*model.Foreshadow, error) {
	foreshadow, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if foreshadow.Status == model.ForeshadowStatusResolved {
		return nil, errors.New("foreshadow already resolved")
	}

	foreshadow.Status = model.ForeshadowStatusResolved
	foreshadow.ResolveChapterID = &chapterID
	foreshadow.ResolveChapterNum = chapterNum
	foreshadow.ResolveContent = content
	now := time.Now()
	foreshadow.ResolvedAt = &now

	if err := s.repo.Update(foreshadow); err != nil {
		return nil, err
	}

	return foreshadow, nil
}

// AddHint 添加暗示
func (s *ForeshadowService) AddHint(foreshadowID string, chapterID string, chapterNum int, content string, hintType string) (*model.ForeshadowHint, error) {
	// 验证伏笔存在
	foreshadow, err := s.repo.GetByID(foreshadowID)
	if err != nil {
		return nil, err
	}

	// 更新伏笔状态为已暗示
	if foreshadow.Status == model.ForeshadowStatusPlanted {
		foreshadow.Status = model.ForeshadowStatusHinted
		s.repo.Update(foreshadow)
	}

	hint := &model.ForeshadowHint{
		ForeshadowID: foreshadowID,
		ChapterID:    chapterID,
		ChapterNum:   chapterNum,
		Content:      content,
		HintType:     hintType,
	}

	if err := s.hintRepo.Create(hint); err != nil {
		return nil, err
	}

	return hint, nil
}

// GetHints 获取伏笔的暗示列表
func (s *ForeshadowService) GetHints(foreshadowID string) ([]model.ForeshadowHint, error) {
	return s.hintRepo.ListByForeshadow(foreshadowID)
}

// DeleteHint 删除暗示
func (s *ForeshadowService) DeleteHint(hintID string) error {
	return s.hintRepo.Delete(hintID)
}

// GetPendingReminders 获取待提醒的伏笔
func (s *ForeshadowService) GetPendingReminders(projectID string, currentChapter int) ([]model.Foreshadow, error) {
	return s.repo.ListPendingReminders(projectID, currentChapter)
}

// GetOverdueForeshadows 获取超期未回收的伏笔
func (s *ForeshadowService) GetOverdueForeshadows(projectID string, currentChapter int) ([]model.Foreshadow, error) {
	return s.repo.ListOverdue(projectID, currentChapter, 50) // 默认50章为超期阈值
}

// CreateReminder 创建提醒
func (s *ForeshadowService) CreateReminder(foreshadowID string, chapterNum int, message string) (*model.ForeshadowReminder, error) {
	reminder := &model.ForeshadowReminder{
		ForeshadowID: foreshadowID,
		ChapterNum:   chapterNum,
		Message:      message,
		IsRead:       false,
	}

	if err := s.reminderRepo.Create(reminder); err != nil {
		return nil, err
	}

	return reminder, nil
}

// GetUnreadReminders 获取未读提醒
func (s *ForeshadowService) GetUnreadReminders(projectID string) ([]model.ForeshadowReminder, error) {
	return s.reminderRepo.ListUnread(projectID)
}

// MarkReminderAsRead 标记提醒为已读
func (s *ForeshadowService) MarkReminderAsRead(reminderID string) error {
	return s.reminderRepo.MarkAsRead(reminderID)
}

// MarkAllRemindersAsRead 标记所有提醒为已读
func (s *ForeshadowService) MarkAllRemindersAsRead(projectID string) error {
	return s.reminderRepo.MarkAllAsRead(projectID)
}

// CheckAndCreateReminders 检查并创建提醒
func (s *ForeshadowService) CheckAndCreateReminders(projectID string, currentChapter int) ([]model.ForeshadowReminder, error) {
	// 获取需要提醒的伏笔
	pending, err := s.GetPendingReminders(projectID, currentChapter)
	if err != nil {
		return nil, err
	}

	var reminders []model.ForeshadowReminder
	for _, f := range pending {
		reminder, err := s.CreateReminder(f.ID, currentChapter, "伏笔「"+f.Title+"」已到达提醒章节，请考虑回收")
		if err == nil {
			reminders = append(reminders, *reminder)
		}
	}

	// 获取超期伏笔
	overdue, err := s.GetOverdueForeshadows(projectID, currentChapter)
	if err != nil {
		return reminders, nil
	}

	for _, f := range overdue {
		reminder, err := s.CreateReminder(f.ID, currentChapter, "伏笔「"+f.Title+"」已超过50章未回收，请尽快处理")
		if err == nil {
			reminders = append(reminders, *reminder)
		}
	}

	return reminders, nil
}
