package model

import (
	"time"

	"github.com/google/uuid"
)

// Character 角色模型
type Character struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"project_id"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	Alias       string     `gorm:"size:200" json:"alias,omitempty"`
	Gender      string     `gorm:"size:20" json:"gender,omitempty"`
	Age         string     `gorm:"size:50" json:"age,omitempty"`
	Birthday    string     `gorm:"size:50" json:"birthday,omitempty"`
	Appearance  string     `gorm:"type:text" json:"appearance,omitempty"`
	Personality string     `gorm:"type:text" json:"personality,omitempty"`
	Background  string     `gorm:"type:text" json:"background,omitempty"`
	Abilities   string     `gorm:"type:text" json:"abilities,omitempty"`
	Goals       string     `gorm:"type:text" json:"goals,omitempty"`
	Role        string     `gorm:"size:50" json:"role,omitempty"` // protagonist, antagonist, supporting, minor
	Status      string     `gorm:"size:20;default:active" json:"status"`
	AvatarURL   string     `gorm:"size:500" json:"avatar_url,omitempty"`
	Metadata    JSON       `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Project      Project              `gorm:"foreignKey:ProjectID" json:"-"`
	Experiences  []CharacterExperience `gorm:"foreignKey:CharacterID" json:"experiences,omitempty"`
	Relationships []CharacterRelationship `gorm:"foreignKey:CharacterID" json:"relationships,omitempty"`
}

// CharacterRelationship 角色关系模型
type CharacterRelationship struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CharacterID     uuid.UUID `gorm:"type:uuid;index;not null" json:"character_id"`
	TargetID        uuid.UUID `gorm:"type:uuid;index;not null" json:"target_id"`
	RelationType    string    `gorm:"size:50;not null" json:"relation_type"` // family, friend, enemy, lover, mentor, etc.
	Description     string    `gorm:"type:text" json:"description,omitempty"`
	StartChapter    *int      `json:"start_chapter,omitempty"`
	EndChapter      *int      `json:"end_chapter,omitempty"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Character Character `gorm:"foreignKey:CharacterID" json:"-"`
	Target    Character `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// CharacterExperience 角色经历模型
type CharacterExperience struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CharacterID uuid.UUID `gorm:"type:uuid;index;not null" json:"character_id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Content     string    `gorm:"type:text" json:"content,omitempty"`
	ChapterRef  *int      `json:"chapter_ref,omitempty"`
	TimePoint   string    `gorm:"size:100" json:"time_point,omitempty"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Character Character `gorm:"foreignKey:CharacterID" json:"-"`
}

// Location 地点模型
type Location struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"project_id"`
	ParentID    *uuid.UUID `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	Type        string     `gorm:"size:50" json:"type,omitempty"` // world, continent, country, city, building, room
	Climate     string     `gorm:"size:100" json:"climate,omitempty"`
	Features    string     `gorm:"type:text" json:"features,omitempty"`
	ImageURL    string     `gorm:"size:500" json:"image_url,omitempty"`
	Metadata    JSON       `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Project  Project    `gorm:"foreignKey:ProjectID" json:"-"`
	Parent   *Location  `gorm:"foreignKey:ParentID" json:"-"`
	Children []Location `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// Organization 组织模型
type Organization struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID   uuid.UUID `gorm:"type:uuid;index;not null" json:"project_id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Type        string    `gorm:"size:50" json:"type,omitempty"` // sect, guild, kingdom, company, etc.
	Description string    `gorm:"type:text" json:"description,omitempty"`
	History     string    `gorm:"type:text" json:"history,omitempty"`
	Structure   string    `gorm:"type:text" json:"structure,omitempty"`
	Goals       string    `gorm:"type:text" json:"goals,omitempty"`
	LocationID  *uuid.UUID `gorm:"type:uuid;index" json:"location_id,omitempty"`
	LogoURL     string    `gorm:"size:500" json:"logo_url,omitempty"`
	Metadata    JSON      `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Project  Project               `gorm:"foreignKey:ProjectID" json:"-"`
	Location *Location             `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Members  []OrganizationMember  `gorm:"foreignKey:OrganizationID" json:"members,omitempty"`
}

// OrganizationMember 组织成员模型
type OrganizationMember struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;index;not null" json:"organization_id"`
	CharacterID    uuid.UUID `gorm:"type:uuid;index;not null" json:"character_id"`
	Position       string    `gorm:"size:100" json:"position,omitempty"`
	Rank           string    `gorm:"size:50" json:"rank,omitempty"`
	JoinChapter    *int      `json:"join_chapter,omitempty"`
	LeaveChapter   *int      `json:"leave_chapter,omitempty"`
	Status         string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"-"`
	Character    Character    `gorm:"foreignKey:CharacterID" json:"character,omitempty"`
}

// WorldSetting 世界设定模型
type WorldSetting struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"project_id"`
	ParentID    *uuid.UUID `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	Category    string     `gorm:"size:50;not null" json:"category"` // power_system, history, culture, technology, magic, etc.
	Title       string     `gorm:"size:200;not null" json:"title"`
	Content     string     `gorm:"type:text" json:"content,omitempty"`
	Metadata    JSON       `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Project  Project        `gorm:"foreignKey:ProjectID" json:"-"`
	Parent   *WorldSetting  `gorm:"foreignKey:ParentID" json:"-"`
	Children []WorldSetting `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName 指定表名
func (Character) TableName() string {
	return "characters"
}

func (CharacterRelationship) TableName() string {
	return "character_relationships"
}

func (CharacterExperience) TableName() string {
	return "character_experiences"
}

func (Location) TableName() string {
	return "locations"
}

func (Organization) TableName() string {
	return "organizations"
}

func (OrganizationMember) TableName() string {
	return "organization_members"
}

func (WorldSetting) TableName() string {
	return "world_settings"
}
