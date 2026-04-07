package auth

import (
	"context"
	"track-selection/internal/domain/admin"
	"track-selection/internal/domain/shared/events"
)

type CreateAdminRegHandler struct {
	adminRepo admin.Repository
}

func NewCreateAdminRegHandler(adminRepo admin.Repository) *CreateAdminRegHandler {
	return &CreateAdminRegHandler{adminRepo: adminRepo}
}

func (h *CreateAdminRegHandler) Handle(ctx context.Context, event events.DomainEvent) error {
	e, ok := event.(admin.AdminRegisteredEvent)
	if !ok {
		return nil
	}

	admin, err := admin.NewAdmin(e.UserID, e.Email)
	if err != nil {
		return err
	}
	return h.adminRepo.Save(ctx, admin)
}
