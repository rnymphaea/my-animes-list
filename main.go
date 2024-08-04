package main

import (
	"anime/router"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"log"
	"net/http"
)

type Anime struct {
	ID       uint64 `json:"mal_id"`
	Title    string `json:"title"`
	Episodes uint16 `json:"episodes"`
}

type Data struct {
	Data Anime `json:"data"`
}

func main() {
	var data Data
	//var anime Anime?
	engine := html.New("./pages", ".html")
	for i := 1000; i < 1010; i++ {
		URL := fmt.Sprintf("https://api.jikan.moe/v4/anime/%d", i)
		//println(URL)
		response, err := http.Get(URL)
		if response.StatusCode != 200 {
			continue
		}
		if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(data)
		}
		response.Body.Close()
	}
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New())
	router.SetupRoutes(app)
	//app.Use("/myanimes", jwtware.New(jwtware.Config{
	//	SigningKey:   []byte(jwtSecret),
	//	TokenLookup:  "cookie:jwt",
	//	ErrorHandler: jwtError}))

	log.Fatal(app.Listen(":3000"))
}

//
//func jwtError(c *fiber.Ctx, err error) error {
//	if err.Error() == "Missing or malformed JWT" {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing or malformed JWT"})
//	}
//	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
//}
