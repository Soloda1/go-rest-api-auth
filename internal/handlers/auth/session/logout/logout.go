package sessionLogout

import (
	"go-rest-api-auth/internal/database/auth"
	"go-rest-api-auth/internal/utils"

	"log/slog"
	"net/http"

	"time"
)

// Response represents the session logout response payload.
// swagger:model
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func New(log *slog.Logger, sessionManager auth.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Session Logout user")

		cookie, err := r.Cookie("session_id")
		if err != nil {
			log.Error("Error getting cookie", slog.String("error", err.Error()))
			utils.SendError(w, "Error getting cookie")
			return
		}

		err = sessionManager.DeleteSession(cookie.Value)
		if err != nil {
			log.Error("Error deleting session", slog.String("error", err.Error()))
			utils.SendError(w, "Error deleting session")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   "",
			Expires: time.Unix(0, 0),
		})

		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
		})
	}
}
