package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/muhafs/go-fiber-gorm/utils"
)

func Auth(c *fiber.Ctx) error {
	// get Header's value
	auth := c.Get("Authorization")

	// validate token format
	tokens := strings.Split(auth, " ")
	if len(tokens) != 2 || tokens[0] != "Token" {
		return utils.ErrorJSON(c, fiber.StatusBadRequest, "invalid token")
	}

	// validate token
	if tokens[1] != "secret" {
		return utils.ErrorJSON(c, fiber.StatusUnauthorized, "unauthorized")
	}

	// pass to next
	return c.Next()
}
