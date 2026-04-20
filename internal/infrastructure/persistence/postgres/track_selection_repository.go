package postgres

import (
	"context"
	"fmt"
	"track-selection/internal/domain/student"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackSelectionRepository struct {
	pool *pgxpool.Pool
}

func NewTrackSelectionRepository(pool *pgxpool.Pool) *TrackSelectionRepository {
	return &TrackSelectionRepository{pool: pool}
}

func (r *TrackSelectionRepository) Save(ctx context.Context, selection *student.TrackSelection) error {
	query := `
        INSERT INTO student_track_selections (id, student_id, track_id, selected_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (student_id, track_id) DO UPDATE SET
            selected_at = EXCLUDED.selected_at,
            updated_at = EXCLUDED.updated_at
    `
	_, err := r.pool.Exec(ctx, query,
		selection.ID, selection.StudentID, selection.TrackID,
		selection.SelectedAt, selection.CreatedAt, selection.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save track selection: %w", err)
	}
	return nil
}

func (r *TrackSelectionRepository) FindByStudentID(ctx context.Context, studentID string) ([]*student.TrackSelection, error) {
	query := `SELECT id, student_id, track_id, selected_at, created_at, updated_at FROM student_track_selections WHERE student_id = $1`

	rows, err := r.pool.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var selections []*student.TrackSelection
	for rows.Next() {
		var s student.TrackSelection
		err := rows.Scan(&s.ID, &s.StudentID, &s.TrackID, &s.SelectedAt, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		selections = append(selections, &s)
	}
	return selections, nil
}

func (r *TrackSelectionRepository) Delete(ctx context.Context, studentID, trackID string) error {
	query := `DELETE FROM student_track_selections WHERE student_id = $1 AND track_id = $2`
	_, err := r.pool.Exec(ctx, query, studentID, trackID)
	return err
}

func (r *TrackSelectionRepository) Exists(ctx context.Context, studentID, trackID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM student_track_selections WHERE student_id = $1 AND track_id = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, studentID, trackID).Scan(&exists)
	return exists, err
}
