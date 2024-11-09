package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log/slog"
	"os"
)

type CacheClient struct {
	Cache *redis.Client
	Ctx   context.Context
	Log   *slog.Logger
}

func NewRedisClient(ctx context.Context, log *slog.Logger, redisUrl string) *CacheClient {
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Error("Error parsing redis url", slog.String("error", err.Error()))
		os.Exit(1)
	}

	rdb := redis.NewClient(opt)
	ok := rdb.Ping(ctx).Err()
	if ok != nil {
		log.Error("Error connecting to redis", slog.String("error", ok.Error()))
		os.Exit(1)
	}

	return &CacheClient{
		Cache: rdb,
		Ctx:   ctx,
		Log:   log,
	}
}
