package handlers

import (
	"github.com/EddieAlvarez01/sist-backend/models"
	"github.com/gofiber/fiber/v2"
)

type InstitutionService struct {
	InstitutionModel *models.ManageInstitution
}

//Register a new institution
func (i *InstitutionService) CreateInstitution(c *fiber.Ctx) error {
	var institution models.Institution
	if err := c.BodyParser(&institution); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{Code: fiber.StatusBadRequest, Message: "Invalid JSON"})
	}
	_, err := i.InstitutionModel.Create(&institution)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(institution)
}

//Get all institutions
func (i *InstitutionService) GetAllInstitutions(c *fiber.Ctx) error {
	institutions, err := i.InstitutionModel.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(institutions)
}
