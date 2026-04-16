package track

import (
	"context"
	"track-selection/internal/domain/shared/errors"
	"track-selection/internal/infrastructure/persistence/postgres"
)

type DeleteTrackUseCase struct {
	trackRepo *postgres.TrackRepository
}

func NewDeleteTrackUseCase(repo *postgres.TrackRepository) *DeleteTrackUseCase {
	return &DeleteTrackUseCase{trackRepo: repo}
}

func (uc *DeleteTrackUseCase) Execute(ctx context.Context, id string) error {
	// Проверяем, существует ли трек
	_, err := uc.trackRepo.FindByID(ctx, id)
	if err != nil {
		return errors.ErrNotFound
	}

	// Удаляем трек
	return uc.trackRepo.Delete(ctx, id)
}
