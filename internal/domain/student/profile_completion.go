package student

import (
	"time"

	"github.com/google/uuid"
)

type ProfileCompletion struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	IsComplete  bool       `json:"is_complete"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewProfileCompletion(userID string) *ProfileCompletion {
	now := time.Now()
	return &ProfileCompletion{
		ID:         uuid.New().String(),
		UserID:     userID,
		IsComplete: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (p *ProfileCompletion) SetComplete(completedAt time.Time) {
	p.IsComplete = true
	p.CompletedAt = &completedAt
	p.UpdatedAt = completedAt
}

func (p *ProfileCompletion) SetIncomplete() {
	p.IsComplete = false
	p.CompletedAt = nil
	p.UpdatedAt = time.Now()
}
