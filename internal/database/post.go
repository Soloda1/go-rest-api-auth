package database

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"slices"
	"strings"
	"time"
)

type PostDTO struct {
	Id        int              `json:"id,omitempty"`
	Title     string           `json:"title,omitempty"`
	Content   string           `json:"content,omitempty"`
	UserId    int              `json:"user_id,omitempty"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Tags      []string         `json:"tags,omitempty"`
}

type PostServiceImplementation struct {
	tagsService TagsService
	pg          *DbPool
}

//go:generate go run github.com/vektra/mockery/v2@v2.46.3 --name PostService --output ../../testing/mocks
type PostService interface {
	CreatePost(post PostDTO) (PostDTO, error)
	DeletePost(postID int) error
	UpdatePost(post PostDTO) error
	GetPost(postID int) (PostDTO, error)
	GetALlPosts() ([]PostDTO, error)
}

func NewPostService(pg *DbPool, tagsService TagsService) PostService {
	return &PostServiceImplementation{
		tagsService: tagsService,
		pg:          pg,
	}
}

func (service *PostServiceImplementation) DeletePost(postID int) error {
	query := `DELETE FROM posts WHERE id = @id`
	args := pgx.NamedArgs{
		"id": postID,
	}
	_, err := service.pg.Db.Exec(service.pg.Ctx, query, args)
	if err != nil {
		service.pg.Log.Error("Error deleting post", slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (service *PostServiceImplementation) CreatePost(post PostDTO) (PostDTO, error) {
	var createdPost PostDTO

	createdAt := pgtype.Timestamp{
		Time:  time.Now(),
		Valid: true,
	}

	args := pgx.NamedArgs{
		"title":      post.Title,
		"content":    post.Content,
		"user_id":    post.UserId,
		"created_at": createdAt,
	}

	query := `INSERT INTO posts (title, content, user_id, created_at) VALUES (@title, @content, @user_id, @created_at) RETURNING id, title, content, user_id, created_at`

	err := service.pg.Db.QueryRow(service.pg.Ctx, query, args).Scan(
		&createdPost.Id,
		&createdPost.Title,
		&createdPost.Content,
		&createdPost.UserId,
		&createdPost.CreatedAt,
	)
	if err != nil {
		service.pg.Log.Error("Error creating post", slog.String("err", err.Error()))
		return PostDTO{}, err
	}

	if post.Tags != nil {
		err = service.CreateTagsToPost(&createdPost, post.Tags, false)
		if err != nil {
			service.pg.Log.Error("Error add tags to post", slog.String("err", err.Error()))
			return PostDTO{}, err
		}
	}

	return createdPost, nil
}

func (service *PostServiceImplementation) UpdatePost(post PostDTO) error {
	query := `UPDATE posts SET `
	var setClauses []string
	args := pgx.NamedArgs{"id": post.Id}

	if post.Title != "" {
		setClauses = append(setClauses, "title = @title")
		args["title"] = post.Title
	}
	if post.Content != "" {
		setClauses = append(setClauses, "content = @content")
		args["content"] = post.Content
	}
	if post.Tags != nil {
		err := service.CreateTagsToPost(&PostDTO{Id: post.Id, Tags: nil}, post.Tags, true)
		if err != nil {
			service.pg.Log.Error("Error adding tags to post update", slog.String("err", err.Error()))
			return err
		}
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = @id"

	_, err := service.pg.Db.Exec(service.pg.Ctx, query, args)
	if err != nil {
		service.pg.Log.Error("Error updating post", slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (service *PostServiceImplementation) GetPost(postID int) (PostDTO, error) {
	query := `SELECT id, title, content, user_id, created_at FROM posts WHERE id = @id`
	args := pgx.NamedArgs{"id": postID}
	row := service.pg.Db.QueryRow(service.pg.Ctx, query, args)
	post := PostDTO{}
	err := row.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt)
	if err != nil {
		service.pg.Log.Error("Error getting post", slog.String("err", err.Error()))
		return PostDTO{}, err
	}

	return post, nil
}

func (service *PostServiceImplementation) GetALlPosts() ([]PostDTO, error) {
	query := `SELECT id, title, content, user_id, created_at FROM posts`

	rows, err := service.pg.Db.Query(service.pg.Ctx, query)
	if err != nil {
		service.pg.Log.Error("Error sql query getting all posts", slog.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var posts []PostDTO
	for rows.Next() {
		var post PostDTO

		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt)
		if err != nil {
			service.pg.Log.Error("Error scanning post", slog.String("err", err.Error()))
			return nil, err
		}

		tagsDto, err := service.tagsService.GetPostTagsRelation(post.Id)
		if err != nil {
			service.pg.Log.Error("Error getting post tags relation", slog.String("err", err.Error()))
			return nil, err
		}
		for _, tag := range tagsDto {
			post.Tags = append(post.Tags, tag.Name)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (service *PostServiceImplementation) CreateTagsToPost(post *PostDTO, tags []string, removeAll bool) error {
	tagsChan := make(chan TagsDTO, len(tags))
	eg, ctx := errgroup.WithContext(service.pg.Ctx)

	if removeAll {
		if err := service.tagsService.DeletePostTagsRelation(post.Id); err != nil {
			service.pg.Log.Error("error deleting tags relation", slog.String("err", err.Error()))
			return err
		}
	}

	for _, tag := range tags {
		tag := tag
		eg.Go(func() error {
			tagDTO, err := service.tagsService.GetTagByName(tag)
			if err != nil {
				tagDTO, err = service.tagsService.CreateTag(TagsDTO{Name: tag})
				if err != nil {
					service.pg.Log.Error("error creating tag", slog.String("err", err.Error()))
					return err
				}
			}

			select {
			case tagsChan <- tagDTO:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})
	}

	go func() {
		_ = eg.Wait()
		close(tagsChan)
	}()

	for tag := range tagsChan {
		tag := tag
		if !slices.Contains(post.Tags, tag.Name) {
			post.Tags = append(post.Tags, tag.Name)
			eg.Go(func() error {
				err := service.tagsService.CreatePostTagsRelation(tag, *post)
				if err != nil {
					service.pg.Log.Error("error creating post-tag relation", slog.String("err", err.Error()))
					return err
				}
				return nil
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
