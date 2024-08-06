package main

import (
	"anime/database"
	"anime/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"log"
)

type Anime struct {
	ID          uint64 `json:"mal_id"`
	Title       string `json:"title"`
	Episodes    uint16 `json:"episodes"`
	Description string `json:"synopsis"`
}

func main() {

	engine := html.New("./pages", ".html")

	conn, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	database.CreateTables(conn)
	//database.AddAllAnimes()
	defer conn.Close()

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New())
	router.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
