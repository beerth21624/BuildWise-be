package repositories

import (
	"boonkosang/internal/domain/models"
	"context"

	"github.com/google/uuid"
)

type ContractRepository interface {
	Create(ctx context.Context, projectID uuid.UUID, fileURL string) error
	Delete(ctx context.Context, projectID uuid.UUID) error
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (*models.Contract, error)
	ValidateProjectStatus(ctx context.Context, projectID uuid.UUID) error
}
