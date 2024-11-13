package main

import (
	"context"
	"fmt"
	"go-rest-api-auth/config"
	"go-rest-api-auth/internal/api"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.DATABASE.Username, cfg.DATABASE.Password, cfg.DATABASE.Host, cfg.DATABASE.Port, cfg.DATABASE.DbName)
	redisUrl := fmt.Sprintf("redis://:%s@%s:%s/%v", cfg.REDIS.Password, cfg.REDIS.Host, cfg.REDIS.Port, cfg.REDIS.DbIndex)
	server := api.NewAPIServer(cfg.HTTPServer.Address, dsn, redisUrl)
	ctx := context.Background()

	done := make(chan bool)
	go func() {
		if err := server.Run(cfg, ctx); err != nil {
			slog.Error("Server error", slog.String("err", err.Error()))
		}
		done <- true
	}()

	//Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown: %v", err)
	}
	<-done
	slog.Info("Server exiting")

}
