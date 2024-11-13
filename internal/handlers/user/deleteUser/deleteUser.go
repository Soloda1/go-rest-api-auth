package deleteUser

import (
	"go-rest-api-auth/internal/database"
	"go-rest-api-auth/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

// Response represents the deletion user response payload.
// swagger:model
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	UserID int    `json:"user_id,omitempty"`
}

func New(log *slog.Logger, service database.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Delete user")

		userID, err := strconv.Atoi(r.PathValue("userID"))
		if err != nil {
			log.Error("Invalid user id", slog.String("user_id", r.PathValue("userID")), slog.String("error", err.Error()))
			utils.SendError(w, "Invalid user id")
			return
		}

		_, err = service.GetUserById(userID)
		if err != nil {
			log.Error("User not found", slog.String("user_id", r.PathValue("userID")), slog.String("Error", err.Error()))
			utils.SendError(w, "User not found")
			return
		}

		err = service.DeleteUser(userID)
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
