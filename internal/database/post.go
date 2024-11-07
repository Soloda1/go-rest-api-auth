package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"log/slog"
	"os"
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

//TODO добавить логи во всю бд все таки

func CreateTagsToPost(ctx context.Context, pg *Dbpool, post *PostDTO, tags []string, removeAll bool) {
	tagsChan := make(chan TagsDTO, len(tags))

	if removeAll == true {
		err := pg.DeletePostTagsRelation(ctx, post.Id)
		if err != nil {
			slog.Error("error in in delete tags gorutine", slog.String("err", err.Error()))
			os.Exit(1) //TODO как то исправить  чтоб при ошибке не вырубался сервер
			return
		}
	}

	for _, tag := range tags {
		go func() {
			tagDTO, err := pg.GetTag(ctx, tag)
			if err != nil {
				tagDTO, err = pg.CreateTag(ctx, TagsDTO{Name: tag})
				if err != nil {
					slog.Error("error in create tag", slog.String("err", err.Error()))
					os.Exit(1) //TODO как то исправить  чтоб при ошибке не вырубался сервер
					return
				}
			}
			tagsChan <- tagDTO
		}()
	}

	for range tags {
		tag := <-tagsChan
		if !slices.Contains(post.Tags, tag.Name) {
			post.Tags = append(post.Tags, tag.Name)
			go func() {
				err := pg.CreatePostTagsRelation(ctx, tag, *post)
				if err != nil {
					slog.Error("error in create relation post tags", slog.String("err", err.Error()))
					os.Exit(1) //TODO как то исправить  чтоб при ошибке не вырубался сервер
					return
				}
			}()
		}
	}
}

func (pg *Dbpool) CreatePost(ctx context.Context, post PostDTO) (PostDTO, error) {
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
		return PostDTO{}, err
	}

	if post.Tags != nil {
		CreateTagsToPost(ctx, pg, &createdPost, post.Tags, false)
	}

	return createdPost, nil
}

func (pg *Dbpool) DeletePost(ctx context.Context, postID int) error {
	query := `DELETE FROM posts WHERE id = @id`
	args := pgx.NamedArgs{
		"id": postID,
	}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (pg *Dbpool) UpdatePost(ctx context.Context, post PostDTO) error {
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
		CreateTagsToPost(ctx, pg, &PostDTO{Id: post.Id, Tags: nil}, post.Tags, true)
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = @id"

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (pg *Dbpool) GetPost(ctx context.Context, postID int) (PostDTO, error) {
	query := `SELECT id, title, content, user_id, created_at FROM posts WHERE id = @id`
	args := pgx.NamedArgs{"id": postID}
	row := pg.db.QueryRow(ctx, query, args)
	post := PostDTO{}
	err := row.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt)
	if err != nil {
		return PostDTO{}, err
	}

	return post, nil
}

func (pg *Dbpool) GetALlPosts(ctx context.Context) ([]PostDTO, error) {
	query := `SELECT id, title, content, user_id, created_at FROM posts`

	rows, err := pg.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[PostDTO])
}
