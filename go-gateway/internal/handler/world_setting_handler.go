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

// ========== 世界设定 Handler ==========

type WorldSettingHandler struct {
	settingService *service.WorldSettingService
}

func NewWorldSettingHandler(settingService *service.WorldSettingService) *WorldSettingHandler {
	return &WorldSettingHandler{
		settingService: settingService,
	}
}

// Create 创建世界设定
// POST /api/v1/projects/:id/world-settings
func (h *WorldSettingHandler) Create(c *gin.Context) {
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

	var req service.CreateWorldSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	setting, err := h.settingService.Create(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to create world setting")
		}
		return
	}

	response.Success(c, http.StatusCreated, "world setting created successfully", setting)
}

// List 获取世界设定列表
// GET /api/v1/projects/:id/world-settings
func (h *WorldSettingHandler) List(c *gin.Context) {
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

	category := c.Query("category")
	sort := c.Query("sort")

	result, err := h.settingService.List(c.Request.Context(), userID, projectID, page, pageSize, category, sort)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to list world settings")
		}
		return
	}

	response.Success(c, http.StatusOK, "world settings retrieved successfully", result)
}

// GetTree 获取世界设定树结构
// GET /api/v1/projects/:id/world-settings/tree
func (h *WorldSettingHandler) GetTree(c *gin.Context) {
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

	category := c.Query("category")

	tree, err := h.settingService.GetTree(c.Request.Context(), userID, projectID, category)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get world setting tree")
		}
		return
	}

	response.Success(c, http.StatusOK, "world setting tree retrieved successfully", tree)
}

// GetByCategory 按分类获取设定
// GET /api/v1/projects/:id/world-settings/category/:category
func (h *WorldSettingHandler) GetByCategory(c *gin.Context) {
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

	category := c.Param("category")
	if category == "" {
		response.Error(c, http.StatusBadRequest, "category is required")
		return
	}

	settings, err := h.settingService.GetByCategory(c.Request.Context(), userID, projectID, category)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get world settings")
		}
		return
	}

	response.Success(c, http.StatusOK, "world settings retrieved successfully", settings)
}

// Get 获取单个世界设定
// GET /api/v1/world-settings/:id
func (h *WorldSettingHandler) Get(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	settingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid setting id")
		return
	}

	setting, err := h.settingService.Get(c.Request.Context(), userID, settingID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrWorldSettingNotFound):
			response.Error(c, http.StatusNotFound, "world setting not found")
		case errors.Is(err, service.ErrWorldSettingNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get world setting")
		}
		return
	}

	response.Success(c, http.StatusOK, "world setting retrieved successfully", setting)
}

// Update 更新世界设定
// PUT /api/v1/world-settings/:id
func (h *WorldSettingHandler) Update(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	settingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid setting id")
		return
	}

	var req service.UpdateWorldSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	setting, err := h.settingService.Update(c.Request.Context(), userID, settingID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrWorldSettingNotFound):
			response.Error(c, http.StatusNotFound, "world setting not found")
		case errors.Is(err, service.ErrWorldSettingNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update world setting")
		}
		return
	}

	response.Success(c, http.StatusOK, "world setting updated successfully", setting)
}

// Delete 删除世界设定
// DELETE /api/v1/world-settings/:id
func (h *WorldSettingHandler) Delete(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	settingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid setting id")
		return
	}

	err = h.settingService.Delete(c.Request.Context(), userID, settingID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrWorldSettingNotFound):
			response.Error(c, http.StatusNotFound, "world setting not found")
		case errors.Is(err, service.ErrWorldSettingNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to delete world setting")
		}
		return
	}

	response.Success(c, http.StatusOK, "world setting deleted successfully", nil)
}
