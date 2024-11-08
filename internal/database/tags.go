package database

import (
	"github.com/jackc/pgx/v5"
	"log/slog"
	"strconv"
)

type TagsDTO struct {
	Id   int
	Name string
}

type TagsServiceImplementation struct {
	pg *DbPool
}

type TagsService interface {
	CreateTag(tag TagsDTO) (TagsDTO, error)
	DeleteTag(tagID int) error
	GetALlTags() ([]TagsDTO, error)
	GetTagByName(tagName string) (TagsDTO, error)
	GetTagByID(tagID int) (TagsDTO, error)
	CreatePostTagsRelation(tag TagsDTO, post PostDTO) error
	DeletePostTagsRelation(postId int) error
	GetPostTagsRelation(postID int) ([]TagsDTO, error)
}

func NewTagService(pg *DbPool) TagsService {
	return &TagsServiceImplementation{
		pg: pg,
	}
}

func (service *TagsServiceImplementation) CreateTag(tag TagsDTO) (TagsDTO, error) {
	args := pgx.NamedArgs{"name": tag.Name}

	query := `INSERT INTO tags (name) VALUES (@name) RETURNING id, name`

	var createdTag TagsDTO
	err := service.pg.Db.QueryRow(service.pg.Ctx, query, args).Scan(
		&createdTag.Id,
		&createdTag.Name,
	)
	if err != nil {
		service.pg.Log.Error("Error inserting tag into database", slog.String("error", err.Error()), slog.String("tag", tag.Name))
		return TagsDTO{}, err

	}

	return createdTag, nil
}

func (service *TagsServiceImplementation) DeleteTag(tagID int) error {
	query := `DELETE FROM tags WHERE id = @id`
	args := pgx.NamedArgs{"id": tagID}
	_, err := service.pg.Db.Exec(service.pg.Ctx, query, args)
	if err != nil {
		service.pg.Log.Error("Error in deleting tag from database", slog.String("error", err.Error()), slog.String("tagid", strconv.Itoa(tagID)))
		return err
	}
	return nil
}

func (service *TagsServiceImplementation) GetALlTags() ([]TagsDTO, error) {
	query := `SELECT id, name FROM tags`

	rows, err := service.pg.Db.Query(service.pg.Ctx, query)
	if err != nil {
		service.pg.Log.Error("Error in get all tags from database", slog.String("error", err.Error()), slog.String("query", query))
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[TagsDTO])
}

func (service *TagsServiceImplementation) GetTagByName(tagName string) (TagsDTO, error) {
	query := `SELECT id, name FROM tags WHERE name = @name`
	args := pgx.NamedArgs{"name": tagName}
	row := service.pg.Db.QueryRow(service.pg.Ctx, query, args)
	tag := TagsDTO{}
	err := row.Scan(&tag.Id, &tag.Name)
	if err != nil {
		service.pg.Log.Error("Error get tag by name from database", slog.String("error", err.Error()), slog.String("query", query))
		return TagsDTO{}, err
	}

	return tag, nil
}

func (service *TagsServiceImplementation) GetTagByID(tagID int) (TagsDTO, error) {
	query := `SELECT id, name FROM tags WHERE id = @id`
	args := pgx.NamedArgs{"id": tagID}
	row := service.pg.Db.QueryRow(service.pg.Ctx, query, args)
	tag := TagsDTO{}
	err := row.Scan(&tag.Id, &tag.Name)
	if err != nil {
		service.pg.Log.Error("Error getting tag by id from database", slog.String("error", err.Error()), slog.String("query", query))
		return TagsDTO{}, err
	}
	return tag, nil
}

func (service *TagsServiceImplementation) CreatePostTagsRelation(tag TagsDTO, post PostDTO) error {
	args := pgx.NamedArgs{
		"post_id": post.Id,
		"tag_id":  tag.Id,
	}

	query := `INSERT INTO posts_tags (post_id, tag_id) VALUES (@post_id, @tag_id)`

	_, err := service.pg.Db.Exec(service.pg.Ctx, query, args)
	if err != nil {
		service.pg.Log.Error("Error creating post tags relation from database", slog.String("error", err.Error()), slog.String("query", query))
		return err
	}

	return nil
}

func (service *TagsServiceImplementation) DeletePostTagsRelation(postId int) error {
	query := `DELETE FROM posts_tags WHERE post_id = @post_id`
	args := pgx.NamedArgs{"post_id": postId}
	_, err := service.pg.Db.Exec(service.pg.Ctx, query, args)
	if err != nil {
		service.pg.Log.Error("Error deleting post tags relation from database", slog.String("error", err.Error()), slog.String("query", query))
		return err
	}
	return nil
}

func (service *TagsServiceImplementation) GetPostTagsRelation(postID int) ([]TagsDTO, error) {
	query := `SELECT tag_id FROM posts_tags WHERE post_id = @post_id`
	rows, err := service.pg.Db.Query(service.pg.Ctx, query, pgx.NamedArgs{"post_id": postID})
	if err != nil {
		service.pg.Log.Error("Error get post tags relation from database", slog.String("error", err.Error()), slog.String("query", query))
		return nil, err
	}
	defer rows.Close()

	var tagIDs []int
	for rows.Next() {
		var tagID int
		if err = rows.Scan(&tagID); err != nil {
			service.pg.Log.Error("Error scanning row post tags relation", slog.String("error", err.Error()))
			return nil, err
		}
		tagIDs = append(tagIDs, tagID)
	}

	if len(tagIDs) == 0 {
		return nil, nil
	}

	var tags []TagsDTO
	for _, tagID := range tagIDs {
		tag, err := service.GetTagByID(tagID)
		if err != nil {
			service.pg.Log.Error("Error add tags to post from post tags relation from database", slog.String("error", err.Error()))
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
