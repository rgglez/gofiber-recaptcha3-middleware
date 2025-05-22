package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recaptcha3 "github.com/rgglez/gofiber-recaptcha3-middleware"
)

func main() {
	var secretKey string = os.Getenv("RECAPTCHA3_SECRET_KEY")

	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))

	app.Use(recaptcha3.New(recaptcha3.Config{
		Secret:         secretKey,
		TokenField:     "recaptcha_token",
		ExpectedAction: "contact",
		MinScore:       0.5,
	}))

	app.Post("/contact", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Form sent correctly",
		})
	})

	app.Listen(":3000")
}
