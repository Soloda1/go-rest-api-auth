package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"sync"
)

type Dbpool struct {
	db *pgxpool.Pool
}

var (
	pgInstance *Dbpool
	pgOnce     sync.Once
)

func NewDbPool(ctx context.Context, dbUrl string, log *slog.Logger) *Dbpool {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, dbUrl)
		if err != nil {
			log.Debug("Failed to connect to database", slog.String("url", dbUrl), slog.String("error", err.Error()))
			os.Exit(1)
		}

		pgInstance = &Dbpool{db}
	})

	go func() {
		query := `
		CREATE TABLE IF NOT EXISTS users (
		    id serial PRIMARY KEY,
		    username VARCHAR(30) NOT NULL UNIQUE ,
		    password VARCHAR(255) NOT NULL
		)
	`
		_, err := pgInstance.db.Exec(ctx, query)
		if err != nil {
			log.Debug("Failed to create user table", slog.String("error", err.Error()))
			os.Exit(1)
		}

		log.Info("Created user table")
	}()

	return pgInstance
}

func (pg *Dbpool) CreateUser(ctx context.Context, log *slog.Logger, username, password string) error {
	query := `INSERT INTO users (username, password) VALUES (@username, @password)`
	args := pgx.NamedArgs{
		"username": username,
		"password": password,
	}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
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
	var query string
	args := pgx.NamedArgs{
		"id": user.UserID,
	}
	if user.Password == "" {
		query = `UPDATE users SET username = @username WHERE id = @id`
		args["username"] = user.Username
	} else if user.Username == "" {
		query = `UPDATE users SET password = @password WHERE id = @id`
		args["password"] = user.Password
	} else {
		query = `UPDATE users SET username = @username, password = @password  WHERE id = @id`
		args["username"] = user.Username
		args["password"] = user.Password
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (pg *Dbpool) GetUser(ctx context.Context, userID int) (UserDTO, error) {
	query := `SELECT id,username,password FROM users WHERE id = @id`
	args := pgx.NamedArgs{
		"id": userID,
	}
	row := pg.db.QueryRow(ctx, query, args)
	user := UserDTO{}
	err := row.Scan(&user.UserID, &user.Username, &user.Password)
	if err != nil {
		return UserDTO{}, err
	}

	return user, nil
}

func (pg *Dbpool) GetALlUsers(ctx context.Context) ([]UserDTO, error) {
	query := `SELECT id,username,password FROM users`

	rows, err := pg.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[UserDTO])
}

func (pg *Dbpool) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *Dbpool) Close() {
	pg.db.Close()
}
