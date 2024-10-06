package main

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/tolubydesign/go-jwt/config"
	"github.com/tolubydesign/go-jwt/controller"
)

func main() {
	// Setup project configuration
	c, err := config.BuildConfiguration()
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		panic(err)
	}

	fmt.Println("environment", c.Environment)
	fmt.Println("configuration issuer", c.JWT.Issuer)
	fmt.Println("configuration secret", c.JWT.Secret)
	port := "3255"
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Go JWT",
		AppName:       "JSON Web Token",
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).SendString(err.Error())
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			// Return from handler
			return nil
		},
	})

	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept",
		// AllowOrigins:     "*",
		// AllowCredentials: true,
		AllowMethods: "GET,POST,HEAD,PUT",
	}))
	controller.SetupMethods(app, c)
	app.Listen(fmt.Sprintf(":%v", port))
}
