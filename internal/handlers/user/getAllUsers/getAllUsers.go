package getAllUsers

import (
	"gocourse/internal/database"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
)

type Response struct {
	Status string             `json:"status"`
	Error  string             `json:"error,omitempty"`
	Users  []database.UserDTO `json:"usernames,omitempty"`
}

func New(log *slog.Logger, service database.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("get all users")

		users, err := service.GetALlUsers()
		if err != nil {
			log.Error("get all users failed", slog.String("error", err.Error()))
			utils.SendError(w, "get all users failed")
			return
		}

		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
			Users:  users,
		})
	}
}
