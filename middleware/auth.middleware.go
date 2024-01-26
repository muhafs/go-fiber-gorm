package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/muhafs/go-fiber-gorm/utils"
)

func Auth(c *fiber.Ctx) error {
	// get Header's value
	auth := c.Get("Authorization")

	// split type and token from auth string
	splittedAuth := strings.Split(auth, " ") // ["Bearer", "secrettokenstring123"]

	// validate auth string format
	if len(splittedAuth) != 2 {
		return utils.ErrorJSON(c, fiber.StatusUnauthorized, "missing authorization header")
	}

	// extract type and token
	hType := splittedAuth[0]  // "Bearer"
	hToken := splittedAuth[1] // "secrettokenstring123"

	// validate token string
	if _, err := utils.VerifyToken(hToken); err != nil || hType != "Bearer" {
		return utils.ErrorJSON(c, fiber.StatusUnauthorized, "invalid token")
	}

	// pass to next
	return c.Next()
}
