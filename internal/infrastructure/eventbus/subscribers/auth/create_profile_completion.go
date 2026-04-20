package auth

import (
	"context"
	"time"
	"track-selection/internal/domain/shared/events"
	"track-selection/internal/domain/student"
	"track-selection/internal/infrastructure/persistence/postgres"

	"github.com/google/uuid"
)

type CreateProfileCompletionHandler struct {
	prefsRepo   *postgres.PreferencesRepository
	profileRepo *postgres.ProfileCompletionRepository
}

func NewCreateProfileCompletionHandler(
	prefsRepo *postgres.PreferencesRepository,
	profileRepo *postgres.ProfileCompletionRepository,
) *CreateProfileCompletionHandler {
	return &CreateProfileCompletionHandler{
		prefsRepo:   prefsRepo,
		profileRepo: profileRepo,
	}
}

func (h *CreateProfileCompletionHandler) Handle(ctx context.Context, event events.DomainEvent) error {
	e, ok := event.(student.StudentRegisteredEvent)
	if !ok {
		return nil
	}
	prefs := &student.Preferences{
		ID:                uuid.New().String(),
		UserID:            e.UserID,
		ProfessionalGoals: []int{},
		Grades:            student.GenerateRandomGrades(),
		Skills:            student.Skills{},
		LearningStyle:     1,
		Certificates:      0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := h.prefsRepo.Save(ctx, prefs); err != nil {
		return err
	}

	// 2. Создаем запись о заполнении профиля
	return h.profileRepo.CreateDefault(ctx, e.UserID)
}
