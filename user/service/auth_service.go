package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"zhasa2.0/user/entities"
	. "zhasa2.0/user/repository"
)

type AuthorizationService interface {
	RequestCode(phone entities.Phone) (entities.OtpId, error)
	Login(otpId entities.OtpId, code entities.OtpCode) (*entities.User, error)
}

type AuthResponse struct {
	UserTokenData
}

type SafeAuthorizationService struct {
	ctx                 context.Context
	recoveryService     RecoveryService
	getUserByPhoneFunc  GetUserByPhoneFunc
	getUserByIdFunc     GetUserByIdFunc
	addUserCodeFunc     AddUserCodeFunc
	getAuthCodeByIdFunc GetAuthCodeByIdFunc
}

func NewAuthorizationService(
	ctx context.Context,
	getUserByPhoneFunc GetUserByPhoneFunc,
	addUserCodeFunc AddUserCodeFunc,
	getUserByIdFunc GetUserByIdFunc,
	getAuthCodeByIdFunc GetAuthCodeByIdFunc,
) AuthorizationService {
	return SafeAuthorizationService{
		ctx:                 ctx,
		recoveryService:     NewRecoveryService(),
		getUserByPhoneFunc:  getUserByPhoneFunc,
		addUserCodeFunc:     addUserCodeFunc,
		getUserByIdFunc:     getUserByIdFunc,
		getAuthCodeByIdFunc: getAuthCodeByIdFunc,
	}
}

func (service SafeAuthorizationService) RequestCode(phone entities.Phone) (entities.OtpId, error) {
	user, err := service.getUserByPhoneFunc(phone)
	if err != nil {
		return 0, errors.New("user not found")
	}

	otp, err := service.recoveryService.GenerateSendRecoveryCode(*user)
	id, err := service.addUserCodeFunc(entities.UserId(user.Id), *otp)
	return id, err
}

func (service SafeAuthorizationService) Login(otpId entities.OtpId, code entities.OtpCode) (*entities.User, error) {

	userId, err := service.VerifyRecoveryCode(otpId, code)
	if err != nil {
		return nil, err
	}

	user, err := service.getUserByIdFunc(int32(userId))
	if err != nil {
		fmt.Println(err, "userId: ", userId)
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (service SafeAuthorizationService) VerifyRecoveryCode(otpId entities.OtpId, code entities.OtpCode) (entities.UserId, error) {

	auth, err := service.getAuthCodeByIdFunc(otpId)

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
