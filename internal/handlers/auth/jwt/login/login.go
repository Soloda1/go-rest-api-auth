package login

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"gocourse/internal/database"
	"gocourse/internal/database/auth"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Status       string `json:"status"`
	Error        string `json:"error,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func New(log *slog.Logger, tokenManager *auth.JwtManager, userService database.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Login user")

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

		accessToken, err := tokenManager.GenerateJWT(strconv.Itoa(user.Id), "access", tokenManager.AccessExpiresAt)
		if err != nil {
			log.Error("failed to generate access token", slog.String("username", req.Username), slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		refreshToken, err := tokenManager.GenerateJWT(strconv.Itoa(user.Id), "refresh", tokenManager.RefreshExpiresAt)
		if err != nil {
			log.Error("failed to generate refresh token", slog.String("username", req.Username), slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		err = tokenManager.SaveRefreshToken(refreshToken)
		if err != nil {
			log.Error("failed to save refresh token", slog.String("username", req.Username), slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		utils.Send(w, Response{
			Status:       http.StatusText(http.StatusOK),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
