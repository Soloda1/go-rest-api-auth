package main

import (
	"gocourse/api"
	"gocourse/config"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()
	server := api.NewAPIServer(cfg.HTTPServer.Address)
	err := server.Run(cfg)
	if err != nil {
		slog.Error("Server error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
