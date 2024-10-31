package usecase

import (
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type ContractUseCase interface {
	GetContract(ctx context.Context, projectID uuid.UUID) (*responses.ContractResponse, error)
	CreateContract(ctx context.Context, projectID uuid.UUID, req requests.UploadContractRequest) error
	DeleteContract(ctx context.Context, projectID uuid.UUID) error
}

type contractUseCase struct {
	contractRepo repositories.ContractRepository
	projectRepo  repositories.ProjectRepository
}

func NewContractUsecase(
	contractRepo repositories.ContractRepository,
	projectRepo repositories.ProjectRepository,
) ContractUseCase {
	return &contractUseCase{
		contractRepo: contractRepo,
		projectRepo:  projectRepo,
	}
}

func (u *contractUseCase) GetContract(ctx context.Context, projectID uuid.UUID) (*responses.ContractResponse, error) {
	contract, err := u.contractRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}
	if contract == nil {
		return nil, errors.New("contract not found")
	}

	// Convert *models.Contract to *responses.ContractResponse
	contractResponse := &responses.ContractResponse{
		ContractID: contract.ContractID,
		ProjectID:  contract.ProjectID,
		FileURL:    contract.FileURL.String,
		CreatedAt:  contract.CreatedAt,
		UpdatedAt:  contract.UpdatedAt.Time,
	}

	return contractResponse, nil
}

func (u *contractUseCase) CreateContract(ctx context.Context, projectID uuid.UUID, req requests.UploadContractRequest) error {

	// Check if project exists and get its status
	project, err := u.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}
	if project == nil {
		return errors.New("project not found")
	}

	// Check if contract already exists
	existingContract, err := u.contractRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to check existing contract: %w", err)
	}
	if existingContract != nil {
		return errors.New("contract already exists for this project")
	}

	// Create new contract
	err = u.contractRepo.Create(ctx, projectID, req.FileURL)
	if err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}

	return nil
}

func (u *contractUseCase) DeleteContract(ctx context.Context, projectID uuid.UUID) error {
	// Get existing contract
	existingContract, err := u.contractRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get existing contract: %w", err)
	}
	if existingContract == nil {
		return errors.New("contract not found")
	}

	// Delete contract record from database
	err = u.contractRepo.Delete(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete contract record: %w", err)
	}

	return nil
}
