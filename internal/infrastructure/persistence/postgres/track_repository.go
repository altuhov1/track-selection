package postgres

import (
	"context"
	"encoding/json"
	"track-selection/internal/domain/track"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackRepository struct {
	pool *pgxpool.Pool
}

func NewTrackRepository(pool *pgxpool.Pool) *TrackRepository {
	return &TrackRepository{pool: pool}
}

func (r *TrackRepository) Save(ctx context.Context, t *track.Track) error {
	curriculumJSON, err := json.Marshal(t.Curriculum)
	if err != nil {
		return err
	}
	requirementsJSON, err := json.Marshal(t.Requirements)
	if err != nil {
		return err
	}
	teachersJSON, err := json.Marshal(t.Teachers)
	if err != nil {
		return err
	}
	goalsJSON, err := json.Marshal(t.ProfessionalGoals)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO tracks (
			id, name, description, curriculum, requirements, teachers,
			difficulty, type, employment_prospects, alumni_reviews, web_link,
			has_certificates, learning_style,
			desired_tech_skills, desired_math_skills, desired_soft_skills,
			professional_goals, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			curriculum = EXCLUDED.curriculum,
			requirements = EXCLUDED.requirements,
			teachers = EXCLUDED.teachers,
			difficulty = EXCLUDED.difficulty,
			type = EXCLUDED.type,
			employment_prospects = EXCLUDED.employment_prospects,
			alumni_reviews = EXCLUDED.alumni_reviews,
			web_link = EXCLUDED.web_link,
			has_certificates = EXCLUDED.has_certificates,
			learning_style = EXCLUDED.learning_style,
			desired_tech_skills = EXCLUDED.desired_tech_skills,
			desired_math_skills = EXCLUDED.desired_math_skills,
			desired_soft_skills = EXCLUDED.desired_soft_skills,
			professional_goals = EXCLUDED.professional_goals,
			updated_at = EXCLUDED.updated_at
	`

	_, err = r.pool.Exec(ctx, query,
		t.ID, t.Name, t.Description, curriculumJSON, requirementsJSON, teachersJSON,
		t.Difficulty, t.Type, t.EmploymentProspects, t.AlumniReviews, t.WebLink,
		t.HasCertificates, t.LearningStyle,
		t.DesiredTechSkills, t.DesiredMathSkills, t.DesiredSoftSkills,
		goalsJSON, t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *TrackRepository) FindByID(ctx context.Context, id string) (*track.Track, error) {
	query := `SELECT * FROM tracks WHERE id = $1`

	var t track.Track
	var curriculumJSON, requirementsJSON, teachersJSON, goalsJSON []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Name, &t.Description, &curriculumJSON, &requirementsJSON,
		&teachersJSON, &t.Difficulty, &t.Type, &t.EmploymentProspects, &t.AlumniReviews, &t.WebLink,
		&t.HasCertificates, &t.LearningStyle,
		&t.DesiredTechSkills, &t.DesiredMathSkills, &t.DesiredSoftSkills,
		&goalsJSON, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(curriculumJSON, &t.Curriculum); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(requirementsJSON, &t.Requirements); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(teachersJSON, &t.Teachers); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(goalsJSON, &t.ProfessionalGoals); err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TrackRepository) FindAll(ctx context.Context) ([]*track.Track, error) {
	query := `SELECT * FROM tracks ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*track.Track
	for rows.Next() {
		var t track.Track
		var curriculumJSON, requirementsJSON, teachersJSON, goalsJSON []byte

		err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &curriculumJSON, &requirementsJSON,
			&teachersJSON, &t.Difficulty, &t.Type, &t.EmploymentProspects, &t.AlumniReviews, &t.WebLink,
			&t.HasCertificates, &t.LearningStyle,
			&t.DesiredTechSkills, &t.DesiredMathSkills, &t.DesiredSoftSkills,
			&goalsJSON, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(curriculumJSON, &t.Curriculum); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(requirementsJSON, &t.Requirements); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(teachersJSON, &t.Teachers); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(goalsJSON, &t.ProfessionalGoals); err != nil {
			return nil, err
		}

		tracks = append(tracks, &t)
	}

	return tracks, nil
}

func (r *TrackRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tracks WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
