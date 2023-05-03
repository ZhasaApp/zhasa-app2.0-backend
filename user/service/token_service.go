package service

import (
	"fmt"
	"github.com/o1egl/paseto"
	"strconv"
	"time"
)

type Token string

type UserTokenData struct {
	Id        int32
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// TokenService responsible for conversion sensitive user data into secure view
type TokenService interface {
	GenerateToken(user *UserTokenData) (Token, error)
	VerifyToken(token Token) (*UserTokenData, error)
}

// PasetoTokenService implements TokenService by using Paseto specification
type PasetoTokenService struct {
	encryptionKey *[]byte
}

func NewTokenService(encryptionKey *[]byte) TokenService {
	return PasetoTokenService{
		encryptionKey: encryptionKey,
	}
}

func (ts PasetoTokenService) GenerateToken(user *UserTokenData) (Token, error) {
	v2 := paseto.NewV2()

	jsonToken := paseto.JSONToken{}
	jsonToken.Set("Id", fmt.Sprintf("%d", user.Id))
	jsonToken.Set("first_name", user.FirstName)
	jsonToken.Set("last_name", user.LastName)
	jsonToken.Set("phone", user.Phone)

	footer := map[string]interface{}{
		"issued_at": time.Now(),
	}

	tokenString, err := v2.Encrypt(*ts.encryptionKey, jsonToken, footer)
	if err != nil {
		return "", err
	}
	return Token(tokenString), nil
}

func (ts PasetoTokenService) VerifyToken(token Token) (*UserTokenData, error) {
	v2 := paseto.NewV2()

	var extractedData map[string]interface{}
	var footer map[string]interface{}
	err := v2.Decrypt(string(token), *ts.encryptionKey, &extractedData, &footer)
	if err != nil {
		return nil, err
	}
	userId, _ := strconv.Atoi(extractedData["Id"].(string))

	// Create a user object from the extracted data
	user := UserTokenData{
		Id:        int32(userId),
		FirstName: extractedData["first_name"].(string),
		LastName:  extractedData["last_name"].(string),
		Phone:     extractedData["phone"].(string),
	}

	return &user, nil
}
