package service

import (
	"context"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
)

// ========== 地点服务 ==========

type LocationService struct {
	locationRepo *repository.LocationRepository
	projectRepo  *repository.ProjectRepository
}

func NewLocationService() *LocationService {
	return &LocationService{
		locationRepo: repository.NewLocationRepository(),
		projectRepo:  repository.NewProjectRepository(),
	}
}

// CreateLocationRequest 创建地点请求
type CreateLocationRequest struct {
	ParentID    *uuid.UUID             `json:"parent_id"`
	Name        string                 `json:"name" binding:"required,min=1,max=100"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Climate     string                 `json:"climate"`
	Features    string                 `json:"features"`
	ImageURL    string                 `json:"image_url"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateLocationRequest 更新地点请求
type UpdateLocationRequest struct {
	ParentID    *uuid.UUID `json:"parent_id"`
	Name        *string    `json:"name" binding:"omitempty,min=1,max=100"`
	Description *string    `json:"description"`
	Type        *string    `json:"type"`
	Climate     *string    `json:"climate"`
	Features    *string    `json:"features"`
	ImageURL    *string    `json:"image_url"`
	SortOrder   *int       `json:"sort_order"`
}

// LocationListResponse 地点列表响应
type LocationListResponse struct {
	Items      []model.Location `json:"items"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// Create 创建地点
func (s *LocationService) Create(ctx context.Context, userID, projectID uuid.UUID, req *CreateLocationRequest) (*model.Location, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	location := &model.Location{
		ProjectID:   projectID,
		ParentID:    req.ParentID,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Climate:     req.Climate,
		Features:    req.Features,
		ImageURL:    req.ImageURL,
		Metadata:    model.JSON(req.Metadata),
	}

	if err := s.locationRepo.Create(ctx, location); err != nil {
		return nil, err
	}

	return location, nil
}

// Get 获取地点
func (s *LocationService) Get(ctx context.Context, userID, locationID uuid.UUID) (*model.Location, error) {
	location, err := s.locationRepo.FindByID(ctx, locationID)
	if err != nil {
		return nil, err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, location.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrLocationNotOwned
	}

	return location, nil
}

// List 获取地点列表
func (s *LocationService) List(ctx context.Context, userID, projectID uuid.UUID, page, pageSize int, locationType, sort string) (*LocationListResponse, error) {
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
		Status:   locationType,
		Sort:     sort,
	}

	locations, total, err := s.locationRepo.FindByProjectID(ctx, projectID, opts)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &LocationListResponse{
		Items:      locations,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetTree 获取地点树结构
func (s *LocationService) GetTree(ctx context.Context, userID, projectID uuid.UUID) ([]model.Location, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	return s.locationRepo.FindRootsByProjectID(ctx, projectID)
}

// Update 更新地点
func (s *LocationService) Update(ctx context.Context, userID, locationID uuid.UUID, req *UpdateLocationRequest) (*model.Location, error) {
	location, err := s.locationRepo.FindByID(ctx, locationID)
	if err != nil {
		return nil, err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, location.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrLocationNotOwned
	}

	// 更新字段
	if req.ParentID != nil {
		location.ParentID = req.ParentID
	}
	if req.Name != nil {
		location.Name = *req.Name
	}
	if req.Description != nil {
		location.Description = *req.Description
	}
	if req.Type != nil {
		location.Type = *req.Type
	}
	if req.Climate != nil {
		location.Climate = *req.Climate
	}
	if req.Features != nil {
		location.Features = *req.Features
	}
	if req.ImageURL != nil {
		location.ImageURL = *req.ImageURL
	}
	if req.SortOrder != nil {
		location.SortOrder = *req.SortOrder
	}

	if err := s.locationRepo.Update(ctx, location); err != nil {
		return nil, err
	}

	return location, nil
}

// Delete 删除地点
func (s *LocationService) Delete(ctx context.Context, userID, locationID uuid.UUID) error {
	location, err := s.locationRepo.FindByID(ctx, locationID)
	if err != nil {
		return err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, location.ProjectID)
	if err != nil {
		return err
	}
	if project.UserID != userID {
		return ErrLocationNotOwned
	}

	return s.locationRepo.Delete(ctx, locationID)
}

// ========== 组织服务 ==========

type OrganizationService struct {
	orgRepo    *repository.OrganizationRepository
	memberRepo *repository.OrganizationMemberRepository
	projectRepo *repository.ProjectRepository
}

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{
		orgRepo:    repository.NewOrganizationRepository(),
		memberRepo: repository.NewOrganizationMemberRepository(),
		projectRepo: repository.NewProjectRepository(),
	}
}

// CreateOrganizationRequest 创建组织请求
type CreateOrganizationRequest struct {
	Name        string                 `json:"name" binding:"required,min=1,max=100"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	History     string                 `json:"history"`
	Structure   string                 `json:"structure"`
	Goals       string                 `json:"goals"`
	LocationID  *uuid.UUID             `json:"location_id"`
	LogoURL     string                 `json:"logo_url"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateOrganizationRequest 更新组织请求
type UpdateOrganizationRequest struct {
	Name        *string    `json:"name" binding:"omitempty,min=1,max=100"`
	Type        *string    `json:"type"`
	Description *string    `json:"description"`
	History     *string    `json:"history"`
	Structure   *string    `json:"structure"`
	Goals       *string    `json:"goals"`
	LocationID  *uuid.UUID `json:"location_id"`
	LogoURL     *string    `json:"logo_url"`
	SortOrder   *int       `json:"sort_order"`
}

// OrganizationListResponse 组织列表响应
type OrganizationListResponse struct {
	Items      []model.Organization `json:"items"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// Create 创建组织
func (s *OrganizationService) Create(ctx context.Context, userID, projectID uuid.UUID, req *CreateOrganizationRequest) (*model.Organization, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	org := &model.Organization{
		ProjectID:   projectID,
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		History:     req.History,
		Structure:   req.Structure,
		Goals:       req.Goals,
		LocationID:  req.LocationID,
		LogoURL:     req.LogoURL,
		Metadata:    model.JSON(req.Metadata),
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

// Get 获取组织
func (s *OrganizationService) Get(ctx context.Context, userID, orgID uuid.UUID) (*model.Organization, error) {
	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, org.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrOrganizationNotOwned
	}

	return org, nil
}

// List 获取组织列表
func (s *OrganizationService) List(ctx context.Context, userID, projectID uuid.UUID, page, pageSize int, orgType, sort string) (*OrganizationListResponse, error) {
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
		Status:   orgType,
		Sort:     sort,
	}

	orgs, total, err := s.orgRepo.FindByProjectID(ctx, projectID, opts)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &OrganizationListResponse{
		Items:      orgs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update 更新组织
func (s *OrganizationService) Update(ctx context.Context, userID, orgID uuid.UUID, req *UpdateOrganizationRequest) (*model.Organization, error) {
	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, org.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrOrganizationNotOwned
	}

	// 更新字段
	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.Type != nil {
		org.Type = *req.Type
	}
	if req.Description != nil {
		org.Description = *req.Description
	}
	if req.History != nil {
		org.History = *req.History
	}
	if req.Structure != nil {
		org.Structure = *req.Structure
	}
	if req.Goals != nil {
		org.Goals = *req.Goals
	}
	if req.LocationID != nil {
		org.LocationID = req.LocationID
	}
	if req.LogoURL != nil {
		org.LogoURL = *req.LogoURL
	}
	if req.SortOrder != nil {
		org.SortOrder = *req.SortOrder
	}

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

// Delete 删除组织
func (s *OrganizationService) Delete(ctx context.Context, userID, orgID uuid.UUID) error {
	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return err
	}

	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, org.ProjectID)
	if err != nil {
		return err
	}
	if project.UserID != userID {
		return ErrOrganizationNotOwned
	}

	return s.orgRepo.Delete(ctx, orgID)
}

// --- 组织成员 ---

// AddMemberRequest 添加成员请求
type AddMemberRequest struct {
	CharacterID  uuid.UUID `json:"character_id" binding:"required"`
	Position     string    `json:"position"`
	Rank         string    `json:"rank"`
	JoinChapter  *int      `json:"join_chapter"`
	LeaveChapter *int      `json:"leave_chapter"`
}

// AddMember 添加组织成员
func (s *OrganizationService) AddMember(ctx context.Context, userID, orgID uuid.UUID, req *AddMemberRequest) (*model.OrganizationMember, error) {
	// 验证组织所有权
	_, err := s.Get(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}

	member := &model.OrganizationMember{
		OrganizationID: orgID,
		CharacterID:    req.CharacterID,
		Position:       req.Position,
		Rank:           req.Rank,
		JoinChapter:    req.JoinChapter,
		LeaveChapter:   req.LeaveChapter,
		Status:         "active",
	}

	if err := s.memberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}

// GetMembers 获取组织成员列表
func (s *OrganizationService) GetMembers(ctx context.Context, userID, orgID uuid.UUID) ([]model.OrganizationMember, error) {
	// 验证组织所有权
	_, err := s.Get(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}

	return s.memberRepo.FindByOrganizationID(ctx, orgID)
}

// RemoveMember 移除组织成员
func (s *OrganizationService) RemoveMember(ctx context.Context, userID, memberID uuid.UUID) error {
	member, err := s.memberRepo.FindByID(ctx, memberID)
	if err != nil {
		return err
	}

	// 验证组织所有权
	_, err = s.Get(ctx, userID, member.OrganizationID)
	if err != nil {
		return err
	}

	return s.memberRepo.Delete(ctx, memberID)
}
