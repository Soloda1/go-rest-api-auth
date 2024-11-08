package logout

import (
	"gocourse/internal/database/auth"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger, tokenManager *auth.JwtManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Logout user")

	}
}
