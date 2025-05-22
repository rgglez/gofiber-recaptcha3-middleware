/*

   Copyright 2025 Rodolfo González González

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recaptcha3 "github.com/rgglez/gofiber-recaptcha3-middleware"
)

func RecaptchaActionMiddleware(action string) fiber.Handler {
	return recaptcha3.New(recaptcha3.Config{
		Secret:         os.Getenv("RECAPTCHA3_SECRET_KEY"),
		ExpectedAction: action,
		MinScore:       0.5,
	})
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))

	app.Post("/contact", RecaptchaActionMiddleware("contact"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Form sent correctly",
		})
	})

	app.Post("/help", RecaptchaActionMiddleware("help"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Form sent correctly",
		})
	})

	app.Post("/wrong", RecaptchaActionMiddleware("xxx"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Form sent correctly",
		})
	})

	app.Listen(":3000")
}
