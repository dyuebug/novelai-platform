package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-gateway/internal/model"
	"go-gateway/internal/repository"
	"go-gateway/internal/service"
	"go-gateway/pkg/response"
)

// ChapterGenerateServiceInterface 章节生成服务接口
type ChapterGenerateServiceInterface interface {
	Get(ctx context.Context, userID, chapterID uuid.UUID) (*model.Chapter, error)
}

// ProjectGenerateServiceInterface 项目服务接口
type ProjectGenerateServiceInterface interface {
	Get(ctx context.Context, userID, projectID uuid.UUID) (*model.Project, error)
}

// AIGenerateServiceInterface AI 服务接口
type AIGenerateServiceInterface interface {
	GenerateStream(ctx context.Context, req *service.GenerateRequest) (service.Stream, error)
}

type ChapterGenerateHandler struct {
	ai             AIGenerateServiceInterface
	chapterService ChapterGenerateServiceInterface
	projectService ProjectGenerateServiceInterface
	chapterRepo    *repository.ChapterRepository
	versionRepo    *repository.ChapterVersionRepository
}

func NewChapterGenerateHandler(
	ai AIGenerateServiceInterface,
	chapterService ChapterGenerateServiceInterface,
	projectService ProjectGenerateServiceInterface,
) *ChapterGenerateHandler {
	return &ChapterGenerateHandler{
		ai:             ai,
		chapterService: chapterService,
		projectService: projectService,
		chapterRepo:    repository.NewChapterRepository(),
		versionRepo:    repository.NewChapterVersionRepository(),
	}
}

// GenerateChapterRequest 章节生成请求
type GenerateChapterRequest struct {
	Model       string  `json:"model"`
	Provider    string  `json:"provider"`
	Instruction string  `json:"instruction"` // 用户指令
	Context     string  `json:"context"`     // 额外上下文
	Temperature float32 `json:"temperature"`
	MaxTokens   int32   `json:"max_tokens"`
}

// PartialRegenerateRequest 局部重写请求
type PartialRegenerateRequest struct {
	Selection   string  `json:"selection" binding:"required"` // 选中的文本
	Instruction string  `json:"instruction"`                  // 重写指令
	Model       string  `json:"model"`
	Provider    string  `json:"provider"`
	Temperature float32 `json:"temperature"`
}

// PolishRequest 润色请求
type PolishRequest struct {
	Model       string  `json:"model"`
	Provider    string  `json:"provider"`
	Instruction string  `json:"instruction"` // 润色指令
	Temperature float32 `json:"temperature"`
}

// GenerateStream 流式生成章节内容
// POST /api/v1/chapters/:id/generate-stream
func (h *ChapterGenerateHandler) GenerateStream(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	chapterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid chapter id")
		return
	}

	var req GenerateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	// 验证章节所有权
	chapter, err := h.chapterService.Get(c.Request.Context(), userID, chapterID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrChapterNotFound):
			response.Error(c, http.StatusNotFound, "chapter not found")
		case errors.Is(err, service.ErrChapterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get chapter")
		}
		return
	}

	// 获取项目信息
	project, err := h.projectService.Get(c.Request.Context(), userID, chapter.ProjectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get project")
		return
	}

	// 构建系统提示词
	systemPrompt := buildChapterGeneratePrompt(project, chapter, req.Instruction)

	// 设置默认值
	model := req.Model
	if model == "" {
		model = "gpt-4o"
	}
	provider := req.Provider
	if provider == "" {
		provider = "openai"
	}

	// 调用 AI 服务
	stream, err := h.ai.GenerateStream(c.Request.Context(), &service.GenerateRequest{
		Prompt:       req.Instruction,
		Model:        model,
		Provider:     provider,
		Temperature:  req.Temperature,
		MaxTokens:    req.MaxTokens,
		SystemPrompt: systemPrompt,
		UserID:       userID.String(),
		ProjectID:    project.ID.String(),
		ChapterID:    chapterID.String(),
		TaskType:     "generate",
		Context:      req.Context,
	})
	if err != nil {
		response.Error(c, http.StatusServiceUnavailable, fmt.Sprintf("AI service unavailable: %v", err))
		return
	}
	defer stream.CloseSend()

	// 收集生成的内容
	var generatedContent string

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	c.Stream(func(w io.Writer) bool {
		resp, err := stream.Recv()
		if err == io.EOF {
			// 保存生成的内容
			if generatedContent != "" {
				h.saveGeneratedContent(c.Request.Context(), chapter, generatedContent, model, provider)
			}
			c.SSEvent("done", map[string]interface{}{
				"status":     "completed",
				"word_count": utf8.RuneCountInString(generatedContent),
			})
			return false
		}
		if err != nil {
			c.SSEvent("error", map[string]string{"message": err.Error()})
			return false
		}
		if resp.Error != "" {
			c.SSEvent("error", map[string]string{"message": resp.Error})
			return false
		}

		generatedContent += resp.Content
		c.SSEvent("content", map[string]string{"text": resp.Content})

		if resp.Done {
			if generatedContent != "" {
				h.saveGeneratedContent(c.Request.Context(), chapter, generatedContent, model, provider)
			}
			c.SSEvent("done", map[string]interface{}{
				"status":     "completed",
				"word_count": utf8.RuneCountInString(generatedContent),
			})
			return false
		}
		return true
	})
}

// PartialRegenerateStream 局部重写流式生成
// POST /api/v1/chapters/:id/partial-regenerate-stream
func (h *ChapterGenerateHandler) PartialRegenerateStream(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	chapterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid chapter id")
		return
	}

	var req PartialRegenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	// 验证章节所有权
	chapter, err := h.chapterService.Get(c.Request.Context(), userID, chapterID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrChapterNotFound):
			response.Error(c, http.StatusNotFound, "chapter not found")
		case errors.Is(err, service.ErrChapterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get chapter")
		}
		return
	}

	// 构建重写提示词
	systemPrompt := buildPartialRegeneratePrompt(chapter, req.Selection, req.Instruction)

	model := req.Model
	if model == "" {
		model = "gpt-4o"
	}
	provider := req.Provider
	if provider == "" {
		provider = "openai"
	}

	stream, err := h.ai.GenerateStream(c.Request.Context(), &service.GenerateRequest{
		Prompt:       req.Instruction,
		Model:        model,
		Provider:     provider,
		Temperature:  req.Temperature,
		SystemPrompt: systemPrompt,
		UserID:       userID.String(),
		ChapterID:    chapterID.String(),
		TaskType:     "partial_regenerate",
		Selection:    req.Selection,
	})
	if err != nil {
		response.Error(c, http.StatusServiceUnavailable, fmt.Sprintf("AI service unavailable: %v", err))
		return
	}
	defer stream.CloseSend()

	var generatedContent string

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	c.Stream(func(w io.Writer) bool {
		resp, err := stream.Recv()
		if err == io.EOF {
			c.SSEvent("done", map[string]interface{}{
				"status":           "completed",
				"original":         req.Selection,
				"replacement":      generatedContent,
				"word_count_delta": utf8.RuneCountInString(generatedContent) - utf8.RuneCountInString(req.Selection),
			})
			return false
		}
		if err != nil {
			c.SSEvent("error", map[string]string{"message": err.Error()})
			return false
		}
		if resp.Error != "" {
			c.SSEvent("error", map[string]string{"message": resp.Error})
			return false
		}

		generatedContent += resp.Content
		c.SSEvent("content", map[string]string{"text": resp.Content})

		if resp.Done {
			c.SSEvent("done", map[string]interface{}{
				"status":           "completed",
				"original":         req.Selection,
				"replacement":      generatedContent,
				"word_count_delta": utf8.RuneCountInString(generatedContent) - utf8.RuneCountInString(req.Selection),
			})
			return false
		}
		return true
	})
}

// PolishStream 润色流式生成
// POST /api/v1/chapters/:id/polish-stream
func (h *ChapterGenerateHandler) PolishStream(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	chapterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid chapter id")
		return
	}

	var req PolishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	// 验证章节所有权
	chapter, err := h.chapterService.Get(c.Request.Context(), userID, chapterID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrChapterNotFound):
			response.Error(c, http.StatusNotFound, "chapter not found")
		case errors.Is(err, service.ErrChapterNotOwned):
			response.Error(c, http.StatusForbidden, "access denied")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get chapter")
		}
		return
	}

	// 构建润色提示词
	systemPrompt := buildPolishPrompt(chapter, req.Instruction)

	model := req.Model
	if model == "" {
		model = "gpt-4o"
	}
	provider := req.Provider
	if provider == "" {
		provider = "openai"
	}

	stream, err := h.ai.GenerateStream(c.Request.Context(), &service.GenerateRequest{
		Prompt:       chapter.Content,
		Model:        model,
		Provider:     provider,
		Temperature:  req.Temperature,
		SystemPrompt: systemPrompt,
		UserID:       userID.String(),
		ChapterID:    chapterID.String(),
		TaskType:     "polish",
		Instruction:  req.Instruction,
	})
	if err != nil {
		response.Error(c, http.StatusServiceUnavailable, fmt.Sprintf("AI service unavailable: %v", err))
		return
	}
	defer stream.CloseSend()

	var generatedContent string

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	c.Stream(func(w io.Writer) bool {
		resp, err := stream.Recv()
		if err == io.EOF {
			c.SSEvent("done", map[string]interface{}{
				"status":     "completed",
				"word_count": utf8.RuneCountInString(generatedContent),
			})
			return false
		}
		if err != nil {
			c.SSEvent("error", map[string]string{"message": err.Error()})
			return false
		}
		if resp.Error != "" {
			c.SSEvent("error", map[string]string{"message": resp.Error})
			return false
		}

		generatedContent += resp.Content
		c.SSEvent("content", map[string]string{"text": resp.Content})

		if resp.Done {
			c.SSEvent("done", map[string]interface{}{
				"status":     "completed",
				"word_count": utf8.RuneCountInString(generatedContent),
			})
			return false
		}
		return true
	})
}

// saveGeneratedContent 保存生成的内容
func (h *ChapterGenerateHandler) saveGeneratedContent(ctx context.Context, chapter *model.Chapter, content, aiModel, aiProvider string) {
	// 更新章节内容
	chapter.Content = content
	chapter.WordCount = utf8.RuneCountInString(content)
	_ = h.chapterRepo.Update(ctx, chapter)

	// 创建版本记录
	maxVersion, _ := h.versionRepo.GetMaxVersionNumber(ctx, chapter.ID)
	version := &model.ChapterVersion{
		ChapterID:     chapter.ID,
		VersionNumber: maxVersion + 1,
		Content:       content,
		WordCount:     chapter.WordCount,
		Source:        "ai",
		AIProvider:    aiProvider,
		AIModel:       aiModel,
	}
	_ = h.versionRepo.Create(ctx, version)
}

// buildChapterGeneratePrompt 构建章节生成提示词
func buildChapterGeneratePrompt(project *model.Project, chapter *model.Chapter, instruction string) string {
	prompt := fmt.Sprintf(`你是一位专业的网络小说作家，正在创作一部%s类型的小说《%s》。

小说简介：%s

当前任务：撰写第%d章《%s》的内容。

写作要求：
1. 保持文风一致，符合网络小说的阅读习惯
2. 注意情节的连贯性和节奏感
3. 人物对话要生动自然
4. 适当设置悬念和爽点
5. 字数控制在2000-4000字左右

`, project.Genre, project.Title, project.Description, chapter.ChapterNumber, chapter.Title)

	if instruction != "" {
		prompt += fmt.Sprintf("用户特别要求：%s\n\n", instruction)
	}

	prompt += "请直接输出章节正文内容，不需要输出章节标题。"

	return prompt
}

// buildPartialRegeneratePrompt 构建局部重写提示词
func buildPartialRegeneratePrompt(chapter *model.Chapter, selection, instruction string) string {
	prompt := fmt.Sprintf(`你是一位专业的网络小说编辑，需要对以下选中的文本进行重写。

章节标题：%s

选中的原文：
%s

`, chapter.Title, selection)

	if instruction != "" {
		prompt += fmt.Sprintf("重写要求：%s\n\n", instruction)
	} else {
		prompt += "请保持原意，优化表达，使文字更加流畅生动。\n\n"
	}

	prompt += "请直接输出重写后的内容，不需要任何解释。"

	return prompt
}

// buildPolishPrompt 构建润色提示词
func buildPolishPrompt(chapter *model.Chapter, instruction string) string {
	prompt := fmt.Sprintf(`你是一位专业的网络小说编辑，需要对以下章节内容进行润色。

章节标题：%s

润色要求：
1. 修正错别字和语法错误
2. 优化句子结构，使表达更流畅
3. 增强描写的生动性
4. 保持原有的情节和人物设定不变

`, chapter.Title)

	if instruction != "" {
		prompt += fmt.Sprintf("特别要求：%s\n\n", instruction)
	}

	prompt += `请以JSON格式输出润色建议，格式如下：
{
  "suggestions": [
    {
      "original": "原文片段",
      "revised": "修改后的片段",
      "reason": "修改原因"
    }
  ],
  "overall_score": 85,
  "summary": "整体评价"
}

请分析以下内容并给出润色建议：`

	return prompt
}
