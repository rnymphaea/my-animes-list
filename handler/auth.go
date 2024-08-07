package handler

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"anime/config"
	"anime/database"
)

func Login(c *fiber.Ctx) error {
	type Credentials struct {
		Email    string `form:"email" json:"email" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	creds := new(Credentials)
	if err := c.BodyParser(creds); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if ok, err := database.CheckPassword(creds.Email, creds.Password); !ok || err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = creds.Email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	secret, ok := config.Get("jwtsecret")
	if !ok {
		log.Println("JWT secret not set")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	t, err := token.SignedString([]byte(secret))
	log.Println("TOKEN: ", t)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		HTTPOnly: true,
	})

	return c.Redirect("/", fiber.StatusSeeOther)
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:  "jwt",
		Value: "",
	})
	log.Println(c.Request())
	return c.Redirect("/", fiber.StatusSeeOther)
}
