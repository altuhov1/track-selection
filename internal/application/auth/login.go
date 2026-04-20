package auth

import (
	"context"
	"track-selection/internal/domain/auth"
	"track-selection/internal/domain/shared/errors"
	"track-selection/internal/domain/shared/value_objects"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
}

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
	email, err := value_objects.NewEmail(input.Email)
	if err != nil {
		return nil, errors.ErrInvalidEmail
	}

	authUser, err := uc.authRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	if authUser == nil {
		return nil, errors.ErrUnauthorized
	}
	if !authUser.CheckPassword(input.Password) {
		return nil, errors.ErrUnauthorized
	}

	token, err := uc.jwtService.GenerateToken(
		authUser.ID,
		authUser.Role,
		authUser.FirstName,
		authUser.LastName,
		authUser.Email.String(),
	)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{Token: token}, nil
}
