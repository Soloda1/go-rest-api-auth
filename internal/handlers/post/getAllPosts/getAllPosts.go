package getAllPosts

import (
	"gocourse/internal/database"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
)

type Response struct {
	Status string             `json:"status"`
	Error  string             `json:"error,omitempty"`
	Posts  []database.PostDTO `json:"posts"`
}

func New(log *slog.Logger, storage *database.DbPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("get all posts")

		posts, err := storage.GetALlPosts()
		if err != nil {
			log.Error("get all posts failed", slog.String("error", err.Error()))
			utils.SendError(w, "get all posts failed")
			return
		}

		utils.Send(w, Response{
			Status: http.StatusText(http.StatusOK),
			Posts:  posts,
		})
	}
}
