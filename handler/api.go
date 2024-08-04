package handler

import (
	"anime/middleware"
	"github.com/gofiber/fiber/v2"
)

func MyAnimeHandler(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.SendFile("pages/myanimes.html")
}

func IndexPage(c *fiber.Ctx) error {
	if err := middleware.ValidateToken(c); err != nil {
		return c.Render("index", fiber.Map{
			"IsAuthenticated": true,
		})
	} else {
		return c.Render("index", fiber.Map{
			"IsAuthenticated": false,
		})
	}
}
