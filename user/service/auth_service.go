package service

import (
	"context"
	"errors"
	"zhasa2.0/user/entities"

	"zhasa2.0/user/repository"
)

type AuthorizationService interface {
	Login(phone entities.Phone, code recoveryCode) (*entities.User, error)
}

type SafeAuthorizationService struct {
	ctx             context.Context
	repo            repository.UserRepository
	recoveryService RecoveryService
}

func (service SafeAuthorizationService) SendCode(phone entities.Phone) error {
	user, err := service.repo.GetUserByPhone(phone)
	if err != nil {
		return err
	}

	return service.recoveryService.SendRecoveryCode(*user)
}

func (service SafeAuthorizationService) Login(phone entities.Phone, code recoveryCode) (*entities.User, error) {

	user, err := service.repo.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}

	verified := service.recoveryService.VerifyRecoveryCode(*user, code)
	if !verified {
		return nil, errors.New("wrong code")
	}
	return user, nil
}
