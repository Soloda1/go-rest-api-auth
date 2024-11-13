package deletePost

import (
	"go-rest-api-auth/internal/database"
	"go-rest-api-auth/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

// Response represents the deleting post response payload.
// swagger:model
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	PostID int    `json:"post_id,omitempty"`
}

func New(log *slog.Logger, service database.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Delete post")

		postID, err := strconv.Atoi(r.PathValue("postID"))
		if err != nil {
			log.Error("Invalid post id", slog.String("post_id", r.PathValue("postID")), slog.String("error", err.Error()))
			utils.SendError(w, "Invalid post id")
			return
		}

		_, err = service.GetPost(postID)
		if err != nil {
			log.Error("post not found", slog.String("post_id", r.PathValue("postID")), slog.String("Error", err.Error()))
			utils.SendError(w, "post not found")
			return
		}

		err = service.DeletePost(postID)
		if err != nil {
			log.Error("Error deleting post", slog.String("post_id", r.PathValue("postID")), slog.String("error", err.Error()))
			utils.SendError(w, "Error deleting post")
			return
		}

		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
			PostID: postID,
		})
	}
}
