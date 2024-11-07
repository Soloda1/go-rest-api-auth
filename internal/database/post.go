package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"slices"
	"strings"
	"time"
)

type PostDTO struct {
	Id        int
	Title     string
	Content   string
	UserId    int
	CreatedAt pgtype.Timestamp
	Tags      []string
}

func CreateTagsToPost(ctx context.Context, log *slog.Logger, pg *Dbpool, post *PostDTO, tags []string, removeAll bool) error {
	tagsChan := make(chan TagsDTO, len(tags))
	g, ctx := errgroup.WithContext(ctx)

	if removeAll {
		if err := pg.DeletePostTagsRelation(ctx, log, post.Id); err != nil {
			log.Error("error deleting tags relation", slog.String("err", err.Error()))
			return err
		}
	}

	for _, tag := range tags {
		tag := tag
		g.Go(func() error {
			tagDTO, err := pg.GetTagByName(ctx, log, tag)
			if err != nil {
				tagDTO, err = pg.CreateTag(ctx, log, TagsDTO{Name: tag})
				if err != nil {
					log.Error("error creating tag", slog.String("err", err.Error()))
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
		_ = g.Wait()
		close(tagsChan)
	}()

	for tag := range tagsChan {
		tag := tag
		if !slices.Contains(post.Tags, tag.Name) {
			post.Tags = append(post.Tags, tag.Name)
			g.Go(func() error {
				err := pg.CreatePostTagsRelation(ctx, log, tag, *post)
				if err != nil {
					log.Error("error creating post-tag relation", slog.String("err", err.Error()))
					return err
				}
				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (pg *Dbpool) CreatePost(ctx context.Context, log *slog.Logger, post PostDTO) (PostDTO, error) {
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

	err := pg.db.QueryRow(ctx, query, args).Scan(
		&createdPost.Id,
		&createdPost.Title,
		&createdPost.Content,
		&createdPost.UserId,
		&createdPost.CreatedAt,
	)
	if err != nil {
		log.Error("Error creating post", slog.String("err", err.Error()))
		return PostDTO{}, err
	}

	if post.Tags != nil {
		err = CreateTagsToPost(ctx, log, pg, &createdPost, post.Tags, false)
		if err != nil {
			log.Error("Error add tags to post", slog.String("err", err.Error()))
			return PostDTO{}, err
		}
	}

	return createdPost, nil
}

func (pg *Dbpool) DeletePost(ctx context.Context, log *slog.Logger, postID int) error {
	query := `DELETE FROM posts WHERE id = @id`
	args := pgx.NamedArgs{
		"id": postID,
	}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		log.Error("Error deleting post", slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (pg *Dbpool) UpdatePost(ctx context.Context, log *slog.Logger, post PostDTO) error {
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
		err := CreateTagsToPost(ctx, log, pg, &PostDTO{Id: post.Id, Tags: nil}, post.Tags, true)
		if err != nil {
			log.Error("Error adding tags to post update", slog.String("err", err.Error()))
			return err
		}
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = @id"

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		log.Error("Error updating post", slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (pg *Dbpool) GetPost(ctx context.Context, log *slog.Logger, postID int) (PostDTO, error) {
	query := `SELECT id, title, content, user_id, created_at FROM posts WHERE id = @id`
	args := pgx.NamedArgs{"id": postID}
	row := pg.db.QueryRow(ctx, query, args)
	post := PostDTO{}
	err := row.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt)
	if err != nil {
		log.Error("Error getting post", slog.String("err", err.Error()))
		return PostDTO{}, err
	}

	return post, nil
}

func (pg *Dbpool) GetALlPosts(ctx context.Context, log *slog.Logger) ([]PostDTO, error) {
	query := `SELECT id, title, content, user_id, created_at FROM posts`

	rows, err := pg.db.Query(ctx, query)
	if err != nil {
		log.Error("Error sql query getting all posts", slog.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var posts []PostDTO
	for rows.Next() {
		var post PostDTO

		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt)
		if err != nil {
			log.Error("Error scanning post", slog.String("err", err.Error()))
			return nil, err
		}

		tagsDto, err := pg.GetPostTagsRelation(ctx, log, post.Id)
		if err != nil {
			log.Error("Error getting post tags relation", slog.String("err", err.Error()))
			return nil, err
		}
		for _, tag := range tagsDto {
			post.Tags = append(post.Tags, tag.Name)
		}

		posts = append(posts, post)
	}

	return posts, nil
}
