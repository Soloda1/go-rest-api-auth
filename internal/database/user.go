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

type UserServiceImplementation struct {
	pg *DbPool
}

type UserService interface {
	CreateUser(user UserDTO) (UserDTO, error)
	DeleteUser(userID int) error
	UpdateUser(user UserDTO) error
	GetUserById(userID int) (UserDTO, error)
	GetALlUsers() ([]UserDTO, error)
	GetUserByName(username string) (UserDTO, error)
}

func NewUserService(pg *DbPool) UserService {
	return &UserServiceImplementation{
		pg: pg,
	}
}

func (service *UserServiceImplementation) CreateUser(user UserDTO) (UserDTO, error) {
	dateJoined := pgtype.Date{
		Time:  time.Now(), // Текущая дата, время будет проигнорировано
		Valid: true,       // Отмечаем, что значение установлено
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		service.pg.Log.Error("Error hashing password in create user in database", slog.String("password", user.Password))
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
	err = service.pg.Db.QueryRow(service.pg.Ctx, query, args).Scan(
		&createdUser.Id,
		&createdUser.Username,
		&createdUser.Password,
		&createdUser.Description,
		&createdUser.DateJoined,
	)
	if err != nil {
		service.pg.Log.Error("Error creating new user in database", slog.String("username", user.Username))
		return UserDTO{}, err
	}

	return createdUser, nil
}

func (service *UserServiceImplementation) DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = @id`
	args := pgx.NamedArgs{
		"id": userID,
	}
	_, err := service.pg.Db.Exec(service.pg.Ctx, query, args)
	if err != nil {
		service.pg.Log.Error("Error deleting user from database", slog.String("user_id", strconv.Itoa(userID)))
		return err
	}
	return nil
}

func (service *UserServiceImplementation) UpdateUser(user UserDTO) error {
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
			service.pg.Log.Error("Error hashing password in update user in database", slog.String("password", user.Password))
			return err
		}
		args["password"] = hashedPassword
	}
	if user.Description != "" {
		setClauses = append(setClauses, "description = @description")
		args["description"] = user.Description
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = @id"

	_, err := service.pg.Db.Exec(service.pg.Ctx, query, args)
	if err != nil {
		service.pg.Log.Error("Error updating user in database", slog.String("user_id", strconv.Itoa(user.Id)))
		return err
	}
	return nil
}

func (service *UserServiceImplementation) GetUserById(userID int) (UserDTO, error) {
	query := `SELECT id,username,password,description,date_joined FROM users WHERE id = @id`
	args := pgx.NamedArgs{
		"id": userID,
	}
	row := service.pg.Db.QueryRow(service.pg.Ctx, query, args)
	user := UserDTO{}
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Description, &user.DateJoined)
	if err != nil {
		service.pg.Log.Error("Error getting user by id from database", slog.String("user_id", strconv.Itoa(userID)), slog.String("error", err.Error()))
		return UserDTO{}, err
	}

	return user, nil
}

func (service *UserServiceImplementation) GetUserByName(username string) (UserDTO, error) {
	query := `SELECT id,username,password,description,date_joined FROM users WHERE username = @username`
	args := pgx.NamedArgs{
		"username": username,
	}
	row := service.pg.Db.QueryRow(service.pg.Ctx, query, args)
	user := UserDTO{}
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Description, &user.DateJoined)
	if err != nil {
		service.pg.Log.Error("Error getting user by name from database", slog.String("username", username), slog.String("error", err.Error()))
		return UserDTO{}, err
	}

	return user, nil
}

func (service *UserServiceImplementation) GetALlUsers() ([]UserDTO, error) {
	query := `SELECT id,username,password,description,date_joined FROM users`

	rows, err := service.pg.Db.Query(service.pg.Ctx, query)
	if err != nil {
		service.pg.Log.Error("Error getting all users from database")
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[UserDTO])
}
