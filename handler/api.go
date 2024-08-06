package handler

import (
	"anime/database"
	"anime/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
)

func MyAnimeHandler(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	email := user.(*jwt.Token).Claims.(jwt.MapClaims)["email"].(string)
	log.Println(email)
	animes, err := database.GetAnimeList(email)
	if err != nil {
		log.Println(err)
	}
	return c.Render("myanimes", fiber.Map{"Animes": animes})
}

func AddAnime(c *fiber.Ctx) error {
	type RequestData struct {
		Title string `json:"data"`
	}
	user := c.Locals("user")
	if user == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	email := user.(*jwt.Token).Claims.(jwt.MapClaims)["email"].(string)

	var rq RequestData
	if err := c.BodyParser(&rq); err != nil {
		log.Println(err)
	}

	if err := database.AddAnime(email, rq.Title); err != nil {
		log.Println(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
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
