package rest

import (
	"boonkosang/internal/requests"
	"boonkosang/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ContractHandler struct {
	contractUseCase usecase.ContractUseCase
}

func NewContractHandler(contractUseCase usecase.ContractUseCase) *ContractHandler {
	return &ContractHandler{
		contractUseCase: contractUseCase,
	}
}

func (h *ContractHandler) ContractRoutes(app *fiber.App) {
	contract := app.Group("/contracts/:projectId")

	contract.Get("/", h.GetContract)
	contract.Post("/", h.CreateContract)
	contract.Delete("/", h.DeleteContract)
}

func (h *ContractHandler) GetContract(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	contract, err := h.contractUseCase.GetContract(c.Context(), projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Contract retrieved successfully",
		"data":    contract,
	})

}

func (h *ContractHandler) CreateContract(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	var req requests.UploadContractRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.contractUseCase.CreateContract(c.Context(), projectID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Contract created successfully",
	})
}

func (h *ContractHandler) DeleteContract(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	if err := h.contractUseCase.DeleteContract(c.Context(), projectID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Contract deleted successfully",
	})
}
