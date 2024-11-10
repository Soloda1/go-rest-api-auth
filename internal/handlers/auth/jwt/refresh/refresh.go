package refresh

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"gocourse/internal/database/auth"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

// Request represents the jwt refresh request payload.
// swagger:model
type Request struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Response represents the jwt refresh response payload.
// swagger:model
type Response struct {
	Status       string `json:"status"`
	Error        string `json:"error,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func New(log *slog.Logger, tokenManager *auth.JwtManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Refresh user's tokens")

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

		refreshTokenClaim, err := tokenManager.ValidateJWT(req.RefreshToken, "refresh")
		if err != nil {
			log.Error("failed to validate token", slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		exist, err := tokenManager.IsRefreshTokenValid(req.RefreshToken)
		if err != nil {
			log.Error("failed to check refresh token", slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		} else if !exist && err == nil {
			log.Error("refresh token not found", slog.String("error", "refresh token is not valid"))
			utils.SendError(w, "refresh token not found")
			return
		}

		userID, err := strconv.Atoi(refreshTokenClaim["sub"].(string))
		if err != nil {
			log.Error("failed to convert userID to int", slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		err = tokenManager.DeleteRefreshToken(userID)
		if err != nil {
			log.Error("failed to delete refresh token from database", slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		accessToken, err := tokenManager.GenerateJWT(strconv.Itoa(userID), "access", tokenManager.AccessExpiresAt)
		if err != nil {
			log.Error("failed to generate access token", slog.String("userID", strconv.Itoa(userID)), slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		refreshToken, err := tokenManager.GenerateJWT(strconv.Itoa(userID), "refresh", tokenManager.RefreshExpiresAt)
		if err != nil {
			log.Error("failed to generate refresh token", slog.String("userID", strconv.Itoa(userID)), slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		err = tokenManager.SaveRefreshToken(refreshToken)
		if err != nil {
			log.Error("failed to save refresh token", slog.String("userID", strconv.Itoa(userID)), slog.String("error", err.Error()))
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
