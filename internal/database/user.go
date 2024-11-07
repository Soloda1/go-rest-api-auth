package database

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"gocourse/internal/utils"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type UserDTO struct {
	Id          int
	Username    string
	Password    string
	Description string
	DateJoined  pgtype.Date
}

func (pg *DbPool) CreateUser(user UserDTO) (UserDTO, error) {
	dateJoined := pgtype.Date{
		Time:  time.Now(), // Текущая дата, время будет проигнорировано
		Valid: true,       // Отмечаем, что значение установлено
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		pg.log.Error("Error hashing password in create user in database", slog.String("password", user.Password))
		return UserDTO{}, err
	}

	args := pgx.NamedArgs{
		"username":   user.Username,
		"password":   hashedPassword,
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
	err = pg.db.QueryRow(pg.ctx, query, args).Scan(
		&createdUser.Id,
		&createdUser.Username,
		&createdUser.Password,
		&createdUser.Description,
		&createdUser.DateJoined,
	)
	if err != nil {
		pg.log.Error("Error creating new user in database", slog.String("username", user.Username))
		return UserDTO{}, err
	}

	return createdUser, nil
}

func (pg *DbPool) DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = @id`
	args := pgx.NamedArgs{
		"id": userID,
	}
	_, err := pg.db.Exec(pg.ctx, query, args)
	if err != nil {
		pg.log.Error("Error deleting user from database", slog.String("user_id", strconv.Itoa(userID)))
		return err
	}
	return nil
}

func (pg *DbPool) UpdateUser(user UserDTO) error {
	query := `UPDATE users SET `
	var setClauses []string
	args := pgx.NamedArgs{"id": user.Id}

	if user.Username != "" {
		setClauses = append(setClauses, "username = @username")
		args["username"] = user.Username
	}
	if user.Password != "" {
		setClauses = append(setClauses, "password = @password")
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			pg.log.Error("Error hashing password in update user in database", slog.String("password", user.Password))
			return err
		}
		args["password"] = hashedPassword
	}
	if user.Description != "" {
		setClauses = append(setClauses, "description = @description")
		args["description"] = user.Description
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = @id"

	_, err := pg.db.Exec(pg.ctx, query, args)
	if err != nil {
		pg.log.Error("Error updating user in database", slog.String("user_id", strconv.Itoa(user.Id)))
		return err
	}
	return nil
}

func (pg *DbPool) GetUser(userID int) (UserDTO, error) {
	query := `SELECT id,username,password,description,date_joined FROM users WHERE id = @id`
	args := pgx.NamedArgs{
		"id": userID,
	}
	row := pg.db.QueryRow(pg.ctx, query, args)
	user := UserDTO{}
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Description, &user.DateJoined)
	if err != nil {
		pg.log.Error("Error getting user from database", slog.String("user_id", strconv.Itoa(userID)))
		return UserDTO{}, err
	}

	return user, nil
}

func (pg *DbPool) GetALlUsers() ([]UserDTO, error) {
	query := `SELECT id,username,password,description,date_joined FROM users`

	rows, err := pg.db.Query(pg.ctx, query)
	if err != nil {
		pg.log.Error("Error getting all users from database")
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[UserDTO])
}
