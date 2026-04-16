package track

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, track *Track) error
	FindByID(ctx context.Context, id string) (*Track, error)
	FindAll(ctx context.Context) ([]*Track, error)
	Delete(ctx context.Context, id string) error
}
