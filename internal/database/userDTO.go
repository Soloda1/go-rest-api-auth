package database

import "github.com/jackc/pgx/v5/pgtype"

type UserDTO struct {
	Id          int
	Username    string
	Password    string
	Description string
	DateJoined  pgtype.Date
}
