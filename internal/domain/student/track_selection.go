package student

import (
	"time"

	"github.com/google/uuid"
)

type TrackSelection struct {
	ID         string    `json:"id"`
	StudentID  string    `json:"student_id"`
	TrackID    string    `json:"track_id"`
	SelectedAt time.Time `json:"selected_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewTrackSelection(studentID, trackID string) *TrackSelection {
	now := time.Now()
	return &TrackSelection{
		ID:         uuid.New().String(),
		StudentID:  studentID,
		TrackID:    trackID,
		SelectedAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
