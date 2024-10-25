package api

import (
	"gocourse/config"
	"gocourse/handlers/createUser"
	"gocourse/handlers/deleteUser"
	"gocourse/handlers/getAllUsers"
	"gocourse/handlers/getUser"
	"gocourse/handlers/updateUser"
	"gocourse/middleware"
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
	dbURI   string
}

func NewAPIServer(address string) *APIServer {
	return &APIServer{address: address}
}

func (s *APIServer) Run(cfg *config.Config) error {
	log := setupLogger("local")

	router := http.NewServeMux()

	mainMiddlewareStack := middleware.CreateStack(
		middleware.RequestLoggerMiddleware(log),
	)

	router.HandleFunc("GET /users/{userID}", getUser.New(log))
	router.HandleFunc("GET /users", getAllUsers.New(log))
	router.HandleFunc("POST /users", createUser.New(log))
	router.HandleFunc("DELETE /users/{userID}", deleteUser.New(log))
	router.HandleFunc("PUT /users/{userID}", updateUser.New(log))

	v1 := http.NewServeMux()
	v1MiddlewareStack := middleware.CreateStack(
		middleware.RequireAuthMiddleware(log), // Мидлвейр для авторизации
	)
	v1.HandleFunc("POST /create-user", createUser.New(log))

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
