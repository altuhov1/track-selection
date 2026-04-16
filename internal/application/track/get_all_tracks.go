package track

import (
	"context"
	"track-selection/internal/domain/track"
)

type GetAllTracksUseCase struct {
	trackRepo track.Repository
}

func NewGetAllTracksUseCase(repo track.Repository) *GetAllTracksUseCase {
	return &GetAllTracksUseCase{trackRepo: repo}
}

func (uc *GetAllTracksUseCase) Execute(ctx context.Context) ([]*track.Track, error) {
	return uc.trackRepo.FindAll(ctx)
}
