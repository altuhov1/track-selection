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

// Save сохраняет или обновляет студента
func (r *AdminRepository) Save(ctx context.Context, s *admin.Admin) error {
	query := `
        INSERT INTO admins (id, auth_user_id, email, username, rating, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE SET
            email = EXCLUDED.email,
            username = EXCLUDED.username,
            rating = EXCLUDED.rating,
            updated_at = EXCLUDED.updated_at
    `

	_, err := r.pool.Exec(ctx, query,
		s.ID().String(),
		s.AuthUserID(),
		s.Email().String(),
		s.Username(),
		s.Rating(),
		s.CreatedAt(),
		s.UpdatedAt(),
	)
	fmt.Println(err)
	if err != nil {
		return fmt.Errorf("failed to save admin: %w", err)
	}

	return nil
}

// FindByID находит студента по ID
func (r *AdminRepository) FindByID(ctx context.Context, id admin.AdminID) (*admin.Admin, error) {
	query := `
        SELECT id, auth_user_id, email, username, rating, created_at, updated_at
        FROM admins
        WHERE id = $1
    `

	var (
		idStr      string
		authUserID string
		emailStr   string
		username   string
		rating     int
		createdAt  time.Time
		updatedAt  time.Time
	)

	err := r.pool.QueryRow(ctx, query, id.String()).Scan(
		&idStr,
		&authUserID,
		&emailStr,
		&username,
		&rating,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find admin by id: %w", err)
	}

	// Парсим Email
	email, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	adminID, err := admin.AdminIDFromString(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid admin id in database: %w", err)
	}

	// Создаем студента через конструктор для восстановления из БД
	return admin.NewAdminFromDB(
		adminID,
		authUserID,
		email,
		username,
		rating,
		createdAt,
		updatedAt,
	), nil
}

// FindByEmail находит студента по email
func (r *AdminRepository) FindByEmail(ctx context.Context, email value_objects.Email) (*admin.Admin, error) {
	query := `
        SELECT id, auth_user_id, email, username, rating, created_at, updated_at
        FROM admins
        WHERE email = $1
    `

	var (
		idStr      string
		authUserID string
		emailStr   string
		username   string
		rating     int
		createdAt  time.Time
		updatedAt  time.Time
	)

	err := r.pool.QueryRow(ctx, query, email.String()).Scan(
		&idStr,
		&authUserID,
		&emailStr,
		&username,
		&rating,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find admin by email: %w", err)
	}

	// Парсим AdminID
	adminID, err := admin.AdminIDFromString(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid admin id in database: %w", err)
	}

	// Email уже есть, но для консистенции парсим заново
	parsedEmail, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	return admin.NewAdminFromDB(
		adminID,
		authUserID,
		parsedEmail,
		username,
		rating,
		createdAt,
		updatedAt,
	), nil
}

// FindByAuthAdminID находит студента по ID учетной записи
func (r *AdminRepository) FindByAuthAdminID(ctx context.Context, authUserID string) (*admin.Admin, error) {
	query := `
        SELECT id, auth_user_id, email, username, rating, created_at, updated_at
        FROM admins
        WHERE auth_user_id = $1
    `

	var (
		idStr         string
		authUserIDOut string
		emailStr      string
		username      string
		rating        int
		createdAt     time.Time
		updatedAt     time.Time
	)

	err := r.pool.QueryRow(ctx, query, authUserID).Scan(
		&idStr,
		&authUserIDOut,
		&emailStr,
		&username,
		&rating,
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
		username,
		rating,
		createdAt,
		updatedAt,
	), nil
}

// ExistsByEmail проверяет существование студента по email
func (r *AdminRepository) ExistsByEmail(ctx context.Context, email value_objects.Email) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM admins WHERE email = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, email.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// UpdateRating обновляет рейтинг студента
func (r *AdminRepository) UpdateRating(ctx context.Context, adminID admin.AdminID, rating int) error {
	query := `UPDATE admins SET rating = $1, updated_at = $2 WHERE id = $3`

	_, err := r.pool.Exec(ctx, query, rating, time.Now(), adminID.String())
	if err != nil {
		return fmt.Errorf("failed to update rating: %w", err)
	}

	return nil
}
