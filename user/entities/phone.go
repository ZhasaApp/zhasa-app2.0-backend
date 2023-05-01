package entities

import (
	"errors"
	"regexp"
)

type Phone string

func NewPhone(phoneString string) (*Phone, error) {
	err := validatePhoneNumber(phoneString)
	if err != nil {
		return nil, err
	}
	phone := Phone(phoneString)
	return &phone, nil
}

func validatePhoneNumber(phoneNumber string) error {
	// Check that the phone number matches the pattern for a Kazakhstan phone number
	match, err := regexp.MatchString(`^\+7\d{10}$`, phoneNumber)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("phone number is not valid for Kazakhstan")
	}

	// Return nil if the phone number is valid
	return nil
}
