package student

import (
	"context"
	"encoding/json"
	"track-selection/internal/domain/selection"
	"track-selection/internal/domain/shared/errors"
	"track-selection/internal/domain/student"
	"track-selection/internal/infrastructure/persistence/postgres"
)

type GetRecommendationsUseCase struct {
	prefsRepo      *postgres.PreferencesRepository
	trackRepo      *postgres.TrackRepository
	profileChecker *student.ProfileChecker
}

func NewGetRecommendationsUseCase(
	prefsRepo *postgres.PreferencesRepository,
	trackRepo *postgres.TrackRepository,
	profileChecker *student.ProfileChecker,
) *GetRecommendationsUseCase {
	return &GetRecommendationsUseCase{
		prefsRepo:      prefsRepo,
		trackRepo:      trackRepo,
		profileChecker: profileChecker,
	}
}

type GetRecommendationsOutput struct {
	Recommendations []selection.TrackScore `json:"recommendations"`
}

func (uc *GetRecommendationsUseCase) Execute(ctx context.Context, userID string) (*GetRecommendationsOutput, error) {
	prefs, err := uc.prefsRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !uc.profileChecker.IsProfileComplete(prefs) {
		return nil, errors.ErrProfileNotComplete
	}

	tracks, err := uc.trackRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	studentData := selection.StudentData{
		ProfessionalGoals: prefs.ProfessionalGoals,
		Grades: selection.Grades{
			Informatics:              prefs.Grades.Informatics,
			Programming:              prefs.Grades.Programming,
			ForeignLanguage:          prefs.Grades.ForeignLanguage,
			Physics:                  prefs.Grades.Physics,
			AIG:                      prefs.Grades.AIG,
			MathAnalysis:             prefs.Grades.MathAnalysis,
			AlgorithmsDataStructures: prefs.Grades.AlgorithmsDataStructures,
			Databases:                prefs.Grades.Databases,
			DiscreteMath:             prefs.Grades.DiscreteMath,
		},
		Skills: selection.Skills{
			Databases:              prefs.Skills.Databases,
			SystemArchitecture:     prefs.Skills.SystemArchitecture,
			AlgorithmicProgramming: prefs.Skills.AlgorithmicProgramming,
			PublicSpeaking:         prefs.Skills.PublicSpeaking,
			Testing:                prefs.Skills.Testing,
			Analytics:              prefs.Skills.Analytics,
			MachineLearning:        prefs.Skills.MachineLearning,
			OSKnowledge:            prefs.Skills.OSKnowledge,
			ResearchProjects:       prefs.Skills.ResearchProjects,
		},
		LearningStyle:     prefs.LearningStyle,
		DesiredTechSkills: prefs.Skills.Databases,
		DesiredMathSkills: prefs.Skills.Analytics,
		DesiredSoftSkills: prefs.Skills.PublicSpeaking,
	}

	var prometheeInputs []selection.PrometheeInput
	for _, t := range tracks {
		var requirements []selection.Requirement
		if len(t.Requirements) > 0 {
			reqJSON, _ := json.Marshal(t.Requirements)
			json.Unmarshal(reqJSON, &requirements)
		}

		prometheeInputs = append(prometheeInputs, selection.PrometheeInput{
			TrackID:           t.ID,
			TrackName:         t.Name,
			ProfessionalGoals: t.ProfessionalGoals,
			Employment:        t.EmploymentProspects,
			AlumniReviews:     t.AlumniReviews,
			Difficulty:        t.Difficulty,
			HasCertificates:   t.HasCertificates,
			LearningStyle:     t.LearningStyle,
			DesiredTechSkills: t.DesiredTechSkills,
			DesiredMathSkills: t.DesiredMathSkills,
			DesiredSoftSkills: t.DesiredSoftSkills,
			Requirements:      requirements,
		})
	}

	calculator := selection.NewPrometheeCalculator(selection.DefaultWeights())
	scores := calculator.CalculateScores(prometheeInputs, studentData)

	return &GetRecommendationsOutput{
		Recommendations: scores,
	}, nil
}
