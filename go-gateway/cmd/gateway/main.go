package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-gateway/internal/app"
	"go-gateway/internal/config"
	"go-gateway/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Log.Level)
	defer log.Sync()

	server := app.NewServer(cfg, log)

	go func() {
		log.Info("server starting", logger.String("addr", cfg.Server.Addr))
		if err := server.Run(cfg.Server.Addr); err != nil {
			log.Fatal("server run failed", logger.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server shutdown error", logger.Error(err))
	}
	log.Info("server exited")
}
