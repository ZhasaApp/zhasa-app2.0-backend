package service

import (
	"crypto/rand"
	"math/big"
	"zhasa2.0/user/entities"
)

type CodeSender interface {
	SendRecoveryCode(phone entities.Phone, code entities.OtpCode) error
}

type RecoveryCodeGenerator interface {
	GenerateRecoveryCode() (entities.OtpCode, error)
}

type FourDigitsRecoveryCodeGenerator struct{}

func (ep FourDigitsRecoveryCodeGenerator) SendRecoveryCode(phone entities.Phone, code entities.OtpCode) error {
	return nil
}

type TestPasswordRecoveryCodeGenerator struct{}

type RecoveryService struct {
	CodeSender
	RecoveryCodeGenerator
}

func NewRecoveryService() RecoveryService {
	return RecoveryService{
		FourDigitsRecoveryCodeGenerator{},
		TestPasswordRecoveryCodeGenerator{},
	}
}

type PhoneCodeSender struct {
	phone entities.Phone
}

func (ep FourDigitsRecoveryCodeGenerator) GenerateRecoveryCode() (entities.OtpCode, error) {
	code, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		return 0, err
	}
	codeResult := entities.OtpCode(code.Int64()) + 1000

	return codeResult, nil
}

func (t TestPasswordRecoveryCodeGenerator) GenerateRecoveryCode() (entities.OtpCode, error) {
	return 7777, nil
}

func (e PhoneCodeSender) SendRecoveryCode(entities.OtpCode) error {
	return nil
}

func (p RecoveryService) GenerateSendRecoveryCode(user entities.User) (*entities.OtpCode, error) {
	recoveryCode, err := p.GenerateRecoveryCode()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &recoveryCode, p.SendRecoveryCode(user.Phone, recoveryCode)
}
