package model

import (
	"time"

	"github.com/google/uuid"
)

// Project 项目模型
type Project struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;index;not null" json:"user_id"`
	Title         string     `gorm:"size:200;not null" json:"title"`
	Description   string     `gorm:"type:text" json:"description,omitempty"`
	Genre         string     `gorm:"size:50" json:"genre,omitempty"`
	Status        string     `gorm:"size:20;default:draft" json:"status"`
	TotalChapters int        `gorm:"default:0" json:"total_chapters"`
	TotalWords    int        `gorm:"default:0" json:"total_words"`
	CoverURL      string     `gorm:"size:500" json:"cover_url,omitempty"`
	Metadata      JSON       `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	IsDeleted     bool       `gorm:"default:false" json:"is_deleted"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	User     User      `gorm:"foreignKey:UserID" json:"-"`
	Chapters []Chapter `gorm:"foreignKey:ProjectID" json:"chapters,omitempty"`
	Outlines []Outline `gorm:"foreignKey:ProjectID" json:"outlines,omitempty"`
}

// Outline 大纲模型
type Outline struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID    uuid.UUID  `gorm:"type:uuid;index;not null" json:"project_id"`
	Title        string     `gorm:"size:200" json:"title,omitempty"`
	Content      string     `gorm:"type:text;not null" json:"content"`
	OutlineType  string     `gorm:"size:50;default:main" json:"outline_type"`
	ParentID     *uuid.UUID `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	SortOrder    int        `gorm:"default:0" json:"sort_order"`
	StartChapter *int       `json:"start_chapter,omitempty"`
	EndChapter   *int       `json:"end_chapter,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Project  Project   `gorm:"foreignKey:ProjectID" json:"-"`
	Parent   *Outline  `gorm:"foreignKey:ParentID" json:"-"`
	Children []Outline `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// Chapter 章节模型
type Chapter struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID       uuid.UUID  `gorm:"type:uuid;index;not null" json:"project_id"`
	ChapterNumber   int        `gorm:"not null" json:"chapter_number"`
	Title           string     `gorm:"size:200;not null" json:"title"`
	Content         string     `gorm:"type:text" json:"content,omitempty"`
	Summary         string     `gorm:"type:text" json:"summary,omitempty"`
	Status          string     `gorm:"size:20;default:draft" json:"status"`
	WordCount       int        `gorm:"default:0" json:"word_count"`
	POVCharacter    string     `gorm:"size:100" json:"pov_character,omitempty"`
	Location        string     `gorm:"size:200" json:"location,omitempty"`
	TimeSetting     string     `gorm:"size:200" json:"time_setting,omitempty"`
	HookType        string     `gorm:"size:50" json:"hook_type,omitempty"`
	HookStrength    string     `gorm:"size:20" json:"hook_strength,omitempty"`
	CoolpointPatterns JSON     `gorm:"type:jsonb;default:'[]'" json:"coolpoint_patterns,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	PublishedAt     *time.Time `json:"published_at,omitempty"`

	// 关联
	Project  Project          `gorm:"foreignKey:ProjectID" json:"-"`
	Versions []ChapterVersion `gorm:"foreignKey:ChapterID" json:"versions,omitempty"`
}

// ChapterVersion 章节版本模型
type ChapterVersion struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ChapterID     uuid.UUID  `gorm:"type:uuid;index;not null" json:"chapter_id"`
	VersionNumber int        `gorm:"not null" json:"version_number"`
	Content       string     `gorm:"type:text;not null" json:"content"`
	WordCount     int        `gorm:"default:0" json:"word_count"`
	Source        string     `gorm:"size:50;default:manual" json:"source"`
	AIProvider    string     `gorm:"size:50" json:"ai_provider,omitempty"`
	AIModel       string     `gorm:"size:100" json:"ai_model,omitempty"`
	PromptUsed    string     `gorm:"type:text" json:"prompt_used,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy     *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`

	// 关联
	Chapter Chapter `gorm:"foreignKey:ChapterID" json:"-"`
}

// JSON 类型用于 JSONB 字段
type JSON map[string]interface{}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}

func (Outline) TableName() string {
	return "outlines"
}

func (Chapter) TableName() string {
	return "chapters"
}

func (ChapterVersion) TableName() string {
	return "chapter_versions"
}
