package handlers

import (
	"github.com/EddieAlvarez01/sist-backend/models"
	"github.com/gofiber/fiber/v2"
)

type OperationService struct {
	OperationModel *models.ManageOperation
}

//Create a new operation
func (o *OperationService) CreateOperation(c *fiber.Ctx) error {
	var operationDTO models.DTOOperation
	if err := c.BodyParser(&operationDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{Code: fiber.StatusBadRequest, Message: "Invalid JSON"})
	}
	account, err := o.OperationModel.Create(&operationDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(account)
}

//Get all operations by account holder id
func (o *OperationService) GetAllOperationsByAccountHolderID(c *fiber.Ctx) error {
	id := c.Params("id")
	operations, err := o.OperationModel.GetAllByAccountHolderID(id,"",1)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(operations)
}

func (o *OperationService) GetAllOperationsByAccountHolderIDAndMonth(c *fiber.Ctx) error {
	id := c.Params("id")
	month := c.Params("month")
	operations, err := o.OperationModel.GetAllByAccountHolderID(id, month, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(operations)
}

//Get total by institution ID
func (o *OperationService) GetTotalsByInstitutionID(c *fiber.Ctx) error {
	id := c.Params("id")
	total, err := o.OperationModel.GetInstitutionTotals(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(total)
}
