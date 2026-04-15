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
	// Получаем предпочтения
	prefs, err := uc.prefsRepo.FindByUserID(ctx, userID)
	if err != nil {
		// Нет предпочтений - профиль не заполнен
		return uc.updateCompletionStatus(ctx, userID, false)
	}

	// Проверяем полноту
	isComplete := uc.profileChecker.IsProfileComplete(prefs)

	return uc.updateCompletionStatus(ctx, userID, isComplete)
}

func (uc *CheckProfileCompletionUseCase) updateCompletionStatus(ctx context.Context, userID string, isComplete bool) error {
	// Получаем текущий статус
	current, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		// Создаем новый
		current = student.NewProfileCompletion(userID)
	}

	// Если статус не изменился - ничего не делаем
	if current.IsComplete == isComplete {
		return nil
	}

	// Обновляем статус
	if isComplete {
		current.SetComplete(time.Now())
	} else {
		current.SetIncomplete()
	}

	// Сохраняем
	if err := uc.profileRepo.Save(ctx, current); err != nil {
		return err
	}

	// Публикуем событие
	event := student.NewProfileCompletedEvent(userID, isComplete)
	return uc.eventBus.Publish(ctx, event)
}
