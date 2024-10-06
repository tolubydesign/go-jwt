package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tolubydesign/go-jwt/config"
)

// Verify jwt token provided.
// Confirm that it came from us.
// Confirm that it hasn't expired.
func JWTVerification(token string) (*jwt.Token, *fiber.Error) {
	// var tkn jwt.Token
	fmt.Println("JWTVerification . verifying jwt token.")
	// Initialize a new instance of `Claims`
	fiberErr := &fiber.Error{}
	configuration, err := config.GetConfiguration()
	if err != nil {
		fiberErr.Code = fiber.StatusInternalServerError
		fiberErr.Message = err.Error()
		return nil, fiberErr
	}

	key := configuration.JWT.Secret
	secretString := []byte(key)

	err = ValidateString(token)
	if err != nil {
		fiberErr.Code = fiber.StatusBadRequest
		fiberErr.Message = err.Error()

		return nil, fiberErr
	}

	fmt.Println("JWTVerification . attempting to parsing jwt token.")
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return secretString, nil
	})

	// parsing errors result
	if err != nil {
		fmt.Println("JWTVerification . token parse error", err.Error())
		fiberErr.Code = fiber.StatusUnauthorized
		fiberErr.Message = err.Error()
		return nil, fiberErr
	}

	fmt.Println("JWTVerification . token validation check.")
	if !tkn.Valid {
		log.Println("in tkn invalid", tkn)
		fiberErr.Code = fiber.StatusUnauthorized
		fiberErr.Message = "invalid token"
		return nil, fiberErr
	}

	log.Println("JWTVerification . token parsing complete.")
	return tkn, nil
}

// Construct JWT token based on user information. Requires an id, email and name.
func BuildUserJWT(u struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}) (string, error) {
	// TODO: log security event
	// Get secret key from config
	configuration, err := config.GetConfiguration()
	if err != nil {
		return "", err
	}

	jwtConfig := configuration.JWT
	secret := configuration.JWT.Secret
	// fmt.Println("jwt config", configuration.JWT.Issuer)

	// Declare the expiration time of the token.
	now := time.Now()
	expirationTime := now.Add(time.Duration(jwtConfig.ExpiresAtAmount) * time.Hour)
	notBeforeTime := now.Add(time.Duration(jwtConfig.NotBeforeDateAmount) * time.Hour)

	// Create the JWT claims, which includes the username and expiry time
	// jwt.MapClaims
	claims := struct {
		Id                   string `json:"id"`
		Email                string `json:"email"`
		Name                 string `json:"name"`
		jwt.RegisteredClaims `json:"username"`
	}{
		Id:               u.Id,
		Email:            u.Email,
		Name:             u.Name,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	// In JWT, the expiry time is expressed as unix milliseconds
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	claims.Issuer = jwtConfig.Issuer
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.Audience = []string{jwtConfig.Audience}
	// Claim cant be older than the set date, in the past.
	claims.NotBefore = jwt.NewNumericDate(notBeforeTime)
	claims.Subject = jwtConfig.Subject

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println("created claims")

	// Create the JWT token string
	ts, err := token.SignedString(secret)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", err
	}

	return ts, err
}

func GetJWTSecretKey() ([]byte, error) {
	configuration, err := config.GetConfiguration()
	if err != nil {
		return nil, err
	}

	var secret = []byte(configuration.JWT.Secret)
	return secret, nil
}
