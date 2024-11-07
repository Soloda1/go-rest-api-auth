package getUser

import (
	"gocourse/internal/database"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	Status string           `json:"status"`
	Error  string           `json:"error,omitempty"`
	User   database.UserDTO `json:"user,omitempty"`
}

func New(log *slog.Logger, storage *database.DbPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("get one user")

		userID, err := strconv.Atoi(r.PathValue("userID"))
		if err != nil {
			log.Error("Invalid user id", slog.String("user_id", r.PathValue("userID")), slog.String("Error", err.Error()))
			utils.SendError(w, "Invalid user id")
			return
		}

		user, err := storage.GetUser(userID)
		if err != nil {
			log.Error("User not found", slog.String("user_id", r.PathValue("userID")), slog.String("Error", err.Error()))
			utils.SendError(w, "User not found")
			return
		}

		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
			User:   user,
		})
	}
}
