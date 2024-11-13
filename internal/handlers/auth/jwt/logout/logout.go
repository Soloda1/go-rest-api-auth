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

		userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
		if err != nil {
			log.Error("Invalid user ID")
			utils.SendError(w, err.Error())
			return
		}

		err = tokenManager.DeleteRefreshToken(userID)
		if err != nil {
			log.Error("Error deleting refresh token")
			utils.SendError(w, err.Error())
			return
		}

		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
		})
	}
}
