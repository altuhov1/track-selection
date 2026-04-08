package student

import (
	"context"
	"track-selection/internal/domain/shared/value_objects"
)

type Repository interface {
	Save(ctx context.Context, student *Student) error
	FindByID(ctx context.Context, id StudentID) (*Student, error)
	FindByEmail(ctx context.Context, email value_objects.Email) (*Student, error)
	FindByAuthUserID(ctx context.Context, authUserID string) (*Student, error)
	ExistsByEmail(ctx context.Context, email value_objects.Email) (bool, error)
	UpdateRating(ctx context.Context, studentID StudentID, rating int) error
}
