package auth

import (
	"context"
	"track-selection/internal/domain/shared/events"
	"track-selection/internal/domain/student"
)

type CreateStudentRegHandler struct {
	studentRepo student.Repository
}

func NewCreateStudentRegHandler(studentRepo student.Repository) *CreateStudentRegHandler {
	return &CreateStudentRegHandler{studentRepo: studentRepo}
}

func (h *CreateStudentRegHandler) Handle(ctx context.Context, event events.DomainEvent) error {
	e, ok := event.(student.StudentRegisteredEvent)
	if !ok {
		return nil
	}

	student, err := student.NewStudent(
		e.UserID,
		e.Email,
		e.FirstName,
		e.LastName,
	)
	if err != nil {
		return err
	}

	return h.studentRepo.Save(ctx, student)
}
