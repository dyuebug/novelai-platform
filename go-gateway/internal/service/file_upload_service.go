package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

var (
	ErrFileNotOwned       = errors.New("file not owned by user")
	ErrInvalidFileType    = errors.New("invalid file type")
	ErrFileTooLarge       = errors.New("file too large")
	ErrFileUploadFailed   = errors.New("file upload failed")
	ErrThumbnailFailed    = errors.New("thumbnail generation failed")
	ErrStoragePathInvalid = errors.New("storage path invalid")
)

const (
	MaxImageSize      = 5 * 1024 * 1024  // 5MB
	MaxDocumentSize   = 10 * 1024 * 1024 // 10MB
	ThumbnailWidth    = 200
	ThumbnailHeight   = 200
	ImageQuality      = 80
	UploadDir         = "uploads"
	ThumbnailDir      = "thumbnails"
)

// 允许的图片类型
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// FileRepository 文件仓库接口
type FileRepository interface {
	Create(ctx context.Context, file *model.File) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.File, error)
	FindByHash(ctx context.Context, hash string) (*model.File, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, opts *repository.ListOptions) ([]model.File, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	GetTotalSizeByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

// FileUploadService 文件上传服务
type FileUploadService struct {
	fileRepo FileRepository
}

// NewFileUploadService 创建文件上传服务
func NewFileUploadService() *FileUploadService {
	return &FileUploadService{
		fileRepo: repository.NewFileRepository(),
	}
}

// UploadFileRequest 上传文件请求
type UploadFileRequest struct {
	File     *multipart.FileHeader
	FileType string // image, document, etc.
}

// UploadFile 上传文件
func (s *FileUploadService) UploadFile(ctx context.Context, userID uuid.UUID, req *UploadFileRequest) (*model.File, error) {
	// 打开上传的文件
	src, err := req.File.Open()
	if err != nil {
		return nil, ErrFileUploadFailed
	}
	defer src.Close()

	// 读取文件内容用于计算哈希
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, ErrFileUploadFailed
	}

	// 计算文件哈希
	hash := sha256.Sum256(fileBytes)
	hashStr := hex.EncodeToString(hash[:])

	// 检查文件是否已存在（去重）
	existingFile, err := s.fileRepo.FindByHash(ctx, hashStr)
	if err == nil && existingFile != nil {
		// 文件已存在，直接返回
		return existingFile, nil
	}

	// 验证文件类型和大小
	mimeType := req.File.Header.Get("Content-Type")
	fileSize := req.File.Size

	if req.FileType == "image" {
		if !allowedImageTypes[mimeType] {
			return nil, ErrInvalidFileType
		}
		if fileSize > MaxImageSize {
			return nil, ErrFileTooLarge
		}
	}

	// 生成存储路径
	ext := filepath.Ext(req.File.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	datePath := time.Now().Format("2006/01/02")
	storagePath := filepath.Join(UploadDir, datePath, fileName)

	// 确保目录存在
	dirPath := filepath.Dir(storagePath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, ErrFileUploadFailed
	}

	// 保存文件
	dst, err := os.Create(storagePath)
	if err != nil {
		return nil, ErrFileUploadFailed
	}
	defer dst.Close()

	if _, err := dst.Write(fileBytes); err != nil {
		return nil, ErrFileUploadFailed
	}

	// 创建文件记录
	file := &model.File{
		UserID:      userID,
		FileName:    req.File.Filename,
		StoragePath: storagePath,
		FileSize:    fileSize,
		MimeType:    mimeType,
		FileType:    req.FileType,
		Hash:        hashStr,
		Metadata:    make(map[string]interface{}),
	}

	// 如果是图片，处理图片信息和缩略图
	if req.FileType == "image" {
		if err := s.processImage(file, fileBytes); err != nil {
			// 图片处理失败不影响上传，只记录错误
			file.Metadata = map[string]interface{}{
				"thumbnail_error": err.Error(),
			}
		}
	}

	// 保存文件记录到数据库
	if err := s.fileRepo.Create(ctx, file); err != nil {
		// 删除已上传的文件
		os.Remove(storagePath)
		return nil, err
	}

	return file, nil
}

// processImage 处理图片（获取尺寸、生成缩略图）
func (s *FileUploadService) processImage(file *model.File, fileBytes []byte) error {
	// 解码图片
	img, format, err := image.Decode(strings.NewReader(string(fileBytes)))
	if err != nil {
		return err
	}

	// 获取图片尺寸
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	file.Width = &width
	file.Height = &height

	// 生成缩略图
	thumbnail := resize.Thumbnail(ThumbnailWidth, ThumbnailHeight, img, resize.Lanczos3)

	// 生成缩略图路径
	ext := filepath.Ext(file.FileName)
	thumbnailFileName := fmt.Sprintf("%s_thumb%s", uuid.New().String(), ext)
	datePath := time.Now().Format("2006/01/02")
	thumbnailPath := filepath.Join(ThumbnailDir, datePath, thumbnailFileName)

	// 确保缩略图目录存在
	thumbnailDirPath := filepath.Dir(thumbnailPath)
	if err := os.MkdirAll(thumbnailDirPath, 0755); err != nil {
		return ErrThumbnailFailed
	}

	// 保存缩略图
	thumbnailFile, err := os.Create(thumbnailPath)
	if err != nil {
		return ErrThumbnailFailed
	}
	defer thumbnailFile.Close()

	// 根据格式编码缩略图
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(thumbnailFile, thumbnail, &jpeg.Options{Quality: ImageQuality})
	case "png":
		err = png.Encode(thumbnailFile, thumbnail)
	default:
		// 默认使用 JPEG
		err = jpeg.Encode(thumbnailFile, thumbnail, &jpeg.Options{Quality: ImageQuality})
	}

	if err != nil {
		return ErrThumbnailFailed
	}

	file.ThumbnailURL = &thumbnailPath
	return nil
}

// GetFile 获取文件信息
func (s *FileUploadService) GetFile(ctx context.Context, userID, fileID uuid.UUID) (*model.File, error) {
	file, err := s.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	// 验证所有权
	if file.UserID != userID {
		return nil, ErrFileNotOwned
	}

	return file, nil
}

// DeleteFile 删除文件
func (s *FileUploadService) DeleteFile(ctx context.Context, userID, fileID uuid.UUID) error {
	file, err := s.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		return err
	}

	// 验证所有权
	if file.UserID != userID {
		return ErrFileNotOwned
	}

	// 删除物理文件
	if err := os.Remove(file.StoragePath); err != nil && !os.IsNotExist(err) {
		// 文件删除失败，记录但不阻止数据库删除
	}

	// 删除缩略图
	if file.ThumbnailURL != nil {
		if err := os.Remove(*file.ThumbnailURL); err != nil && !os.IsNotExist(err) {
			// 缩略图删除失败，记录但不阻止数据库删除
		}
	}

	// 删除数据库记录
	return s.fileRepo.Delete(ctx, fileID)
}

// ListFiles 获取用户文件列表
func (s *FileUploadService) ListFiles(ctx context.Context, userID uuid.UUID, page, pageSize int, fileType, sort string) (*FileListResponse, error) {
	// 限制 pageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page <= 0 {
		page = 1
	}

	opts := &repository.ListOptions{
		Page:     page,
		PageSize: pageSize,
		Status:   fileType,
		Sort:     sort,
	}

	files, total, err := s.fileRepo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &FileListResponse{
		Items:      files,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// FileListResponse 文件列表响应
type FileListResponse struct {
	Items      []model.File `json:"items"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}

// GetUserStorageStats 获取用户存储统计
func (s *FileUploadService) GetUserStorageStats(ctx context.Context, userID uuid.UUID) (*StorageStats, error) {
	count, err := s.fileRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalSize, err := s.fileRepo.GetTotalSizeByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &StorageStats{
		FileCount: count,
		TotalSize: totalSize,
	}, nil
}

// StorageStats 存储统计
type StorageStats struct {
	FileCount int64 `json:"file_count"`
	TotalSize int64 `json:"total_size"`
}
