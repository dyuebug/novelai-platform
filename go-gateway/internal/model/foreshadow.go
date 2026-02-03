package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ForeshadowStatus 伏笔状态
type ForeshadowStatus string

const (
	ForeshadowStatusPlanted   ForeshadowStatus = "planted"   // 已埋下
	ForeshadowStatusHinted    ForeshadowStatus = "hinted"    // 已暗示
	ForeshadowStatusResolved  ForeshadowStatus = "resolved"  // 已回收
	ForeshadowStatusAbandoned ForeshadowStatus = "abandoned" // 已放弃
)

// ForeshadowPriority 伏笔优先级
type ForeshadowPriority string

const (
	ForeshadowPriorityHigh   ForeshadowPriority = "high"   // 高优先级 - 主线伏笔
	ForeshadowPriorityMedium ForeshadowPriority = "medium" // 中优先级 - 支线伏笔
	ForeshadowPriorityLow    ForeshadowPriority = "low"    // 低优先级 - 细节伏笔
)

// Foreshadow 伏笔
type Foreshadow struct {
	ID          string             `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID   string             `gorm:"type:uuid;not null;index" json:"project_id"`
	Title       string             `gorm:"size:200;not null" json:"title"`
	Description string             `gorm:"type:text" json:"description"`
	Content     string             `gorm:"type:text" json:"content"`      // 伏笔内容/原文
	Priority    ForeshadowPriority `gorm:"size:20;default:medium" json:"priority"`
	Status      ForeshadowStatus   `gorm:"size:20;default:planted" json:"status"`

	// 埋设信息
	PlantChapterID  *string `gorm:"type:uuid" json:"plant_chapter_id"`
	PlantChapterNum int     `gorm:"default:0" json:"plant_chapter_num"`
	PlantPosition   string  `gorm:"size:500" json:"plant_position"` // 埋设位置描述

	// 回收信息
	ResolveChapterID  *string    `gorm:"type:uuid" json:"resolve_chapter_id"`
	ResolveChapterNum int        `gorm:"default:0" json:"resolve_chapter_num"`
	ResolveContent    string     `gorm:"type:text" json:"resolve_content"` // 回收内容
	ResolvedAt        *time.Time `json:"resolved_at"`

	// 提醒设置
	RemindChapterNum int  `gorm:"default:0" json:"remind_chapter_num"` // 提醒章节号
	RemindEnabled    bool `gorm:"default:true" json:"remind_enabled"`

	// 关联
	RelatedCharacterIDs []string `gorm:"type:text;serializer:json" json:"related_character_ids"`
	Tags                []string `gorm:"type:text;serializer:json" json:"tags"`

	// 元数据
	Metadata  map[string]interface{} `gorm:"type:text;serializer:json" json:"metadata"`
	SortOrder int                    `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// BeforeCreate 创建前钩子
func (f *Foreshadow) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}

// ForeshadowHint 伏笔暗示记录
type ForeshadowHint struct {
	ID            string    `gorm:"type:uuid;primaryKey" json:"id"`
	ForeshadowID  string    `gorm:"type:uuid;not null;index" json:"foreshadow_id"`
	ChapterID     string    `gorm:"type:uuid;not null" json:"chapter_id"`
	ChapterNum    int       `gorm:"default:0" json:"chapter_num"`
	Content       string    `gorm:"type:text" json:"content"` // 暗示内容
	HintType      string    `gorm:"size:50" json:"hint_type"` // direct, indirect, symbolic
	CreatedAt     time.Time `json:"created_at"`

	// 关联
	Foreshadow *Foreshadow `gorm:"foreignKey:ForeshadowID" json:"foreshadow,omitempty"`
}

// BeforeCreate 创建前钩子
func (h *ForeshadowHint) BeforeCreate(tx *gorm.DB) error {
	if h.ID == "" {
		h.ID = uuid.New().String()
	}
	return nil
}

// ForeshadowReminder 伏笔提醒
type ForeshadowReminder struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	ForeshadowID string    `gorm:"type:uuid;not null;index" json:"foreshadow_id"`
	ChapterNum   int       `gorm:"not null" json:"chapter_num"`
	Message      string    `gorm:"type:text" json:"message"`
	IsRead       bool      `gorm:"default:false" json:"is_read"`
	CreatedAt    time.Time `json:"created_at"`

	// 关联
	Foreshadow *Foreshadow `gorm:"foreignKey:ForeshadowID" json:"foreshadow,omitempty"`
}

// BeforeCreate 创建前钩子
func (r *ForeshadowReminder) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

// ForeshadowStats 伏笔统计
type ForeshadowStats struct {
	Total       int `json:"total"`
	Planted     int `json:"planted"`
	Hinted      int `json:"hinted"`
	Resolved    int `json:"resolved"`
	Abandoned   int `json:"abandoned"`
	HighPriority   int `json:"high_priority"`
	MediumPriority int `json:"medium_priority"`
	LowPriority    int `json:"low_priority"`
	OverdueCount   int `json:"overdue_count"` // 超期未回收数量
}
