package requests

import (
	"github.com/google/uuid"
)

type UploadContractRequest struct {
	FileURL string `json:"file_url" validate:"required,url"`
}

type DeleteContractRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
}
