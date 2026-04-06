package auth

import (
	"context"
	"track-selection/internal/domain/shared/value_objects"
)

type AuthUserRepository interface {
	Save(ctx context.Context, user *AuthUser) error
	FindByEmail(ctx context.Context, email value_objects.Email) (*AuthUser, error)
	FindByID(ctx context.Context, id string) (*AuthUser, error)
	ExistsByEmail(ctx context.Context, email value_objects.Email) (bool, error)
}
