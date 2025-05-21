# gofiber-recaptcha3-middleware

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/gofiber-recaptcha3-middleware/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/gofiber-recaptcha3-middleware)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/gofiber-recaptcha3-middleware)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/gofiber-recaptcha3-middleware)](https://goreportcard.com/report/github.com/rgglez/gofiber-recaptcha3-middleware)
[![GitHub release](https://img.shields.io/github/release/rgglez/gofiber-recaptcha3-middleware.svg)](https://github.com/rgglez/gofiber-recaptcha3-middleware/releases/)

**gofiber-recaptcha3-middleware** is a [gofiber](https://gofiber.io/) [middleware](https://docs.gofiber.io/category/-middleware/) implementing [Google ReCaptcha 3](https://developers.google.com/recaptcha/docs/v3?hl=es-419) validation.

## Installation

```bash
go get github.com/rgglez/gofiber-recaptcha3-middleware
```

## Usage

```go
import gofiberip "github.com/rgglez/gofiber-recaptcha3-middleware/recaptcha3"
package main

import (
	"github.com/gofiber/fiber/v2"
	recaptcha3 "github.com/rgglez/gofiber-recaptcha3-middleware/recaptcha3"
)

func main() {
	app := fiber.New()

  var minScore float = 0.5

	app.Post("/endpoint", recaptcha3.RecaptchaMiddleware(middleware.RecaptchaConfig{
		Secret:         "YOUR_SECRET_KEY",
		ExpectedAction: "submit",
		MinScore:       minScore,
	}), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Form sent correctly",
		})
	})

	app.Listen(":3000")
}
```

## Configuration


## Example

An example is included in the [example](example/) directory.

## Dependencies

* go get github.com/gofiber/fiber/v2
* go get github.com/go-playground/validator/v10

## License

Copyright (c) 2025 Rodolfo González González

Licensed under the [Apache 2.0](LICENSE) license. Read the [LICENSE](LICENSE) file.

