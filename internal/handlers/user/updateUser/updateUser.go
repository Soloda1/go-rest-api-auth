package updateUser

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"gocourse/internal/database"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func New(log *slog.Logger, storage *database.Dbpool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Update user")

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

		//get request body info
		var req Request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		//validating request body info
		if req.Username == "" && req.Password == "" {
			log.Error("Empty request data")
			utils.SendError(w, "Empty request data")
			return
		} else {
			err = validator.New().Struct(req)
			if err != nil {
				log.Error("failed to validate request", slog.String("error", err.Error()))
				utils.SendError(w, err.Error())
				return
			}
		}

		user := database.UserDTO{userID, req.Username, req.Password}
		err = storage.UpdateUser(context.Background(), user)
		if err != nil {
			log.Error("failed to update user", slog.String("user_id", r.PathValue("userID")), slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		utils.Send(w, Response{
			Status:   http.StatusText(http.StatusOK),
			Username: req.Username,
			Password: req.Password,
		})
	}
}
