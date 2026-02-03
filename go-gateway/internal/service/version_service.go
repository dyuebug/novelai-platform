package service

import (
	"context"
	"errors"
	"strings"

	"go-gateway/internal/model"
	"go-gateway/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrVersionNotFound = errors.New("version not found")
)

type VersionService struct {
	chapterRepo        *repository.ChapterRepository
	chapterVersionRepo *repository.ChapterVersionRepository
	projectRepo        *repository.ProjectRepository
}

func NewVersionService() *VersionService {
	return &VersionService{
		chapterRepo:        repository.NewChapterRepository(),
		chapterVersionRepo: repository.NewChapterVersionRepository(),
		projectRepo:        repository.NewProjectRepository(),
	}
}

// VersionListResponse 版本列表响应
type VersionListResponse struct {
	Items      []model.ChapterVersion `json:"items"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

// DiffResponse 版本对比响应
type DiffResponse struct {
	Version1      int        `json:"version1"`
	Version2      int        `json:"version2"`
	Content1      string     `json:"content1"`
	Content2      string     `json:"content2"`
	Additions     int        `json:"additions"`
	Deletions     int        `json:"deletions"`
	DiffLines     []DiffLine `json:"diff_lines"`
}

// DiffLine 差异行
type DiffLine struct {
	Type    string `json:"type"` // "add", "delete", "equal"
	Content string `json:"content"`
	LineNum int    `json:"line_num,omitempty"`
}

// ListVersions 获取章节版本列表
func (s *VersionService) ListVersions(ctx context.Context, userID, chapterID uuid.UUID, page, pageSize int) (*VersionListResponse, error) {
	// 验证章节所有权
	chapter, err := s.chapterRepo.FindByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	project, err := s.projectRepo.FindByID(ctx, chapter.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrChapterNotOwned
	}

	// 限制 pageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page <= 0 {
		page = 1
	}

	opts := &repository.ListOptions{
		Page:     page,
		PageSize: pageSize,
	}

	versions, total, err := s.chapterVersionRepo.FindByChapterID(ctx, chapterID, opts)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &VersionListResponse{
		Items:      versions,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetVersion 获取指定版本
func (s *VersionService) GetVersion(ctx context.Context, userID, chapterID uuid.UUID, versionNumber int) (*model.ChapterVersion, error) {
	// 验证章节所有权
	chapter, err := s.chapterRepo.FindByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	project, err := s.projectRepo.FindByID(ctx, chapter.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrChapterNotOwned
	}

	return s.chapterVersionRepo.FindByVersionNumber(ctx, chapterID, versionNumber)
}

// RestoreVersion 恢复到指定版本
func (s *VersionService) RestoreVersion(ctx context.Context, userID, chapterID uuid.UUID, versionNumber int) (*model.Chapter, error) {
	// 验证章节所有权
	chapter, err := s.chapterRepo.FindByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	project, err := s.projectRepo.FindByID(ctx, chapter.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrChapterNotOwned
	}

	// 获取目标版本
	version, err := s.chapterVersionRepo.FindByVersionNumber(ctx, chapterID, versionNumber)
	if err != nil {
		return nil, err
	}

	// 更新章节内容
	chapter.Content = version.Content
	chapter.WordCount = version.WordCount

	if err := s.chapterRepo.Update(ctx, chapter); err != nil {
		return nil, err
	}

	// 创建新版本记录 (恢复操作也产生新版本)
	maxVersion, _ := s.chapterVersionRepo.GetMaxVersionNumber(ctx, chapterID)
	newVersion := &model.ChapterVersion{
		ChapterID:     chapterID,
		VersionNumber: maxVersion + 1,
		Content:       version.Content,
		WordCount:     version.WordCount,
		Source:        "restore",
	}
	_ = s.chapterVersionRepo.Create(ctx, newVersion)

	return chapter, nil
}

// DiffVersions 对比两个版本
func (s *VersionService) DiffVersions(ctx context.Context, userID, chapterID uuid.UUID, v1, v2 int) (*DiffResponse, error) {
	// 验证章节所有权
	chapter, err := s.chapterRepo.FindByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	project, err := s.projectRepo.FindByID(ctx, chapter.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrChapterNotOwned
	}

	// 获取两个版本
	version1, version2, err := s.chapterVersionRepo.GetTwoVersions(ctx, chapterID, v1, v2)
	if err != nil {
		return nil, err
	}

	// 简单的行级 diff
	lines1 := strings.Split(version1.Content, "\n")
	lines2 := strings.Split(version2.Content, "\n")

	diffLines, additions, deletions := computeSimpleDiff(lines1, lines2)

	return &DiffResponse{
		Version1:  v1,
		Version2:  v2,
		Content1:  version1.Content,
		Content2:  version2.Content,
		Additions: additions,
		Deletions: deletions,
		DiffLines: diffLines,
	}, nil
}

// computeSimpleDiff 简单的行级差异计算
func computeSimpleDiff(lines1, lines2 []string) ([]DiffLine, int, int) {
	additions := 0
	deletions := 0

	// 使用简单的 LCS 算法进行对比
	m, n := len(lines1), len(lines2)

	// 构建 LCS 表
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if lines1[i-1] == lines2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] > dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}

	// 回溯生成 diff
	i, j := m, n
	var result []DiffLine

	for i > 0 || j > 0 {
		if i > 0 && j > 0 && lines1[i-1] == lines2[j-1] {
			result = append(result, DiffLine{Type: "equal", Content: lines1[i-1]})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			result = append(result, DiffLine{Type: "add", Content: lines2[j-1], LineNum: j})
			additions++
			j--
		} else if i > 0 {
			result = append(result, DiffLine{Type: "delete", Content: lines1[i-1], LineNum: i})
			deletions++
			i--
		}
	}

	// 反转结果
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, additions, deletions
}
