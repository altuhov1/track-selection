package admin

import (
	"context"
	"track-selection/internal/domain/shared/value_objects"
)

type Repository interface {
	Save(ctx context.Context, admin *Admin) error
	FindByID(ctx context.Context, id AdminID) (*Admin, error)
	FindByEmail(ctx context.Context, email value_objects.Email) (*Admin, error)
	FindByAuthUserID(ctx context.Context, authUserID string) (*Admin, error)
	ExistsByEmail(ctx context.Context, email value_objects.Email) (bool, error)
}
