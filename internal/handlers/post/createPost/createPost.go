package createPost

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"gocourse/internal/database"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

// Request represents the creation post request payload.
// swagger:model
type Request struct {
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

// Response represents the creation post response payload.
// swagger:model
type Response struct {
	Status string           `json:"status"`
	Error  string           `json:"error,omitempty"`
	Post   database.PostDTO `json:"post"`
}

func New(log *slog.Logger, service database.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Create Post")

		userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
		if err != nil {
			log.Debug("User ID not found")
			utils.SendError(w, "User ID not found")
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
		err = validator.New().Struct(req)
		if err != nil {
			log.Error("failed to validate request", slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		//create user in db
		postDto := database.PostDTO{
			Title:   req.Title,
			Content: req.Content,
			UserId:  userID,
			Tags:    req.Tags,
		}
		createdPost, err := service.CreatePost(postDto)
		if err != nil {
			log.Error("failed to create post", slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		//send response
		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
			Post: database.PostDTO{
				Id:        createdPost.Id,
				Title:     createdPost.Title,
				Content:   createdPost.Content,
				UserId:    userID,
				CreatedAt: createdPost.CreatedAt,
				Tags:      createdPost.Tags,
			},
		})
	}
}
