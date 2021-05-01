package handlers

import (
	"fmt"
	"github.com/EddieAlvarez01/sist-backend/models"
	"github.com/gofiber/fiber/v2"
)

type AccountService struct {
	AccountModel *models.ManageAccount
}

//Get all accounts by Account holder ID
func (a *AccountService) GetAllAccountsHolderAccounts(c *fiber.Ctx) error {
	id := c.Params("id")
	fmt.Println(id)
	accounts, err := a.AccountModel.GetAllByAccountHolder(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(accounts)
}
