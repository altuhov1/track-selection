package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"track-selection/internal/config"

	_ "track-selection/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.MustLoadConfigMigrate()
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: cfg.GetLogLevel(),
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("Файл .env не найден, использую системные переменные окружения")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PG_DBUser, cfg.PG_DBPassword, cfg.PG_DBHost, cfg.PG_PORT, cfg.PG_DBName, cfg.PG_DBSSLMode)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		slog.Error("Did not connect to pgx", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		slog.Error("Did not up goose", "err", err)
		os.Exit(1)
	}
}
