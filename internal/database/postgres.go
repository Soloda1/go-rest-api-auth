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

	go SetupTables(ctx, log)

	return pgInstance
}

func SetupTables(ctx context.Context, log *slog.Logger) {
	done := make(chan bool)

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
		done <- true
	}()

	go func() {
		<-done
		query := `
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT DEFAULT '',
			user_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`
		_, err := pgInstance.db.Exec(ctx, query)
		if err != nil {
			log.Debug("Failed to create user table", slog.String("error", err.Error()))
			os.Exit(1)
		}

		log.Info("Created posts table")
		done <- true
	}()

	go func() {
		query := `
		CREATE TABLE IF NOT EXISTS tags (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) UNIQUE NOT NULL
		)
	`
		_, err := pgInstance.db.Exec(ctx, query)
		if err != nil {
			log.Debug("Failed to create user table", slog.String("error", err.Error()))
			os.Exit(1)
		}

		log.Info("Created tags table")
		done <- true
	}()

	go func() {
		<-done
		<-done
		query := `
		CREATE TABLE IF NOT EXISTS posts_tags (
			post_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, tag_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		)
	`
		_, err := pgInstance.db.Exec(ctx, query)
		if err != nil {
			log.Debug("Failed to create user table", slog.String("error", err.Error()))
			os.Exit(1)
		}

		log.Info("Created ManyToMany tags <=> posts table")
		done <- true
	}()

}

func (pg *Dbpool) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *Dbpool) Close() {
	pg.db.Close()
}
