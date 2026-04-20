package student

import (
	"context"
	"track-selection/internal/domain/student"
	"track-selection/internal/infrastructure/persistence/postgres"
)

type GetProfileCompletionUseCase struct {
	profileRepo *postgres.ProfileCompletionRepository
}

func NewGetProfileCompletionUseCase(repo *postgres.ProfileCompletionRepository) *GetProfileCompletionUseCase {
	return &GetProfileCompletionUseCase{profileRepo: repo}
}

func (uc *GetProfileCompletionUseCase) Execute(ctx context.Context, userID string) (*student.ProfileCompletion, error) {
	status, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return &student.ProfileCompletion{
			UserID:     userID,
			IsComplete: false,
		}, nil
	}
	return status, nil
}
