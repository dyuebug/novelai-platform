package service

import (
	"context"
	"errors"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrCharacterNotOwned    = errors.New("character not owned by user")
	ErrLocationNotOwned     = errors.New("location not owned by user")
	ErrOrganizationNotOwned = errors.New("organization not owned by user")
	ErrWorldSettingNotOwned = errors.New("world setting not owned by user")
)

// ========== 角色服务 ==========

type CharacterService struct {
	characterRepo    *repository.CharacterRepository
	relationshipRepo *repository.CharacterRelationshipRepository
	experienceRepo   *repository.CharacterExperienceRepository
	projectRepo      *repository.ProjectRepository
}

func NewCharacterService() *CharacterService {
	return &CharacterService{
		characterRepo:    repository.NewCharacterRepository(),
		relationshipRepo: repository.NewCharacterRelationshipRepository(),
		experienceRepo:   repository.NewCharacterExperienceRepository(),
		projectRepo:      repository.NewProjectRepository(),
	}
}

// CreateCharacterRequest 创建角色请求
type CreateCharacterRequest struct {
	Name        string                 `json:"name" binding:"required,min=1,max=100"`
	Alias       string                 `json:"alias"`
	Gender      string                 `json:"gender"`
	Age         string                 `json:"age"`
	Birthday    string                 `json:"birthday"`
	Appearance  string                 `json:"appearance"`
	Personality string                 `json:"personality"`
	Background  string                 `json:"background"`
	Abilities   string                 `json:"abilities"`
	Goals       string                 `json:"goals"`
	Role        string                 `json:"role"`
	AvatarURL   string                 `json:"avatar_url"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateCharacterRequest 更新角色请求
type UpdateCharacterRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=100"`
	Alias       *string `json:"alias"`
	Gender      *string `json:"gender"`
	Age         *string `json:"age"`
	Birthday    *string `json:"birthday"`
	Appearance  *string `json:"appearance"`
	Personality *string `json:"personality"`
	Background  *string `json:"background"`
	Abilities   *string `json:"abilities"`
	Goals       *string `json:"goals"`
	Role        *string `json:"role"`
	Status      *string `json:"status"`
	AvatarURL   *string `json:"avatar_url"`
	SortOrder   *int    `json:"sort_order"`
}

// CharacterListResponse 角色列表响应
type CharacterListResponse struct {
	Items      []model.Character `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// Create 创建角色
func (s *CharacterService) Create(ctx context.Context, userID, projectID uuid.UUID, req *CreateCharacterRequest) (*model.Character, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	character := &model.Character{
		ProjectID:   projectID,
		Name:        req.Name,
		Alias:       req.Alias,
		Gender:      req.Gender,
		Age:         req.Age,
		Birthday:    req.Birthday,
		Appearance:  req.Appearance,
		Personality: req.Personality,
		Background:  req.Background,
		Abilities:   req.Abilities,
		Goals:       req.Goals,
		Role:        req.Role,
		Status:      "active",
		AvatarURL:   req.AvatarURL,
		Metadata:    model.JSON(req.Metadata),
	}

	if err := s.characterRepo.Create(ctx, character); err != nil {
		return nil, err
	}

	return character, nil
}

// Get 获取角色
func (s *CharacterService) Get(ctx context.Context, userID, characterID uuid.UUID) (*model.Character, error) {
	character, err := s.characterRepo.FindByID(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, character.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrCharacterNotOwned
	}

	return character, nil
}

// List 获取角色列表
func (s *CharacterService) List(ctx context.Context, userID, projectID uuid.UUID, page, pageSize int, role, sort string) (*CharacterListResponse, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}
	if page <= 0 {
		page = 1
	}

	opts := &repository.ListOptions{
		Page:     page,
		PageSize: pageSize,
		Status:   role,
		Sort:     sort,
	}

	characters, total, err := s.characterRepo.FindByProjectID(ctx, projectID, opts)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &CharacterListResponse{
		Items:      characters,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update 更新角色
func (s *CharacterService) Update(ctx context.Context, userID, characterID uuid.UUID, req *UpdateCharacterRequest) (*model.Character, error) {
	character, err := s.characterRepo.FindByID(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, character.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrCharacterNotOwned
	}

	// 更新字段
	if req.Name != nil {
		character.Name = *req.Name
	}
	if req.Alias != nil {
		character.Alias = *req.Alias
	}
	if req.Gender != nil {
		character.Gender = *req.Gender
	}
	if req.Age != nil {
		character.Age = *req.Age
	}
	if req.Birthday != nil {
		character.Birthday = *req.Birthday
	}
	if req.Appearance != nil {
		character.Appearance = *req.Appearance
	}
	if req.Personality != nil {
		character.Personality = *req.Personality
	}
	if req.Background != nil {
		character.Background = *req.Background
	}
	if req.Abilities != nil {
		character.Abilities = *req.Abilities
	}
	if req.Goals != nil {
		character.Goals = *req.Goals
	}
	if req.Role != nil {
		character.Role = *req.Role
	}
	if req.Status != nil {
		character.Status = *req.Status
	}
	if req.AvatarURL != nil {
		character.AvatarURL = *req.AvatarURL
	}
	if req.SortOrder != nil {
		character.SortOrder = *req.SortOrder
	}

	if err := s.characterRepo.Update(ctx, character); err != nil {
		return nil, err
	}

	return character, nil
}

// Delete 删除角色
func (s *CharacterService) Delete(ctx context.Context, userID, characterID uuid.UUID) error {
	character, err := s.characterRepo.FindByID(ctx, characterID)
	if err != nil {
		return err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, character.ProjectID)
	if err != nil {
		return err
	}
	if project.UserID != userID {
		return ErrCharacterNotOwned
	}

	return s.characterRepo.Delete(ctx, characterID)
}

// --- 角色关系 ---

// CreateRelationshipRequest 创建关系请求
type CreateRelationshipRequest struct {
	TargetID     uuid.UUID `json:"target_id" binding:"required"`
	RelationType string    `json:"relation_type" binding:"required"`
	Description  string    `json:"description"`
	StartChapter *int      `json:"start_chapter"`
	EndChapter   *int      `json:"end_chapter"`
}

// CreateRelationship 创建角色关系
func (s *CharacterService) CreateRelationship(ctx context.Context, userID, characterID uuid.UUID, req *CreateRelationshipRequest) (*model.CharacterRelationship, error) {
	// 验证角色所有权
	_, err := s.Get(ctx, userID, characterID)
	if err != nil {
		return nil, err
	}

	rel := &model.CharacterRelationship{
		CharacterID:  characterID,
		TargetID:     req.TargetID,
		RelationType: req.RelationType,
		Description:  req.Description,
		StartChapter: req.StartChapter,
		EndChapter:   req.EndChapter,
	}

	if err := s.relationshipRepo.Create(ctx, rel); err != nil {
		return nil, err
	}

	return rel, nil
}

// GetRelationships 获取角色关系列表
func (s *CharacterService) GetRelationships(ctx context.Context, userID, characterID uuid.UUID) ([]model.CharacterRelationship, error) {
	// 验证角色所有权
	_, err := s.Get(ctx, userID, characterID)
	if err != nil {
		return nil, err
	}

	return s.relationshipRepo.FindByCharacterID(ctx, characterID)
}

// DeleteRelationship 删除角色关系
func (s *CharacterService) DeleteRelationship(ctx context.Context, userID, relationshipID uuid.UUID) error {
	rel, err := s.relationshipRepo.FindByID(ctx, relationshipID)
	if err != nil {
		return err
	}

	// 验证角色所有权
	_, err = s.Get(ctx, userID, rel.CharacterID)
	if err != nil {
		return err
	}

	return s.relationshipRepo.Delete(ctx, relationshipID)
}

// --- 角色经历 ---

// CreateExperienceRequest 创建经历请求
type CreateExperienceRequest struct {
	Title      string `json:"title" binding:"required,min=1,max=200"`
	Content    string `json:"content"`
	ChapterRef *int   `json:"chapter_ref"`
	TimePoint  string `json:"time_point"`
	SortOrder  int    `json:"sort_order"`
}

// CreateExperience 创建角色经历
func (s *CharacterService) CreateExperience(ctx context.Context, userID, characterID uuid.UUID, req *CreateExperienceRequest) (*model.CharacterExperience, error) {
	// 验证角色所有权
	_, err := s.Get(ctx, userID, characterID)
	if err != nil {
		return nil, err
	}

	exp := &model.CharacterExperience{
		CharacterID: characterID,
		Title:       req.Title,
		Content:     req.Content,
		ChapterRef:  req.ChapterRef,
		TimePoint:   req.TimePoint,
		SortOrder:   req.SortOrder,
	}

	if err := s.experienceRepo.Create(ctx, exp); err != nil {
		return nil, err
	}

	return exp, nil
}

// GetExperiences 获取角色经历列表
func (s *CharacterService) GetExperiences(ctx context.Context, userID, characterID uuid.UUID) ([]model.CharacterExperience, error) {
	// 验证角色所有权
	_, err := s.Get(ctx, userID, characterID)
	if err != nil {
		return nil, err
	}

	return s.experienceRepo.FindByCharacterID(ctx, characterID)
}

// DeleteExperience 删除角色经历
func (s *CharacterService) DeleteExperience(ctx context.Context, userID, experienceID uuid.UUID) error {
	exp, err := s.experienceRepo.FindByID(ctx, experienceID)
	if err != nil {
		return err
	}

	// 验证角色所有权
	_, err = s.Get(ctx, userID, exp.CharacterID)
	if err != nil {
		return err
	}

	return s.experienceRepo.Delete(ctx, experienceID)
}
