package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-gateway/internal/repository"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"
)

// SearchServiceInterface 搜索服务接口
type SearchServiceInterface interface {
	GlobalSearch(ctx context.Context, userID uuid.UUID, req *service.GlobalSearchRequest) (*service.SearchResponse, error)
	ProjectSearch(ctx context.Context, userID, projectID uuid.UUID, req *service.ProjectSearchRequest) (*service.SearchResponse, error)
	AdvancedFilter(ctx context.Context, userID, projectID uuid.UUID, req *service.AdvancedFilterRequest) (*service.SearchResponse, error)
}

// SearchHandler 搜索处理器
type SearchHandler struct {
	searchService SearchServiceInterface
}

// NewSearchHandler 创建搜索处理器
func NewSearchHandler(searchService SearchServiceInterface) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// GlobalSearch 全局搜索
// GET /api/v1/search
func (h *SearchHandler) GlobalSearch(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req service.GlobalSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	result, err := h.searchService.GlobalSearch(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "search failed")
		return
	}

	response.Success(c, http.StatusOK, "search completed", result)
}

// ProjectSearch 项目内搜索
// GET /api/v1/projects/:id/search
func (h *SearchHandler) ProjectSearch(c *gin.Context) {
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

	var req service.ProjectSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	result, err := h.searchService.ProjectSearch(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "search failed")
		}
		return
	}

	response.Success(c, http.StatusOK, "search completed", result)
}

// AdvancedFilter 高级筛选
// POST /api/v1/projects/:id/advanced-filter
func (h *SearchHandler) AdvancedFilter(c *gin.Context) {
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

	var req service.AdvancedFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	result, err := h.searchService.AdvancedFilter(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProjectNotFound):
			response.Error(c, http.StatusNotFound, "project not found")
		case errors.Is(err, service.ErrProjectNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "filter failed")
		}
		return
	}

	response.Success(c, http.StatusOK, "filter completed", result)
}
