package track

import (
	"context"
	"encoding/json"
	"track-selection/internal/domain/shared/errors"
	"track-selection/internal/domain/track"
)

type UpdateTrackUseCase struct {
	trackRepo track.Repository
}

func NewUpdateTrackUseCase(repo track.Repository) *UpdateTrackUseCase {
	return &UpdateTrackUseCase{trackRepo: repo}
}

func (uc *UpdateTrackUseCase) Execute(ctx context.Context, id string, updates json.RawMessage) error {
	track, err := uc.trackRepo.FindByID(ctx, id)
	if err != nil {
		return errors.ErrNotFound
	}

	var updatesMap map[string]interface{}
	if err := json.Unmarshal(updates, &updatesMap); err != nil {
		return errors.ErrInvalidRequest
	}

	track.Update(updatesMap)

	return uc.trackRepo.Save(ctx, track)
}
