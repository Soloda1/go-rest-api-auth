package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"log/slog"
	"strings"
	"time"
)

func (pg *Dbpool) CreateUser(ctx context.Context, user UserDTO) (UserDTO, error) {
	dateJoined := pgtype.Date{
		Time:  time.Now(), // Текущая дата, время будет проигнорировано
		Valid: true,       // Отмечаем, что значение установлено
	}

	args := pgx.NamedArgs{
		"username":   user.Username,
		"password":   user.Password,
		"dateJoined": dateJoined,
	}
	var query string
	if user.Description == "" {
		query = `INSERT INTO users (username, password, date_joined) VALUES (@username, @password, @dateJoined) RETURNING id, username, password, description, date_joined`
	} else {
		args["description"] = user.Description
		query = `INSERT INTO users (username, password, description, date_joined) VALUES (@username, @password, @description, @dateJoined) RETURNING id, username, password, description, date_joined`
	}
	var createdUser UserDTO
	err := pg.db.QueryRow(ctx, query, args).Scan(
		&createdUser.Id,
		&createdUser.Username,
		&createdUser.Password,
		&createdUser.Description,
		&createdUser.DateJoined,
	)
	if err != nil {
		return UserDTO{}, err
	}

	return createdUser, nil
}

func (pg *Dbpool) DeleteUser(ctx context.Context, userID int) error {
	query := `DELETE FROM users WHERE id = @id`
	args := pgx.NamedArgs{
		"id": userID,
	}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (pg *Dbpool) UpdateUser(ctx context.Context, user UserDTO) error {
	query := `UPDATE users SET `
	var setClauses []string
	args := pgx.NamedArgs{"id": user.Id}

	if user.Username != "" {
		setClauses = append(setClauses, "username = @username")
		args["username"] = user.Username
	}
	if user.Password != "" {
		setClauses = append(setClauses, "password = @password")
		args["password"] = user.Password
	}
	if user.Description != "" {
		setClauses = append(setClauses, "description = @description")
		args["description"] = user.Description
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = @id"
	slog.Info(query)

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (pg *Dbpool) GetUser(ctx context.Context, userID int) (UserDTO, error) {
	query := `SELECT id,username,password,description,date_joined FROM users WHERE id = @id`
	args := pgx.NamedArgs{
		"id": userID,
	}
	row := pg.db.QueryRow(ctx, query, args)
	user := UserDTO{}
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Description, &user.DateJoined)
	if err != nil {
		return UserDTO{}, err
	}

	return user, nil
}

func (pg *Dbpool) GetALlUsers(ctx context.Context) ([]UserDTO, error) {
	query := `SELECT id,username,password,description,date_joined FROM users`

	rows, err := pg.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[UserDTO])
}
