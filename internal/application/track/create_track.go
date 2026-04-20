package track

import (
	"context"
	"track-selection/internal/domain/track"
	"track-selection/internal/infrastructure/persistence/postgres"
)

type CreateTrackUseCase struct {
	trackRepo *postgres.TrackRepository
}

func NewCreateTrackUseCase(repo *postgres.TrackRepository) *CreateTrackUseCase {
	return &CreateTrackUseCase{trackRepo: repo}
}

type CreateTrackInput struct {
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	Curriculum          track.Curriculum    `json:"curriculum"`
	Requirements        []track.Requirement `json:"requirements"`
	Teachers            []string            `json:"teachers"`
	Difficulty          int                 `json:"difficulty"`
	Type                int                 `json:"type"`
	EmploymentProspects int                 `json:"employment_prospects"`
	AlumniReviews       int                 `json:"alumni_reviews"`
	WebLink             string              `json:"web_link"`
	HasCertificates     int                 `json:"has_certificates"`
	LearningStyle       int                 `json:"learning_style"`
	DesiredTechSkills   int                 `json:"desired_tech_skills"`
	DesiredMathSkills   int                 `json:"desired_math_skills"`
	DesiredSoftSkills   int                 `json:"desired_soft_skills"`
	ProfessionalGoals   []int               `json:"professional_goals"`
}

func (uc *CreateTrackUseCase) Execute(ctx context.Context, input CreateTrackInput) (*track.Track, error) {
	t, err := track.NewTrack(
		input.Name,
		input.Description,
		input.Curriculum,
		input.Requirements,
		input.Teachers,
		input.Difficulty,
		input.Type,
		input.EmploymentProspects,
		input.AlumniReviews,
		input.WebLink,
		input.HasCertificates,
		input.LearningStyle,
		input.DesiredTechSkills,
		input.DesiredMathSkills,
		input.DesiredSoftSkills,
		input.ProfessionalGoals,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.trackRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}
