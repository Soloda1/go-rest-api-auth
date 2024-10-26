package database

import (
	"context"
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
		    password VARCHAR(255) NOT NULL,
		    description TEXT DEFAULT '',
		    date_joined DATE DEFAULT CURRENT_DATE
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

func (pg *Dbpool) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *Dbpool) Close() {
	pg.db.Close()
}
