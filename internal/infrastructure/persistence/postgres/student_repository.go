package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"
	"track-selection/internal/domain/shared/value_objects"
	"track-selection/internal/domain/student"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository struct {
	pool *pgxpool.Pool
}

func NewStudentRepository(pool *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{pool: pool}
}

// Save сохраняет или обновляет студента
func (r *StudentRepository) Save(ctx context.Context, s *student.Student) error {
	query := `
        INSERT INTO students (id, auth_user_id, email, first_name, last_name, username, rating, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (id) DO UPDATE SET
            email = EXCLUDED.email,
            first_name = EXCLUDED.first_name,
            last_name = EXCLUDED.last_name,
            username = EXCLUDED.username,
            rating = EXCLUDED.rating,
            updated_at = EXCLUDED.updated_at
    `

	_, err := r.pool.Exec(ctx, query,
		s.ID().String(),
		s.AuthUserID(),
		s.Email().String(),
		s.FirstName(),
		s.LastName(),
		s.Username(),
		s.Rating(),
		s.CreatedAt(),
		s.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to save student: %w", err)
	}

	return nil
}

// FindByID находит студента по ID
func (r *StudentRepository) FindByID(ctx context.Context, id student.StudentID) (*student.Student, error) {
	query := `
        SELECT id, auth_user_id, email, first_name, last_name, username, rating, created_at, updated_at
        FROM students
        WHERE id = $1
    `

	var (
		idStr      string
		authUserID string
		emailStr   string
		firstName  string
		lastName   string
		username   string
		rating     int
		createdAt  time.Time
		updatedAt  time.Time
	)

	err := r.pool.QueryRow(ctx, query, id.String()).Scan(
		&idStr,
		&authUserID,
		&emailStr,
		&firstName,
		&lastName,
		&username,
		&rating,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find student by id: %w", err)
	}

	// Парсим Email
	email, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	// Парсим StudentID
	studentID, err := student.StudentIDFromString(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid student id in database: %w", err)
	}

	// Создаем студента через конструктор для восстановления из БД
	return student.NewStudentFromDB(
		studentID,
		authUserID,
		email,
		firstName,
		lastName,
		username,
		rating,
		createdAt,
		updatedAt,
	), nil
}

// FindByEmail находит студента по email
func (r *StudentRepository) FindByEmail(ctx context.Context, email value_objects.Email) (*student.Student, error) {
	query := `
        SELECT id, auth_user_id, email, first_name, last_name, username, rating, created_at, updated_at
        FROM students
        WHERE email = $1
    `

	var (
		idStr      string
		authUserID string
		emailStr   string
		firstName  string
		lastName   string
		username   string
		rating     int
		createdAt  time.Time
		updatedAt  time.Time
	)

	err := r.pool.QueryRow(ctx, query, email.String()).Scan(
		&idStr,
		&authUserID,
		&emailStr,
		&firstName,
		&lastName,
		&username,
		&rating,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find student by email: %w", err)
	}

	// Парсим StudentID
	studentID, err := student.StudentIDFromString(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid student id in database: %w", err)
	}

	// Email уже есть, но для консистенции парсим заново
	parsedEmail, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	return student.NewStudentFromDB(
		studentID,
		authUserID,
		parsedEmail,
		firstName,
		lastName,
		username,
		rating,
		createdAt,
		updatedAt,
	), nil
}

// FindByAuthUserID находит студента по ID учетной записи
func (r *StudentRepository) FindByAuthUserID(ctx context.Context, authUserID string) (*student.Student, error) {
	query := `
        SELECT id, auth_user_id, email, first_name, last_name, username, rating, created_at, updated_at
        FROM students
        WHERE auth_user_id = $1
    `

	var (
		idStr         string
		authUserIDOut string
		emailStr      string
		firstName     string
		lastName      string
		username      string
		rating        int
		createdAt     time.Time
		updatedAt     time.Time
	)

	err := r.pool.QueryRow(ctx, query, authUserID).Scan(
		&idStr,
		&authUserIDOut,
		&emailStr,
		&firstName,
		&lastName,
		&username,
		&rating,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find student by auth_user_id: %w", err)
	}

	studentID, err := student.StudentIDFromString(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid student id in database: %w", err)
	}

	parsedEmail, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	return student.NewStudentFromDB(
		studentID,
		authUserIDOut,
		parsedEmail,
		firstName,
		lastName,
		username,
		rating,
		createdAt,
		updatedAt,
	), nil
}

// ExistsByEmail проверяет существование студента по email
func (r *StudentRepository) ExistsByEmail(ctx context.Context, email value_objects.Email) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM students WHERE email = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, email.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// UpdateRating обновляет рейтинг студента
func (r *StudentRepository) UpdateRating(ctx context.Context, studentID student.StudentID, rating int) error {
	query := `UPDATE students SET rating = $1, updated_at = $2 WHERE id = $3`

	_, err := r.pool.Exec(ctx, query, rating, time.Now(), studentID.String())
	if err != nil {
		return fmt.Errorf("failed to update rating: %w", err)
	}

	return nil
}
