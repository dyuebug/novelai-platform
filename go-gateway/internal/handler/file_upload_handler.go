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

// FileUploadServiceInterface 文件上传服务接口
type FileUploadServiceInterface interface {
	UploadFile(ctx context.Context, userID uuid.UUID, req *service.UploadFileRequest) (*model.File, error)
	GetFile(ctx context.Context, userID, fileID uuid.UUID) (*model.File, error)
	DeleteFile(ctx context.Context, userID, fileID uuid.UUID) error
	ListFiles(ctx context.Context, userID uuid.UUID, page, pageSize int, fileType, sort string) (*service.FileListResponse, error)
	GetUserStorageStats(ctx context.Context, userID uuid.UUID) (*service.StorageStats, error)
}

// FileUploadHandler 文件上传处理器
type FileUploadHandler struct {
	fileService FileUploadServiceInterface
}

// NewFileUploadHandler 创建文件上传处理器
func NewFileUploadHandler(fileService FileUploadServiceInterface) *FileUploadHandler {
	return &FileUploadHandler{
		fileService: fileService,
	}
}

// Upload 上传文件
// POST /api/v1/upload
func (h *FileUploadHandler) Upload(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "file is required")
		return
	}

	// 获取文件类型（默认为 image）
	fileType := c.DefaultPostForm("file_type", "image")

	req := &service.UploadFileRequest{
		File:     file,
		FileType: fileType,
	}

	uploadedFile, err := h.fileService.UploadFile(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidFileType):
			response.Error(c, http.StatusBadRequest, "invalid file type")
		case errors.Is(err, service.ErrFileTooLarge):
			response.Error(c, http.StatusBadRequest, "file too large")
		case errors.Is(err, service.ErrFileUploadFailed):
			response.Error(c, http.StatusInternalServerError, "file upload failed")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to upload file")
		}
		return
	}

	response.Success(c, http.StatusCreated, "file uploaded successfully", uploadedFile)
}

// GetFile 获取文件元数据
// GET /api/v1/files/:id
func (h *FileUploadHandler) GetFile(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid file id")
		return
	}

	file, err := h.fileService.GetFile(c.Request.Context(), userID, fileID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrFileNotFound):
			response.Error(c, http.StatusNotFound, "file not found")
		case errors.Is(err, service.ErrFileNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get file")
		}
		return
	}

	response.Success(c, http.StatusOK, "file retrieved successfully", file)
}

// DeleteFile 删除文件
// DELETE /api/v1/files/:id
func (h *FileUploadHandler) DeleteFile(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid file id")
		return
	}

	err = h.fileService.DeleteFile(c.Request.Context(), userID, fileID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrFileNotFound):
			response.Error(c, http.StatusNotFound, "file not found")
		case errors.Is(err, service.ErrFileNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to delete file")
		}
		return
	}

	response.Success(c, http.StatusOK, "file deleted successfully", nil)
}

// ListFiles 获取文件列表
// GET /api/v1/files
func (h *FileUploadHandler) ListFiles(c *gin.Context) {
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

	fileType := c.Query("file_type")
	sort := c.Query("sort")

	result, err := h.fileService.ListFiles(c.Request.Context(), userID, page, pageSize, fileType, sort)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to list files")
		return
	}

	response.Success(c, http.StatusOK, "files retrieved successfully", result)
}

// GetStorageStats 获取存储统计
// GET /api/v1/files/stats
func (h *FileUploadHandler) GetStorageStats(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	stats, err := h.fileService.GetUserStorageStats(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get storage stats")
		return
	}

	response.Success(c, http.StatusOK, "storage stats retrieved successfully", stats)
}

// DownloadFile 下载文件
// GET /api/v1/files/:id/download
func (h *FileUploadHandler) DownloadFile(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid file id")
		return
	}

	file, err := h.fileService.GetFile(c.Request.Context(), userID, fileID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrFileNotFound):
			response.Error(c, http.StatusNotFound, "file not found")
		case errors.Is(err, service.ErrFileNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get file")
		}
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+file.FileName)
	c.Header("Content-Type", file.MimeType)

	// 发送文件
	c.File(file.StoragePath)
}

// GetThumbnail 获取缩略图
// GET /api/v1/files/:id/thumbnail
func (h *FileUploadHandler) GetThumbnail(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid file id")
		return
	}

	file, err := h.fileService.GetFile(c.Request.Context(), userID, fileID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrFileNotFound):
			response.Error(c, http.StatusNotFound, "file not found")
		case errors.Is(err, service.ErrFileNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get file")
		}
		return
	}

	// 检查是否有缩略图
	if file.ThumbnailURL == nil {
		response.Error(c, http.StatusNotFound, "thumbnail not found")
		return
	}

	// 设置响应头
	c.Header("Content-Type", file.MimeType)
	c.Header("Cache-Control", "public, max-age=31536000")

	// 发送缩略图
	c.File(*file.ThumbnailURL)
}
