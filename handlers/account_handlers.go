package handlers

import (
	"github.com/EddieAlvarez01/sist-backend/models"
	"github.com/gofiber/fiber/v2"
)

type AccountService struct {
	AccountModel *models.ManageAccount
}

//Get all accounts by Account holder ID
func (a *AccountService) GetAllAccountsHolderAccounts(c *fiber.Ctx) error {
	id := c.Params("id")
	accounts, err := a.AccountModel.GetAllByAccountHolder(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(accounts)
}

//Create a new account for an account holder
func (a *AccountService) CreateNewAccount(c *fiber.Ctx) error {
	var account models.Account
	if err := c.BodyParser(&account); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{Code: fiber.StatusBadRequest, Message: "Invalid JSON"})
	}
	_, err := a.AccountModel.Create(&account)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(account)
}
