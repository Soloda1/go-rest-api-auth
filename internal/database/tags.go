package database

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type TagsDTO struct {
	Id   int
	Name string
}

func (pg *Dbpool) CreateTag(ctx context.Context, tag TagsDTO) (TagsDTO, error) {
	args := pgx.NamedArgs{"name": tag.Name}

	query := `INSERT INTO tags (name) VALUES (@name) RETURNING id, name`

	var createdTag TagsDTO
	err := pg.db.QueryRow(ctx, query, args).Scan(
		&createdTag.Id,
		&createdTag.Name,
	)
	if err != nil {
		return TagsDTO{}, err
	}

	return createdTag, nil
}

func (pg *Dbpool) DeleteTag(ctx context.Context, tagID int) error {
	query := `DELETE FROM tags WHERE id = @id`
	args := pgx.NamedArgs{"id": tagID}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (pg *Dbpool) GetALlTags(ctx context.Context) ([]TagsDTO, error) {
	query := `SELECT id, name FROM tags`

	rows, err := pg.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[TagsDTO])
}

func (pg *Dbpool) GetTag(ctx context.Context, tagName string) (TagsDTO, error) {
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

func (pg *Dbpool) CreatePostTagsRelation(ctx context.Context, tag TagsDTO, post PostDTO) error {
	args := pgx.NamedArgs{
		"post_id": post.Id,
		"tag_id":  tag.Id,
	}

	query := `INSERT INTO posts_tags (post_id, tag_id) VALUES (@post_id, @tag_id)`

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (pg *Dbpool) DeletePostTagsRelation(ctx context.Context, postId int) error {
	query := `DELETE FROM posts_tags WHERE post_id = @post_id`
	args := pgx.NamedArgs{"post_id": postId}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}
