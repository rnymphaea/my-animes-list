package router

import (
	"anime/handler"
	"anime/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	//app.Get("/", func(c *fiber.Ctx) error {
	//	return c.SendFile("pages/index.html")
	//})
	app.Get("/", handler.IndexPage)
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.SendFile("pages/login.html")
	})
	app.Get("/myanimes", middleware.Authenticated(), handler.MyAnimeHandler)
	app.Get("/signup", func(c *fiber.Ctx) error {
		return c.SendFile("pages/signup.html")
	})
	app.Post("/login", handler.Login)
}
