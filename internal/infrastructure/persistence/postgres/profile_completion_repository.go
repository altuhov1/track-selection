package postgres

import (
	"context"
	"database/sql"
	"track-selection/internal/domain/student"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileCompletionRepository struct {
	pool *pgxpool.Pool
}

func NewProfileCompletionRepository(pool *pgxpool.Pool) *ProfileCompletionRepository {
	return &ProfileCompletionRepository{pool: pool}
}

func (r *ProfileCompletionRepository) Save(ctx context.Context, pc *student.ProfileCompletion) error {
	query := `
        INSERT INTO profile_completion (id, user_id, is_complete, completed_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (id) DO UPDATE SET
            is_complete = EXCLUDED.is_complete,
            completed_at = EXCLUDED.completed_at,
            updated_at = EXCLUDED.updated_at
    `

	var completedAt interface{}
	if pc.CompletedAt != nil {
		completedAt = *pc.CompletedAt
	} else {
		completedAt = nil
	}

	_, err := r.pool.Exec(ctx, query,
		pc.ID, pc.UserID, pc.IsComplete, completedAt, pc.CreatedAt, pc.UpdatedAt,
	)
	return err
}

func (r *ProfileCompletionRepository) FindByUserID(ctx context.Context, userID string) (*student.ProfileCompletion, error) {
	query := `
        SELECT id, user_id, is_complete, completed_at, created_at, updated_at
        FROM profile_completion
        WHERE user_id = $1
    `

	var pc student.ProfileCompletion
	var completedAt sql.NullTime

	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&pc.ID, &pc.UserID, &pc.IsComplete, &completedAt, &pc.CreatedAt, &pc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		pc.CompletedAt = &completedAt.Time
	}

	return &pc, nil
}

func (r *ProfileCompletionRepository) CreateDefault(ctx context.Context, userID string) error {
	pc := student.NewProfileCompletion(userID)
	return r.Save(ctx, pc)
}
