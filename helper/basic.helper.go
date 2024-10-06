package helper

import "github.com/gofiber/fiber/v2"

// Generic method of getting a Fiber Context request header
func GetRequestHeader(ctx *fiber.Ctx, name string) string {
	var header string
	headers := ctx.GetReqHeaders()
	header = headers[name][0]
	return header
}
