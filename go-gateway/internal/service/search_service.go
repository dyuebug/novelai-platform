package service

import (
	"context"
	"fmt"
	"strings"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
)

// SearchProjectRepository 项目搜索仓库接口
type SearchProjectRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.Project, error)
	Search(ctx context.Context, userID uuid.UUID, query string, page, pageSize int) ([]model.Project, int64, error)
}

// SearchChapterRepository 章节搜索仓库接口
type SearchChapterRepository interface {
	Search(ctx context.Context, userID, projectID uuid.UUID, query string, page, pageSize int) ([]model.Chapter, int64, error)
	AdvancedFilter(ctx context.Context, projectID uuid.UUID, req interface{}) ([]model.Chapter, int64, error)
}

// SearchCharacterRepository 角色搜索仓库接口
type SearchCharacterRepository interface {
	Search(ctx context.Context, projectID uuid.UUID, query string, page, pageSize int) ([]model.Character, int64, error)
}

// SearchLocationRepository 地点搜索仓库接口
type SearchLocationRepository interface {
	Search(ctx context.Context, projectID uuid.UUID, query string, page, pageSize int) ([]model.Location, int64, error)
}

// SearchService 搜索服务
type SearchService struct {
	projectRepo   SearchProjectRepository
	chapterRepo   SearchChapterRepository
	characterRepo SearchCharacterRepository
	locationRepo  SearchLocationRepository
}

// NewSearchService 创建搜索服务
func NewSearchService() *SearchService {
	return &SearchService{
		projectRepo:   repository.NewProjectRepository(),
		chapterRepo:   repository.NewChapterRepository(),
		characterRepo: repository.NewCharacterRepository(),
		locationRepo:  repository.NewLocationRepository(),
	}
}

// GlobalSearchRequest 全局搜索请求
type GlobalSearchRequest struct {
	Query    string `form:"q" binding:"required,min=1"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

// ProjectSearchRequest 项目内搜索请求
type ProjectSearchRequest struct {
	Query    string `form:"q" binding:"required,min=1"`
	Type     string `form:"type"` // chapter, character, location, all
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

// AdvancedFilterRequest 高级筛选请求
type AdvancedFilterRequest struct {
	Tags         []string `json:"tags"`
	Status       []string `json:"status"`
	DateFrom     string   `json:"date_from"`
	DateTo       string   `json:"date_to"`
	WordCountMin int      `json:"word_count_min"`
	WordCountMax int      `json:"word_count_max"`
	Page         int      `json:"page"`
	PageSize     int      `json:"page_size"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Type      string                 `json:"type"`      // project, chapter, character, location
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	Snippet   string                 `json:"snippet"`   // 匹配片段
	Highlight string                 `json:"highlight"` // 高亮文本
	Metadata  map[string]interface{} `json:"metadata"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Results    []SearchResult `json:"results"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// GlobalSearch 全局搜索
func (s *SearchService) GlobalSearch(ctx context.Context, userID uuid.UUID, req *GlobalSearchRequest) (*SearchResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	results := []SearchResult{}
	var total int64

	// 搜索项目
	projects, projectTotal, err := s.searchProjects(ctx, userID, req.Query, req.Page, req.PageSize)
	if err == nil {
		for _, project := range projects {
			results = append(results, SearchResult{
				Type:      "project",
				ID:        project.ID.String(),
				Title:     project.Title,
				Snippet:   truncate(project.Description, 200),
				Highlight: highlightText(project.Title, req.Query),
				Metadata: map[string]interface{}{
					"genre":  project.Genre,
					"status": project.Status,
				},
			})
		}
		total += projectTotal
	}

	// 搜索章节
	chapters, chapterTotal, err := s.searchChapters(ctx, userID, uuid.Nil, req.Query, req.Page, req.PageSize)
	if err == nil {
		for _, chapter := range chapters {
			results = append(results, SearchResult{
				Type:      "chapter",
				ID:        chapter.ID.String(),
				Title:     chapter.Title,
				Snippet:   truncate(chapter.Summary, 200),
				Highlight: highlightText(chapter.Title, req.Query),
				Metadata: map[string]interface{}{
					"chapter_number": chapter.ChapterNumber,
					"status":         chapter.Status,
					"word_count":     chapter.WordCount,
				},
			})
		}
		total += chapterTotal
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &SearchResponse{
		Results:    results,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// ProjectSearch 项目内搜索
func (s *SearchService) ProjectSearch(ctx context.Context, userID, projectID uuid.UUID, req *ProjectSearchRequest) (*SearchResponse, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	results := []SearchResult{}
	var total int64

	// 根据类型搜索
	switch req.Type {
	case "chapter":
		chapters, chapterTotal, err := s.searchChapters(ctx, userID, projectID, req.Query, req.Page, req.PageSize)
		if err != nil {
			return nil, err
		}
		for _, chapter := range chapters {
			results = append(results, SearchResult{
				Type:      "chapter",
				ID:        chapter.ID.String(),
				Title:     chapter.Title,
				Snippet:   truncate(chapter.Summary, 200),
				Highlight: highlightText(chapter.Title, req.Query),
				Metadata: map[string]interface{}{
					"chapter_number": chapter.ChapterNumber,
					"status":         chapter.Status,
					"word_count":     chapter.WordCount,
				},
			})
		}
		total = chapterTotal

	case "character":
		characters, characterTotal, err := s.searchCharacters(ctx, projectID, req.Query, req.Page, req.PageSize)
		if err != nil {
			return nil, err
		}
		for _, character := range characters {
			results = append(results, SearchResult{
				Type:      "character",
				ID:        character.ID.String(),
				Title:     character.Name,
				Snippet:   truncate(character.Background, 200),
				Highlight: highlightText(character.Name, req.Query),
				Metadata: map[string]interface{}{
					"role":   character.Role,
					"status": character.Status,
				},
			})
		}
		total = characterTotal

	case "location":
		locations, locationTotal, err := s.searchLocations(ctx, projectID, req.Query, req.Page, req.PageSize)
		if err != nil {
			return nil, err
		}
		for _, location := range locations {
			results = append(results, SearchResult{
				Type:      "location",
				ID:        location.ID.String(),
				Title:     location.Name,
				Snippet:   truncate(location.Description, 200),
				Highlight: highlightText(location.Name, req.Query),
				Metadata: map[string]interface{}{
					"type": location.Type,
				},
			})
		}
		total = locationTotal

	default: // "all" or empty
		// 搜索所有类型
		chapters, _, _ := s.searchChapters(ctx, userID, projectID, req.Query, 1, 10)
		for _, chapter := range chapters {
			results = append(results, SearchResult{
				Type:      "chapter",
				ID:        chapter.ID.String(),
				Title:     chapter.Title,
				Snippet:   truncate(chapter.Summary, 200),
				Highlight: highlightText(chapter.Title, req.Query),
			})
		}

		characters, _, _ := s.searchCharacters(ctx, projectID, req.Query, 1, 10)
		for _, character := range characters {
			results = append(results, SearchResult{
				Type:      "character",
				ID:        character.ID.String(),
				Title:     character.Name,
				Snippet:   truncate(character.Background, 200),
				Highlight: highlightText(character.Name, req.Query),
			})
		}

		locations, _, _ := s.searchLocations(ctx, projectID, req.Query, 1, 10)
		for _, location := range locations {
			results = append(results, SearchResult{
				Type:      "location",
				ID:        location.ID.String(),
				Title:     location.Name,
				Snippet:   truncate(location.Description, 200),
				Highlight: highlightText(location.Name, req.Query),
			})
		}

		total = int64(len(results))
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &SearchResponse{
		Results:    results,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// AdvancedFilter 高级筛选
func (s *SearchService) AdvancedFilter(ctx context.Context, userID, projectID uuid.UUID, req *AdvancedFilterRequest) (*SearchResponse, error) {
	// 验证项目所有权
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrProjectNotOwned
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 使用 ChapterRepository 的高级筛选
	chapters, total, err := s.chapterRepo.AdvancedFilter(ctx, projectID, req)
	if err != nil {
		return nil, err
	}

	results := []SearchResult{}
	for _, chapter := range chapters {
		results = append(results, SearchResult{
			Type:    "chapter",
			ID:      chapter.ID.String(),
			Title:   chapter.Title,
			Snippet: truncate(chapter.Summary, 200),
			Metadata: map[string]interface{}{
				"chapter_number": chapter.ChapterNumber,
				"status":         chapter.Status,
				"word_count":     chapter.WordCount,
			},
		})
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &SearchResponse{
		Results:    results,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// 辅助方法

func (s *SearchService) searchProjects(ctx context.Context, userID uuid.UUID, query string, page, pageSize int) ([]model.Project, int64, error) {
	return s.projectRepo.Search(ctx, userID, query, page, pageSize)
}

func (s *SearchService) searchChapters(ctx context.Context, userID, projectID uuid.UUID, query string, page, pageSize int) ([]model.Chapter, int64, error) {
	return s.chapterRepo.Search(ctx, userID, projectID, query, page, pageSize)
}

func (s *SearchService) searchCharacters(ctx context.Context, projectID uuid.UUID, query string, page, pageSize int) ([]model.Character, int64, error) {
	return s.characterRepo.Search(ctx, projectID, query, page, pageSize)
}

func (s *SearchService) searchLocations(ctx context.Context, projectID uuid.UUID, query string, page, pageSize int) ([]model.Location, int64, error) {
	return s.locationRepo.Search(ctx, projectID, query, page, pageSize)
}

// truncate 截断文本
func truncate(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// highlightText 高亮匹配文本
func highlightText(text, query string) string {
	if query == "" {
		return text
	}
	// 简单的高亮实现，使用 <mark> 标签
	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)

	if strings.Contains(lowerText, lowerQuery) {
		// 找到匹配位置
		index := strings.Index(lowerText, lowerQuery)
		before := text[:index]
		match := text[index : index+len(query)]
		after := text[index+len(query):]
		return fmt.Sprintf("%s<mark>%s</mark>%s", before, match, after)
	}

	return text
}
