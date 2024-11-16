package jwtLogout

import (
	"go-rest-api-auth/internal/database/auth"
	"go-rest-api-auth/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

// Response represents the jwt logout response payload.
// swagger:model
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func New(log *slog.Logger, tokenManager auth.JwtManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("JWT  Logout user")

		userIDValue := r.Context().Value("user_id")
		if userIDValue == nil {
			log.Error("user_id not found in context")
			utils.SendError(w, "Invalid user context")
			return
		}

		userID, err := strconv.Atoi(userIDValue.(string))
		if err != nil {
			log.Error("Invalid user ID format")
			utils.SendError(w, "Invalid user ID")
			return
		}

		err = tokenManager.DeleteRefreshToken(userID)
		if err != nil {
			log.Error("Error deleting refresh token")
			utils.SendError(w, "Error deleting refresh token")
			return
		}

		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
		})
	}
}
