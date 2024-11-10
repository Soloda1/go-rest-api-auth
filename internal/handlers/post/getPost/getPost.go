package getPost

import (
	"gocourse/internal/database"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

// Response represents the get post response payload.
// swagger:model
type Response struct {
	Status string           `json:"status"`
	Error  string           `json:"error,omitempty"`
	Post   database.PostDTO `json:"post"`
}

func New(log *slog.Logger, service database.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("get one post")

		postID, err := strconv.Atoi(r.PathValue("postID"))
		if err != nil {
			log.Error("Invalid post id", slog.String("post_id", r.PathValue("postID")), slog.String("Error", err.Error()))
			utils.SendError(w, "Invalid post id")
			return
		}

		post, err := service.GetPost(postID)
		if err != nil {
			log.Error("post not found", slog.String("post_id", r.PathValue("postID")), slog.String("Error", err.Error()))
			utils.SendError(w, "post not found")
			return
		}

		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
			Post:   post,
		})
	}
}
