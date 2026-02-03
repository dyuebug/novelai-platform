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

// ========== 地点 Handler ==========

type LocationHandler struct {
	locationService *service.LocationService
}

func NewLocationHandler(locationService *service.LocationService) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
	}
}

// Create 创建地点
// POST /api/v1/projects/:id/locations
func (h *LocationHandler) Create(c *gin.Context) {
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

	var req service.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	location, err := h.locationService.Create(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to create location")
		}
		return
	}

	response.Success(c, http.StatusCreated, "location created successfully", location)
}

// List 获取地点列表
// GET /api/v1/projects/:id/locations
func (h *LocationHandler) List(c *gin.Context) {
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

	locationType := c.Query("type")
	sort := c.Query("sort")

	result, err := h.locationService.List(c.Request.Context(), userID, projectID, page, pageSize, locationType, sort)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to list locations")
		}
		return
	}

	response.Success(c, http.StatusOK, "locations retrieved successfully", result)
}

// GetTree 获取地点树结构
// GET /api/v1/projects/:id/locations/tree
func (h *LocationHandler) GetTree(c *gin.Context) {
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

	tree, err := h.locationService.GetTree(c.Request.Context(), userID, projectID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get location tree")
		}
		return
	}

	response.Success(c, http.StatusOK, "location tree retrieved successfully", tree)
}

// Get 获取单个地点
// GET /api/v1/locations/:id
func (h *LocationHandler) Get(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	locationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid location id")
		return
	}

	location, err := h.locationService.Get(c.Request.Context(), userID, locationID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrLocationNotFound):
			response.Error(c, http.StatusNotFound, "location not found")
		case errors.Is(err, service.ErrLocationNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get location")
		}
		return
	}

	response.Success(c, http.StatusOK, "location retrieved successfully", location)
}

// Update 更新地点
// PUT /api/v1/locations/:id
func (h *LocationHandler) Update(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	locationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid location id")
		return
	}

	var req service.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	location, err := h.locationService.Update(c.Request.Context(), userID, locationID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrLocationNotFound):
			response.Error(c, http.StatusNotFound, "location not found")
		case errors.Is(err, service.ErrLocationNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update location")
		}
		return
	}

	response.Success(c, http.StatusOK, "location updated successfully", location)
}

// Delete 删除地点
// DELETE /api/v1/locations/:id
func (h *LocationHandler) Delete(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	locationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid location id")
		return
	}

	err = h.locationService.Delete(c.Request.Context(), userID, locationID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrLocationNotFound):
			response.Error(c, http.StatusNotFound, "location not found")
		case errors.Is(err, service.ErrLocationNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to delete location")
		}
		return
	}

	response.Success(c, http.StatusOK, "location deleted successfully", nil)
}

// ========== 组织 Handler ==========

type OrganizationHandler struct {
	orgService *service.OrganizationService
}

func NewOrganizationHandler(orgService *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

// Create 创建组织
// POST /api/v1/projects/:id/organizations
func (h *OrganizationHandler) Create(c *gin.Context) {
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

	var req service.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	org, err := h.orgService.Create(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to create organization")
		}
		return
	}

	response.Success(c, http.StatusCreated, "organization created successfully", org)
}

// List 获取组织列表
// GET /api/v1/projects/:id/organizations
func (h *OrganizationHandler) List(c *gin.Context) {
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

	orgType := c.Query("type")
	sort := c.Query("sort")

	result, err := h.orgService.List(c.Request.Context(), userID, projectID, page, pageSize, orgType, sort)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to list organizations")
		}
		return
	}

	response.Success(c, http.StatusOK, "organizations retrieved successfully", result)
}

// Get 获取单个组织
// GET /api/v1/organizations/:id
func (h *OrganizationHandler) Get(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid organization id")
		return
	}

	org, err := h.orgService.Get(c.Request.Context(), userID, orgID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrOrganizationNotFound):
			response.Error(c, http.StatusNotFound, "organization not found")
		case errors.Is(err, service.ErrOrganizationNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get organization")
		}
		return
	}

	response.Success(c, http.StatusOK, "organization retrieved successfully", org)
}

// Update 更新组织
// PUT /api/v1/organizations/:id
func (h *OrganizationHandler) Update(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid organization id")
		return
	}

	var req service.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	org, err := h.orgService.Update(c.Request.Context(), userID, orgID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrOrganizationNotFound):
			response.Error(c, http.StatusNotFound, "organization not found")
		case errors.Is(err, service.ErrOrganizationNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update organization")
		}
		return
	}

	response.Success(c, http.StatusOK, "organization updated successfully", org)
}

// Delete 删除组织
// DELETE /api/v1/organizations/:id
func (h *OrganizationHandler) Delete(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid organization id")
		return
	}

	err = h.orgService.Delete(c.Request.Context(), userID, orgID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrOrganizationNotFound):
			response.Error(c, http.StatusNotFound, "organization not found")
		case errors.Is(err, service.ErrOrganizationNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to delete organization")
		}
		return
	}

	response.Success(c, http.StatusOK, "organization deleted successfully", nil)
}

// --- 组织成员 ---

// AddMember 添加组织成员
// POST /api/v1/organizations/:id/members
func (h *OrganizationHandler) AddMember(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid organization id")
		return
	}

	var req service.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	member, err := h.orgService.AddMember(c.Request.Context(), userID, orgID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to add member")
		return
	}

	response.Success(c, http.StatusCreated, "member added successfully", member)
}

// GetMembers 获取组织成员列表
// GET /api/v1/organizations/:id/members
func (h *OrganizationHandler) GetMembers(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid organization id")
		return
	}

	members, err := h.orgService.GetMembers(c.Request.Context(), userID, orgID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get members")
		return
	}

	response.Success(c, http.StatusOK, "members retrieved successfully", members)
}

// RemoveMember 移除组织成员
// DELETE /api/v1/members/:id
func (h *OrganizationHandler) RemoveMember(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	memberID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid member id")
		return
	}

	err = h.orgService.RemoveMember(c.Request.Context(), userID, memberID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to remove member")
		return
	}

	response.Success(c, http.StatusOK, "member removed successfully", nil)
}
