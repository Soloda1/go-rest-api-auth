package main

import (
	"fmt"
	"gocourse/config"
	"gocourse/internal/api"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.DATABASE.Username, cfg.DATABASE.Password, cfg.DATABASE.Host, cfg.DATABASE.Port, cfg.DATABASE.DbName)
	server := api.NewAPIServer(cfg.HTTPServer.Address, dsn)
	err := server.Run(cfg)
	if err != nil {
		slog.Error("Server error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
