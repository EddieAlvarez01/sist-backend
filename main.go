package main

import (
	"fmt"
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

	//Server
	app := fiber.New()

	//Middlewares
	app.Use(cors.New(cors.ConfigDefault))
	app.Use(recover.New())

	//Routes
	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.SendString("Hello world")
	})

	//Run server
	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))))

}
