package request

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tolubydesign/go-jwt/config"
	"github.com/tolubydesign/go-jwt/helper"
	"github.com/tolubydesign/go-jwt/utils"
)

// Confirm user identity an their JWT token
func AuthenticateUserToken(ctx *fiber.Ctx, c *config.Config) error {
	// Get token from header
	log.Println("AuthenticateUserToken . capturing jwt token")
	tokenHeader := helper.GetRequestHeader(ctx, "Authentication")
	if tokenHeader == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "no token found")
	}

	// find auth substring
	found := strings.Contains(tokenHeader, "Bearer ")
	if !found {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid token provide")
	}

	// Remove token string start
	tokenString := strings.Replace(tokenHeader, "Bearer ", "", -1)

	// Verify that token
	jwtToken, fiberErr := utils.JWTVerification(tokenString)
	if fiberErr != nil {
		fmt.Errorf("AuthenticateUserToken . jwt verification. fiber error: %s", fiberErr.Message)
		return fiber.NewError(fiberErr.Code, fiberErr.Message)
	}

	if jwtToken == nil {
		return fiber.NewError(fiber.StatusBadRequest, "token untranslatable")
	}

	log.Println("AuthenticateUserToken . json web token valid. checking token claims")

	// Get Claims
	claims := jwtToken.Claims.(jwt.MapClaims)
	iss, err := claims.GetIssuer()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// check server issuer matches request issuer
	if iss != c.JWT.Issuer {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	log.Printf("AuthenticateUserToken . claims: %v", claims)
	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(fiber.Map{
		"message": "authorized",
		"valid":   true,
	})
}

// Request that builds and returns a JSON Web Token, based on user details the requester provides
func BuildJSONWebToken(ctx *fiber.Ctx, c *config.Config) error {
	// get user information from request body
	var err error
	var body struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Get data from fiber context
	byteBody := ctx.Body()
	// Convert Struct to JSON
	err = json.Unmarshal(byteBody, &body)
	// json, err := json.Marshal(body.Content)
	if err != nil {
		return err
	}

	// validate email and name
	err = utils.EmailValidation(body.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	err = utils.ValidateString(body.Name)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user := struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Name:  body.Name,
		Email: body.Email,
	}

	token, err := utils.BuildUserJWT(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(fiber.Map{
		"message": "token generated",
		"token":   token,
	})
}

// Request to expose the content of a json web token.
func TranslateJsonWebToken(ctx *fiber.Ctx, c *config.Config) error {
	// Get token from header
	log.Println("TranslateJsonWebToken . capturing jwt token")
	tokenHeader := helper.GetRequestHeader(ctx, "Authentication")
	if tokenHeader == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "no token found")
	}

	// find auth substring
	found := strings.Contains(tokenHeader, "Bearer ")
	if !found {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid token provide")
	}

	// Remove token string start
	tokenString := strings.Replace(tokenHeader, "Bearer ", "", -1)

	// Verify that token
	jwtToken, fiberErr := utils.JWTVerification(tokenString)
	if fiberErr != nil {
		fmt.Errorf("TranslateJsonWebToken . jwt verification. fiber error: %s", fiberErr.Message)
		return fiber.NewError(fiberErr.Code, fiberErr.Message)
	}

	if jwtToken == nil {
		return fiber.NewError(fiber.StatusBadRequest, "token untranslatable")
	}

	// Get Claims
	claims := jwtToken.Claims.(jwt.MapClaims)
	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(fiber.Map{
		"claims": claims,
	})
}
