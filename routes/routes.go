package routes

import (
	"github.com/EddieAlvarez01/sist-backend/handlers"
	"github.com/gofiber/fiber/v2"
)

type HandlersRoutes struct {
	*handlers.AccountHolderService
	*handlers.InstitutionService
	*handlers.AccountService
}

func (h *HandlersRoutes) RoutesUp(app fiber.Router) {
	//Routes Account Holder
	routesAccount := app.Group("/account_holder")
	routesAccount.Post("/", h.RegisterNewUser)
	routesAccount.Get("/", h.GetAllAccountHolders)
	routesAccount.Post("/login", h.LoginUser)

	//Routes institution
	routesInstitution := app.Group("/institution")
	routesInstitution.Post("/", h.CreateInstitution)
	routesInstitution.Get("/", h.GetAllInstitutions)

	//Routes account
	routesAccount2 := app.Group("/account")
	routesAccount2.Get("/account_holder/:id", h.GetAllAccountsHolderAccounts)
}
