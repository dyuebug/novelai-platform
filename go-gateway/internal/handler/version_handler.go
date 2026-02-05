package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-gateway/internal/model"
	"go-gateway/internal/repository"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"
)

// VersionServiceInterface 版本服务接口
type VersionServiceInterface interface {
	ListVersions(ctx context.Context, userID, chapterID uuid.UUID, page, pageSize int) (*service.VersionListResponse, error)
	GetVersion(ctx context.Context, userID, chapterID uuid.UUID, versionNumber int) (*model.ChapterVersion, error)
	RestoreVersion(ctx context.Context, userID, chapterID uuid.UUID, versionNumber int) (*model.Chapter, error)
	DiffVersions(ctx context.Context, userID, chapterID uuid.UUID, v1, v2 int) (*service.DiffResponse, error)
}

type VersionHandler struct {
	versionService VersionServiceInterface
}

func NewVersionHandler(versionService VersionServiceInterface) *VersionHandler {
	return &VersionHandler{
		versionService: versionService,
	}
}

// ListVersions 获取章节版本列表
// GET /api/v1/chapters/:id/versions
func (h *VersionHandler) ListVersions(c *gin.Context) {
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

	result, err := h.versionService.ListVersions(c.Request.Context(), userID, chapterID, page, pageSize)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrChapterNotFound):
			response.Error(c, http.StatusNotFound, "chapter not found")
		case errors.Is(err, service.ErrChapterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to list versions")
		}
		return
	}

	response.Success(c, http.StatusOK, "versions retrieved successfully", result)
}

// GetVersion 获取指定版本
// GET /api/v1/chapters/:id/versions/:n
func (h *VersionHandler) GetVersion(c *gin.Context) {
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

	versionNumber, err := strconv.Atoi(c.Param("n"))
	if err != nil || versionNumber < 1 {
		response.Error(c, http.StatusBadRequest, "invalid version number")
		return
	}

	version, err := h.versionService.GetVersion(c.Request.Context(), userID, chapterID, versionNumber)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrChapterNotFound):
			response.Error(c, http.StatusNotFound, "chapter not found")
		case errors.Is(err, repository.ErrVersionNotFound):
			response.Error(c, http.StatusNotFound, "version not found")
		case errors.Is(err, service.ErrChapterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get version")
		}
		return
	}

	response.Success(c, http.StatusOK, "version retrieved successfully", version)
}

// RestoreVersion 恢复到指定版本
// POST /api/v1/chapters/:id/versions/:n/restore
func (h *VersionHandler) RestoreVersion(c *gin.Context) {
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

	versionNumber, err := strconv.Atoi(c.Param("n"))
	if err != nil || versionNumber < 1 {
		response.Error(c, http.StatusBadRequest, "invalid version number")
		return
	}

	chapter, err := h.versionService.RestoreVersion(c.Request.Context(), userID, chapterID, versionNumber)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrChapterNotFound):
			response.Error(c, http.StatusNotFound, "chapter not found")
		case errors.Is(err, repository.ErrVersionNotFound):
			response.Error(c, http.StatusNotFound, "version not found")
		case errors.Is(err, service.ErrChapterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to restore version")
		}
		return
	}

	response.Success(c, http.StatusOK, "version restored successfully", chapter)
}

// DiffVersions 对比两个版本
// GET /api/v1/chapters/:id/versions/diff
func (h *VersionHandler) DiffVersions(c *gin.Context) {
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

	// 解析版本号参数
	v1, err := strconv.Atoi(c.Query("v1"))
	if err != nil || v1 < 1 {
		response.Error(c, http.StatusBadRequest, "invalid v1 parameter")
		return
	}

	v2, err := strconv.Atoi(c.Query("v2"))
	if err != nil || v2 < 1 {
		response.Error(c, http.StatusBadRequest, "invalid v2 parameter")
		return
	}

	diff, err := h.versionService.DiffVersions(c.Request.Context(), userID, chapterID, v1, v2)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrChapterNotFound):
			response.Error(c, http.StatusNotFound, "chapter not found")
		case errors.Is(err, repository.ErrVersionNotFound):
			response.Error(c, http.StatusNotFound, "version not found")
		case errors.Is(err, service.ErrChapterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to diff versions")
		}
		return
	}

	response.Success(c, http.StatusOK, "diff generated successfully", diff)
}
