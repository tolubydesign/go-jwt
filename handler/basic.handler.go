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

// Basic method to handle request responses
func HandleResponse(context ResponseHandlerParameters) error {
	c := context.Context
	if c == nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("Request context was not provided. Message: ", context.Message))
	}

	c.Response().StatusCode()
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return c.Status(context.Code).JSON(fiber.Map{
		"error":   context.Error,
		"code":    context.Code,
		"message": context.Message,
		"data":    context.Data,
	})
}

// func HandleError(context ResponseHandlerParameters) error {
// 	c := context.Context

// 	return fiber.NewError(fiber.StatusInternalServerError, error.Error(err))
// }
