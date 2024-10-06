package utils

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// Confirm that the provide uuid is a valid one
func ValidUuid(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// NOTE: The Laws of Reflection - https://go.dev/blog/laws-of-reflection
func ValidateLimitedStringVariable(str string) error {
	value := reflect.TypeOf(str)
	if value.Kind() != reflect.String {
		return errors.New("variable is not a valid string")
	}

	if strings.ReplaceAll(str, " ", "") == "" {
		return errors.New("valid string must be provided")
	}

	if len(str) > 255 {
		return errors.New("string is too long")
	}

	return nil
}

func AdditionalStringValidation(str string) error {
	value := reflect.TypeOf(str)
	if value.Kind() != reflect.String {
		return errors.New("additional validation found. Invalid value provided")
	}

	return nil
}

func ValidateString(str string) error {
	value := reflect.TypeOf(str)
	if value.Kind() != reflect.String {
		return errors.New("variable is not a valid string")
	}

	if strings.ReplaceAll(str, " ", "") == "" {
		return errors.New("valid string must be provided")
	}

	return nil
}

// Validate that email provide follows the expected email pattern.
// Prevent invalid email strings from being added to the database.
func EmailValidation(str string) error {
	matched, err := regexp.Match(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, []byte(str))
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("invalid email provided")
	}

	return nil
}

// Validate that role provided is correct
// Roles can only be `admin` or `user`. Default is `user`
// func UserRoleCheck(role string) string {
// 	if role == "admin" {
// 		return "admin"
// 	} else {
// 		return "user"
// 	}
// }
