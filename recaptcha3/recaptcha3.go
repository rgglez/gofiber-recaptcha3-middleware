package recaptcha3

import (
	"encoding/json"
	"net/http"
	"net/url"

	fiber "github.com/gofiber/fiber/v2"
)

//-----------------------------------------------------------------------------

// Middleware configuration
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

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
}

//-----------------------------------------------------------------------------

// Set the default configuration.
var ConfigDefault = Config{
	Next:       nil,
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
	}

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
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

		// Verify with Google ReCaptcha API
		resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
			url.Values{
				"secret":   {cfg.Secret},
				"response": {token},
			})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error contacting reCAPTCHA servers",
			})
		}
		defer resp.Body.Close()

		var recaptchaResp RecaptchaResponse
		if err := json.NewDecoder(resp.Body).Decode(&recaptchaResp); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error parsing reCAPTCHA response",
			})
		}

		// Validation
		if !recaptchaResp.Success {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":       "Failed reCAPTCHA validation",
				"error_codes": recaptchaResp.ErrorCodes,
			})
		}
		if recaptchaResp.Score < cfg.MinScore {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Too low reCAPTCHA score",
				"score": recaptchaResp.Score,
			})
		}
		if recaptchaResp.Action != cfg.ExpectedAction {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "Unexpected reCAPTCHA actiom",
				"action": recaptchaResp.Action,
			})
		}

		// OK, next handler
		return c.Next()
	}
}
