package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// QualityAIClientInterface AI 客户端接口
type QualityAIClientInterface interface {
	AnalyzeReadingPower(ctx context.Context, chapterID, content, provider, model string) (map[string]interface{}, error)
	CheckConsistency(ctx context.Context, chapterID, content string, context map[string]interface{}, provider, model string) (map[string]interface{}, error)
	MultiAgentReview(ctx context.Context, chapterID, content, provider, model string) (map[string]interface{}, error)
	EvaluateQuality(ctx context.Context, chapterID, content string, context map[string]interface{}, provider, model string) (map[string]interface{}, error)
}

// QualityHandler 质量评估处理器
type QualityHandler struct {
	aiClient QualityAIClientInterface
	logger   *zap.Logger
}

// NewQualityHandler 创建质量评估处理器
func NewQualityHandler(aiClient QualityAIClientInterface, logger *zap.Logger) *QualityHandler {
	return &QualityHandler{
		aiClient: aiClient,
		logger:   logger,
	}
}

// AnalyzeReadingPowerRequest 追读力分析请求
type AnalyzeReadingPowerRequest struct {
	Content  string `json:"content" binding:"required"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
}

// CheckConsistencyRequest 一致性检查请求
type CheckConsistencyRequest struct {
	Content         string                   `json:"content" binding:"required"`
	Characters      []map[string]interface{} `json:"characters"`
	WorldSettings   []map[string]interface{} `json:"world_settings"`
	PreviousSummary string                   `json:"previous_summary"`
	Provider        string                   `json:"provider"`
	Model           string                   `json:"model"`
}

// MultiAgentReviewRequest 多Agent审查请求
type MultiAgentReviewRequest struct {
	Content  string `json:"content" binding:"required"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
}

// EvaluateQualityRequest 综合质量评估请求
type EvaluateQualityRequest struct {
	Content         string                   `json:"content" binding:"required"`
	Characters      []map[string]interface{} `json:"characters"`
	WorldSettings   []map[string]interface{} `json:"world_settings"`
	PreviousSummary string                   `json:"previous_summary"`
	Provider        string                   `json:"provider"`
	Model           string                   `json:"model"`
}

// AnalyzeReadingPower 追读力分析
// POST /api/v1/chapters/:id/analyze-reading-power
func (h *QualityHandler) AnalyzeReadingPower(c *gin.Context) {
	chapterID := c.Param("id")

	var req AnalyzeReadingPowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Analyzing reading power",
		zap.String("chapter_id", chapterID),
		zap.Int("content_length", len(req.Content)),
	)

	// 调用 Python AI 服务
	result, err := h.aiClient.AnalyzeReadingPower(c.Request.Context(), chapterID, req.Content, req.Provider, req.Model)
	if err != nil {
		h.logger.Error("Failed to analyze reading power", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "分析失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CheckConsistency 一致性检查
// POST /api/v1/chapters/:id/check-consistency
func (h *QualityHandler) CheckConsistency(c *gin.Context) {
	chapterID := c.Param("id")

	var req CheckConsistencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Checking consistency",
		zap.String("chapter_id", chapterID),
		zap.Int("content_length", len(req.Content)),
	)

	// 构建上下文
	context := map[string]interface{}{
		"characters":       req.Characters,
		"world_settings":   req.WorldSettings,
		"previous_summary": req.PreviousSummary,
	}

	// 调用 Python AI 服务
	result, err := h.aiClient.CheckConsistency(c.Request.Context(), chapterID, req.Content, context, req.Provider, req.Model)
	if err != nil {
		h.logger.Error("Failed to check consistency", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// MultiAgentReview 多Agent审查
// POST /api/v1/chapters/:id/multi-agent-review
func (h *QualityHandler) MultiAgentReview(c *gin.Context) {
	chapterID := c.Param("id")

	var req MultiAgentReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Starting multi-agent review",
		zap.String("chapter_id", chapterID),
		zap.Int("content_length", len(req.Content)),
	)

	// 调用 Python AI 服务
	result, err := h.aiClient.MultiAgentReview(c.Request.Context(), chapterID, req.Content, req.Provider, req.Model)
	if err != nil {
		h.logger.Error("Failed to multi-agent review", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "审查失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// EvaluateQuality 综合质量评估
// POST /api/v1/chapters/:id/evaluate-quality
func (h *QualityHandler) EvaluateQuality(c *gin.Context) {
	chapterID := c.Param("id")

	var req EvaluateQualityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Evaluating quality",
		zap.String("chapter_id", chapterID),
		zap.Int("content_length", len(req.Content)),
	)

	// 构建上下文
	context := map[string]interface{}{
		"characters":       req.Characters,
		"world_settings":   req.WorldSettings,
		"previous_summary": req.PreviousSummary,
	}

	// 调用 Python AI 服务
	result, err := h.aiClient.EvaluateQuality(c.Request.Context(), chapterID, req.Content, context, req.Provider, req.Model)
	if err != nil {
		h.logger.Error("Failed to evaluate quality", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评估失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}
