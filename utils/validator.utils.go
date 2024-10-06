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
