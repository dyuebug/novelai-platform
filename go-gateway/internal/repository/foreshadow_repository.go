package repository

import (
	"go-gateway/internal/database"
	"go-gateway/internal/model"
)

// ForeshadowRepository 伏笔仓库
type ForeshadowRepository struct{}

// NewForeshadowRepository 创建伏笔仓库
func NewForeshadowRepository() *ForeshadowRepository {
	return &ForeshadowRepository{}
}

// Create 创建伏笔
func (r *ForeshadowRepository) Create(foreshadow *model.Foreshadow) error {
	return database.DB.Create(foreshadow).Error
}

// GetByID 根据ID获取伏笔
func (r *ForeshadowRepository) GetByID(id string) (*model.Foreshadow, error) {
	var foreshadow model.Foreshadow
	err := database.DB.Where("id = ?", id).First(&foreshadow).Error
	if err != nil {
		return nil, err
	}
	return &foreshadow, nil
}

// Update 更新伏笔
func (r *ForeshadowRepository) Update(foreshadow *model.Foreshadow) error {
	return database.DB.Save(foreshadow).Error
}

// Delete 删除伏笔
func (r *ForeshadowRepository) Delete(id string) error {
	return database.DB.Delete(&model.Foreshadow{}, "id = ?", id).Error
}

// ListByProject 获取项目的伏笔列表
func (r *ForeshadowRepository) ListByProject(projectID string, status string, priority string, page, pageSize int) ([]model.Foreshadow, int64, error) {
	var foreshadows []model.Foreshadow
	var total int64

	query := database.DB.Model(&model.Foreshadow{}).Where("project_id = ?", projectID)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("sort_order ASC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&foreshadows).Error
	if err != nil {
		return nil, 0, err
	}

	return foreshadows, total, nil
}

// ListByChapter 获取章节相关的伏笔
func (r *ForeshadowRepository) ListByChapter(chapterID string) ([]model.Foreshadow, error) {
	var foreshadows []model.Foreshadow
	err := database.DB.Where("plant_chapter_id = ? OR resolve_chapter_id = ?", chapterID, chapterID).
		Find(&foreshadows).Error
	return foreshadows, err
}

// ListPendingReminders 获取待提醒的伏笔
func (r *ForeshadowRepository) ListPendingReminders(projectID string, currentChapter int) ([]model.Foreshadow, error) {
	var foreshadows []model.Foreshadow
	err := database.DB.Where("project_id = ? AND status = ? AND remind_enabled = ? AND remind_chapter_num <= ? AND remind_chapter_num > 0",
		projectID, model.ForeshadowStatusPlanted, true, currentChapter).
		Find(&foreshadows).Error
	return foreshadows, err
}

// ListOverdue 获取超期未回收的伏笔
func (r *ForeshadowRepository) ListOverdue(projectID string, currentChapter int, threshold int) ([]model.Foreshadow, error) {
	var foreshadows []model.Foreshadow
	// 超过 threshold 章未回收的伏笔
	err := database.DB.Where("project_id = ? AND status = ? AND plant_chapter_num > 0 AND plant_chapter_num + ? < ?",
		projectID, model.ForeshadowStatusPlanted, threshold, currentChapter).
		Find(&foreshadows).Error
	return foreshadows, err
}

// GetStats 获取伏笔统计
func (r *ForeshadowRepository) GetStats(projectID string, currentChapter int) (*model.ForeshadowStats, error) {
	stats := &model.ForeshadowStats{}
	var count int64

	// 总数
	database.DB.Model(&model.Foreshadow{}).Where("project_id = ?", projectID).Count(&count)
	stats.Total = int(count)

	// 按状态统计
	database.DB.Model(&model.Foreshadow{}).Where("project_id = ? AND status = ?", projectID, model.ForeshadowStatusPlanted).Count(&count)
	stats.Planted = int(count)
	database.DB.Model(&model.Foreshadow{}).Where("project_id = ? AND status = ?", projectID, model.ForeshadowStatusHinted).Count(&count)
	stats.Hinted = int(count)
	database.DB.Model(&model.Foreshadow{}).Where("project_id = ? AND status = ?", projectID, model.ForeshadowStatusResolved).Count(&count)
	stats.Resolved = int(count)
	database.DB.Model(&model.Foreshadow{}).Where("project_id = ? AND status = ?", projectID, model.ForeshadowStatusAbandoned).Count(&count)
	stats.Abandoned = int(count)

	// 按优先级统计
	database.DB.Model(&model.Foreshadow{}).Where("project_id = ? AND priority = ?", projectID, model.ForeshadowPriorityHigh).Count(&count)
	stats.HighPriority = int(count)
	database.DB.Model(&model.Foreshadow{}).Where("project_id = ? AND priority = ?", projectID, model.ForeshadowPriorityMedium).Count(&count)
	stats.MediumPriority = int(count)
	database.DB.Model(&model.Foreshadow{}).Where("project_id = ? AND priority = ?", projectID, model.ForeshadowPriorityLow).Count(&count)
	stats.LowPriority = int(count)

	// 超期数量 (超过50章未回收)
	database.DB.Model(&model.Foreshadow{}).
		Where("project_id = ? AND status = ? AND plant_chapter_num > 0 AND plant_chapter_num + 50 < ?",
			projectID, model.ForeshadowStatusPlanted, currentChapter).
		Count(&count)
	stats.OverdueCount = int(count)

	return stats, nil
}

// ForeshadowHintRepository 伏笔暗示仓库
type ForeshadowHintRepository struct{}

// NewForeshadowHintRepository 创建伏笔暗示仓库
func NewForeshadowHintRepository() *ForeshadowHintRepository {
	return &ForeshadowHintRepository{}
}

// Create 创建暗示
func (r *ForeshadowHintRepository) Create(hint *model.ForeshadowHint) error {
	return database.DB.Create(hint).Error
}

// ListByForeshadow 获取伏笔的暗示列表
func (r *ForeshadowHintRepository) ListByForeshadow(foreshadowID string) ([]model.ForeshadowHint, error) {
	var hints []model.ForeshadowHint
	err := database.DB.Where("foreshadow_id = ?", foreshadowID).
		Order("chapter_num ASC").
		Find(&hints).Error
	return hints, err
}

// Delete 删除暗示
func (r *ForeshadowHintRepository) Delete(id string) error {
	return database.DB.Delete(&model.ForeshadowHint{}, "id = ?", id).Error
}

// ForeshadowReminderRepository 伏笔提醒仓库
type ForeshadowReminderRepository struct{}

// NewForeshadowReminderRepository 创建伏笔提醒仓库
func NewForeshadowReminderRepository() *ForeshadowReminderRepository {
	return &ForeshadowReminderRepository{}
}

// Create 创建提醒
func (r *ForeshadowReminderRepository) Create(reminder *model.ForeshadowReminder) error {
	return database.DB.Create(reminder).Error
}

// ListUnread 获取未读提醒
func (r *ForeshadowReminderRepository) ListUnread(projectID string) ([]model.ForeshadowReminder, error) {
	var reminders []model.ForeshadowReminder
	err := database.DB.Joins("JOIN foreshadows ON foreshadows.id = foreshadow_reminders.foreshadow_id").
		Where("foreshadows.project_id = ? AND foreshadow_reminders.is_read = ?", projectID, false).
		Preload("Foreshadow").
		Find(&reminders).Error
	return reminders, err
}

// MarkAsRead 标记为已读
func (r *ForeshadowReminderRepository) MarkAsRead(id string) error {
	return database.DB.Model(&model.ForeshadowReminder{}).Where("id = ?", id).Update("is_read", true).Error
}

// MarkAllAsRead 标记所有为已读
func (r *ForeshadowReminderRepository) MarkAllAsRead(projectID string) error {
	return database.DB.Model(&model.ForeshadowReminder{}).
		Joins("JOIN foreshadows ON foreshadows.id = foreshadow_reminders.foreshadow_id").
		Where("foreshadows.project_id = ?", projectID).
		Update("is_read", true).Error
}
