package routes

import (
	"github.com/EddieAlvarez01/sist-backend/handlers"
	"github.com/gofiber/fiber/v2"
)

type HandlersRoutes struct {
	*handlers.AccountHolderService
}

func (h *HandlersRoutes) RoutesUp(app fiber.Router) {
	//Routes Account Holder
	routesAccount := app.Group("/account_holder")
	routesAccount.Post("/", h.RegisterNewUser)
	routesAccount.Post("/login", h.LoginUser)
}
