# gofiber-recaptcha3-middleware

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/gofiber-recaptcha3-middleware/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/gofiber-recaptcha3-middleware)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/gofiber-recaptcha3-middleware)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/gofiber-recaptcha3-middleware)](https://goreportcard.com/report/github.com/rgglez/gofiber-recaptcha3-middleware)
[![GitHub release](https://img.shields.io/github/release/rgglez/gofiber-recaptcha3-middleware.svg)](https://github.com/rgglez/gofiber-recaptcha3-middleware/releases/)

**gofiber-recaptcha3-middleware** is a [gofiber](https://gofiber.io/) [middleware](https://docs.gofiber.io/category/-middleware/) implementing [Google reCAPTCHA v3](https://developers.google.com/recaptcha/docs/v3?hl=es-419) validation.

## Installation

```bash
go get github.com/rgglez/gofiber-recaptcha3-middleware
```

## Usage

Basic usage:

```go
import recaptcha3 "github.com/rgglez/gofiber-recaptcha3-middleware"

// ...

app.Use(recaptcha3.New(recaptcha3.Config{
	Secret:         "YOUR_SECRET_KEY",
	TokenField:     "recaptcha_token",
	ExpectedAction: "contact",
	MinScore:       0.5,
}))
```

See the [examples](examples/) directory for more usage examples.

## Configuration

* **Next** defines a function to skip this middleware when returned true. Optional. Default: nil
* **NoNext** defines a function to "force" this middleware when returned true.  Optional. Default: nil
* **Secret** reCAPTCHA secret key. Required string.
* **ExpectedAction** is the expected action. Example: "login" o "submit". Required string.
* **MinScore** the lowest acceptable score. Example: 0.5. Required float.
* **TokenField** the token field in the JSON body. Default: "recaptcha_token"	

## Example

An example is included in the [example](example/) directory. Both backend and frontend are included. 

* You should copy the ```frontend/index.html``` file to a web server.
* You can run the backend as a normal Go app:

```bash
go run .
```

## Dependencies

* go get github.com/gofiber/fiber/v2
* go get github.com/go-playground/validator/v10
* go get resty.dev/v3

## Security

* Remember that reCAPTCHA v3 does not stop bots directly, but just gives a risk score. It is responsability of the application to accept the request or deny it.
* Another captcha (such as reCAPTCHA v2) could be shown with a classic challenge as a "plan B", if the score is lower than expected. 

## License

Copyright (c) 2025 Rodolfo González González

Licensed under the [Apache 2.0](LICENSE) license. Read the [LICENSE](LICENSE) file.

