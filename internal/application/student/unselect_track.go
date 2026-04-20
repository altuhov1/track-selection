package student

import (
	"context"
	"track-selection/internal/domain/shared/errors"
	"track-selection/internal/infrastructure/persistence/postgres"
)

type UnselectTrackUseCase struct {
	selectionRepo *postgres.TrackSelectionRepository
	studentRepo   *postgres.StudentRepository
}

func NewUnselectTrackUseCase(
	selectionRepo *postgres.TrackSelectionRepository,
	studentRepo *postgres.StudentRepository,
) *UnselectTrackUseCase {
	return &UnselectTrackUseCase{
		selectionRepo: selectionRepo,
		studentRepo:   studentRepo,
	}
}

func (uc *UnselectTrackUseCase) Execute(ctx context.Context, authUserID, trackID string) error {
	// Находим студента по auth_user_id
	stud, err := uc.studentRepo.FindByAuthUserID(ctx, authUserID)
	if err != nil || stud == nil {
		return errors.ErrNotFound
	}

	// Проверяем, существует ли выбор
	exists, err := uc.selectionRepo.Exists(ctx, stud.ID().String(), trackID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.ErrNotFound
	}

	return uc.selectionRepo.Delete(ctx, stud.ID().String(), trackID)
}
