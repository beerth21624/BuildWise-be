package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type contractRepository struct {
	db *sqlx.DB
}

func NewContractRepository(db *sqlx.DB) repositories.ContractRepository {
	return &contractRepository{db: db}
}

func (r *contractRepository) ValidateProjectStatus(ctx context.Context, projectID uuid.UUID) error {
	query := `
		SELECT 
			p.status as project_status,
			b.status as boq_status,
			q.status as quotation_status
		FROM project p
		LEFT JOIN boq b ON b.project_id = p.project_id
		LEFT JOIN quotation q ON q.project_id = p.project_id
		WHERE p.project_id = $1`

	type StatusCheck struct {
		ProjectStatus   string         `db:"project_status"`
		BOQStatus       sql.NullString `db:"boq_status"`
		QuotationStatus sql.NullString `db:"quotation_status"`
	}

	var status StatusCheck
	err := r.db.GetContext(ctx, &status, query, projectID)
	if err != nil {
		return fmt.Errorf("failed to get project status: %w", err)
	}

	if status.ProjectStatus == "completed" {
		return errors.New("project is already completed")
	}

	if !status.BOQStatus.Valid || status.BOQStatus.String != "approved" {
		return errors.New("BOQ must be approved")
	}

	if !status.QuotationStatus.Valid || status.QuotationStatus.String != "approved" {
		return errors.New("quotation must be approved")
	}

	return nil
}

func (r *contractRepository) Create(ctx context.Context, projectID uuid.UUID, fileURL string) error {
	if err := r.ValidateProjectStatus(ctx, projectID); err != nil {
		return err
	}

	query := `
		INSERT INTO contract (
			contract_id,
			project_id,
			file_url,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)`

	_, err := r.db.ExecContext(ctx, query, uuid.New(), projectID, fileURL)
	if err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}

	return nil
}

func (r *contractRepository) Delete(ctx context.Context, projectID uuid.UUID) error {
	if err := r.ValidateProjectStatus(ctx, projectID); err != nil {
		return err
	}

	query := `DELETE FROM contract WHERE project_id = $1`
	result, err := r.db.ExecContext(ctx, query, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("contract not found")
	}

	return nil
}

func (r *contractRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*models.Contract, error) {
	var contract models.Contract
	query := `SELECT * FROM contract WHERE project_id = $1`
	err := r.db.GetContext(ctx, &contract, query, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}
	return &contract, nil
}
