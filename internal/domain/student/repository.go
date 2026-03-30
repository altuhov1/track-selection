package student

import "context"

type Repository interface {
	Save(ctx context.Context, student *Student) error
	FindByID(ctx context.Context, id StudentID) (*Student, error)
	FindByEmail(ctx context.Context, email Email) (*Student, error)
}
