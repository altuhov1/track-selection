package auth

import (
	"context"
	"track-selection/internal/domain/shared/events"
	"track-selection/internal/domain/student"
)

type CreateStudentHandler struct {
	studentRepo student.Repository
}

func NewCreateStudentHandler(studentRepo student.Repository) *CreateStudentHandler {
	return &CreateStudentHandler{studentRepo: studentRepo}
}

func (h *CreateStudentHandler) Handle(ctx context.Context, event events.DomainEvent) error {
	e, ok := event.(events.UserRegisteredEvent)
	if !ok {
		return nil
	}

	if e.Role != "student" {
		return nil
	}

	student, err := student.NewStudent(e.UserID, e.Email)
	if err != nil {
		return err
	}
	return h.studentRepo.Save(ctx, student)
}
