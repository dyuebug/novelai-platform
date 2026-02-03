package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-gateway/internal/repository"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"
)

// ========== 角色 Handler ==========

type CharacterHandler struct {
	characterService *service.CharacterService
}

func NewCharacterHandler(characterService *service.CharacterService) *CharacterHandler {
	return &CharacterHandler{
		characterService: characterService,
	}
}

// Create 创建角色
// POST /api/v1/projects/:id/characters
func (h *CharacterHandler) Create(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	var req service.CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	character, err := h.characterService.Create(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to create character")
		}
		return
	}

	response.Success(c, http.StatusCreated, "character created successfully", character)
}

// List 获取角色列表
// GET /api/v1/projects/:id/characters
func (h *CharacterHandler) List(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	page := 1
	pageSize := 50
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

	role := c.Query("role")
	sort := c.Query("sort")

	result, err := h.characterService.List(c.Request.Context(), userID, projectID, page, pageSize, role, sort)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to list characters")
		}
		return
	}

	response.Success(c, http.StatusOK, "characters retrieved successfully", result)
}

// Get 获取单个角色
// GET /api/v1/characters/:id
func (h *CharacterHandler) Get(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	characterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid character id")
		return
	}

	character, err := h.characterService.Get(c.Request.Context(), userID, characterID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrCharacterNotFound):
			response.Error(c, http.StatusNotFound, "character not found")
		case errors.Is(err, service.ErrCharacterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get character")
		}
		return
	}

	response.Success(c, http.StatusOK, "character retrieved successfully", character)
}

// Update 更新角色
// PUT /api/v1/characters/:id
func (h *CharacterHandler) Update(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	characterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid character id")
		return
	}

	var req service.UpdateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	character, err := h.characterService.Update(c.Request.Context(), userID, characterID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrCharacterNotFound):
			response.Error(c, http.StatusNotFound, "character not found")
		case errors.Is(err, service.ErrCharacterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update character")
		}
		return
	}

	response.Success(c, http.StatusOK, "character updated successfully", character)
}

// Delete 删除角色
// DELETE /api/v1/characters/:id
func (h *CharacterHandler) Delete(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	characterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid character id")
		return
	}

	err = h.characterService.Delete(c.Request.Context(), userID, characterID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrCharacterNotFound):
			response.Error(c, http.StatusNotFound, "character not found")
		case errors.Is(err, service.ErrCharacterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to delete character")
		}
		return
	}

	response.Success(c, http.StatusOK, "character deleted successfully", nil)
}

// --- 角色关系 ---

// CreateRelationship 创建角色关系
// POST /api/v1/characters/:id/relationships
func (h *CharacterHandler) CreateRelationship(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	characterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid character id")
		return
	}

	var req service.CreateRelationshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	rel, err := h.characterService.CreateRelationship(c.Request.Context(), userID, characterID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create relationship")
		return
	}

	response.Success(c, http.StatusCreated, "relationship created successfully", rel)
}

// GetRelationships 获取角色关系列表
// GET /api/v1/characters/:id/relationships
func (h *CharacterHandler) GetRelationships(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	characterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid character id")
		return
	}

	rels, err := h.characterService.GetRelationships(c.Request.Context(), userID, characterID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get relationships")
		return
	}

	response.Success(c, http.StatusOK, "relationships retrieved successfully", rels)
}

// DeleteRelationship 删除角色关系
// DELETE /api/v1/relationships/:id
func (h *CharacterHandler) DeleteRelationship(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	relID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid relationship id")
		return
	}

	err = h.characterService.DeleteRelationship(c.Request.Context(), userID, relID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to delete relationship")
		return
	}

	response.Success(c, http.StatusOK, "relationship deleted successfully", nil)
}

// --- 角色经历 ---

// CreateExperience 创建角色经历
// POST /api/v1/characters/:id/experiences
func (h *CharacterHandler) CreateExperience(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	characterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid character id")
		return
	}

	var req service.CreateExperienceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	exp, err := h.characterService.CreateExperience(c.Request.Context(), userID, characterID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create experience")
		return
	}

	response.Success(c, http.StatusCreated, "experience created successfully", exp)
}

// GetExperiences 获取角色经历列表
// GET /api/v1/characters/:id/experiences
func (h *CharacterHandler) GetExperiences(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	characterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid character id")
		return
	}

	exps, err := h.characterService.GetExperiences(c.Request.Context(), userID, characterID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get experiences")
		return
	}

	response.Success(c, http.StatusOK, "experiences retrieved successfully", exps)
}

// DeleteExperience 删除角色经历
// DELETE /api/v1/experiences/:id
func (h *CharacterHandler) DeleteExperience(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	expID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid experience id")
		return
	}

	err = h.characterService.DeleteExperience(c.Request.Context(), userID, expID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to delete experience")
		return
	}

	response.Success(c, http.StatusOK, "experience deleted successfully", nil)
}
