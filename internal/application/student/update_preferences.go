package student

import (
	"context"
	"encoding/json"
	"time"
	"track-selection/internal/domain/shared/errors"
	"track-selection/internal/domain/shared/events"
	"track-selection/internal/domain/student"
	"track-selection/internal/infrastructure/persistence/postgres"

	"github.com/google/uuid"
)

type UpdatePreferencesUseCase struct {
	prefsRepo      *postgres.PreferencesRepository
	profileRepo    *postgres.ProfileCompletionRepository
	profileChecker *student.ProfileChecker
	eventBus       events.EventBus
}

func NewUpdatePreferencesUseCase(
	prefsRepo *postgres.PreferencesRepository,
	profileRepo *postgres.ProfileCompletionRepository,
	profileChecker *student.ProfileChecker,
	eventBus events.EventBus,
) *UpdatePreferencesUseCase {
	return &UpdatePreferencesUseCase{
		prefsRepo:      prefsRepo,
		profileRepo:    profileRepo,
		profileChecker: profileChecker,
		eventBus:       eventBus,
	}
}

func (uc *UpdatePreferencesUseCase) Execute(ctx context.Context, userID string, updates json.RawMessage) error {
	// Получаем существующие или создаём новые
	prefs, err := uc.prefsRepo.FindByUserID(ctx, userID)
	if err != nil {
		prefs = &student.Preferences{
			ID:        uuid.New().String(),
			UserID:    userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	// Парсим обновления
	var updatesMap map[string]interface{}
	if err := json.Unmarshal(updates, &updatesMap); err != nil {
		return errors.ErrInvalidRequest
	}

	// Применяем обновления
	prefs.Merge(updatesMap)

	// Валидация в DOMAIN слое!
	if err := prefs.ValidatePartial(updatesMap); err != nil {
		return err
	}

	// Сохраняем предпочтения
	if err := uc.prefsRepo.Save(ctx, prefs); err != nil {
		return err
	}

	// Проверяем полноту профиля
	isComplete := uc.profileChecker.IsProfileComplete(prefs)

	// Обновляем статус
	if err := uc.updateProfileCompletion(ctx, userID, isComplete); err != nil {
		return err
	}

	// Публикуем событие
	event := student.NewProfileCompletedEvent(userID, isComplete)
	uc.eventBus.Publish(ctx, event)

	return nil
}

func (uc *UpdatePreferencesUseCase) updateProfileCompletion(ctx context.Context, userID string, isComplete bool) error {
	current, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		current = &student.ProfileCompletion{
			ID:         uuid.New().String(),
			UserID:     userID,
			IsComplete: false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}

	if current.IsComplete == isComplete {
		return nil
	}

	current.IsComplete = isComplete
	current.UpdatedAt = time.Now()

	if isComplete {
		now := time.Now()
		current.CompletedAt = &now
	} else {
		current.CompletedAt = nil
	}

	return uc.profileRepo.Save(ctx, current)
}
