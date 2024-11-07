package deletePost

import (
	"gocourse/internal/database"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	PostID int    `json:"post_id,omitempty"`
}

func New(log *slog.Logger, storage *database.DbPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Delete post")

		postID, err := strconv.Atoi(r.PathValue("postID"))
		if err != nil {
			log.Error("Invalid post id", slog.String("post_id", r.PathValue("postID")), slog.String("error", err.Error()))
			utils.SendError(w, "Invalid post id")
			return
		}

		_, err = storage.GetPost(postID)
		if err != nil {
			log.Error("post not found", slog.String("post_id", r.PathValue("postID")), slog.String("Error", err.Error()))
			utils.SendError(w, "post not found")
			return
		}

		err = storage.DeletePost(postID)
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
