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

func CreateTagsToPost(ctx context.Context, log *slog.Logger, pg *DbPool, post *PostDTO, tags []string, removeAll bool) error {
	tagsChan := make(chan TagsDTO, len(tags))
	eg, ctx := errgroup.WithContext(ctx)

	if removeAll {
		if err := pg.DeletePostTagsRelation(post.Id); err != nil {
			log.Error("error deleting tags relation", slog.String("err", err.Error()))
			return err
		}
	}

	for _, tag := range tags {
		tag := tag
		eg.Go(func() error {
			tagDTO, err := pg.GetTagByName(tag)
			if err != nil {
				tagDTO, err = pg.CreateTag(TagsDTO{Name: tag})
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
		_ = eg.Wait()
		close(tagsChan)
	}()

	for tag := range tagsChan {
		tag := tag
		if !slices.Contains(post.Tags, tag.Name) {
			post.Tags = append(post.Tags, tag.Name)
			eg.Go(func() error {
				err := pg.CreatePostTagsRelation(tag, *post)
				if err != nil {
					log.Error("error creating post-tag relation", slog.String("err", err.Error()))
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

func (pg *DbPool) CreatePost(post PostDTO) (PostDTO, error) {
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

	err := pg.db.QueryRow(pg.ctx, query, args).Scan(
		&createdPost.Id,
		&createdPost.Title,
		&createdPost.Content,
		&createdPost.UserId,
		&createdPost.CreatedAt,
	)
	if err != nil {
		pg.log.Error("Error creating post", slog.String("err", err.Error()))
		return PostDTO{}, err
	}

	if post.Tags != nil {
		err = CreateTagsToPost(pg.ctx, pg.log, pg, &createdPost, post.Tags, false)
		if err != nil {
			pg.log.Error("Error add tags to post", slog.String("err", err.Error()))
			return PostDTO{}, err
		}
	}

	return createdPost, nil
}

func (pg *DbPool) DeletePost(postID int) error {
	query := `DELETE FROM posts WHERE id = @id`
	args := pgx.NamedArgs{
		"id": postID,
	}
	_, err := pg.db.Exec(pg.ctx, query, args)
	if err != nil {
		pg.log.Error("Error deleting post", slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (pg *DbPool) UpdatePost(post PostDTO) error {
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
		err := CreateTagsToPost(pg.ctx, pg.log, pg, &PostDTO{Id: post.Id, Tags: nil}, post.Tags, true)
		if err != nil {
			pg.log.Error("Error adding tags to post update", slog.String("err", err.Error()))
			return err
		}
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = @id"

	_, err := pg.db.Exec(pg.ctx, query, args)
	if err != nil {
		pg.log.Error("Error updating post", slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (pg *DbPool) GetPost(postID int) (PostDTO, error) {
	query := `SELECT id, title, content, user_id, created_at FROM posts WHERE id = @id`
	args := pgx.NamedArgs{"id": postID}
	row := pg.db.QueryRow(pg.ctx, query, args)
	post := PostDTO{}
	err := row.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt)
	if err != nil {
		pg.log.Error("Error getting post", slog.String("err", err.Error()))
		return PostDTO{}, err
	}

	return post, nil
}

func (pg *DbPool) GetALlPosts() ([]PostDTO, error) {
	query := `SELECT id, title, content, user_id, created_at FROM posts`

	rows, err := pg.db.Query(pg.ctx, query)
	if err != nil {
		pg.log.Error("Error sql query getting all posts", slog.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var posts []PostDTO
	for rows.Next() {
		var post PostDTO

		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt)
		if err != nil {
			pg.log.Error("Error scanning post", slog.String("err", err.Error()))
			return nil, err
		}

		tagsDto, err := pg.GetPostTagsRelation(post.Id)
		if err != nil {
			pg.log.Error("Error getting post tags relation", slog.String("err", err.Error()))
			return nil, err
		}
		for _, tag := range tagsDto {
			post.Tags = append(post.Tags, tag.Name)
		}

		posts = append(posts, post)
	}

	return posts, nil
}
