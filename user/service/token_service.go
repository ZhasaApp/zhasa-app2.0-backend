package service

import (
	"fmt"
	"github.com/o1egl/paseto"
	"strconv"
	"time"
	"zhasa2.0/user/entities"
)

type Token string

type UserTokenData struct {
	id        int
	firstName string
	lastName  string
	email     string
}

// TokenService responsible for conversion sensitive user data into secure view
type TokenService interface {
	GenerateToken(user *entities.User) (Token, error)
}

// PasetoTokenService implementats TokenService by using Paseto specification
type PasetoTokenService struct {
	encryptionKey *[]byte
}

func (ts *PasetoTokenService) GenerateToken(user *UserTokenData) (Token, error) {
	v2 := paseto.NewV2()

	jsonToken := paseto.JSONToken{}
	jsonToken.Set("id", fmt.Sprintf("%d", user.id))
	jsonToken.Set("email", user.email)
	jsonToken.Set("first_name", user.firstName)
	jsonToken.Set("last_name", user.lastName)

	footer := map[string]interface{}{
		"issued_at": time.Now(),
	}

	tokenString, err := v2.Encrypt(*ts.encryptionKey, jsonToken, footer)
	if err != nil {
		return "", err
	}
	return Token(tokenString), nil
}

func (ts *PasetoTokenService) VerifyToken(token Token) (*UserTokenData, error) {
	v2 := paseto.NewV2()

	var extractedData map[string]interface{}
	var footer map[string]interface{}
	err := v2.Decrypt(string(token), *ts.encryptionKey, &extractedData, &footer)
	if err != nil {
		return nil, err
	}
	userId, _ := strconv.Atoi(extractedData["id"].(string))

	// Create a user object from the extracted data
	user := UserTokenData{
		id:        userId,
		email:     extractedData["email"].(string),
		firstName: extractedData["first_name"].(string),
		lastName:  extractedData["last_name"].(string),
	}

	return &user, nil
}
