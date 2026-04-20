package student

import (
	"context"
	"track-selection/internal/domain/track"
	"track-selection/internal/infrastructure/persistence/postgres"
)

type GetSelectedTracksUseCase struct {
	selectionRepo *postgres.TrackSelectionRepository
	trackRepo     *postgres.TrackRepository
	studentRepo   *postgres.StudentRepository
}

func NewGetSelectedTracksUseCase(
	selectionRepo *postgres.TrackSelectionRepository,
	trackRepo *postgres.TrackRepository,
	studentRepo *postgres.StudentRepository,
) *GetSelectedTracksUseCase {
	return &GetSelectedTracksUseCase{
		selectionRepo: selectionRepo,
		trackRepo:     trackRepo,
		studentRepo:   studentRepo,
	}
}

type GetSelectedTracksOutput struct {
	Tracks []*track.Track `json:"tracks"`
}

func (uc *GetSelectedTracksUseCase) Execute(ctx context.Context, studentID string) (*GetSelectedTracksOutput, error) {
	stud, err := uc.studentRepo.FindByAuthUserID(ctx, studentID)
	if err != nil || stud == nil {
		return &GetSelectedTracksOutput{Tracks: []*track.Track{}}, nil
	}

	// Получаем выборы студента
	selections, err := uc.selectionRepo.FindByStudentID(ctx, stud.ID().String())
	if err != nil {
		return &GetSelectedTracksOutput{Tracks: []*track.Track{}}, nil
	}

	// Получаем треки по ID
	var tracks []*track.Track
	for _, s := range selections {
		t, err := uc.trackRepo.FindByID(ctx, s.TrackID)
		if err != nil {
			continue
		}
		if t != nil {
			tracks = append(tracks, t)
		}
	}

	return &GetSelectedTracksOutput{Tracks: tracks}, nil
}
