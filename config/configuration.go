package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type JWTConfiguration struct {
	Issuer              string `json:"issuer"`
	Audience            string `json:"audience"`
	NotBeforeDateAmount int    `json:"notBeforeDateAmount"`
	ExpiresAtAmount     int    `json:"expiresAtAmount"`
	Subject             string `json:"subject"`
	Secret              []byte `json:"secret"`
}

type Config struct {
	Environment string           `json:"environment"`
	JWT         JWTConfiguration `json:"jwt"`
}

var configSingleton *Config

/*
Build and return the environmental configuration.

Returns Configuration or error, if issues occur.
*/
func BuildConfiguration() (*Config, error) {
	envArg := os.Args[1]
	environmentPath := fmt.Sprintf(".env.%s", envArg)
	log.Println("environment: ", envArg, " | .env file: ", environmentPath)

	// Deny processing if environment argument isn't what we want
	if (envArg == "development") || (envArg == "production") {
		// Note: Alternative method of getting a .env file
		var envs map[string]string
		envs, err := godotenv.Read(environmentPath)
		if err != nil {
			return nil, err
		}

		environment := envs["ENV"]
		secret := envs["JWT_SECRET_KEY"]
		jwtExpiresAtString := envs["JWT_EXPIRES_AT"]
		jwtNotBeforeString := envs["JWT_NOT_BEFORE"]
		jwtExpiresAt, err := strconv.Atoi(jwtExpiresAtString)
		if err != nil {
			return nil, err
		}

		jwtNotBefore, err := strconv.Atoi(jwtNotBeforeString)
		if err != nil {
			return nil, err
		}

		jwt := JWTConfiguration{
			Issuer:              envs["JWT_ISSUER"],
			Audience:            envs["JWT_AUDIENCE"],
			NotBeforeDateAmount: jwtNotBefore - (jwtNotBefore * 2),
			ExpiresAtAmount:     jwtExpiresAt,
			Subject:             envs["JWT_SUBJECT"],
			Secret:              []byte(secret),
		}

		configSingleton = &Config{
			Environment: environment,
			JWT:         jwt,
		}

		return configSingleton, nil
	}

	return nil, errors.New("incorrect environment variables provided")
}

func GetConfiguration() (*Config, error) {
	if configSingleton == nil {
		// Build configuration
		build, e := BuildConfiguration()
		if e != nil {
			return nil, errors.New("project configuration error")
		}

		return build, nil
	}

	return configSingleton, nil
}
