package login

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"gocourse/internal/database"
	"gocourse/internal/database/auth"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
	SessionID string `json:"session_id"`
}

func New(log *slog.Logger, sessionManager *auth.SessionManager, userService database.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Session Login user")

		//get request body info
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		//validating request body info
		err = validator.New().Struct(req)
		if err != nil {
			log.Error("failed to validate request", slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		user, err := userService.GetUserByName(req.Username)
		if err != nil {
			log.Error("failed to get user", slog.String("username", req.Username))
			utils.SendError(w, err.Error())
			return
		}

		if !utils.CheckPasswordHash(req.Password, user.Password) {
			log.Error("invalid password", slog.String("username", req.Username))
			utils.SendError(w, "invalid password")
			return
		}

		_, sessionExists := sessionManager.GetSessionByUserID(strconv.Itoa(user.Id))
		if !errors.Is(sessionExists, sessionManager.ErrSessionNotFound) {
			log.Error("session already exists", slog.String("username", req.Username))
			utils.SendError(w, "session already exists")
			return
		}

		sessionID, err := sessionManager.CreateSession(strconv.Itoa(user.Id))
		if err != nil {
			log.Error("failed to create session", slog.String("username", req.Username))
			utils.SendError(w, err.Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Expires:  time.Now().Add(sessionManager.Ttl),
			HttpOnly: true,
		})

		utils.Send(w, Response{
			Status:    http.StatusText(http.StatusOK),
			SessionID: sessionID,
		})
	}
}
