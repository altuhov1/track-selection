package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"
	"track-selection/internal/domain/admin"
	"track-selection/internal/domain/shared/value_objects"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	pool *pgxpool.Pool
}

func NewAdminRepository(pool *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{pool: pool}
}

func (r *AdminRepository) Save(ctx context.Context, a *admin.Admin) error {
	query := `
        INSERT INTO admins (id, auth_user_id, email, first_name, last_name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE SET
            email = EXCLUDED.email,
            first_name = EXCLUDED.first_name,
            last_name = EXCLUDED.last_name,
            updated_at = EXCLUDED.updated_at
    `

	_, err := r.pool.Exec(ctx, query,
		a.ID().String(),
		a.AuthUserID(),
		a.Email().String(),
		a.FirstName(),
		a.LastName(),
		a.CreatedAt(),
		a.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save admin: %w", err)
	}

	return nil
}

func (r *AdminRepository) FindByID(ctx context.Context, id admin.AdminID) (*admin.Admin, error) {
	query := `
        SELECT id, auth_user_id, email, first_name, last_name, created_at, updated_at
        FROM admins
        WHERE id = $1
    `

	var (
		idStr      string
		authUserID string
		emailStr   string
		firstName  string
		lastName   string
		createdAt  time.Time
		updatedAt  time.Time
	)

	err := r.pool.QueryRow(ctx, query, id.String()).Scan(
		&idStr,
		&authUserID,
		&emailStr,
		&firstName,
		&lastName,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find admin by id: %w", err)
	}

	email, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	adminID, err := admin.AdminIDFromString(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid admin id in database: %w", err)
	}

	return admin.NewAdminFromDB(
		adminID,
		authUserID,
		email,
		firstName,
		lastName,
		createdAt,
		updatedAt,
	), nil
}

func (r *AdminRepository) FindByEmail(ctx context.Context, email value_objects.Email) (*admin.Admin, error) {
	query := `
        SELECT id, auth_user_id, email, first_name, last_name, created_at, updated_at
        FROM admins
        WHERE email = $1
    `

	var (
		idStr      string
		authUserID string
		emailStr   string
		firstName  string
		lastName   string
		createdAt  time.Time
		updatedAt  time.Time
	)

	err := r.pool.QueryRow(ctx, query, email.String()).Scan(
		&idStr,
		&authUserID,
		&emailStr,
		&firstName,
		&lastName,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find admin by email: %w", err)
	}

	adminID, err := admin.AdminIDFromString(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid admin id in database: %w", err)
	}

	parsedEmail, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	return admin.NewAdminFromDB(
		adminID,
		authUserID,
		parsedEmail,
		firstName,
		lastName,
		createdAt,
		updatedAt,
	), nil
}

func (r *AdminRepository) FindByAuthUserID(ctx context.Context, authUserID string) (*admin.Admin, error) {
	query := `
        SELECT id, auth_user_id, email, first_name, last_name, created_at, updated_at
        FROM admins
        WHERE auth_user_id = $1
    `

	var (
		idStr         string
		authUserIDOut string
		emailStr      string
		firstName     string
		lastName      string
		createdAt     time.Time
		updatedAt     time.Time
	)

	err := r.pool.QueryRow(ctx, query, authUserID).Scan(
		&idStr,
		&authUserIDOut,
		&emailStr,
		&firstName,
		&lastName,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find admin by auth_user_id: %w", err)
	}

	adminID, err := admin.AdminIDFromString(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid admin id in database: %w", err)
	}

	parsedEmail, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	return admin.NewAdminFromDB(
		adminID,
		authUserIDOut,
		parsedEmail,
		firstName,
		lastName,
		createdAt,
		updatedAt,
	), nil
}

func (r *AdminRepository) ExistsByEmail(ctx context.Context, email value_objects.Email) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM admins WHERE email = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, email.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}
