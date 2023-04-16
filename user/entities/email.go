package entities

import (
	"errors"
	"regexp"
)

type Email string

func NewEmail(email string) (Email, error) {
	if !isValidEmail(email) {
		return "", errors.New("invalid email format")
	}
	return Email(email), nil
}

func isValidEmail(email string) bool {
	regex := `^\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	emailValidator := regexp.MustCompile(regex)
	return emailValidator.MatchString(email)
}
