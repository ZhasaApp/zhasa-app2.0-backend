package entities

import (
	"errors"
	"regexp"
	"time"
)

type Avatar struct {
	Url string
}

type CreateUserRequest struct {
	Phone     Phone
	FirstName Name
	LastName  Name
}

type User struct {
	Id        int32
	Phone     Phone
	Avatar    Avatar
	FirstName Name
	LastName  Name
}

type UserAuth struct {
	Code      OtpCode
	UserId    int32
	CreatedAt time.Time
}

type Name string

type OtpCode int

func NewName(name string) (*Name, error) {
	// Check that the name is not empty
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	// Check that the name matches the pattern for a valid name
	match, err := regexp.MatchString(`^[A-Za-z][A-Za-z'-]*[A-Za-z]$`, name)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, errors.New("name is not valid")
	}
	res := Name(name)
	return &res, nil
}
