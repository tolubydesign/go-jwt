package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ResponseHandlerParameters struct {
	Context *fiber.Ctx
	Error   bool
	Code    int
	Message string
	Data    interface{}
}

// Override default error handler
var ErrorHandler = func(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Send custom error page
	err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
	if err != nil {
		// In case the SendFile fails
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Return from handler
	return nil
}
