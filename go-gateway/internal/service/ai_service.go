package service

import (
	"context"

	"go-gateway/internal/config"
	"go-gateway/internal/grpcclient"
	"go-gateway/pkg/logger"
)

// 类型别名，方便外部使用
type GenerateRequest = grpcclient.GenerateRequest
type GenerateResponse = grpcclient.GenerateResponse
type Stream = grpcclient.Stream

// AIService AI 服务接口
type AIService interface {
	GenerateStream(ctx context.Context, req *GenerateRequest) (Stream, error)
}

type aiService struct {
	client *grpcclient.AIClient
	log    *logger.Logger
}

// NewAIService 创建 AI 服务实例
func NewAIService(cfg config.GRPCConfig, log *logger.Logger) AIService {
	return &aiService{
		client: grpcclient.NewAIClient(cfg, log),
		log:    log,
	}
}

func (s *aiService) GenerateStream(ctx context.Context, req *GenerateRequest) (Stream, error) {
	s.log.Debug("starting generate stream",
		logger.String("model", req.Model),
		logger.String("provider", req.Provider),
	)
	return s.client.GenerateStream(ctx, req)
}
