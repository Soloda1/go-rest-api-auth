package api

import (
	"context"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "gocourse/cmd/main/docs"
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

// swagger:model
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

// Run Function
// @title API DOCUMENTATION
// @version 1.0
// @description This is a sample server.
// @host localhost:8000
// @BasePath /
func (s *APIServer) Run(cfg *config.Config, ctx context.Context) error {
	log := setupLogger(cfg.Env)

	router := http.NewServeMux()

	router.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	log.Info("Swagger enbaled", slog.String("url", s.address+"/swagger/"+"index.html"))

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
	//
	// @Summary JWT Login
	// @Description Login using JWT authentication
	// @Tags Auth
	// @Accept json
	// @Produce json
	// @Param request body jwtLogin.Request true "JWT login request"
	// @Success 200 {object} jwtLogin.Response
	// @Router /jwt_login [post]
	router.HandleFunc("POST /jwt_login", jwtLogin.New(log, TokenManager, UserService))

	// @Summary Refresh JWT
	// @Description Refresh JWT token
	// @Tags Auth
	// @Accept json
	// @Produce json
	// @Param request body refresh.Request true "JWT refresh request"
	// @Success 200 {object} refresh.Response
	// @Router /refresh [post]
	router.HandleFunc("POST /refresh", refresh.New(log, TokenManager))

	//Session auth
	// @Summary Session Login
	// @Description Login using session-based authentication
	// @Tags Auth
	// @Accept json
	// @Produce json
	// @Param request body sessionLogin.Request true "Session login request"
	// @Success 200 {object} sessionLogin.Response
	// @Router /session_login [post]
	router.HandleFunc("POST /session_login", sessionLogin.New(log, SessionManager, UserService))

	// @Summary Get User
	// @Description Get user by ID
	// @Tags Users
	// @Produce json
	// @Param userID path string true "User ID"
	// @Success 200 {object} getUser.Response
	// @Router /users/{userID} [get]
	router.HandleFunc("GET /users/{userID}", getUser.New(log, UserService))

	// @Summary Get All Users
	// @Description Get a list of all users
	// @Tags Users
	// @Produce json
	// @Success 200 {array} getAllUsers.Response
	// @Router /users [get]
	router.HandleFunc("GET /users", getAllUsers.New(log, UserService))

	// @Summary Create User
	// @Description Create a new user
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Param request body createUser.Request true "Create user request"
	// @Success 201 {object} createUser.Response
	// @Router /users [post]
	router.HandleFunc("POST /users", createUser.New(log, UserService))

	// @Summary Delete User
	// @Description Delete user by ID
	// @Tags Users
	// @Produce json
	// @Param userID path string true "User ID"
	// @Success 204
	// @Router /users/{userID} [delete]
	router.HandleFunc("DELETE /users/{userID}", deleteUser.New(log, UserService))

	// @Summary Update User
	// @Description Update user information by ID
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Param userID path string true "User ID"
	// @Param request body updateUser.Request true "Update user request"
	// @Success 200 {object} updateUser.Response
	// @Router /users/{userID} [put]
	router.HandleFunc("PUT /users/{userID}", updateUser.New(log, UserService))

	v1 := http.NewServeMux()
	v1MiddlewareStack := middleware.CreateStack(
		//middleware.TestAuthMiddleware(log),
		middleware.JWTAuthMiddleware(log, TokenManager),
	)

	// @Summary JWT Logout
	// @Description Logout and invalidate JWT token
	// @Tags Auth
	// @Produce json
	// @Success 200
	// @Router /v1/logout [get]
	v1.HandleFunc("GET /logout", jwtLogout.New(log, TokenManager))

	// @Summary Create Post
	// @Description Create a new post
	// @Tags Posts
	// @Accept json
	// @Produce json
	// @Param request body createPost.Request true "Create post request"
	// @Success 201 {object} createPost.Response
	// @Router /v1/posts [post]
	v1.HandleFunc("POST /posts", createPost.New(log, PostService))

	// @Summary Get All Posts
	// @Description Get a list of all posts
	// @Tags Posts
	// @Produce json
	// @Success 200 {array} getAllPosts.Response
	// @Router /v1/posts [get]
	v1.HandleFunc("GET /posts", getAllPosts.New(log, PostService))

	// @Summary Get Post
	// @Description Get post by ID
	// @Tags Posts
	// @Produce json
	// @Param postID path string true "Post ID"
	// @Success 200 {object} getPost.Response
	// @Router /v1/posts/{postID} [get]
	v1.HandleFunc("GET /posts/{postID}", getPost.New(log, PostService))

	// @Summary Update Post
	// @Description Update post information by ID
	// @Tags Posts
	// @Accept json
	// @Produce json
	// @Param postID path string true "Post ID"
	// @Param request body updatePost.Request true "Update post request"
	// @Success 200 {object} updatePost.Response
	// @Router /v1/posts/{postID} [put]
	v1.HandleFunc("PUT /posts/{postID}", updatePost.New(log, PostService))

	// @Summary Delete Post
	// @Description Delete post by ID
	// @Tags Posts
	// @Produce json
	// @Param postID path string true "Post ID"
	// @Success 204
	// @Router /v1/posts/{postID} [delete]
	v1.HandleFunc("DELETE /posts/{postID}", deletePost.New(log, PostService))

	v2 := http.NewServeMux()
	v2MiddlewareStack := middleware.CreateStack(
		middleware.SessionAuthMiddleware(log, SessionManager),
	)

	// @Summary Logout from session-based authentication
	// @Description Logs out the user by clearing the session stored in the cookie "session_id".
	// @Tags session_auth
	// @Produce json
	// @Success 200 {string} string "Successfully logged out"
	// @Router /v2/logout [get]
	v2.HandleFunc("GET /logout", sessionLogout.New(log, SessionManager))

	// @Summary Create a new post
	// @Description Create a post in the system with session-based authentication (requires "session_id" cookie).
	// @Tags posts
	// @Accept json
	// @Produce json
	// @Param post body createPost.Request true "Post details"
	// @Success 201 {object} createPost.Response "Post created successfully"
	// @Router /v2/posts [post]
	v2.HandleFunc("POST /posts", createPost.New(log, PostService))

	// @Summary Get all posts
	// @Description Retrieve all posts in the system with session-based authentication (requires "session_id" cookie).
	// @Tags posts
	// @Produce json
	// @Success 200 {array} getAllPosts.Response "List of posts"
	// @Router /v2/posts [get]
	v2.HandleFunc("GET /posts", getAllPosts.New(log, PostService))

	// @Summary Get post by ID
	// @Description Retrieve a specific post by its ID with session-based authentication (requires "session_id" cookie).
	// @Tags posts
	// @Produce json
	// @Param postID path string true "ID of the post"
	// @Success 200 {object} getPost.Response "Post details"
	// @Router /v2/posts/{postID} [get]
	v2.HandleFunc("GET /posts/{postID}", getPost.New(log, PostService))

	// @Summary Update a post by ID
	// @Description Update a specific post by its ID with session-based authentication (requires "session_id" cookie).
	// @Tags posts
	// @Accept json
	// @Produce json
	// @Param postID path string true "ID of the post"
	// @Param post body updatePost.Request true "Updated post details"
	// @Success 200 {object} updatePost.Response "Post updated successfully"
	// @Router /v2/posts/{postID} [put]
	v2.HandleFunc("PUT /posts/{postID}", updatePost.New(log, PostService))

	// @Summary Delete a post by ID
	// @Description Delete a specific post by its ID with session-based authentication (requires "session_id" cookie).
	// @Tags posts
	// @Produce json
	// @Param postID path string true "ID of the post"
	// @Success 200 {string} string "Post deleted successfully"
	// @Router /v2/posts/{postID} [delete]
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
