package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"

	"anime/database"
	"anime/router"
)

func main() {

	engine := html.New("./assets/pages", ".html")

	conn, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	database.CreateTables(conn)

	//database.AddAllAnimes() # uncomment if table is not filled

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
