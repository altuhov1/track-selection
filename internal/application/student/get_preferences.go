package student

import (
	"context"
	"track-selection/internal/domain/student"
	"track-selection/internal/infrastructure/persistence/postgres"
)

type GetPreferencesUseCase struct {
	prefsRepo *postgres.PreferencesRepository
}

func NewGetPreferencesUseCase(repo *postgres.PreferencesRepository) *GetPreferencesUseCase {
	return &GetPreferencesUseCase{prefsRepo: repo}
}

func (uc *GetPreferencesUseCase) Execute(ctx context.Context, userID string) (*student.Preferences, error) {
	prefs, err := uc.prefsRepo.FindByUserID(ctx, userID)
	if err != nil {
		// Если нет предпочтений, возвращаем пустые (дефолтные)
		return &student.Preferences{
			ProfessionalGoals: []int{},
			Skills:            student.Skills{},
			LearningStyle:     1,
			Certificates:      0,
		}, nil
	}
	return prefs, nil
}
