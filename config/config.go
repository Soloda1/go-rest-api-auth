package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type DATABASE struct {
	Username string `env:"DATABASE_USERNAME" env-default:"postgres"`
	Password string `env:"DATABASE_PASSWORD" env-default:"admin"`
	Host     string `env:"DATABASE_HOST" env-default:"localhost"`
	Port     string `env:"DATABASE_PORT" env-default:"5432"`
	DbName   string `env:"DATABASE_DB_NAME" env-default:"postgres"`
}

type Config struct {
	Env        string `env:"ENV" env-default:"local"`
	HTTPServer `env-required:"true"`
	DATABASE   `env-required:"true"`
	JWT        `env-required:"true"`
	REDIS      `env-required:"true"`
}

type HTTPServer struct {
	Address     string        `env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8000"`
	Timeout     time.Duration `env:"HTTP_SERVER_TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

type JWT struct {
	Secret           string        `env:"JWT_SECRET" env-default:"test-secret-key"`
	AccessExpiresAt  time.Duration `env:"JWT_ACCESS_EXPIRES_AT" env-default:"15m"`
	RefreshExpiresAt time.Duration `env:"JWT_REFRESH_EXPIRES_AT" env-default:"720h"`
}

type REDIS struct {
	Host     string        `env:"REDIS_HOST" env-default:"localhost"`
	Port     string        `env:"REDIS_PORT" env-default:"6379"`
	Password string        `env:"REDIS_PASSWORD" env-default:"admin"`
	DbIndex  int           `env:"REDIS_DB_INDEX" env-default:"0"`
	TTL      time.Duration `env:"REDIS_TTL" env-default:"360h"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	}

	if err := godotenv.Load(configPath); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	var config Config
	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalf("Can't read env config: %s", err)
	}

	return &config
}
