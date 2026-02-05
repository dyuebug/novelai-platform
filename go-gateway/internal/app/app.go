package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/config"
	"go-gateway/internal/database"
	"go-gateway/internal/router"
	"go-gateway/internal/service"
	"go-gateway/pkg/logger"
)

type Server struct {
	engine *gin.Engine
	server *http.Server
}

func NewServer(cfg *config.Config, log *logger.Logger) *Server {
	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()

	// 初始化数据库连接
	if err := database.Connect(cfg.Database); err != nil {
		log.Fatal("failed to connect database", logger.Error(err))
	}

	// 初始化 EmailService
	emailService := service.NewEmailService(&cfg.Email)

	// 初始化 AuthService
	authService, err := service.NewAuthService(cfg.Auth, emailService)
	if err != nil {
		log.Fatal("failed to initialize auth service", logger.Error(err))
	}

	router.Register(engine, cfg, log, authService)

	return &Server{
		engine: engine,
		server: &http.Server{
			Addr:    cfg.Server.Addr,
			Handler: engine,
		},
	}
}

func (s *Server) Run(addr string) error {
	s.server.Addr = addr
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
