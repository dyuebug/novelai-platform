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

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// getUserID 从 JWT claims 中提取用户 ID
func getUserID(c *gin.Context) (uuid.UUID, error) {
	claims, exists := c.Get("claims")
	if !exists {
		return uuid.Nil, errors.New("unauthorized")
	}

	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		return uuid.Nil, errors.New("invalid claims")
	}

	userIDStr, ok := claimsMap["sub"].(string)
	if !ok {
		return uuid.Nil, errors.New("invalid user id")
	}

	return uuid.Parse(userIDStr)
}

// Create 创建项目
// POST /api/v1/projects
func (h *ProjectHandler) Create(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req service.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	project, err := h.projectService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create project")
		return
	}

	response.Success(c, http.StatusCreated, "project created successfully", project)
}

// List 获取项目列表
// GET /api/v1/projects
func (h *ProjectHandler) List(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// 解析查询参数
	page := 1
	pageSize := 20
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
	genre := c.Query("genre")
	sort := c.Query("sort")

	result, err := h.projectService.List(c.Request.Context(), userID, page, pageSize, status, genre, sort)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to list projects")
		return
	}

	response.Success(c, http.StatusOK, "projects retrieved successfully", result)
}

// Get 获取单个项目
// GET /api/v1/projects/:id
func (h *ProjectHandler) Get(c *gin.Context) {
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

	project, err := h.projectService.Get(c.Request.Context(), userID, projectID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get project")
		}
		return
	}

	response.Success(c, http.StatusOK, "project retrieved successfully", project)
}

// Update 更新项目
// PUT /api/v1/projects/:id
func (h *ProjectHandler) Update(c *gin.Context) {
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

	var req service.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	project, err := h.projectService.Update(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update project")
		}
		return
	}

	response.Success(c, http.StatusOK, "project updated successfully", project)
}

// Delete 删除项目 (软删除)
// DELETE /api/v1/projects/:id
func (h *ProjectHandler) Delete(c *gin.Context) {
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

	err = h.projectService.Delete(c.Request.Context(), userID, projectID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to delete project")
		}
		return
	}

	response.Success(c, http.StatusOK, "project deleted successfully", nil)
}

// Restore 恢复已删除的项目
// POST /api/v1/projects/:id/restore
func (h *ProjectHandler) Restore(c *gin.Context) {
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

	project, err := h.projectService.Restore(c.Request.Context(), userID, projectID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		case errors.Is(err, service.ErrProjectNotDeleted):
			response.Error(c, http.StatusBadRequest, "project is not deleted")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to restore project")
		}
		return
	}

	response.Success(c, http.StatusOK, "project restored successfully", project)
}

// ListChapters 获取项目章节列表 (占位，将在 T4.2 实现)
// GET /api/v1/projects/:id/chapters
func (h *ProjectHandler) ListChapters(c *gin.Context) {
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

	// 验证项目所有权
	_, err = h.projectService.Get(c.Request.Context(), userID, projectID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get project")
		}
		return
	}

	// TODO: T4.2 实现章节列表
	response.Success(c, http.StatusOK, "chapters retrieved successfully", gin.H{
		"items":       []interface{}{},
		"total":       0,
		"page":        1,
		"page_size":   20,
		"total_pages": 0,
	})
}
