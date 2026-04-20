package student

import (
	"context"
	"time"
	"track-selection/internal/domain/shared/events"
	"track-selection/internal/domain/student"
	"track-selection/internal/infrastructure/persistence/postgres"
)

type CheckProfileCompletionUseCase struct {
	prefsRepo      *postgres.PreferencesRepository
	profileRepo    *postgres.ProfileCompletionRepository
	eventBus       events.EventBus
	profileChecker *student.ProfileChecker
}

func NewCheckProfileCompletionUseCase(
	prefsRepo *postgres.PreferencesRepository,
	profileRepo *postgres.ProfileCompletionRepository,
	eventBus events.EventBus,
) *CheckProfileCompletionUseCase {
	return &CheckProfileCompletionUseCase{
		prefsRepo:      prefsRepo,
		profileRepo:    profileRepo,
		eventBus:       eventBus,
		profileChecker: student.NewProfileChecker(),
	}
}

func (uc *CheckProfileCompletionUseCase) Execute(ctx context.Context, userID string) error {
	prefs, err := uc.prefsRepo.FindByUserID(ctx, userID)
	if err != nil {
		return uc.updateCompletionStatus(ctx, userID, false)
	}

	isComplete := uc.profileChecker.IsProfileComplete(prefs)

	return uc.updateCompletionStatus(ctx, userID, isComplete)
}

func (uc *CheckProfileCompletionUseCase) updateCompletionStatus(ctx context.Context, userID string, isComplete bool) error {
	current, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		current = student.NewProfileCompletion(userID)
	}

	if current.IsComplete == isComplete {
		return nil
	}

	if isComplete {
		current.SetComplete(time.Now())
	} else {
		current.SetIncomplete()
	}

	if err := uc.profileRepo.Save(ctx, current); err != nil {
		return err
	}

	event := student.NewProfileCompletedEvent(userID, isComplete)
	return uc.eventBus.Publish(ctx, event)
}
