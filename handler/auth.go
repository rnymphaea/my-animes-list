package handler

import (
	"anime/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"time"
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

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = creds.Email
	claims["exp"] = time.Now().Add(time.Second * 30).Unix()

	t, err := token.SignedString(config.Get("jwtsecret"))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{"data": t})
}
