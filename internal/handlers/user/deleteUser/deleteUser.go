package deleteUser

import (
	"context"
	"gocourse/internal/database"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	UserID int    `json:"user_id,omitempty"`
}

func New(log *slog.Logger, storage *database.Dbpool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Delete user")

		userID, err := strconv.Atoi(r.PathValue("userID"))
		if err != nil {
			log.Error("Invalid user id", slog.String("user_id", r.PathValue("userID")), slog.String("error", err.Error()))
			utils.SendError(w, "Invalid user id")
			return
		}

		_, err = storage.GetUser(context.Background(), userID)
		if err != nil {
			log.Error("User not found", slog.String("user_id", r.PathValue("userID")), slog.String("Error", err.Error()))
			utils.SendError(w, "User not found")
			return
		}

		err = storage.DeleteUser(context.Background(), userID)
		if err != nil {
			log.Error("Error deleting user", slog.String("user_id", r.PathValue("userID")), slog.String("error", err.Error()))
			utils.SendError(w, "Error deleting user")
			return
		}

		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
			UserID: userID,
		})
	}
}
