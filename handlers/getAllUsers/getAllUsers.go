package getAllUsers

import (
	"gocourse/handlers/utils"
	"log/slog"
	"net/http"
)

type Response struct {
	Status    string   `json:"status"`
	Error     string   `json:"error,omitempty"`
	Usernames []string `json:"usernames,omitempty"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("get all users")

		// TODO бизнес логика получения всех юзеров
		users := []string{
			"user1", "user2", "user3", "user4", "user5", "user6", "user7", "user8", "user9", "user10",
		}

		utils.Send(w, Response{
			Status:    http.StatusText(http.StatusOK),
			Usernames: users,
		})
	}
}
