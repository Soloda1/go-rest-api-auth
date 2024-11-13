package jwtLogin

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go-rest-api-auth/internal/database"
	"go-rest-api-auth/internal/database/auth"
	"go-rest-api-auth/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

// Request represents the jwt login request payload.
// swagger:model
type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Response represents the jwt login response payload.
// swagger:model
type Response struct {
	Status       string `json:"status"`
	Error        string `json:"error,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func New(log *slog.Logger, tokenManager auth.JwtManager, userService database.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("JWT Login user")

		//get request body info
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			utils.SendError(w, "failed to decode request body")
			return
		}

		//validating request body info
		err = validator.New().Struct(req)
		if err != nil {
			log.Error("failed to validate request", slog.String("error", err.Error()))
			utils.SendError(w, "failed to validate request")
			return
		}

		user, err := userService.GetUserByName(req.Username)
		if err != nil {
			log.Error("failed to get user", slog.String("username", req.Username))
			utils.SendError(w, "failed to get user")
			return
		}

		if !utils.CheckPasswordHash(req.Password, user.Password) {
			log.Error("invalid password", slog.String("username", req.Username))
			utils.SendError(w, "invalid password")
			return
		}

		accessToken, err := tokenManager.GenerateJWT(strconv.Itoa(user.Id), "access", tokenManager.GetterAccessExpiresAt())
		if err != nil {
			log.Error("failed to generate access token", slog.String("username", req.Username), slog.String("error", err.Error()))
			utils.SendError(w, "failed to generate access token")
			return
		}

		refreshToken, err := tokenManager.GenerateJWT(strconv.Itoa(user.Id), "refresh", tokenManager.GetterRefreshExpiresAt())
		if err != nil {
			log.Error("failed to generate refresh token", slog.String("username", req.Username), slog.String("error", err.Error()))
			utils.SendError(w, "failed to generate refresh token")
			return
		}

		err = tokenManager.SaveRefreshToken(refreshToken)
		if err != nil {
			log.Error("failed to save refresh token", slog.String("username", req.Username), slog.String("error", err.Error()))
			utils.SendError(w, "failed to save refresh token")
			return
		}

		utils.Send(w, Response{
			Status:       http.StatusText(http.StatusOK),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
