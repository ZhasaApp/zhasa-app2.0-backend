package service

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTokenService(t *testing.T) {
	encKey := []byte("YELLOW SUBMARINE, BLACK WIZARDRY")
	tokenService := PasetoTokenService{
		encryptionKey: &encKey,
	}

	testUser := &UserTokenData{
		id:        1,
		firstName: "Test",
		lastName:  "Tested",
		email:     "test@test.com",
	}

	token, err := tokenService.GenerateToken(testUser)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	token2, err := tokenService.GenerateToken(testUser)
	require.NoError(t, err)
	require.NotEmpty(t, token2)

	user, err := tokenService.VerifyToken(token)
	require.NoError(t, err)

	require.Equal(t, testUser, user)
}
