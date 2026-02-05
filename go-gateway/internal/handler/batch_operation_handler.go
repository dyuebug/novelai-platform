package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-gateway/internal/repository"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"
)

// BatchOperationServiceInterface 批量操作服务接口
type BatchOperationServiceInterface interface {
	BatchUpdate(ctx context.Context, userID, projectID uuid.UUID, req *service.BatchUpdateRequest) (*service.BatchOperationResult, error)
	BatchDelete(ctx context.Context, userID, projectID uuid.UUID, req *service.BatchDeleteRequest) (*service.BatchOperationResult, error)
	BatchStatusUpdate(ctx context.Context, userID, projectID uuid.UUID, req *service.BatchStatusUpdateRequest) (*service.BatchOperationResult, error)
	Reorder(ctx context.Context, userID, projectID uuid.UUID, req *service.ReorderRequest) error
}

// BatchOperationHandler 批量操作处理器
type BatchOperationHandler struct {
	batchService   BatchOperationServiceInterface
	projectService *service.ProjectService
}

// NewBatchOperationHandler 创建批量操作处理器
func NewBatchOperationHandler(batchService BatchOperationServiceInterface, projectService *service.ProjectService) *BatchOperationHandler {
	return &BatchOperationHandler{
		batchService:   batchService,
		projectService: projectService,
	}
}

// BatchUpdate 批量更新章节
// POST /api/v1/projects/:id/chapters/batch-update
func (h *BatchOperationHandler) BatchUpdate(c *gin.Context) {
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

	var req service.BatchUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	result, err := h.batchService.BatchUpdate(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		case errors.Is(err, service.ErrBatchLimitExceeded):
			response.Error(c, http.StatusBadRequest, "batch limit exceeded (max 100)")
		case errors.Is(err, service.ErrBatchOperationFailed):
			response.Error(c, http.StatusInternalServerError, "all operations failed")
		default:
			response.Error(c, http.StatusInternalServerError, "batch update failed")
		}
		return
	}

	response.Success(c, http.StatusOK, "batch update completed", result)
}

// BatchDelete 批量删除章节
// POST /api/v1/projects/:id/chapters/batch-delete
func (h *BatchOperationHandler) BatchDelete(c *gin.Context) {
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

	var req service.BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	result, err := h.batchService.BatchDelete(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		case errors.Is(err, service.ErrBatchLimitExceeded):
			response.Error(c, http.StatusBadRequest, "batch limit exceeded (max 100)")
		case errors.Is(err, service.ErrBatchOperationFailed):
			response.Error(c, http.StatusInternalServerError, "all operations failed")
		default:
			response.Error(c, http.StatusInternalServerError, "batch delete failed")
		}
		return
	}

	response.Success(c, http.StatusOK, "batch delete completed", result)
}

// BatchStatusUpdate 批量更新章节状态
// POST /api/v1/projects/:id/chapters/batch-status
func (h *BatchOperationHandler) BatchStatusUpdate(c *gin.Context) {
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

	var req service.BatchStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	result, err := h.batchService.BatchStatusUpdate(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		case errors.Is(err, service.ErrBatchLimitExceeded):
			response.Error(c, http.StatusBadRequest, "batch limit exceeded (max 100)")
		case errors.Is(err, service.ErrBatchOperationFailed):
			response.Error(c, http.StatusInternalServerError, "all operations failed")
		default:
			response.Error(c, http.StatusInternalServerError, "batch status update failed")
		}
		return
	}

	response.Success(c, http.StatusOK, "batch status update completed", result)
}

// Reorder 重排序章节
// POST /api/v1/projects/:id/chapters/reorder
func (h *BatchOperationHandler) Reorder(c *gin.Context) {
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

	var req service.ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	err = h.batchService.Reorder(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		case errors.Is(err, service.ErrBatchLimitExceeded):
			response.Error(c, http.StatusBadRequest, "batch limit exceeded (max 100)")
		case errors.Is(err, service.ErrChapterNotInProject):
			response.Error(c, http.StatusBadRequest, "chapter not in project")
		default:
			response.Error(c, http.StatusInternalServerError, "reorder failed")
		}
		return
	}

	response.Success(c, http.StatusOK, "chapters reordered successfully", nil)
}
