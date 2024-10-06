package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tolubydesign/go-jwt/config"
	"github.com/tolubydesign/go-jwt/request"
	"github.com/tolubydesign/go-jwt/utils"
)

// Setting up API endpoints and connecting relevant function
func SetupMethods(app *fiber.App, configuration *config.Config) {
	app.Get(utils.GET_authenticate, func(ctx *fiber.Ctx) error {
		return request.AuthenticateUserToken(ctx, configuration)
	})

	app.Get(utils.GET_build, func(ctx *fiber.Ctx) error {
		return request.BuildJSONWebToken(ctx, configuration)
	})

	app.Get(utils.GET_verify, func(ctx *fiber.Ctx) error {
		return request.VerifyJSONWebToken(ctx, configuration)
	})
}
