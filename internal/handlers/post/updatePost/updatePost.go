package updatePost

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
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Response struct {
	Status string           `json:"status"`
	Error  string           `json:"error,omitempty"`
	Post   database.PostDTO `json:"post"`
}

func New(log *slog.Logger, storage *database.Dbpool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Update post")

		postID, err := strconv.Atoi(r.PathValue("postID"))
		if err != nil {
			log.Error("Invalid post id", slog.String("post_id", r.PathValue("postID")), slog.String("error", err.Error()))
			utils.SendError(w, "Invalid post id")
			return
		}

		post, err := storage.GetPost(context.Background(), postID)
		if err != nil {
			log.Error("post not found", slog.String("post_id", r.PathValue("postID")), slog.String("Error", err.Error()))
			utils.SendError(w, "post not found")
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
		if req.Title == "" && req.Content == "" {
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

		postDto := database.PostDTO{
			Id:      postID,
			Title:   req.Title,
			Content: req.Content,
		}
		err = storage.UpdatePost(context.Background(), postDto)
		if err != nil {
			log.Error("failed to update post", slog.String("post_id", r.PathValue("postID")), slog.String("error", err.Error()))
			utils.SendError(w, err.Error())
			return
		}

		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
			Post: database.PostDTO{
				Id:        postID,
				Title:     utils.CoalesceString(req.Title, post.Title),
				Content:   utils.CoalesceString(req.Content, post.Content),
				UserId:    post.UserId,
				CreatedAt: post.CreatedAt,
			},
		})
	}
}
