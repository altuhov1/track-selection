package auth

import (
	"context"
	"track-selection/internal/domain/auth"
	"track-selection/internal/domain/shared/errors"
	"track-selection/internal/domain/shared/value_objects"
)

// LoginInput — входные данные
type LoginInput struct {
	Email    string
	Password string
}

// LoginOutput — выходные данные (токен)
type LoginOutput struct {
	Token string
}

// LoginUseCase — Use Case логина
type LoginUseCase struct {
	authRepo   auth.AuthUserRepository
	jwtService auth.JWTService
}

func NewLoginUseCase(
	authRepo auth.AuthUserRepository,
	jwtService auth.JWTService,
) *LoginUseCase {
	return &LoginUseCase{
		authRepo:   authRepo,
		jwtService: jwtService,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// 1. Создаем Email Value Object
	email, err := value_objects.NewEmail(input.Email)
	if err != nil {
		return nil, errors.ErrInvalidEmail
	}

	// 2. Ищем пользователя
	authUser, err := uc.authRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	// 3. Проверяем пароль (бизнес-логика внутри AuthUser)
	if !authUser.CheckPassword(input.Password) {
		return nil, errors.ErrUnauthorized
	}

	// 4. Генерируем JWT токен
	token, err := uc.jwtService.GenerateToken(authUser.ID, authUser.Role)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{Token: token}, nil
}
