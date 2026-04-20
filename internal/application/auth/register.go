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

type RegisterInput struct {
	Email       string
	Password    string
	Role        string
	FirstName   string
	LastName    string
	AdminSecret string
}

type RegisterUseCase struct {
	authRepo       auth.AuthUserRepository
	eventBus       events.EventBus
	adminSecretKey string
}

func NewRegisterUseCase(
	authRepo auth.AuthUserRepository,
	eventBus events.EventBus,
	adminSecretKey string,
) *RegisterUseCase {
	return &RegisterUseCase{
		authRepo:       authRepo,
		eventBus:       eventBus,
		adminSecretKey: adminSecretKey,
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
	if role == auth.RoleAdmin {
		if input.AdminSecret == "" {
			return errors.ErrAdminSecretRequired
		}
		if input.AdminSecret != uc.adminSecretKey {
			return errors.ErrInvalidAdminSecret
		}
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

	authUser, err := auth.NewAuthUser(
		input.Email,
		input.Password,
		input.FirstName,
		input.LastName,
		role,
	)
	if err != nil {
		return err
	}

	if err := uc.authRepo.Save(ctx, authUser); err != nil {
		return err
	}

	var event events.DomainEvent
	switch role {
	case auth.RoleStudent:
		event = student.NewStudentRegisteredEvent(
			authUser.ID,
			authUser.Email.String(),
			input.FirstName,
			input.LastName,
		)
	case auth.RoleAdmin:
		event = admin.NewAdminRegisteredEvent(
			authUser.ID,
			authUser.Email.String(),
			input.FirstName,
			input.LastName,
		)
	}

	return uc.eventBus.Publish(ctx, event)
}
