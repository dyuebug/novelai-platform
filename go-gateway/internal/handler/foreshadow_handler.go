package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/service"
)

// ForeshadowHandler 伏笔处理器
type ForeshadowHandler struct {
	svc *service.ForeshadowService
}

// NewForeshadowHandler 创建伏笔处理器
func NewForeshadowHandler(svc *service.ForeshadowService) *ForeshadowHandler {
	return &ForeshadowHandler{svc: svc}
}

// Create 创建伏笔
func (h *ForeshadowHandler) Create(c *gin.Context) {
	projectID := c.Param("id")

	var req service.CreateForeshadowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foreshadow, err := h.svc.Create(projectID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, foreshadow)
}

// Get 获取伏笔
func (h *ForeshadowHandler) Get(c *gin.Context) {
	id := c.Param("id")

	foreshadow, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "foreshadow not found"})
		return
	}

	c.JSON(http.StatusOK, foreshadow)
}

// Update 更新伏笔
func (h *ForeshadowHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateForeshadowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foreshadow, err := h.svc.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, foreshadow)
}

// Delete 删除伏笔
func (h *ForeshadowHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// List 获取伏笔列表
func (h *ForeshadowHandler) List(c *gin.Context) {
	projectID := c.Param("id")
	status := c.Query("status")
	priority := c.Query("priority")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	foreshadows, total, err := h.svc.List(projectID, status, priority, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": foreshadows,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// GetStats 获取伏笔统计
func (h *ForeshadowHandler) GetStats(c *gin.Context) {
	projectID := c.Param("id")
	currentChapter, _ := strconv.Atoi(c.DefaultQuery("current_chapter", "0"))

	stats, err := h.svc.GetStats(projectID, currentChapter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ResolveRequest 回收伏笔请求
type ResolveRequest struct {
	ChapterID  string `json:"chapter_id"`
	ChapterNum int    `json:"chapter_num"`
	Content    string `json:"content"`
}

// Resolve 回收伏笔
func (h *ForeshadowHandler) Resolve(c *gin.Context) {
	id := c.Param("id")

	var req ResolveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foreshadow, err := h.svc.Resolve(id, req.ChapterID, req.ChapterNum, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, foreshadow)
}

// AddHintRequest 添加暗示请求
type AddHintRequest struct {
	ChapterID  string `json:"chapter_id" binding:"required"`
	ChapterNum int    `json:"chapter_num"`
	Content    string `json:"content" binding:"required"`
	HintType   string `json:"hint_type"`
}

// AddHint 添加暗示
func (h *ForeshadowHandler) AddHint(c *gin.Context) {
	foreshadowID := c.Param("id")

	var req AddHintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hint, err := h.svc.AddHint(foreshadowID, req.ChapterID, req.ChapterNum, req.Content, req.HintType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, hint)
}

// GetHints 获取伏笔的暗示列表
func (h *ForeshadowHandler) GetHints(c *gin.Context) {
	foreshadowID := c.Param("id")

	hints, err := h.svc.GetHints(foreshadowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hints)
}

// DeleteHint 删除暗示
func (h *ForeshadowHandler) DeleteHint(c *gin.Context) {
	hintID := c.Param("id")

	if err := h.svc.DeleteHint(hintID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetPendingReminders 获取待提醒的伏笔
func (h *ForeshadowHandler) GetPendingReminders(c *gin.Context) {
	projectID := c.Param("id")
	currentChapter, _ := strconv.Atoi(c.Query("current_chapter"))

	foreshadows, err := h.svc.GetPendingReminders(projectID, currentChapter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, foreshadows)
}

// GetOverdueForeshadows 获取超期未回收的伏笔
func (h *ForeshadowHandler) GetOverdueForeshadows(c *gin.Context) {
	projectID := c.Param("id")
	currentChapter, _ := strconv.Atoi(c.Query("current_chapter"))

	foreshadows, err := h.svc.GetOverdueForeshadows(projectID, currentChapter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, foreshadows)
}

// CreateReminderRequest 创建提醒请求
type CreateReminderRequest struct {
	ChapterNum int    `json:"chapter_num" binding:"required"`
	Message    string `json:"message" binding:"required"`
}

// CreateReminder 创建提醒
func (h *ForeshadowHandler) CreateReminder(c *gin.Context) {
	foreshadowID := c.Param("id")

	var req CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminder, err := h.svc.CreateReminder(foreshadowID, req.ChapterNum, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reminder)
}

// GetUnreadReminders 获取未读提醒
func (h *ForeshadowHandler) GetUnreadReminders(c *gin.Context) {
	projectID := c.Param("id")

	reminders, err := h.svc.GetUnreadReminders(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

// MarkReminderAsRead 标记提醒为已读
func (h *ForeshadowHandler) MarkReminderAsRead(c *gin.Context) {
	reminderID := c.Param("id")

	if err := h.svc.MarkReminderAsRead(reminderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

// MarkAllRemindersAsRead 标记所有提醒为已读
func (h *ForeshadowHandler) MarkAllRemindersAsRead(c *gin.Context) {
	projectID := c.Param("id")

	if err := h.svc.MarkAllRemindersAsRead(projectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all marked as read"})
}

// CheckAndCreateReminders 检查并创建提醒
func (h *ForeshadowHandler) CheckAndCreateReminders(c *gin.Context) {
	projectID := c.Param("id")
	currentChapter, _ := strconv.Atoi(c.Query("current_chapter"))

	reminders, err := h.svc.CheckAndCreateReminders(projectID, currentChapter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"created_reminders": reminders,
		"count":             len(reminders),
	})
}
