package postgres

import (
	"context"
	"errors"
	"fmt"
	"track-selection/internal/domain/auth"
	"track-selection/internal/domain/shared/value_objects"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}

func (r *AuthRepository) Save(ctx context.Context, user *auth.AuthUser) error {
	query := `
        INSERT INTO auth_users (id, email, password_hash, role, first_name, last_name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Email.String(),
		user.PasswordHash,
		string(user.Role),
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save auth user: %w", err)
	}

	return nil
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email value_objects.Email) (*auth.AuthUser, error) {
	query := `
		SELECT id, email, password_hash, role, first_name, last_name, created_at, updated_at
		FROM auth_users
		WHERE email = $1
	`

	var user auth.AuthUser
	var emailStr string
	var roleStr string

	err := r.pool.QueryRow(ctx, query, email.String()).Scan(
		&user.ID,
		&emailStr,
		&user.PasswordHash,
		&roleStr,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find auth user by email: %w", err)
	}

	parsedEmail, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}
	user.Email = parsedEmail
	user.Role = auth.UserRole(roleStr)

	return &user, nil
}

func (r *AuthRepository) FindByID(ctx context.Context, id string) (*auth.AuthUser, error) {
	query := `
		SELECT id, email, password_hash, role, first_name, last_name, created_at, updated_at
		FROM auth_users
		WHERE id = $1
	`

	var user auth.AuthUser
	var emailStr string
	var roleStr string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&emailStr,
		&user.PasswordHash,
		&roleStr,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find auth user by id: %w", err)
	}

	parsedEmail, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}
	user.Email = parsedEmail
	user.Role = auth.UserRole(roleStr)

	return &user, nil
}

func (r *AuthRepository) ExistsByEmail(ctx context.Context, email value_objects.Email) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM auth_users WHERE email = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, email.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}
