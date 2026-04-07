package auth

import (
	"context"
	"track-selection/internal/domain/admin"
	"track-selection/internal/domain/auth"
	"track-selection/internal/domain/shared/errors"
	"track-selection/internal/domain/shared/events"
	"track-selection/internal/domain/shared/value_objects"
	"track-selection/internal/domain/student"
)

// RegisterInput — входные данные
type RegisterInput struct {
	Email    string
	Password string
	Role     string
}

// RegisterUseCase — Use Case регистрации
type RegisterUseCase struct {
	authRepo auth.AuthUserRepository
	eventBus events.EventBus
}

func NewRegisterUseCase(
	authRepo auth.AuthUserRepository,
	eventBus events.EventBus,
) *RegisterUseCase {
	return &RegisterUseCase{
		authRepo: authRepo,
		eventBus: eventBus,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) error {
	var role auth.UserRole
	switch input.Role {
	case "student":
		role = auth.RoleStudent
	case "admin":
		role = auth.RoleAdmin
	default:
		return errors.ErrInvalidRole
	}

	email, err := value_objects.NewEmail(input.Email)
	if err != nil {
		return err
	}

	exists, err := uc.authRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrAlreadyExists
	}

	authUser, err := auth.NewAuthUser(input.Email, input.Password, role)
	if err != nil {
		return err
	}

	if err := uc.authRepo.Save(ctx, authUser); err != nil {
		return err
	}
	var event events.DomainEvent
	switch role {
	case "student":
		event = student.NewStudentRegisteredEvent(authUser.ID, authUser.Email.String(), string(role))
	case "admin":
		event = admin.NewAdminRegisteredEvent(authUser.ID, authUser.Email.String(), string(role))
	}
	return uc.eventBus.Publish(ctx, event)
}
