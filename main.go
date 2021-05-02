package main

import (
	"fmt"
	"github.com/EddieAlvarez01/sist-backend/handlers"
	"github.com/EddieAlvarez01/sist-backend/models"
	"github.com/EddieAlvarez01/sist-backend/routes"
	"github.com/EddieAlvarez01/sist-backend/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	//Dotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error starting the .env file: ", err)
	}

	//Database
	sistStorage, err := storage.NewSistStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer sistStorage.Session.Close()

	//Server
	app := fiber.New()

	//Middlewares
	app.Use(cors.New(cors.ConfigDefault))
	app.Use(recover.New())

	//Routes
	group := app.Group("/api/v1")
	handlersRoutes := routes.HandlersRoutes{
		AccountHolderService: &handlers.AccountHolderService{AccountHolderModel: &models.ManageAccountHolder{SistStorage: sistStorage}},
		InstitutionService: &handlers.InstitutionService{InstitutionModel: &models.ManageInstitution{SistStorage: sistStorage}},
		AccountService: &handlers.AccountService{AccountModel: &models.ManageAccount{SistStorage: sistStorage}},
		OperationService: &handlers.OperationService{OperationModel: &models.ManageOperation{SistStorage: sistStorage}},
	}
	handlersRoutes.RoutesUp(group)

	//Run server
	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))))

}
