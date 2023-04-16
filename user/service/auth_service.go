package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"zhasa2.0/user/entities"

	"zhasa2.0/user/repository"
)

type AuthorizationService interface {
	Login(body LoginBody) (*entities.User, error)
	ChangePassword(body ChangePasswordBody) error
}

type SafeAuthorizationService struct {
	ctx  context.Context
	repo repository.UserRepository
	enc  entities.PasswordEncryptor
}

type LoginBody struct {
	Email    string
	Password string
}

type ChangePasswordBody struct {
	Email       string
	NewPassword string
}

func (service SafeAuthorizationService) Login(body LoginBody) (*entities.User, error) {
	email, err := entities.NewEmail(body.Email)
	if err != nil {
		return nil, err
	}

	user, err := service.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.Encrypted), []byte(body.Password)); err != nil {
		return nil, errors.New("password doesn't not match")
	}

	return user, nil
}

func (service SafeAuthorizationService) ChangePassword(body ChangePasswordBody) error {
	email, err := entities.NewEmail(body.Email)
	if err != nil {
		return nil
	}

	password, err := entities.NewPassword(body.NewPassword, service.enc)
	if err != nil {
		return err
	}

	err = service.repo.ChangePassword(email, *password)
	return err
}
