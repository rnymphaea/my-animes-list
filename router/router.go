package router

import (
	"github.com/gofiber/fiber/v2"

	"anime/handler"
	"anime/middleware"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", handler.IndexPage)

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.SendFile("pages/login.html")
	})
	app.Post("/login", handler.Login)

	app.Get("/signup", func(c *fiber.Ctx) error {
		return c.SendFile("pages/signup.html")
	})
	app.Post("/signup", handler.SignUp)

	app.Get("/logout", handler.Logout)

	app.Get("/myanimes", middleware.Authenticated(), handler.MyAnimeHandler)
	app.Post("/myanimes", middleware.Authenticated(), handler.AddAnime)

	app.Get("/myanimes/:num", middleware.Authenticated(), handler.ShowAnimeInfo)
}
