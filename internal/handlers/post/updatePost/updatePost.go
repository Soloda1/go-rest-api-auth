package updatePost

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go-rest-api-auth/internal/database"
	"go-rest-api-auth/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

// Request represents the updating post request payload.
// swagger:model
type Request struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// Response represents the updating post response payload.
// swagger:model
type Response struct {
	Status string           `json:"status"`
	Error  string           `json:"error,omitempty"`
	Post   database.PostDTO `json:"post"`
}

func New(log *slog.Logger, service database.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Update post")

		postID, err := strconv.Atoi(r.PathValue("postID"))
		if err != nil {
			log.Error("Invalid post id", slog.String("post_id", r.PathValue("postID")), slog.String("error", err.Error()))
			utils.SendError(w, "Invalid post id")
			return
		}

		post, err := service.GetPost(postID)
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
			utils.SendError(w, "failed to decode request body")
			return
		}

		//validating request body info
		if req.Title == "" && req.Content == "" && req.Tags == nil {
			log.Error("Empty request data")
			utils.SendError(w, "Empty request data")
			return
		} else {
			err = validator.New().Struct(req)
			if err != nil {
				log.Error("failed to validate request", slog.String("error", err.Error()))
				utils.SendError(w, "failed to validate request")
				return
			}
		}

		postDto := database.PostDTO{
			Id:      postID,
			Title:   req.Title,
			Content: req.Content,
			Tags:    req.Tags,
		}
		err = service.UpdatePost(postDto)
		if err != nil {
			log.Error("failed to update post", slog.String("post_id", r.PathValue("postID")), slog.String("error", err.Error()))
			utils.SendError(w, "failed to update post")
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
				Tags:      utils.CoalesceSliceStrings(req.Tags, post.Tags),
			},
		})
	}
}
