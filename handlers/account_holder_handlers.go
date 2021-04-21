package handlers

import (
	"github.com/EddieAlvarez01/sist-backend/models"
	"github.com/gofiber/fiber/v2"
	"log"
)

type AccountHolderService struct {
	AccountHolderModel *models.ManageAccountHolder
}

func (a *AccountHolderService) RegisterNewUser(c *fiber.Ctx) error {
	accountHolder := models.AccountHolder{Role: "USER"}
	if err := c.BodyParser(&accountHolder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{Code: fiber.StatusBadRequest, Message: "Invalid json"})
	}
	_, err := a.AccountHolderModel.Create(&accountHolder)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{Code: fiber.StatusInternalServerError, Message: "Error on create user"})
	}
	return c.Status(fiber.StatusOK).JSON(accountHolder)
}
