package repository

import (
	"context"
	"errors"

	"go-gateway/internal/database"
	"go-gateway/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrLocationNotFound     = errors.New("location not found")
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrMemberNotFound       = errors.New("member not found")
)

// --- 地点 ---

type LocationRepository struct{}

func NewLocationRepository() *LocationRepository {
	return &LocationRepository{}
}

// Create 创建地点
func (r *LocationRepository) Create(ctx context.Context, location *model.Location) error {
	return database.DB.WithContext(ctx).Create(location).Error
}

// FindByID 通过 ID 查找地点
func (r *LocationRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Location, error) {
	var location model.Location
	result := database.DB.WithContext(ctx).
		Where("id = ?", id).
		First(&location)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrLocationNotFound
		}
		return nil, result.Error
	}
	return &location, nil
}

// FindByProjectID 查找项目的所有地点
func (r *LocationRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts *ListOptions) ([]model.Location, int64, error) {
	var locations []model.Location
	var total int64

	query := database.DB.WithContext(ctx).
		Model(&model.Location{}).
		Where("project_id = ?", projectID)

	// 类型过滤
	if opts.Status != "" {
		query = query.Where("type = ?", opts.Status)
	}

	// 计算总数
	query.Count(&total)

	// 排序
	if opts.Sort != "" {
		query = query.Order(opts.Sort)
	} else {
		query = query.Order("sort_order ASC, created_at DESC")
	}

	// 分页
	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	result := query.Find(&locations)
	return locations, total, result.Error
}

// FindRootsByProjectID 查找项目的根地点（无父级）
func (r *LocationRepository) FindRootsByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Location, error) {
	var locations []model.Location
	result := database.DB.WithContext(ctx).
		Where("project_id = ? AND parent_id IS NULL", projectID).
		Order("sort_order ASC, created_at DESC").
		Find(&locations)
	return locations, result.Error
}

// FindChildren 查找子地点
func (r *LocationRepository) FindChildren(ctx context.Context, parentID uuid.UUID) ([]model.Location, error) {
	var locations []model.Location
	result := database.DB.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("sort_order ASC, created_at DESC").
		Find(&locations)
	return locations, result.Error
}

// FindWithChildren 查找地点及其子地点（递归）
func (r *LocationRepository) FindWithChildren(ctx context.Context, id uuid.UUID) (*model.Location, error) {
	var location model.Location
	result := database.DB.WithContext(ctx).
		Preload("Children").
		Where("id = ?", id).
		First(&location)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrLocationNotFound
		}
		return nil, result.Error
	}
	return &location, nil
}

// Update 更新地点
func (r *LocationRepository) Update(ctx context.Context, location *model.Location) error {
	return database.DB.WithContext(ctx).Save(location).Error
}

// Delete 删除地点
func (r *LocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Location{}).Error
}

// Search 搜索地点
func (r *LocationRepository) Search(ctx context.Context, projectID uuid.UUID, query string, page, pageSize int) ([]model.Location, int64, error) {
	var locations []model.Location
	var total int64

	// 使用 ILIKE 进行模糊搜索
	searchPattern := "%" + query + "%"

	queryDB := database.DB.WithContext(ctx).
		Model(&model.Location{}).
		Where("project_id = ?", projectID).
		Where("name ILIKE ? OR description ILIKE ? OR features ILIKE ?", searchPattern, searchPattern, searchPattern)

	// 计算总数
	queryDB.Count(&total)

	// 分页
	offset := (page - 1) * pageSize
	result := queryDB.
		Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&locations)

	return locations, total, result.Error
}

// --- 组织 ---

type OrganizationRepository struct{}

func NewOrganizationRepository() *OrganizationRepository {
	return &OrganizationRepository{}
}

// Create 创建组织
func (r *OrganizationRepository) Create(ctx context.Context, org *model.Organization) error {
	return database.DB.WithContext(ctx).Create(org).Error
}

// FindByID 通过 ID 查找组织
func (r *OrganizationRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	var org model.Organization
	result := database.DB.WithContext(ctx).
		Where("id = ?", id).
		First(&org)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrOrganizationNotFound
		}
		return nil, result.Error
	}
	return &org, nil
}

// FindByProjectID 查找项目的所有组织
func (r *OrganizationRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts *ListOptions) ([]model.Organization, int64, error) {
	var orgs []model.Organization
	var total int64

	query := database.DB.WithContext(ctx).
		Model(&model.Organization{}).
		Where("project_id = ?", projectID)

	// 类型过滤
	if opts.Status != "" {
		query = query.Where("type = ?", opts.Status)
	}

	// 计算总数
	query.Count(&total)

	// 排序
	if opts.Sort != "" {
		query = query.Order(opts.Sort)
	} else {
		query = query.Order("sort_order ASC, created_at DESC")
	}

	// 分页
	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	result := query.Find(&orgs)
	return orgs, total, result.Error
}

// Update 更新组织
func (r *OrganizationRepository) Update(ctx context.Context, org *model.Organization) error {
	return database.DB.WithContext(ctx).Save(org).Error
}

// Delete 删除组织
func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Organization{}).Error
}

// --- 组织成员 ---

type OrganizationMemberRepository struct{}

func NewOrganizationMemberRepository() *OrganizationMemberRepository {
	return &OrganizationMemberRepository{}
}

// Create 创建成员
func (r *OrganizationMemberRepository) Create(ctx context.Context, member *model.OrganizationMember) error {
	return database.DB.WithContext(ctx).Create(member).Error
}

// FindByID 通过 ID 查找成员
func (r *OrganizationMemberRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.OrganizationMember, error) {
	var member model.OrganizationMember
	result := database.DB.WithContext(ctx).
		Preload("Character").
		Where("id = ?", id).
		First(&member)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, result.Error
	}
	return &member, nil
}

// FindByOrganizationID 查找组织的所有成员
func (r *OrganizationMemberRepository) FindByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]model.OrganizationMember, error) {
	var members []model.OrganizationMember
	result := database.DB.WithContext(ctx).
		Preload("Character").
		Where("organization_id = ?", orgID).
		Find(&members)
	return members, result.Error
}

// Update 更新成员
func (r *OrganizationMemberRepository) Update(ctx context.Context, member *model.OrganizationMember) error {
	return database.DB.WithContext(ctx).Save(member).Error
}

// Delete 删除成员
func (r *OrganizationMemberRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.OrganizationMember{}).Error
}
