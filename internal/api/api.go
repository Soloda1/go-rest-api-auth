package api

import (
	"context"
	"gocourse/config"
	"gocourse/internal/database"
	"gocourse/internal/database/auth"
	jwtLogin "gocourse/internal/handlers/auth/jwt/login"
	jwtLogout "gocourse/internal/handlers/auth/jwt/logout"
	"gocourse/internal/handlers/auth/jwt/refresh"
	sessionLogin "gocourse/internal/handlers/auth/session/login"
	sessionLogout "gocourse/internal/handlers/auth/session/logout"
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
	address  string
	dbUrl    string
	redisUrl string
	server   *http.Server
}

func NewAPIServer(address string, dbUrl string, redisUrl string) *APIServer {
	return &APIServer{
		address:  address,
		dbUrl:    dbUrl,
		redisUrl: redisUrl,
	}
}

func (s *APIServer) Run(cfg *config.Config, ctx context.Context) error {
	log := setupLogger(cfg.Env)

	router := http.NewServeMux()

	storage := database.NewDbPool(ctx, s.dbUrl, log)
	log.Info("Database connected")
	defer func() {
		slog.Info("Database disconnected")
		storage.Close()
	}()

	cache := database.NewRedisClient(ctx, log, s.redisUrl)
	log.Info("Redis connected")
	defer func() {
		err := cache.Cache.Close()
		if err != nil {
			slog.Error("Failed to close Redis connection")
		}
		slog.Info("Redis disconnected")
	}()

	TagsService := database.NewTagService(storage)
	PostService := database.NewPostService(storage, TagsService)
	UserService := database.NewUserService(storage)
	TokenManager := auth.NewJwtManager(cfg, storage)
	SessionManager := auth.NewSessionManager(cache, cfg.REDIS.TTL)

	mainMiddlewareStack := middleware.CreateStack(
		middleware.RequestLoggerMiddleware(log),
	)

	//JWT auth
	router.HandleFunc("POST /jwt_login", jwtLogin.New(log, TokenManager, UserService))
	router.HandleFunc("POST /refresh", refresh.New(log, TokenManager))

	//Session auth
	router.HandleFunc("POST /session_login", sessionLogin.New(log, SessionManager, UserService))

	router.HandleFunc("GET /users/{userID}", getUser.New(log, UserService))
	router.HandleFunc("GET /users", getAllUsers.New(log, UserService))
	router.HandleFunc("POST /users", createUser.New(log, UserService))
	router.HandleFunc("DELETE /users/{userID}", deleteUser.New(log, UserService))
	router.HandleFunc("PUT /users/{userID}", updateUser.New(log, UserService))

	v1 := http.NewServeMux()
	v1MiddlewareStack := middleware.CreateStack(
		//middleware.TestAuthMiddleware(log),
		middleware.JWTAuthMiddleware(log, TokenManager),
	)

	v1.HandleFunc("GET /logout", jwtLogout.New(log, TokenManager))

	v1.HandleFunc("POST /posts", createPost.New(log, PostService))
	v1.HandleFunc("GET /posts", getAllPosts.New(log, PostService))
	v1.HandleFunc("GET /posts/{postID}", getPost.New(log, PostService))
	v1.HandleFunc("PUT /posts/{postID}", updatePost.New(log, PostService))
	v1.HandleFunc("DELETE /posts/{postID}", deletePost.New(log, PostService))

	v2 := http.NewServeMux()
	v2MiddlewareStack := middleware.CreateStack(
		middleware.SessionAuthMiddleware(log, SessionManager),
	)

	v2.HandleFunc("GET /logout", sessionLogout.New(log, SessionManager))

	v2.HandleFunc("POST /posts", createPost.New(log, PostService))
	v2.HandleFunc("GET /posts", getAllPosts.New(log, PostService))
	v2.HandleFunc("GET /posts/{postID}", getPost.New(log, PostService))
	v2.HandleFunc("PUT /posts/{postID}", updatePost.New(log, PostService))
	v2.HandleFunc("DELETE /posts/{postID}", deletePost.New(log, PostService))

	router.Handle("/v1/", http.StripPrefix("/v1", v1MiddlewareStack(v1)))
	router.Handle("/v2/", http.StripPrefix("/v2", v2MiddlewareStack(v2)))

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      mainMiddlewareStack(router),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("Server has started: ", slog.String("address", s.address))
	log.Debug("debug logger enabled")

	return s.server.ListenAndServe()
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
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
