package middleware

import (
	"anime/config"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

func Authenticated() func(c *fiber.Ctx) error {
	jwtConfig := jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: config.Get("jwtsecret")},
		TokenLookup:    "cookie:jwt",
		ContextKey:     "user",
		SuccessHandler: ValidateToken,
		ErrorHandler:   jwtError,
	}
	return jwtware.New(jwtConfig)
}

func ValidateToken(c *fiber.Ctx) error {
	cookieJWT := c.Cookies("jwt")
	if cookieJWT == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookieJWT, claims, func(token *jwt.Token) (interface{}, error) {
		return config.Get("jwtsecret"), nil
	})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JWT"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	exptime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(exptime) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "time is over"})
	}

	c.Locals("user", token)
	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
}
