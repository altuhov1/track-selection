package postgres

import (
	"context"
	"log/slog"
	"time"
	"track-selection/internal/domain/track"

	"github.com/google/uuid"
)

func SeedTracks(ctx context.Context, repo *TrackRepository) {
	existing, _ := repo.FindAll(ctx)
	if len(existing) > 0 {
		slog.Info("Tracks already seeded, skipping")
		return
	}

	now := time.Now()

	// Бэкенд трек
	backendCurriculum := track.Curriculum{
		Semesters: []track.Semester{
			{
				Number: 1,
				Courses: []track.Course{
					{Name: "Введение в программирование", Description: "Основы Go", IsElective: false},
					{Name: "Базы данных", Description: "SQL, PostgreSQL", IsElective: false},
					{Name: "Английский язык", Description: "Технический английский", IsElective: true, Options: []string{"A2", "B1", "B2"}},
				},
			},
			{
				Number: 2,
				Courses: []track.Course{
					{Name: "HTTP серверы", Description: "REST API, Gin", IsElective: false},
					{Name: "Микросервисы", Description: "Docker, K8s", IsElective: true, Options: []string{"Docker", "Kubernetes"}},
				},
			},
		},
	}

	// Data Science трек
	dsCurriculum := track.Curriculum{
		Semesters: []track.Semester{
			{
				Number: 1,
				Courses: []track.Course{
					{Name: "Python для DS", Description: "NumPy, Pandas", IsElective: false},
					{Name: "Статистика", Description: "Теория вероятности", IsElective: false},
				},
			},
			{
				Number: 2,
				Courses: []track.Course{
					{Name: "Машинное обучение", Description: "Scikit-learn", IsElective: false},
					{Name: "Глубокое обучение", Description: "TensorFlow или PyTorch", IsElective: true, Options: []string{"TensorFlow", "PyTorch"}},
				},
			},
		},
	}

	// Фронтенд трек
	frontendCurriculum := track.Curriculum{
		Semesters: []track.Semester{
			{
				Number: 1,
				Courses: []track.Course{
					{Name: "HTML/CSS", Description: "Верстка", IsElective: false},
					{Name: "JavaScript", Description: "Основы JS", IsElective: false},
				},
			},
			{
				Number: 2,
				Courses: []track.Course{
					{Name: "React", Description: "React и экосистема", IsElective: false},
					{Name: "Vue", Description: "Vue.js", IsElective: true, Options: []string{"Vue 2", "Vue 3"}},
				},
			},
		},
	}

	tracks := []*track.Track{
		{
			ID:                  uuid.New().String(),
			Name:                "Backend-разработка",
			Description:         "Изучи создание серверной части веб-приложений",
			Curriculum:          backendCurriculum,
			Requirements:        []track.Requirement{{Subject: "programming", MinGrade: 4}},
			Teachers:            []string{"Алексей Иванов"},
			Difficulty:          3,
			Type:                1,
			EmploymentProspects: 9,
			AlumniReviews:       8,
			WebLink:             "https://example.com/backend",
			HasCertificates:     1,
			LearningStyle:       2,
			DesiredTechSkills:   8,
			DesiredMathSkills:   5,
			DesiredSoftSkills:   6,
			ProfessionalGoals:   []int{1, 2},
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  uuid.New().String(),
			Name:                "Data Science и машинное обучение",
			Description:         "Стань специалистом по анализу данных",
			Curriculum:          dsCurriculum,
			Requirements:        []track.Requirement{{Subject: "math_analysis", MinGrade: 4}, {Subject: "programming", MinGrade: 4}},
			Teachers:            []string{"Елена Козлова"},
			Difficulty:          4,
			Type:                2,
			EmploymentProspects: 9,
			AlumniReviews:       9,
			WebLink:             "https://example.com/datascience",
			HasCertificates:     1,
			LearningStyle:       3,
			DesiredTechSkills:   7,
			DesiredMathSkills:   9,
			DesiredSoftSkills:   5,
			ProfessionalGoals:   []int{2, 3},
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  uuid.New().String(),
			Name:                "Frontend-разработка",
			Description:         "Создавай красивые интерфейсы",
			Curriculum:          frontendCurriculum,
			Requirements:        []track.Requirement{{Subject: "programming", MinGrade: 3}},
			Teachers:            []string{"Анна Павлова"},
			Difficulty:          2,
			Type:                1,
			EmploymentProspects: 8,
			AlumniReviews:       8,
			WebLink:             "https://example.com/frontend",
			HasCertificates:     1,
			LearningStyle:       2,
			DesiredTechSkills:   8,
			DesiredMathSkills:   3,
			DesiredSoftSkills:   7,
			ProfessionalGoals:   []int{1, 4},
			CreatedAt:           now,
			UpdatedAt:           now,
		},
	}

	for _, t := range tracks {
		if err := repo.Save(ctx, t); err != nil {
			slog.Error("Failed to seed track", "name", t.Name, "error", err)
		} else {
			slog.Info("Seeded track", "name", t.Name, "id", t.ID)
		}
	}
}
