package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Send custom error page
	err = ctx.Status(code).JSON(e)
	if err != nil {
		// In case the SendFile fails
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	// Return from handler
	return nil
}
