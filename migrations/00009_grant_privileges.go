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
	goose.AddMigrationContext(upGrantPrivileges, downGrantPrivileges)
}

func quotePostgresIdentifier(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

func upGrantPrivileges(ctx context.Context, tx *sql.Tx) error {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Файл .env не найден, использую системные переменные окружения")
	}
	username := os.Getenv("PG_USERNAME_FOR_APP")
	if username == "" {
		return fmt.Errorf("PG_USERNAME_FOR_APP is not set")
	}
	quotedUser := quotePostgresIdentifier(username)

	_, err := tx.ExecContext(ctx, fmt.Sprintf(`
		GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE auth_users TO %s;
	`, quotedUser))
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE students TO %s;
	`, quotedUser))
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE admins TO %s;
	`, quotedUser))
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE user_preferences TO %s;
	`, quotedUser))
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE profile_completion TO %s;
	`, quotedUser))
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE tracks TO %s;
	`, quotedUser))
	if err != nil {
		return err
	}
	// _, err = tx.ExecContext(ctx, fmt.Sprintf(`
	// 	GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE track_teachers TO %s;
	// `, quotedUser))
	// if err != nil {
	// 	return err
	// }
	fmt.Printf("Granted privileges to user: %s\n", username)
	return nil
}

func downGrantPrivileges(ctx context.Context, tx *sql.Tx) error {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Файл .env не найден, использую системные переменные окружения")
	}
	username := os.Getenv("PG_USERNAME_FOR_APP")
	if username == "" {
		username = "myapp_user"
	}
	quotedUser := quotePostgresIdentifier(username)

	_, err := tx.ExecContext(ctx, fmt.Sprintf(`
		REVOKE ALL PRIVILEGES ON TABLE auth_users FROM %s;
	`, quotedUser))
	if err != nil {
		if strings.Contains(err.Error(), "undefined_object") {
			fmt.Printf("Privileges already revoked for user: %s\n", username)
			return nil
		}
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		REVOKE ALL PRIVILEGES ON TABLE students FROM %s;
	`, quotedUser))
	if err != nil {
		if strings.Contains(err.Error(), "undefined_object") {
			fmt.Printf("Privileges already revoked for user: %s\n", username)
			return nil
		}
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		REVOKE ALL PRIVILEGES ON TABLE admins FROM %s;
	`, quotedUser))
	if err != nil {
		if strings.Contains(err.Error(), "undefined_object") {
			fmt.Printf("Privileges already revoked for user: %s\n", username)
			return nil
		}
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		REVOKE ALL PRIVILEGES ON TABLE user_preferences FROM %s;
	`, quotedUser))
	if err != nil {
		if strings.Contains(err.Error(), "undefined_object") {
			fmt.Printf("Privileges already revoked for bookings: %s\n", username)
			return nil
		}
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		REVOKE ALL PRIVILEGES ON TABLE profile_completion FROM %s;
	`, quotedUser))
	if err != nil {
		if strings.Contains(err.Error(), "undefined_object") {
			fmt.Printf("Privileges already revoked for bookings: %s\n", username)
			return nil
		}
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		REVOKE ALL PRIVILEGES ON TABLE tracks FROM %s;
	`, quotedUser))
	if err != nil {
		if strings.Contains(err.Error(), "undefined_object") {
			fmt.Printf("Privileges already revoked for bookings: %s\n", username)
			return nil
		}
		return err
	}

	fmt.Printf("Revoked privileges from user: %s\n", username)
	return nil
}
