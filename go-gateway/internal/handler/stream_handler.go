package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/grpcclient"
	"go-gateway/internal/service"
)

// StreamInterface 流式响应接口
type StreamInterface interface {
	Recv() (*grpcclient.GenerateResponse, error)
	CloseSend() error
}

// AIServiceInterface AI 服务接口
type AIServiceInterface interface {
	GenerateStream(ctx context.Context, req *service.GenerateRequest) (service.Stream, error)
}

type StreamHandler struct {
	ai AIServiceInterface
}

func NewStreamHandler(ai AIServiceInterface) *StreamHandler {
	return &StreamHandler{ai: ai}
}

func (h *StreamHandler) Generate(c *gin.Context) {
	req := &grpcclient.GenerateRequest{
		Prompt:       c.Query("prompt"),
		Model:        c.DefaultQuery("model", "gpt-4o"),
		Provider:     c.DefaultQuery("provider", "openai"),
		SystemPrompt: c.Query("system_prompt"),
		UserID:       c.GetString("user_id"),
		ProjectID:    c.Query("project_id"),
	}

	if req.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prompt is required"})
		return
	}

	stream, err := h.ai.GenerateStream(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": fmt.Sprintf("AI service unavailable: %v", err)})
		return
	}
	defer stream.CloseSend()

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	c.Stream(func(w io.Writer) bool {
		resp, err := stream.Recv()
		if err == io.EOF {
			// 发送完成事件
			c.SSEvent("done", gin.H{"status": "completed"})
			return false
		}
		if err != nil {
			c.SSEvent("error", gin.H{"message": err.Error()})
			return false
		}
		if resp.Error != "" {
			c.SSEvent("error", gin.H{"message": resp.Error})
			return false
		}

		c.SSEvent("message", gin.H{"content": resp.Content})

		if resp.Done {
			c.SSEvent("done", gin.H{"status": "completed"})
			return false
		}
		return true
	})
}
