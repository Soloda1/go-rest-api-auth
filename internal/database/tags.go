package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"strconv"
)

type TagsDTO struct {
	Id   int
	Name string
}

func (pg *Dbpool) CreateTag(ctx context.Context, log *slog.Logger, tag TagsDTO) (TagsDTO, error) {
	args := pgx.NamedArgs{"name": tag.Name}

	query := `INSERT INTO tags (name) VALUES (@name) RETURNING id, name`

	var createdTag TagsDTO
	err := pg.db.QueryRow(ctx, query, args).Scan(
		&createdTag.Id,
		&createdTag.Name,
	)
	if err != nil {
		log.Error("Error inserting tag into database", slog.String("error", err.Error()), slog.String("tag", tag.Name))
		return TagsDTO{}, err

	}

	return createdTag, nil
}

func (pg *Dbpool) DeleteTag(ctx context.Context, log *slog.Logger, tagID int) error {
	query := `DELETE FROM tags WHERE id = @id`
	args := pgx.NamedArgs{"id": tagID}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		log.Error("Error in deleting tag from database", slog.String("error", err.Error()), slog.String("tagid", strconv.Itoa(tagID)))
		return err
	}
	return nil
}

func (pg *Dbpool) GetALlTags(ctx context.Context, log *slog.Logger) ([]TagsDTO, error) {
	query := `SELECT id, name FROM tags`

	rows, err := pg.db.Query(ctx, query)
	if err != nil {
		log.Error("Error in get all tags from database", slog.String("error", err.Error()), slog.String("query", query))
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[TagsDTO])
}

func (pg *Dbpool) GetTagByName(ctx context.Context, log *slog.Logger, tagName string) (TagsDTO, error) {
	query := `SELECT id, name FROM tags WHERE name = @name`
	args := pgx.NamedArgs{"name": tagName}
	row := pg.db.QueryRow(ctx, query, args)
	tag := TagsDTO{}
	err := row.Scan(&tag.Id, &tag.Name)
	if err != nil {
		return TagsDTO{}, err
	}

	return tag, nil
}

func (pg *Dbpool) GetTagByID(ctx context.Context, log *slog.Logger, tagID int) (TagsDTO, error) {
	query := `SELECT id, name FROM tags WHERE id = @id`
	args := pgx.NamedArgs{"id": tagID}
	row := pg.db.QueryRow(ctx, query, args)
	tag := TagsDTO{}
	err := row.Scan(&tag.Id, &tag.Name)
	if err != nil {
		log.Error("Error getting tag by id from database", slog.String("error", err.Error()), slog.String("query", query))
		return TagsDTO{}, err
	}
	return tag, nil
}

func (pg *Dbpool) CreatePostTagsRelation(ctx context.Context, log *slog.Logger, tag TagsDTO, post PostDTO) error {
	args := pgx.NamedArgs{
		"post_id": post.Id,
		"tag_id":  tag.Id,
	}

	query := `INSERT INTO posts_tags (post_id, tag_id) VALUES (@post_id, @tag_id)`

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		log.Error("Error creating post tags relation from database", slog.String("error", err.Error()), slog.String("query", query))
		return err
	}

	return nil
}

func (pg *Dbpool) DeletePostTagsRelation(ctx context.Context, log *slog.Logger, postId int) error {
	query := `DELETE FROM posts_tags WHERE post_id = @post_id`
	args := pgx.NamedArgs{"post_id": postId}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		log.Error("Error deleting post tags relation from database", slog.String("error", err.Error()), slog.String("query", query))
		return err
	}
	return nil
}

func (pg *Dbpool) GetPostTagsRelation(ctx context.Context, log *slog.Logger, postID int) ([]TagsDTO, error) {
	query := `SELECT tag_id FROM posts_tags WHERE post_id = @post_id`
	rows, err := pg.db.Query(ctx, query, pgx.NamedArgs{"post_id": postID})
	if err != nil {
		log.Error("Error get post tags relation from database", slog.String("error", err.Error()), slog.String("query", query))
		return nil, err
	}
	defer rows.Close()

	// Сохраняем ID тегов в слайс
	var tagIDs []int
	for rows.Next() {
		var tagID int
		if err = rows.Scan(&tagID); err != nil {
			log.Error("Error scanning row post tags relation", slog.String("error", err.Error()))
			return nil, err
		}
		tagIDs = append(tagIDs, tagID)
	}

	if len(tagIDs) == 0 {
		return nil, nil
	}

	var tags []TagsDTO
	for _, tagID := range tagIDs {
		tag, err := pg.GetTagByID(ctx, log, tagID)
		if err != nil {
			log.Error("Error add tags to post from post tags relation from database", slog.String("error", err.Error()))
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
