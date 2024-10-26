package api

import (
	"context"
	"gocourse/config"
	"gocourse/internal/database"
	"gocourse/internal/handlers/user/createUser"
	"gocourse/internal/handlers/user/deleteUser"
	"gocourse/internal/handlers/user/getAllUsers"
	"gocourse/internal/handlers/user/getUser"
	"gocourse/internal/handlers/user/updateUser"
	"gocourse/internal/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	endDev   = "dev"
	endProd  = "prod"
)

type APIServer struct {
	address string
	dbUrl   string
}

func NewAPIServer(address string, dbUrl string) *APIServer {
	return &APIServer{address: address, dbUrl: dbUrl}
}

func (s *APIServer) Run(cfg *config.Config) error {
	log := setupLogger(cfg.Env)

	router := http.NewServeMux()

	storage := database.NewDbPool(context.Background(), s.dbUrl, log)
	log.Info("Database connected")
	defer log.Info("Database disconnected")
	defer storage.Close()

	mainMiddlewareStack := middleware.CreateStack(
		middleware.RequestLoggerMiddleware(log),
	)

	router.HandleFunc("GET /users/{userID}", getUser.New(log, storage))
	router.HandleFunc("GET /users", getAllUsers.New(log, storage))
	router.HandleFunc("POST /users", createUser.New(log, storage))
	router.HandleFunc("DELETE /users/{userID}", deleteUser.New(log, storage))
	router.HandleFunc("PUT /users/{userID}", updateUser.New(log, storage))

	v1 := http.NewServeMux()
	v1MiddlewareStack := middleware.CreateStack(
		middleware.RequireAuthMiddleware(log), // Мидлвейр для авторизации
	)
	v1.HandleFunc("POST /create-user", createUser.New(log, storage))

	router.Handle("/v1/", http.StripPrefix("/v1", v1MiddlewareStack(v1)))

	server := http.Server{
		Addr:         s.address,
		Handler:      mainMiddlewareStack(router),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("Server has started: ", slog.String("address", s.address))
	log.Debug("debug logger enabled")

	return server.ListenAndServe()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case endDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case endProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
