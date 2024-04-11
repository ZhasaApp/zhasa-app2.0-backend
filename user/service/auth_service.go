package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
	"zhasa2.0/user/entities"
	. "zhasa2.0/user/repository"
)

type AuthorizationService interface {
	RequestCode(phone entities.Phone) (entities.OtpId, error)
	Login(otpId entities.OtpId, code entities.OtpCode) (*entities.User, error)
	AdminLogin(phone entities.Phone, password string) (*entities.AuthUser, error)
}

type AuthResponse struct {
	UserTokenData
}

type SafeAuthorizationService struct {
	ctx                            context.Context
	recoveryService                RecoveryService
	getUserByPhoneFunc             GetUserByPhoneFunc
	getUserByPhoneWithPasswordFunc GetUserByPhoneWithPasswordFunc
	getUserByIdFunc                GetUserByIdFunc
	addUserCodeFunc                AddUserCodeFunc
	getAuthCodeByIdFunc            GetAuthCodeByIdFunc
	checkDisabledUserFunc          CheckDisabledUserFunc
}

func NewAuthorizationService(
	ctx context.Context,
	getUserByPhoneFunc GetUserByPhoneFunc,
	getUserByPhoneWithPasswordFunc GetUserByPhoneWithPasswordFunc,
	addUserCodeFunc AddUserCodeFunc,
	getUserByIdFunc GetUserByIdFunc,
	getAuthCodeByIdFunc GetAuthCodeByIdFunc,
	checkDisabledUserFunc CheckDisabledUserFunc,
) AuthorizationService {
	return SafeAuthorizationService{
		ctx:                            ctx,
		recoveryService:                NewRecoveryService(),
		getUserByPhoneFunc:             getUserByPhoneFunc,
		getUserByPhoneWithPasswordFunc: getUserByPhoneWithPasswordFunc,
		addUserCodeFunc:                addUserCodeFunc,
		getUserByIdFunc:                getUserByIdFunc,
		getAuthCodeByIdFunc:            getAuthCodeByIdFunc,
		checkDisabledUserFunc:          checkDisabledUserFunc,
	}
}

func (service SafeAuthorizationService) RequestCode(phone entities.Phone) (entities.OtpId, error) {
	user, err := service.getUserByPhoneFunc(phone)
	if err != nil {
		return 0, errors.New("user not found")
	}

	disabled, err := service.checkDisabledUserFunc(user.Id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	if disabled {
		return 0, errors.New("user is disabled")
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

func (service SafeAuthorizationService) AdminLogin(phone entities.Phone, password string) (*entities.AuthUser, error) {
	user, err := service.getUserByPhoneWithPasswordFunc(phone)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Password == "" {
		return nil, errors.New("no password set for user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
