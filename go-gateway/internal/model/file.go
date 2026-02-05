package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// File 文件模型
type File struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index;column:user_id" json:"user_id"`
	FileName     string         `gorm:"size:255;not null;column:file_name" json:"file_name"`           // 原始文件名
	StoragePath  string         `gorm:"size:500;not null;column:storage_path" json:"storage_path"`     // 存储路径
	FileSize     int64          `gorm:"not null;column:file_size" json:"file_size"`                    // 文件大小（字节）
	MimeType     string         `gorm:"size:100;not null;column:mime_type" json:"mime_type"`           // MIME类型
	FileType     string         `gorm:"size:50;not null;index;column:file_type" json:"file_type"`      // 文件类型（image, document, etc.）
	Width        *int           `gorm:"column:width" json:"width,omitempty"`                           // 图片宽度
	Height       *int           `gorm:"column:height" json:"height,omitempty"`                         // 图片高度
	ThumbnailURL *string        `gorm:"size:500;column:thumbnail_url" json:"thumbnail_url,omitempty"`  // 缩略图URL
	Hash         string         `gorm:"size:64;index;column:hash" json:"hash"`                         // 文件哈希（用于去重）
	Metadata     JSON           `gorm:"type:jsonb;column:metadata" json:"metadata"`                    // 额外元数据
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (File) TableName() string {
	return "files"
}

// BeforeCreate GORM钩子：创建前设置UUID
func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// IsImage 判断是否为图片
func (f *File) IsImage() bool {
	return f.FileType == "image"
}

// GetURL 获取文件访问URL
func (f *File) GetURL(baseURL string) string {
	return baseURL + "/api/v1/files/" + f.ID.String() + "/download"
}
