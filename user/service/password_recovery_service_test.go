package service

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"zhasa2.0/user/entities"
)

type TestRecoveryCodeSender struct {
	codesSent []recoveryCode
}

func (t *TestRecoveryCodeSender) sendRecoveryCode(code recoveryCode) error {
	t.codesSent = append(t.codesSent, code)
	return nil
}

func TestPasswordRecoveryService(t *testing.T) {
	// Create a test user
	user := entities.User{
		Id:    1,
		Email: "test@example.com",
	}

	// Initialize the test recovery code sender
	testRCS := &TestRecoveryCodeSender{
		codesSent: make([]recoveryCode, 0),
	}

	// Create a password recovery service with the test recovery code sender and generator
	prs := PasswordRecoveryService{
		rcs: testRCS,
		prg: TestPasswordRecoveryCodeGenerator{},
	}

	// Test sending a recovery code
	err := prs.SendRecoveryCode(user)
	require.NoError(t, err)

	// Test verifying wrong code
	wrongCode := testRCS.codesSent[0] + 1
	valid := prs.VerifyRecoveryCode(user, wrongCode)
	require.False(t, valid)

	// Test verifying the recovery code
	code := testRCS.codesSent[0]
	valid = prs.VerifyRecoveryCode(user, code)

	require.True(t, valid)

	// Test verifying the recovery code again, it should be marked as used and invalid
	valid = prs.VerifyRecoveryCode(user, code)
	require.False(t, valid)

	// Test expiration by setting the expiration to the past and verifying the recovery code
	recoveryCodes.Store(user.Id, recovery{
		Code:       code,
		Expiration: time.Now().Add(-1 * time.Hour),
		Used:       false,
	})

	valid = prs.VerifyRecoveryCode(user, code)

	require.False(t, valid)
}
