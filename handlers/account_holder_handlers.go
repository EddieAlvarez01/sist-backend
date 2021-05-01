package handlers

import (
	"github.com/EddieAlvarez01/sist-backend/models"
	"github.com/gofiber/fiber/v2"
	"log"
)

type AccountHolderService struct {
	AccountHolderModel *models.ManageAccountHolder
}

//Register a new user
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

//Login a user
func (a *AccountHolderService) LoginUser(c *fiber.Ctx) error {
	payload := struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{Code: fiber.StatusBadRequest, Message: "Invalid json"})
	}
	accountHolder, err := a.AccountHolderModel.GetUserByEmailAndPassword(payload.Email, payload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(accountHolder)
}

//Get all account holder
func (a *AccountHolderService) GetAllAccountHolders(c *fiber.Ctx) error {
	accountHolders, err := a.AccountHolderModel.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(accountHolders)
}
