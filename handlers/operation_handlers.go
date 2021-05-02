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
