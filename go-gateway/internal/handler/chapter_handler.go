package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-gateway/internal/repository"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"
)

type ChapterHandler struct {
	chapterService *service.ChapterService
	projectService *service.ProjectService
}

func NewChapterHandler(chapterService *service.ChapterService, projectService *service.ProjectService) *ChapterHandler {
	return &ChapterHandler{
		chapterService: chapterService,
		projectService: projectService,
	}
}

// Create 创建章节
// POST /api/v1/projects/:id/chapters
func (h *ChapterHandler) Create(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	var req service.CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	chapter, err := h.chapterService.Create(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to create chapter")
		}
		return
	}

	response.Success(c, http.StatusCreated, "chapter created successfully", chapter)
}

// List 获取章节列表
// GET /api/v1/projects/:id/chapters
func (h *ChapterHandler) List(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	// 解析查询参数
	page := 1
	pageSize := 50
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}

	status := c.Query("status")
	sort := c.Query("sort")

	result, err := h.chapterService.List(c.Request.Context(), userID, projectID, page, pageSize, status, sort)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to list chapters")
		}
		return
	}

	response.Success(c, http.StatusOK, "chapters retrieved successfully", result)
}

// Get 获取单个章节
// GET /api/v1/chapters/:id
func (h *ChapterHandler) Get(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	chapterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid chapter id")
		return
	}

	chapter, err := h.chapterService.Get(c.Request.Context(), userID, chapterID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrChapterNotFound):
			response.Error(c, http.StatusNotFound, "chapter not found")
		case errors.Is(err, service.ErrChapterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get chapter")
		}
		return
	}

	response.Success(c, http.StatusOK, "chapter retrieved successfully", chapter)
}

// Update 更新章节
// PUT /api/v1/chapters/:id
func (h *ChapterHandler) Update(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	chapterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid chapter id")
		return
	}

	var req service.UpdateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	chapter, err := h.chapterService.Update(c.Request.Context(), userID, chapterID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrChapterNotFound):
			response.Error(c, http.StatusNotFound, "chapter not found")
		case errors.Is(err, service.ErrChapterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update chapter")
		}
		return
	}

	response.Success(c, http.StatusOK, "chapter updated successfully", chapter)
}

// Delete 删除章节
// DELETE /api/v1/chapters/:id
func (h *ChapterHandler) Delete(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	chapterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid chapter id")
		return
	}

	err = h.chapterService.Delete(c.Request.Context(), userID, chapterID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrChapterNotFound):
			response.Error(c, http.StatusNotFound, "chapter not found")
		case errors.Is(err, service.ErrChapterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to delete chapter")
		}
		return
	}

	response.Success(c, http.StatusOK, "chapter deleted successfully", nil)
}

// Reorder 重排序章节
// POST /api/v1/projects/:id/chapters/reorder
func (h *ChapterHandler) Reorder(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	var req service.ReorderChaptersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	err = h.chapterService.Reorder(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to reorder chapters")
		}
		return
	}

	response.Success(c, http.StatusOK, "chapters reordered successfully", nil)
}
