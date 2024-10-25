package deleteUser

import (
	"gocourse/handlers/utils"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	UserID int    `json:"user_id,omitempty"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Delete user")

		userID, err := strconv.Atoi(r.PathValue("userID"))
		if err != nil {
			log.Error("Invalid user id", slog.String("user_id", r.PathValue("userID")))
			utils.SendError(w, "Invalid user id")
			return
		}

		// TODO бизнес логика удаление юзера

		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
			UserID: userID,
		})
	}
}
