package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateAppUser, downCreateAppUser)
}

func quoteIdentifier(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

func quoteLiteral(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `''`) + `'`
}

func upCreateAppUser(ctx context.Context, tx *sql.Tx) error {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Файл .env не найден, использую системные переменные окружения")
	}
	username := os.Getenv("PG_USERNAME_FOR_APP")
	password := os.Getenv("PG_USERPASS_FOR_APP")

	if username == "" || password == "" {
		return fmt.Errorf("PG_USERNAME_FOR_APP and PG_USERPASS_FOR_APP must be set")
	}

	var exists bool
	err := tx.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_roles WHERE rolname = $1)",
		username,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if !exists {
		safeUser := quoteIdentifier(username)
		safePass := quoteLiteral(password)
		_, err = tx.ExecContext(
			ctx,
			fmt.Sprintf("CREATE USER %s WITH PASSWORD %s", safeUser, safePass),
		)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}
	return nil
}

func downCreateAppUser(ctx context.Context, tx *sql.Tx) error {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Файл .env не найден, использую системные переменные окружения", "err", err)
	}
	username := os.Getenv("PG_USERNAME_FOR_APP")
	if username == "" {
		username = "myapp_user"
	}

	safeUser := quoteIdentifier(username)
	_, err := tx.ExecContext(
		ctx,
		fmt.Sprintf("DROP USER IF EXISTS %s", safeUser),
	)
	return err
}
