package api

import (
	"context"
	"gocourse/config"
	"gocourse/internal/database"
	"gocourse/internal/handlers/post/createPost"
	"gocourse/internal/handlers/post/deletePost"
	"gocourse/internal/handlers/post/getAllPosts"
	"gocourse/internal/handlers/post/getPost"
	"gocourse/internal/handlers/post/updatePost"
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
	TagsService := database.NewTagService(storage)
	PostService := database.NewPostService(storage, TagsService)
	UserService := database.NewUserService(storage)

	log.Info("Database connected")
	defer log.Info("Database disconnected")
	defer storage.Close()

	mainMiddlewareStack := middleware.CreateStack(
		middleware.RequestLoggerMiddleware(log),
	)

	router.HandleFunc("GET /users/{userID}", getUser.New(log, UserService))
	router.HandleFunc("GET /users", getAllUsers.New(log, UserService))
	router.HandleFunc("POST /users", createUser.New(log, UserService))
	router.HandleFunc("DELETE /users/{userID}", deleteUser.New(log, UserService))
	router.HandleFunc("PUT /users/{userID}", updateUser.New(log, UserService))

	v1 := http.NewServeMux()
	v1MiddlewareStack := middleware.CreateStack(
		middleware.RequireAuthMiddleware(log), // Мидлвейр для авторизации
	)

	v1.HandleFunc("POST /posts", createPost.New(log, PostService))
	v1.HandleFunc("GET /posts", getAllPosts.New(log, PostService))
	v1.HandleFunc("GET /posts/{postID}", getPost.New(log, PostService))
	v1.HandleFunc("PUT /posts/{postID}", updatePost.New(log, PostService))
	v1.HandleFunc("DELETE /posts/{postID}", deletePost.New(log, PostService))

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
