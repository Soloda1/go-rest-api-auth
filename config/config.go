package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type DATABASE struct {
	Username string `yaml:"username" env-default:"postgres"`
	Password string `yaml:"password" env-default:"admin"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	DbName   string `yaml:"db_name" env-default:"postgres"`
}

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server" env-required:"true"`
	DATABASE   `yaml:"database" env-required:"true"`
	JWT        `yaml:"jwt" env-required:"true"`
	REDIS      `yaml:"redis" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type JWT struct {
	Secret           string        `yaml:"secret" env-default:"test-secret-key"`
	AccessExpiresAt  time.Duration `yaml:"access_expires_at" env-default:"15m"`
	RefreshExpiresAt time.Duration `yaml:"refresh_expires_at" env-default:"720h"`
}

type REDIS struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"6379"`
	Password string `yaml:"password" env-default:"admin"`
	DbIndex  int    `yaml:"db_index" env-default:"0"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Can't read config: %s", err)
	}

	return &config
}
