package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-gateway/internal/model"
	"go-gateway/internal/repository"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"
)

// UserLayoutServiceInterface 用户布局配置服务接口
type UserLayoutServiceInterface interface {
	GetLayoutConfig(ctx context.Context, userID uuid.UUID, layoutType string) (*model.UserLayoutConfig, error)
	SaveLayoutConfig(ctx context.Context, userID uuid.UUID, layoutType string, req *service.SaveLayoutConfigRequest) (*model.UserLayoutConfig, error)
	DeleteLayoutConfig(ctx context.Context, userID uuid.UUID, layoutType string) error
	GetAllLayoutConfigs(ctx context.Context, userID uuid.UUID) ([]model.UserLayoutConfig, error)
	ResetLayoutConfig(ctx context.Context, userID uuid.UUID, layoutType string) (*model.UserLayoutConfig, error)
}

// UserLayoutHandler 用户布局配置处理器
type UserLayoutHandler struct {
	layoutService UserLayoutServiceInterface
}

// NewUserLayoutHandler 创建用户布局配置处理器
func NewUserLayoutHandler(layoutService UserLayoutServiceInterface) *UserLayoutHandler {
	return &UserLayoutHandler{
		layoutService: layoutService,
	}
}

// GetLayoutConfig 获取布局配置
// GET /api/v1/user/layout-config/:layoutType
func (h *UserLayoutHandler) GetLayoutConfig(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	layoutType := c.Param("layoutType")
	if layoutType == "" {
		response.Error(c, http.StatusBadRequest, "layout type is required")
		return
	}

	config, err := h.layoutService.GetLayoutConfig(c.Request.Context(), userID, layoutType)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidLayoutType):
			response.Error(c, http.StatusBadRequest, "invalid layout type")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get layout config")
		}
		return
	}

	response.Success(c, http.StatusOK, "layout config retrieved successfully", config)
}

// SaveLayoutConfig 保存布局配置
// POST /api/v1/user/layout-config/:layoutType
func (h *UserLayoutHandler) SaveLayoutConfig(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	layoutType := c.Param("layoutType")
	if layoutType == "" {
		response.Error(c, http.StatusBadRequest, "layout type is required")
		return
	}

	var req service.SaveLayoutConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	config, err := h.layoutService.SaveLayoutConfig(c.Request.Context(), userID, layoutType, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidLayoutType):
			response.Error(c, http.StatusBadRequest, "invalid layout type")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to save layout config")
		}
		return
	}

	response.Success(c, http.StatusOK, "layout config saved successfully", config)
}

// DeleteLayoutConfig 删除布局配置
// DELETE /api/v1/user/layout-config/:layoutType
func (h *UserLayoutHandler) DeleteLayoutConfig(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	layoutType := c.Param("layoutType")
	if layoutType == "" {
		response.Error(c, http.StatusBadRequest, "layout type is required")
		return
	}

	err = h.layoutService.DeleteLayoutConfig(c.Request.Context(), userID, layoutType)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidLayoutType):
			response.Error(c, http.StatusBadRequest, "invalid layout type")
		case errors.Is(err, repository.ErrLayoutConfigNotFound):
			response.Error(c, http.StatusNotFound, "layout config not found")
		case errors.Is(err, service.ErrLayoutConfigNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to delete layout config")
		}
		return
	}

	response.Success(c, http.StatusOK, "layout config deleted successfully", nil)
}

// GetAllLayoutConfigs 获取所有布局配置
// GET /api/v1/user/layout-configs
func (h *UserLayoutHandler) GetAllLayoutConfigs(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	configs, err := h.layoutService.GetAllLayoutConfigs(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get layout configs")
		return
	}

	response.Success(c, http.StatusOK, "layout configs retrieved successfully", configs)
}

// ResetLayoutConfig 重置布局配置为默认值
// POST /api/v1/user/layout-config/:layoutType/reset
func (h *UserLayoutHandler) ResetLayoutConfig(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	layoutType := c.Param("layoutType")
	if layoutType == "" {
		response.Error(c, http.StatusBadRequest, "layout type is required")
		return
	}

	config, err := h.layoutService.ResetLayoutConfig(c.Request.Context(), userID, layoutType)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidLayoutType):
			response.Error(c, http.StatusBadRequest, "invalid layout type")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to reset layout config")
		}
		return
	}

	response.Success(c, http.StatusOK, "layout config reset successfully", config)
}
