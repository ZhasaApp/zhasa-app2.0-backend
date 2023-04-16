package service

import (
	"github.com/stretchr/testify/require"
	"testing"
	"zhasa2.0/user/entities"
	_ "zhasa2.0/user/repository"
)

type UserRepositoryStub struct {
	GetUserByEmailFn func(email entities.Email) (*entities.User, error)
	ChangePasswordFn func(email entities.Email, password entities.Password) error
}

func (s *UserRepositoryStub) GetUserByEmail(email entities.Email) (*entities.User, error) {
	return s.GetUserByEmailFn(email)
}

func (s *UserRepositoryStub) ChangePassword(email entities.Email, password entities.Password) error {
	return s.ChangePasswordFn(email, password)
}

func TestYourFunction(t *testing.T) {
	encKey := []byte("YELLOW SUBMARINE, BLACK WIZARDRY")

	pEncryptor := entities.Base64Encryptor{
		EncKey: encKey,
	}

	testPassword, err := entities.NewPassword("testPassword", pEncryptor)
	require.NoError(t, err)
	testCases := []struct {
		name           string
		userRepository *UserRepositoryStub
		loginBody      LoginBody
		assert         func(user *entities.User, err error)
		// Add more fields if necessary, depending on your test case
	}{
		{
			name: "Test case user exist",
			userRepository: &UserRepositoryStub{
				GetUserByEmailFn: func(email entities.Email) (*entities.User, error) {
					return &entities.User{
						Email:    "test@gmail.com",
						Password: *testPassword,
					}, nil
				},
				ChangePasswordFn: func(email entities.Email, password entities.Password) error {
					return nil
				},
			},
			loginBody: LoginBody{
				Email:    "test@gmail.com",
				Password: "testPassword",
			},
			assert: func(user *entities.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, user)
			},
		},
		{
			name: "Test case password doesnt not match",
			userRepository: &UserRepositoryStub{
				GetUserByEmailFn: func(email entities.Email) (*entities.User, error) {
					return &entities.User{
						Email:    "test@gmail.com",
						Password: *testPassword,
					}, nil
				},
				ChangePasswordFn: func(email entities.Email, password entities.Password) error {
					return nil
				},
			},
			loginBody: LoginBody{
				Email:    "test@gmail.com",
				Password: "wrong password",
			},
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
				enc:  pEncryptor,
			}

			result, err := authService.Login(LoginBody{
				Email:    tc.loginBody.Email,
				Password: tc.loginBody.Password,
			})

			tc.assert(result, err)
		})
	}
}
