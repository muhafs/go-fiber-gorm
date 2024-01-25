package utils

import "github.com/gofiber/fiber/v2"

func SuccessJSON(ctx *fiber.Ctx, code int, msg string, payload any) error {
	return ctx.Status(code).JSON(fiber.Map{
		"success": true,
		"message": msg,
		"data":    payload,
	})
}

func ErrorJSON(ctx *fiber.Ctx, code int, msg any) error {
	return ctx.Status(code).JSON(fiber.Map{
		"success": false,
		"message": msg,
	})
}
