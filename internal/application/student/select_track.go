package student

import (
	"context"
	"track-selection/internal/domain/shared/errors"
	"track-selection/internal/domain/student"
	"track-selection/internal/infrastructure/persistence/postgres"
)

type SelectTrackUseCase struct {
	selectionRepo *postgres.TrackSelectionRepository
	studentRepo   *postgres.StudentRepository
	trackRepo     *postgres.TrackRepository
}

func NewSelectTrackUseCase(
	selectionRepo *postgres.TrackSelectionRepository,
	studentRepo *postgres.StudentRepository,
	trackRepo *postgres.TrackRepository,
) *SelectTrackUseCase {
	return &SelectTrackUseCase{
		selectionRepo: selectionRepo,
		studentRepo:   studentRepo,
		trackRepo:     trackRepo,
	}
}

type SelectTrackInput struct {
	TrackID string `json:"track_id"`
}

func (uc *SelectTrackUseCase) Execute(ctx context.Context, authUserID string, input SelectTrackInput) error {
	stud, err := uc.studentRepo.FindByAuthUserID(ctx, authUserID)
	if err != nil || stud == nil {
		return errors.ErrNotFound
	}

	track, err := uc.trackRepo.FindByID(ctx, input.TrackID)
	if err != nil || track == nil {
		return errors.ErrNotFound
	}

	selection := student.NewTrackSelection(stud.ID().String(), input.TrackID)

	return uc.selectionRepo.Save(ctx, selection)
}
