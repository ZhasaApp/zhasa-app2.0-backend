package entities

import (
	"errors"
	"regexp"
)

type Phone string

func NewPhone(phoneString string) (Phone, error) {
	if !isValidPhone(phoneString) {
		return "", errors.New("invalid phone format")
	}
	return Phone(phoneString), nil
}

func isValidPhone(email string) bool {
	regex := `^\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	emailValidator := regexp.MustCompile(regex)
	return emailValidator.MatchString(email)
}
