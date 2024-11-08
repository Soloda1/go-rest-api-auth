package logout

import (
	"gocourse/internal/database/auth"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func New(log *slog.Logger, tokenManager *auth.JwtManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Logout user")

		userID, err := strconv.Atoi(r.Context().Value("userID").(string))
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
