package service

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"zhasa2.0/user/entities"
	_ "zhasa2.0/user/repository"
)

type UserRepositoryStub struct {
	GetUserByPhoneFn func(phone entities.Phone) (*entities.User, error)
}

func (s *UserRepositoryStub) GetUserByPhone(phone entities.Phone) (*entities.User, error) {
	return s.GetUserByPhoneFn(phone)
}

func TestAuthorizationService(t *testing.T) {

	testCases := []struct {
		name           string
		userRepository *UserRepositoryStub
		phone          entities.Phone
		code           recoveryCode
		assert         func(user *entities.User, err error)
		// Add more fields if necessary, depending on your test case
	}{
		{
			name: "Test case user exist",
			userRepository: &UserRepositoryStub{
				GetUserByPhoneFn: func(email entities.Phone) (*entities.User, error) {
					return &entities.User{
						Phone: "test@gmail.com",
					}, nil
				},
			},
			phone: entities.Phone("+77777777"),
			assert: func(user *entities.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, user)
			},
			code: 7777,
		},
		{
			name: "Test case code doesnt not match",
			userRepository: &UserRepositoryStub{
				GetUserByPhoneFn: func(email entities.Phone) (*entities.User, error) {
					return &entities.User{
						Phone: "test@gmail.com",
					}, nil
				},
			},

			phone: "test@gmail.com",
			code:  4444,
			assert: func(user *entities.User, err error) {
				expected := errors.New("wrong code")
				require.Equal(t, expected, err)
			},
		},
		{
			name: "Test case user not found",
			userRepository: &UserRepositoryStub{
				GetUserByPhoneFn: func(email entities.Phone) (*entities.User, error) {
					return nil, errors.New("not found")
				},
			},

			phone: "not_exist@gmail.com",
			code:  5555,
			assert: func(user *entities.User, err error) {
				require.Error(t, err)
				require.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authService := SafeAuthorizationService{
				ctx:  nil,
				repo: tc.userRepository,
				recoveryService: RecoveryService{
					rcs: PhoneCodeSender{},
					prg: TestPasswordRecoveryCodeGenerator{},
				},
			}
			err := authService.SendCode(tc.phone)

			result, err := authService.Login(tc.phone, tc.code)

			tc.assert(result, err)
		})
	}
}
