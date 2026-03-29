package config

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type ConfigApp struct {
	PG_DBHost      string `env:"PG_HOST" envDefault:"postgres"`
	PG_DBUser      string `env:"PG_USERNAME_FOR_APP" envDefault:""`
	PG_DBPassword  string `env:"PG_USERPASS_FOR_APP" envDefault:""`
	PG_DBName      string `env:"STANDART_PG_DB_NAME" envDefault:""`
	PG_DBSSLMode   string `env:"PG_SSLMODE" envDefault:"disable"`
	PG_PORT        string `env:"PG_PORT" envDefault:"5432"`
	Logs_Level     string `env:"LOGS_LEVEL_APP" envDefault:"INFO"`
	App_port       string `env:"APP_PORT" envDefault:"8080"`
	Jwt_secret_key string `env:"JWT_SECRET_KEY" envDefault:""`
}

func MustLoadConfigApp() *ConfigApp {

	if err := godotenv.Load(); err != nil {
		fmt.Printf("Failed to load .env file, %v\n", err)
	} else {
		fmt.Println("Loaded configuration from .env file")
	}

	var cfg ConfigApp
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Failed to parse environment variables: %v\n", err)
		panic("configuration error: " + err.Error())
	}

	return &cfg
}

func (c *ConfigApp) GetLogLevel() slog.Level {
	return getLogLevelFromString(c.Logs_Level)
}
