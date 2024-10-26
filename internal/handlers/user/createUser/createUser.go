package createUser

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"gocourse/internal/database"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
)

type Request struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Description string `json:"description"`
}

type Response struct {
	Status string           `json:"status"`
	Error  string           `json:"error,omitempty"`
	User   database.UserDTO `json:"user"`
}

func New(log *slog.Logger, storage *database.Dbpool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Create user")

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

		//create user in db
		userDto := database.UserDTO{
			Username:    req.Username,
			Password:    req.Password,
			Description: req.Description,
		}
		createdUser, err := storage.CreateUser(context.Background(), userDto)
		if err != nil {
			log.Error("failed to create user", slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		//send response
		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
			User: database.UserDTO{
				Id:          createdUser.Id,
				Username:    createdUser.Username,
				Password:    createdUser.Password,
				Description: createdUser.Description,
				DateJoined:  createdUser.DateJoined,
			},
		})
	}
}
