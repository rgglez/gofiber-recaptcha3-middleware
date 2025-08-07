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

package recaptcha3

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"
	resty "resty.dev/v3"
)

//-----------------------------------------------------------------------------

// Middleware configuration
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// NoNext defines a function to "force" this middleware when returned true.
	// Optional. Default: nil
	NoNext func(c *fiber.Ctx) bool

	// reCAPTCHA secret key
	// Required string.
	Secret string

	// Expected action. Example: "login" o "submit"
	// Required string.
	ExpectedAction string

	// Lowest acceptable score. Example: 0.5
	// Required float.
	MinScore float64

	// Token field in the JSON body ()
	// Optional. Default: "recaptcha_token"
	TokenField string

	// Resty client 
	// Optional. Pass it if you want to change the retry or timeout defaults.
	Client *resty.Client
}

//-----------------------------------------------------------------------------

// Set the default configuration.
var ConfigDefault = Config{
	Next:       nil,
	NoNext:     nil,
	TokenField: "recaptcha_token",
}

//-----------------------------------------------------------------------------

// reCAPTCHA response
type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ErrorCodes  []string `json:"error-codes"`
}

//-----------------------------------------------------------------------------

// Global reCAPTCHA middleware with Skip support
func New(config ...Config) fiber.Handler {	
	cfg := ConfigDefault

	if len(config) > 0 {
		cfg = config[0]

		if cfg.TokenField == "" {
			cfg.TokenField = ConfigDefault.TokenField
		}
		if cfg.Client == nil {
			cfg.Client = resty.New().
				SetRetryCount(3).
				SetRetryWaitTime(2 * time.Second).
				SetRetryMaxWaitTime(8 * time.Second).
				AddRetryConditions(
					func(r *resty.Response, err error) bool {
						return err != nil || r.StatusCode() >= 500
					})
		}
	}

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		if cfg.NoNext != nil && !cfg.NoNext(c) {
			return c.Next()
		}

		var body map[string]string
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON",
			})
		}

		token, ok := body[cfg.TokenField]
		if !ok || token == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "reCAPTCHA token missing",
			})
		}

		resp := RecaptchaResponse{}
		_, err := cfg.Client.R().
			SetFormData(map[string]string{
				"secret":   cfg.Secret,
				"response": token,
				"remoteip": c.IP(),
			}).
			SetResult(&resp).
			Post("https://www.google.com/recaptcha/api/siteverify")

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify reCAPTCHA",
			})
		}

		if !resp.Success || resp.Score < cfg.MinScore {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "reCAPTCHA verification failed",
				"details": resp,
			})
		}

		if cfg.ExpectedAction != "" && resp.Action != cfg.ExpectedAction {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "Unexpected reCAPTCHA action",
				"action": resp.Action,
			})
		}

		// OK, next handler
		return c.Next()
	}
}
