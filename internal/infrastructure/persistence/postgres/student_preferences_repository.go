package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"track-selection/internal/domain/shared/errors"
	"track-selection/internal/domain/student"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PreferencesRepository struct {
	pool *pgxpool.Pool
}

func NewPreferencesRepository(pool *pgxpool.Pool) *PreferencesRepository {
	return &PreferencesRepository{pool: pool}
}

func (r *PreferencesRepository) Save(ctx context.Context, prefs *student.Preferences) error {
	// Преобразуем ProfessionalGoals в JSON
	goalsJSON, err := json.Marshal(prefs.ProfessionalGoals)
	if err != nil {
		return errors.ErrInvalidRequest
	}

	query := `
		INSERT INTO user_preferences (
			id, user_id, professional_goals,
			grades_informatics, grades_programming, grades_foreign_language,
			grades_physics, grades_aig, grades_math_analysis,
			grades_algorithms_data_structures, grades_databases, grades_discrete_math,
			skills_databases, skills_system_architecture, skills_algorithmic_programming,
			skills_public_speaking, skills_testing, skills_analytics,
			skills_machine_learning, skills_os_knowledge, skills_research_projects,
			learning_style, certificates, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
		ON CONFLICT (id) DO UPDATE SET
			professional_goals = EXCLUDED.professional_goals,
			grades_informatics = EXCLUDED.grades_informatics,
			grades_programming = EXCLUDED.grades_programming,
			grades_foreign_language = EXCLUDED.grades_foreign_language,
			grades_physics = EXCLUDED.grades_physics,
			grades_aig = EXCLUDED.grades_aig,
			grades_math_analysis = EXCLUDED.grades_math_analysis,
			grades_algorithms_data_structures = EXCLUDED.grades_algorithms_data_structures,
			grades_databases = EXCLUDED.grades_databases,
			grades_discrete_math = EXCLUDED.grades_discrete_math,
			skills_databases = EXCLUDED.skills_databases,
			skills_system_architecture = EXCLUDED.skills_system_architecture,
			skills_algorithmic_programming = EXCLUDED.skills_algorithmic_programming,
			skills_public_speaking = EXCLUDED.skills_public_speaking,
			skills_testing = EXCLUDED.skills_testing,
			skills_analytics = EXCLUDED.skills_analytics,
			skills_machine_learning = EXCLUDED.skills_machine_learning,
			skills_os_knowledge = EXCLUDED.skills_os_knowledge,
			skills_research_projects = EXCLUDED.skills_research_projects,
			learning_style = EXCLUDED.learning_style,
			certificates = EXCLUDED.certificates,
			updated_at = EXCLUDED.updated_at
	`

	_, err = r.pool.Exec(ctx, query,
		prefs.ID, prefs.UserID, goalsJSON,
		prefs.Grades.Informatics, prefs.Grades.Programming, prefs.Grades.ForeignLanguage,
		prefs.Grades.Physics, prefs.Grades.AIG, prefs.Grades.MathAnalysis,
		prefs.Grades.AlgorithmsDataStructures, prefs.Grades.Databases, prefs.Grades.DiscreteMath,
		prefs.Skills.Databases, prefs.Skills.SystemArchitecture, prefs.Skills.AlgorithmicProgramming,
		prefs.Skills.PublicSpeaking, prefs.Skills.Testing, prefs.Skills.Analytics,
		prefs.Skills.MachineLearning, prefs.Skills.OSKnowledge, prefs.Skills.ResearchProjects,
		prefs.LearningStyle, prefs.Certificates, prefs.CreatedAt, prefs.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save preferences: %w", err)
	}

	return nil
}

// FindByUserID находит предпочтения по user_id
func (r *PreferencesRepository) FindByUserID(ctx context.Context, userID string) (*student.Preferences, error) {
	query := `
		SELECT 
			id, user_id, professional_goals,
			grades_informatics, grades_programming, grades_foreign_language,
			grades_physics, grades_aig, grades_math_analysis,
			grades_algorithms_data_structures, grades_databases, grades_discrete_math,
			skills_databases, skills_system_architecture, skills_algorithmic_programming,
			skills_public_speaking, skills_testing, skills_analytics,
			skills_machine_learning, skills_os_knowledge, skills_research_projects,
			learning_style, certificates, created_at, updated_at
		FROM user_preferences
		WHERE user_id = $1
	`

	var prefs student.Preferences
	var goalsJSON []byte

	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&prefs.ID, &prefs.UserID, &goalsJSON,
		&prefs.Grades.Informatics, &prefs.Grades.Programming, &prefs.Grades.ForeignLanguage,
		&prefs.Grades.Physics, &prefs.Grades.AIG, &prefs.Grades.MathAnalysis,
		&prefs.Grades.AlgorithmsDataStructures, &prefs.Grades.Databases, &prefs.Grades.DiscreteMath,
		&prefs.Skills.Databases, &prefs.Skills.SystemArchitecture, &prefs.Skills.AlgorithmicProgramming,
		&prefs.Skills.PublicSpeaking, &prefs.Skills.Testing, &prefs.Skills.Analytics,
		&prefs.Skills.MachineLearning, &prefs.Skills.OSKnowledge, &prefs.Skills.ResearchProjects,
		&prefs.LearningStyle, &prefs.Certificates, &prefs.CreatedAt, &prefs.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Преобразуем JSON обратно в массив
	if err := json.Unmarshal(goalsJSON, &prefs.ProfessionalGoals); err != nil {
		return nil, err
	}

	return &prefs, nil
}
