package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"strings"
	"time"
)

func (pg *Dbpool) CreatePost(ctx context.Context, post PostDTO) (PostDTO, error) {
	createdAt := pgtype.Timestamp{
		Time:  time.Now(), // Текущая дата, время будет проигнорировано
		Valid: true,       // Отмечаем, что значение установлено
	}

	args := pgx.NamedArgs{
		"title":      post.Title,
		"content":    post.Content,
		"user_id":    post.UserId,
		"created_at": createdAt,
	}

	query := `INSERT INTO posts (title, content, user_id, created_at) VALUES (@title, @content, @user_id, @created_at) RETURNING id, title, content, user_id, created_at`

	var createdPost PostDTO
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
