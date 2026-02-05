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
	ErrCharacterNotFound    = errors.New("character not found")
	ErrRelationshipNotFound = errors.New("relationship not found")
	ErrExperienceNotFound   = errors.New("experience not found")
)

type CharacterRepository struct{}

func NewCharacterRepository() *CharacterRepository {
	return &CharacterRepository{}
}

// Create 创建角色
func (r *CharacterRepository) Create(ctx context.Context, character *model.Character) error {
	return database.DB.WithContext(ctx).Create(character).Error
}

// FindByID 通过 ID 查找角色
func (r *CharacterRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Character, error) {
	var character model.Character
	result := database.DB.WithContext(ctx).
		Where("id = ?", id).
		First(&character)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrCharacterNotFound
		}
		return nil, result.Error
	}
	return &character, nil
}

// FindByProjectID 查找项目的所有角色
func (r *CharacterRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts *ListOptions) ([]model.Character, int64, error) {
	var characters []model.Character
	var total int64

	query := database.DB.WithContext(ctx).
		Model(&model.Character{}).
		Where("project_id = ?", projectID)

	// 角色类型过滤
	if opts.Status != "" {
		query = query.Where("role = ?", opts.Status)
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

	result := query.Find(&characters)
	return characters, total, result.Error
}

// Update 更新角色
func (r *CharacterRepository) Update(ctx context.Context, character *model.Character) error {
	return database.DB.WithContext(ctx).Save(character).Error
}

// Delete 删除角色
func (r *CharacterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Character{}).Error
}

// Search 搜索角色
func (r *CharacterRepository) Search(ctx context.Context, projectID uuid.UUID, query string, page, pageSize int) ([]model.Character, int64, error) {
	var characters []model.Character
	var total int64

	// 使用 ILIKE 进行模糊搜索
	searchPattern := "%" + query + "%"

	queryDB := database.DB.WithContext(ctx).
		Model(&model.Character{}).
		Where("project_id = ?", projectID).
		Where("name ILIKE ? OR alias ILIKE ? OR background ILIKE ?", searchPattern, searchPattern, searchPattern)

	// 计算总数
	queryDB.Count(&total)

	// 分页
	offset := (page - 1) * pageSize
	result := queryDB.
		Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&characters)

	return characters, total, result.Error
}

// --- 角色关系 ---

type CharacterRelationshipRepository struct{}

func NewCharacterRelationshipRepository() *CharacterRelationshipRepository {
	return &CharacterRelationshipRepository{}
}

// Create 创建关系
func (r *CharacterRelationshipRepository) Create(ctx context.Context, rel *model.CharacterRelationship) error {
	return database.DB.WithContext(ctx).Create(rel).Error
}

// FindByID 通过 ID 查找关系
func (r *CharacterRelationshipRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.CharacterRelationship, error) {
	var rel model.CharacterRelationship
	result := database.DB.WithContext(ctx).
		Preload("Target").
		Where("id = ?", id).
		First(&rel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRelationshipNotFound
		}
		return nil, result.Error
	}
	return &rel, nil
}

// FindByCharacterID 查找角色的所有关系
func (r *CharacterRelationshipRepository) FindByCharacterID(ctx context.Context, characterID uuid.UUID) ([]model.CharacterRelationship, error) {
	var rels []model.CharacterRelationship
	result := database.DB.WithContext(ctx).
		Preload("Target").
		Where("character_id = ?", characterID).
		Find(&rels)
	return rels, result.Error
}

// Update 更新关系
func (r *CharacterRelationshipRepository) Update(ctx context.Context, rel *model.CharacterRelationship) error {
	return database.DB.WithContext(ctx).Save(rel).Error
}

// Delete 删除关系
func (r *CharacterRelationshipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.CharacterRelationship{}).Error
}

// --- 角色经历 ---

type CharacterExperienceRepository struct{}

func NewCharacterExperienceRepository() *CharacterExperienceRepository {
	return &CharacterExperienceRepository{}
}

// Create 创建经历
func (r *CharacterExperienceRepository) Create(ctx context.Context, exp *model.CharacterExperience) error {
	return database.DB.WithContext(ctx).Create(exp).Error
}

// FindByID 通过 ID 查找经历
func (r *CharacterExperienceRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.CharacterExperience, error) {
	var exp model.CharacterExperience
	result := database.DB.WithContext(ctx).
		Where("id = ?", id).
		First(&exp)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrExperienceNotFound
		}
		return nil, result.Error
	}
	return &exp, nil
}

// FindByCharacterID 查找角色的所有经历
func (r *CharacterExperienceRepository) FindByCharacterID(ctx context.Context, characterID uuid.UUID) ([]model.CharacterExperience, error) {
	var exps []model.CharacterExperience
	result := database.DB.WithContext(ctx).
		Where("character_id = ?", characterID).
		Order("sort_order ASC, created_at ASC").
		Find(&exps)
	return exps, result.Error
}

// Update 更新经历
func (r *CharacterExperienceRepository) Update(ctx context.Context, exp *model.CharacterExperience) error {
	return database.DB.WithContext(ctx).Save(exp).Error
}

// Delete 删除经历
func (r *CharacterExperienceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.CharacterExperience{}).Error
}
