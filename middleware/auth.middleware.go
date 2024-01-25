package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	// get Header's value
	auth := c.Get("Authorization")

	// validate token format
	tokens := strings.Split(auth, " ")
	if len(tokens) != 2 || tokens[0] != "Token" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid token",
		})
	}

	// validate token
	if tokens[1] != "secret" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	// pass to next
	return c.Next()
}
