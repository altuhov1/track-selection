package auth

import (
	"context"
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
	authRepo    auth.AuthUserRepository
	studentRepo student.Repository
	// adminRepo   admin.Repository
	eventBus events.EventBus
}

func NewRegisterUseCase(
	authRepo auth.AuthUserRepository,
	studentRepo student.Repository,
	// adminRepo admin.Repository,
	eventBus events.EventBus,
) *RegisterUseCase {
	return &RegisterUseCase{
		authRepo:    authRepo,
		studentRepo: studentRepo,
		// adminRepo:   adminRepo,
		eventBus: eventBus,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) error {
	// 1. Определяем роль
	var role auth.UserRole
	switch input.Role {
	case "student":
		role = auth.RoleStudent
	case "admin":
		role = auth.RoleAdmin
	default:
		return errors.ErrInvalidRole
	}

	// 2. Проверяем, существует ли email
	email, _ := value_objects.NewEmail(input.Email)
	exists, err := uc.authRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrAlreadyExists
	}

	// 3. Создаем AuthUser (бизнес-логика внутри NewAuthUser)
	authUser, err := auth.NewAuthUser(input.Email, input.Password, role)
	if err != nil {
		return err
	}

	// 4. Сохраняем AuthUser
	if err := uc.authRepo.Save(ctx, authUser); err != nil {
		return err
	}

	// 5. В зависимости от роли создаем Student или Admin
	switch role {
	case auth.RoleStudent:
		student, err := student.NewStudent(authUser.ID, authUser.Email.String())
		if err != nil {
			return err
		}
		if err := uc.studentRepo.Save(ctx, student); err != nil {
			return err
		}
	case auth.RoleAdmin:
		// admin := admin.NewAdmin(authUser.ID, authUser.Email.String())
		// if err := uc.adminRepo.Save(ctx, admin); err != nil {
		// 	return err
		// }
	}

	// 6. Публикуем событие
	event := events.NewUserRegisteredEvent(authUser.ID, authUser.Email.String(), string(role))
	uc.eventBus.Publish(ctx, event)

	return nil
}
