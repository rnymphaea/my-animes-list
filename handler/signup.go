package handler

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"anime/config"
	"anime/database"
)

func SignUp(c *fiber.Ctx) error {
	type Credentials struct {
		Email    string `form:"email" json:"email" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	creds := new(Credentials)
	if err := c.BodyParser(creds); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := database.CreateUser(creds.Email, creds.Password)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusConflict)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = creds.Email
	claims["exp"] = time.Now().Add(time.Second * 30).Unix()

	secret, ok := config.Get("jwtsecret")
	if !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	t, err := token.SignedString([]byte(secret))

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
