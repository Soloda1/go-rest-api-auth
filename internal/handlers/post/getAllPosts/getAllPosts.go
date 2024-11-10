package getAllPosts

import (
	"gocourse/internal/database"
	"gocourse/internal/utils"
	"log/slog"
	"net/http"
)

// Response represents the get all posts response payload.
// swagger:model
type Response struct {
	Status string             `json:"status"`
	Error  string             `json:"error,omitempty"`
	Posts  []database.PostDTO `json:"posts"`
}

func New(log *slog.Logger, service database.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("get all posts")

		posts, err := service.GetALlPosts()
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
