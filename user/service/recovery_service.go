package service

import (
	"crypto/rand"
	"math/big"
	"sync"
	"time"
	"zhasa2.0/user/entities"
)

type recoveryCode int

type recovery struct {
	Code       recoveryCode
	Expiration time.Time
	Used       bool
}

type CodeSender interface {
	sendRecoveryCode(code recoveryCode) error
}

type RecoveryCodeGenerator interface {
	generateRecoveryCode() (recoveryCode, error)
}

type FourDigitsRecoveryCodeGenerator struct{}

type TestPasswordRecoveryCodeGenerator struct{}

type RecoveryService struct {
	rcs CodeSender
	prg RecoveryCodeGenerator
}

type PhoneCodeSender struct {
	phone entities.Phone
}

func (ep FourDigitsRecoveryCodeGenerator) generateRecoveryCode() (recoveryCode, error) {
	code, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		return 0, err
	}
	codeResult := recoveryCode(code.Int64()) + 1000

	return codeResult, nil
}

func (t TestPasswordRecoveryCodeGenerator) generateRecoveryCode() (recoveryCode, error) {
	return 7777, nil
}

func (e PhoneCodeSender) sendRecoveryCode(recoveryCode) error {

	return nil
}

var recoveryCodes sync.Map

func (p RecoveryService) SendRecoveryCode(user entities.User) error {
	recoveryCode, err := p.prg.generateRecoveryCode()
	if err != nil {
		return err
	}

	expiration := time.Now().Add(time.Hour)

	recoveryCodes.Store(user.Id, recovery{
		Code:       recoveryCode,
		Expiration: expiration,
		Used:       false,
	})
	return p.rcs.sendRecoveryCode(recoveryCode)
}

func NewPasswordRecoveryService(phone entities.Phone) RecoveryService {
	return RecoveryService{
		rcs: PhoneCodeSender{
			phone: phone,
		},
		prg: TestPasswordRecoveryCodeGenerator{},
	}
}

func (p RecoveryService) VerifyRecoveryCode(user entities.User, code recoveryCode) bool {
	value, ok := recoveryCodes.Load(user.Id)
	if !ok {
		return false
	}

	recovery := value.(recovery)

	if recovery.Used {
		return false
	}

	if time.Now().After(recovery.Expiration) {
		return false
	}

	if recovery.Code != code {
		return false
	}

	// Mark the recovery code as used
	recovery.Used = true
	recoveryCodes.Store(user.Id, recovery)

	return true
}
