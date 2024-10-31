package responses

import (
	"time"

	"github.com/google/uuid"
)

type ContractResponse struct {
	ContractID uuid.UUID `json:"contract_id"`
	ProjectID  uuid.UUID `json:"project_id"`
	FileURL    string    `json:"file_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
