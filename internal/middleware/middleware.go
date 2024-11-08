package middleware

import (
	"context"
	"gocourse/internal/database/auth"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

func RequestLoggerMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(slog.String("component", "middleware/RequestLoggerMiddleware"))
		log.Info("Logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
			)

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.String("duration", time.Since(t1).String()))
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func TestAuthMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(slog.String("component", "middleware/TestAuthMiddleware"))
		log.Info("TEST Auth middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			log.Debug("Token", slog.String("token", token))
			if token != "Bearer token" {
				utils.SendError(w, "Unauthorized token")
				return
			}

			// Ищем куку с user_id
			cookie, err := r.Cookie("user_id")
			if err != nil {
				utils.SendError(w, "Unauthorized cookies")
				return
			}
			//log.Debug("Cookie", slog.Any("cookie", cookie))

			// Преобразуем значение куки в int
			userID, err := strconv.Atoi(cookie.Value)
			if err != nil {
				utils.SendError(w, "Unauthorized cookies user id")
				return
			}
			//log.Debug("UserID", slog.Any("userID", userID))

			// Добавляем user_id в контекст
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func JWTAuthMiddleware(log *slog.Logger, tokenManager *auth.JwtManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(slog.String("component", "middleware/JWTAuthMiddleware"))
		log.Info("JWT Auth middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				utils.SendError(w, "Missing authorization token")
				return
			}

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			accessTokenClaims, err := tokenManager.ValidateJWT(tokenString, "access")
			if err != nil {
				utils.SendError(w, err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), "userID", accessTokenClaims["sub"].(string))
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
