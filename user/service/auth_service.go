package service

import (
	"context"
	"errors"
	"time"
	"zhasa2.0/user/entities"

	"zhasa2.0/user/repository"
)

type AuthorizationService interface {
	RequestCode(phone entities.Phone) (int32, error)
	Login(phone entities.Phone, code entities.OtpCode) (*entities.User, error)
}

type AuthResponse struct {
	UserTokenData
}

type SafeAuthorizationService struct {
	ctx             context.Context
	repo            repository.UserRepository
	recoveryService RecoveryService
}

func NewAuthorizationService(ctx context.Context, repo repository.UserRepository) AuthorizationService {
	return SafeAuthorizationService{
		ctx:             ctx,
		repo:            repo,
		recoveryService: NewRecoveryService(),
	}
}

func (service SafeAuthorizationService) RequestCode(phone entities.Phone) (int32, error) {
	user, err := service.repo.GetUserByPhone(phone)
	if err != nil {
		return 0, errors.New("user not found")
	}

	otp, err := service.recoveryService.GenerateSendRecoveryCode(*user)
	id, err := service.repo.AddUserCode(user.Id, int32(*otp))
	return id, err
}

func (service SafeAuthorizationService) Login(phone entities.Phone, code entities.OtpCode) (*entities.User, error) {

	user, err := service.repo.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}

	err = service.VerifyRecoveryCode(*user, code)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service SafeAuthorizationService) VerifyRecoveryCode(user entities.User, code entities.OtpCode) error {

	auth, err := service.repo.GetActualUserCode(user.Id)

	if err != nil {
		return err
	}

	if code != auth.Code {
		return errors.New("otp code does not match")
	}

	if time.Now().After(auth.CreatedAt.Add(time.Minute)) {
		return errors.New("otp code expired")
	}

	return nil
}
