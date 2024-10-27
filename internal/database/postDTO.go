package database

import "github.com/jackc/pgx/v5/pgtype"

type PostDTO struct {
	Id        int
	Title     string
	Content   string
	UserId    int
	CreatedAt pgtype.Timestamp
}
