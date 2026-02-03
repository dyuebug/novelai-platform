package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/grpcclient"
	"go.uber.org/zap"
)

// ConstraintHandler 约束检查处理器
type ConstraintHandler struct {
	aiClient *grpcclient.AIClient
	logger   *zap.Logger
}

// NewConstraintHandler 创建约束检查处理器
func NewConstraintHandler(aiClient *grpcclient.AIClient, logger *zap.Logger) *ConstraintHandler {
	return &ConstraintHandler{
		aiClient: aiClient,
		logger:   logger,
	}
}

// CheckConstraintsRequest 约束检查请求
type CheckConstraintsRequest struct {
	Content     string                   `json:"content" binding:"required"`
	Constraints []map[string]interface{} `json:"constraints"`
}

// RequestExemptionRequest 请求豁免
type RequestExemptionRequest struct {
	ConstraintID string `json:"constraint_id" binding:"required"`
	Reason       string `json:"reason" binding:"required"`
}

// CheckConstraints 检查内容是否违反宪法约束
// POST /api/v1/projects/:projectId/chapters/:chapterId/check-constraints
func (h *ConstraintHandler) CheckConstraints(c *gin.Context) {
	projectID := c.Param("projectId")
	chapterID := c.Param("chapterId")

	var req CheckConstraintsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Checking constraints",
		zap.String("project_id", projectID),
		zap.String("chapter_id", chapterID),
		zap.Int("content_length", len(req.Content)),
		zap.Int("constraints_count", len(req.Constraints)),
	)

	// 调用 Python AI 服务
	grpcReq := &grpcclient.CheckConstraintsRequest{
		ProjectID:   projectID,
		ChapterID:   chapterID,
		Content:     req.Content,
		Constraints: req.Constraints,
	}

	result, err := h.aiClient.CheckConstraints(c.Request.Context(), grpcReq)
	if err != nil {
		h.logger.Error("Failed to check constraints", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "约束检查失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// RequestExemption 请求约束豁免
// POST /api/v1/projects/:projectId/chapters/:chapterId/exemptions
func (h *ConstraintHandler) RequestExemption(c *gin.Context) {
	projectID := c.Param("projectId")
	chapterID := c.Param("chapterId")
	userID := c.GetString("user_id") // 从认证中间件获取

	var req RequestExemptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Requesting exemption",
		zap.String("project_id", projectID),
		zap.String("chapter_id", chapterID),
		zap.String("constraint_id", req.ConstraintID),
		zap.String("user_id", userID),
	)

	// 调用 Python AI 服务
	grpcReq := &grpcclient.RequestExemptionRequest{
		ProjectID:    projectID,
		ChapterID:    chapterID,
		ConstraintID: req.ConstraintID,
		Reason:       req.Reason,
		RequestedBy:  userID,
	}

	result, err := h.aiClient.RequestExemption(c.Request.Context(), grpcReq)
	if err != nil {
		h.logger.Error("Failed to request exemption", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "请求豁免失败"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// RevokeExemption 撤销约束豁免
// DELETE /api/v1/exemptions/:exemptionId
func (h *ConstraintHandler) RevokeExemption(c *gin.Context) {
	exemptionID := c.Param("exemptionId")

	h.logger.Info("Revoking exemption",
		zap.String("exemption_id", exemptionID),
	)

	// 调用 Python AI 服务
	grpcReq := &grpcclient.RevokeExemptionRequest{
		ExemptionID: exemptionID,
	}

	err := h.aiClient.RevokeExemption(c.Request.Context(), grpcReq)
	if err != nil {
		h.logger.Error("Failed to revoke exemption", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "撤销豁免失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "豁免已撤销"})
}
