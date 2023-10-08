package service

import (
	"context"
	"errors"
	"time"
	"zhasa2.0/user/entities"

	"zhasa2.0/user/repository"
)

type AuthorizationService interface {
	RequestCode(phone entities.Phone) (entities.OtpId, error)
	Login(otpId entities.OtpId, code entities.OtpCode) (*entities.User, error)
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

func (service SafeAuthorizationService) RequestCode(phone entities.Phone) (entities.OtpId, error) {
	user, err := service.repo.GetUserByPhone(phone)
	if err != nil {
		return 0, errors.New("user not found")
	}

	otp, err := service.recoveryService.GenerateSendRecoveryCode(*user)
	id, err := service.repo.AddUserCode(entities.UserId(user.Id), *otp)
	return id, err
}

func (service SafeAuthorizationService) Login(otpId entities.OtpId, code entities.OtpCode) (*entities.User, error) {

	userId, err := service.VerifyRecoveryCode(otpId, code)
	if err != nil {
		return nil, err
	}

	user, err := service.repo.GetUserById(int32(userId))
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (service SafeAuthorizationService) VerifyRecoveryCode(otpId entities.OtpId, code entities.OtpCode) (entities.UserId, error) {

	auth, err := service.repo.GetAuthCodeById(otpId)

	if err != nil {
		return 0, errors.New("otp code not found")
	}

	if code != auth.Code {
		return 0, errors.New("otp code does not match")
	}

	if time.Now().After(auth.CreatedAt.Add(time.Minute)) {
		return 0, errors.New("otp code expired")
	}

	return auth.UserId, nil
}
