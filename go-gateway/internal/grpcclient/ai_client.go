package grpcclient

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/sony/gobreaker"
	"go-gateway/internal/config"
	"go-gateway/pkg/logger"
	"google.golang.org/grpc"
)

// GenerateRequest AI 生成请求
type GenerateRequest struct {
	Prompt       string  `json:"prompt"`
	Model        string  `json:"model"`
	Provider     string  `json:"provider"`
	Temperature  float32 `json:"temperature"`
	MaxTokens    int32   `json:"max_tokens"`
	SystemPrompt string  `json:"system_prompt"`
	UserID       string  `json:"user_id"`
	ProjectID    string  `json:"project_id"`
	ChapterID    string  `json:"chapter_id"`
	TaskType     string  `json:"task_type"` // generate, partial_regenerate, polish
	Context      string  `json:"context"`   // 上下文信息 (大纲、前文等)
	Selection    string  `json:"selection"` // 选中的文本 (用于局部重写)
	Instruction  string  `json:"instruction"` // 用户指令
}

// GenerateResponse AI 生成响应
type GenerateResponse struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
	Error   string `json:"error"`
}

// Stream 流式响应接口
type Stream interface {
	Recv() (*GenerateResponse, error)
	CloseSend() error
}

// AIClient AI 服务 gRPC 客户端
type AIClient struct {
	pool    *ConnPool
	retry   RetryPolicy
	breaker *gobreaker.CircuitBreaker
	log     *logger.Logger
	timeout time.Duration
}

// NewAIClient 创建 AI 客户端
func NewAIClient(cfg config.GRPCConfig, log *logger.Logger) *AIClient {
	return &AIClient{
		pool:    NewConnPool(cfg.AIServiceAddr, cfg.MaxConns),
		retry:   RetryPolicy{MaxRetries: cfg.RetryMax, Backoff: cfg.RetryBackoff},
		breaker: NewBreaker("ai-generate-stream"),
		log:     log,
		timeout: cfg.Timeout,
	}
}

// GenerateStream 流式生成
func (c *AIClient) GenerateStream(ctx context.Context, req *GenerateRequest) (Stream, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		_ = cancel // 由调用方通过 stream.CloseSend() 管理
	}

	var lastErr error
	for i := 0; i <= c.retry.MaxRetries; i++ {
		res, err := c.breaker.Execute(func() (interface{}, error) {
			conn, release, err := c.pool.Acquire(ctx)
			if err != nil {
				return nil, err
			}
			stream, err := newAIStream(ctx, conn, req, release)
			if err != nil {
				release()
				return nil, err
			}
			return stream, nil
		})
		if err == nil {
			return res.(Stream), nil
		}
		lastErr = err
		c.log.Warn("grpc call failed, retrying",
			logger.Int("attempt", i+1),
			logger.Error(err),
		)
		time.Sleep(c.retry.nextBackoff(i + 1))
	}
	if lastErr == nil {
		lastErr = errors.New("generate stream failed")
	}
	return nil, lastErr
}

// aiStream 内部流实现
type aiStream struct {
	grpc.ClientStream
	release func()
}

func newAIStream(ctx context.Context, conn *grpc.ClientConn, req *GenerateRequest, release func()) (Stream, error) {
	desc := &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: false,
	}
	stream, err := conn.NewStream(ctx, desc, "/ai.AIService/GenerateStream")
	if err != nil {
		return nil, err
	}
	if err := stream.SendMsg(req); err != nil {
		return nil, err
	}
	if err := stream.CloseSend(); err != nil {
		return nil, err
	}
	return &aiStream{ClientStream: stream, release: release}, nil
}

func (s *aiStream) Recv() (*GenerateResponse, error) {
	resp := &GenerateResponse{}
	if err := s.ClientStream.RecvMsg(resp); err != nil {
		if err == io.EOF {
			s.release()
		}
		return nil, err
	}
	return resp, nil
}

func (s *aiStream) CloseSend() error {
	s.release()
	return nil
}

// QualityRequest 质量评估请求基础结构
type QualityRequest struct {
	ChapterID       string                   `json:"chapter_id"`
	Content         string                   `json:"content"`
	Provider        string                   `json:"provider"`
	Model           string                   `json:"model"`
	Characters      []map[string]interface{} `json:"characters,omitempty"`
	WorldSettings   []map[string]interface{} `json:"world_settings,omitempty"`
	PreviousSummary string                   `json:"previous_summary,omitempty"`
}

// AnalyzeReadingPower 追读力分析
func (c *AIClient) AnalyzeReadingPower(ctx context.Context, chapterID, content, provider, model string) (map[string]interface{}, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	req := &QualityRequest{
		ChapterID: chapterID,
		Content:   content,
		Provider:  provider,
		Model:     model,
	}

	var result map[string]interface{}
	var lastErr error

	for i := 0; i <= c.retry.MaxRetries; i++ {
		res, err := c.breaker.Execute(func() (interface{}, error) {
			conn, release, err := c.pool.Acquire(ctx)
			if err != nil {
				return nil, err
			}
			defer release()

			resp := make(map[string]interface{})
			err = conn.Invoke(ctx, "/ai.AIService/AnalyzeReadingPower", req, &resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		})
		if err == nil {
			result = res.(map[string]interface{})
			return result, nil
		}
		lastErr = err
		c.log.Warn("AnalyzeReadingPower failed, retrying",
			logger.Int("attempt", i+1),
			logger.Error(err),
		)
		time.Sleep(c.retry.nextBackoff(i + 1))
	}

	return nil, lastErr
}

// CheckConsistency 一致性检查
func (c *AIClient) CheckConsistency(ctx context.Context, chapterID, content string, contextData map[string]interface{}, provider, model string) (map[string]interface{}, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	// 从 contextData 中提取字段
	var characters []map[string]interface{}
	var worldSettings []map[string]interface{}
	var previousSummary string

	if chars, ok := contextData["characters"].([]map[string]interface{}); ok {
		characters = chars
	}
	if settings, ok := contextData["world_settings"].([]map[string]interface{}); ok {
		worldSettings = settings
	}
	if summary, ok := contextData["previous_summary"].(string); ok {
		previousSummary = summary
	}

	req := &QualityRequest{
		ChapterID:       chapterID,
		Content:         content,
		Provider:        provider,
		Model:           model,
		Characters:      characters,
		WorldSettings:   worldSettings,
		PreviousSummary: previousSummary,
	}

	var result map[string]interface{}
	var lastErr error

	for i := 0; i <= c.retry.MaxRetries; i++ {
		res, err := c.breaker.Execute(func() (interface{}, error) {
			conn, release, err := c.pool.Acquire(ctx)
			if err != nil {
				return nil, err
			}
			defer release()

			resp := make(map[string]interface{})
			err = conn.Invoke(ctx, "/ai.AIService/CheckConsistency", req, &resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		})
		if err == nil {
			result = res.(map[string]interface{})
			return result, nil
		}
		lastErr = err
		c.log.Warn("CheckConsistency failed, retrying",
			logger.Int("attempt", i+1),
			logger.Error(err),
		)
		time.Sleep(c.retry.nextBackoff(i + 1))
	}

	return nil, lastErr
}

// MultiAgentReview 多Agent审查
func (c *AIClient) MultiAgentReview(ctx context.Context, chapterID, content, provider, model string) (map[string]interface{}, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	req := &QualityRequest{
		ChapterID: chapterID,
		Content:   content,
		Provider:  provider,
		Model:     model,
	}

	var result map[string]interface{}
	var lastErr error

	for i := 0; i <= c.retry.MaxRetries; i++ {
		res, err := c.breaker.Execute(func() (interface{}, error) {
			conn, release, err := c.pool.Acquire(ctx)
			if err != nil {
				return nil, err
			}
			defer release()

			resp := make(map[string]interface{})
			err = conn.Invoke(ctx, "/ai.AIService/MultiAgentReview", req, &resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		})
		if err == nil {
			result = res.(map[string]interface{})
			return result, nil
		}
		lastErr = err
		c.log.Warn("MultiAgentReview failed, retrying",
			logger.Int("attempt", i+1),
			logger.Error(err),
		)
		time.Sleep(c.retry.nextBackoff(i + 1))
	}

	return nil, lastErr
}

// EvaluateQuality 综合质量评估
func (c *AIClient) EvaluateQuality(ctx context.Context, chapterID, content string, contextData map[string]interface{}, provider, model string) (map[string]interface{}, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	// 从 contextData 中提取字段
	var characters []map[string]interface{}
	var worldSettings []map[string]interface{}
	var previousSummary string

	if chars, ok := contextData["characters"].([]map[string]interface{}); ok {
		characters = chars
	}
	if settings, ok := contextData["world_settings"].([]map[string]interface{}); ok {
		worldSettings = settings
	}
	if summary, ok := contextData["previous_summary"].(string); ok {
		previousSummary = summary
	}

	req := &QualityRequest{
		ChapterID:       chapterID,
		Content:         content,
		Provider:        provider,
		Model:           model,
		Characters:      characters,
		WorldSettings:   worldSettings,
		PreviousSummary: previousSummary,
	}

	var result map[string]interface{}
	var lastErr error

	for i := 0; i <= c.retry.MaxRetries; i++ {
		res, err := c.breaker.Execute(func() (interface{}, error) {
			conn, release, err := c.pool.Acquire(ctx)
			if err != nil {
				return nil, err
			}
			defer release()

			resp := make(map[string]interface{})
			err = conn.Invoke(ctx, "/ai.AIService/EvaluateQuality", req, &resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		})
		if err == nil {
			result = res.(map[string]interface{})
			return result, nil
		}
		lastErr = err
		c.log.Warn("EvaluateQuality failed, retrying",
			logger.Int("attempt", i+1),
			logger.Error(err),
		)
		time.Sleep(c.retry.nextBackoff(i + 1))
	}

	return nil, lastErr
}

// CheckConstraintsRequest 约束检查请求
type CheckConstraintsRequest struct {
	ProjectID   string                   `json:"project_id"`
	ChapterID   string                   `json:"chapter_id"`
	Content     string                   `json:"content"`
	Constraints []map[string]interface{} `json:"constraints"`
}

// CheckConstraintsResponse 约束检查响应
type CheckConstraintsResponse struct {
	Violations []ConstraintViolation `json:"violations"`
	Passed     bool                  `json:"passed"`
}

// ConstraintViolation 约束违规
type ConstraintViolation struct {
	ConstraintID   string `json:"constraint_id"`
	ConstraintType string `json:"constraint_type"` // HARD or SOFT
	Description    string `json:"description"`
	Severity       string `json:"severity"`
	Suggestion     string `json:"suggestion"`
	CanExempt      bool   `json:"can_exempt"`
}

// CheckConstraints 检查内容是否违反宪法约束
func (c *AIClient) CheckConstraints(ctx context.Context, req *CheckConstraintsRequest) (*CheckConstraintsResponse, error) {
	var result *CheckConstraintsResponse
	var lastErr error

	for i := 0; i <= c.retry.MaxRetries; i++ {
		res, err := c.breaker.Execute(func() (interface{}, error) {
			conn, release, err := c.pool.Acquire(ctx)
			if err != nil {
				return nil, err
			}
			defer release()

			resp := &CheckConstraintsResponse{}
			err = conn.Invoke(ctx, "/ai.AIService/CheckConstraints", req, resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		})
		if err == nil {
			result = res.(*CheckConstraintsResponse)
			return result, nil
		}
		lastErr = err
		c.log.Warn("CheckConstraints failed, retrying",
			logger.Int("attempt", i+1),
			logger.Error(err),
		)
		time.Sleep(c.retry.nextBackoff(i + 1))
	}

	return nil, lastErr
}

// RequestExemptionRequest 请求豁免
type RequestExemptionRequest struct {
	ProjectID    string `json:"project_id"`
	ChapterID    string `json:"chapter_id"`
	ConstraintID string `json:"constraint_id"`
	Reason       string `json:"reason"`
	RequestedBy  string `json:"requested_by"`
}

// ExemptionResponse 豁免响应
type ExemptionResponse struct {
	ID           string `json:"id"`
	ConstraintID string `json:"constraint_id"`
	ChapterID    string `json:"chapter_id"`
	Reason       string `json:"reason"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

// RequestExemption 请求约束豁免
func (c *AIClient) RequestExemption(ctx context.Context, req *RequestExemptionRequest) (*ExemptionResponse, error) {
	var result *ExemptionResponse
	var lastErr error

	for i := 0; i <= c.retry.MaxRetries; i++ {
		res, err := c.breaker.Execute(func() (interface{}, error) {
			conn, release, err := c.pool.Acquire(ctx)
			if err != nil {
				return nil, err
			}
			defer release()

			resp := &ExemptionResponse{}
			err = conn.Invoke(ctx, "/ai.AIService/RequestExemption", req, resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		})
		if err == nil {
			result = res.(*ExemptionResponse)
			return result, nil
		}
		lastErr = err
		c.log.Warn("RequestExemption failed, retrying",
			logger.Int("attempt", i+1),
			logger.Error(err),
		)
		time.Sleep(c.retry.nextBackoff(i + 1))
	}

	return nil, lastErr
}

// RevokeExemptionRequest 撤销豁免请求
type RevokeExemptionRequest struct {
	ExemptionID string `json:"exemption_id"`
}

// RevokeExemption 撤销约束豁免
func (c *AIClient) RevokeExemption(ctx context.Context, req *RevokeExemptionRequest) error {
	var lastErr error

	for i := 0; i <= c.retry.MaxRetries; i++ {
		_, err := c.breaker.Execute(func() (interface{}, error) {
			conn, release, err := c.pool.Acquire(ctx)
			if err != nil {
				return nil, err
			}
			defer release()

			resp := make(map[string]interface{})
			err = conn.Invoke(ctx, "/ai.AIService/RevokeExemption", req, &resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		})
		if err == nil {
			return nil
		}
		lastErr = err
		c.log.Warn("RevokeExemption failed, retrying",
			logger.Int("attempt", i+1),
			logger.Error(err),
		)
		time.Sleep(c.retry.nextBackoff(i + 1))
	}

	return lastErr
}
