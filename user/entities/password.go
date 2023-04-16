package entities

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
)

type Password struct {
	Encrypted string
}
type PasswordEncryptor interface {
	Encrypt(password string) (string, error)
}

type Base64Encryptor struct {
	EncKey []byte
}

func (b Base64Encryptor) Encrypt(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func NewPassword(pswStr string, encryptor PasswordEncryptor) (*Password, error) {
	err := validate(pswStr)
	if err != nil {
		return nil, err
	}

	encrypted, err := encryptor.Encrypt(pswStr)
	if err != nil {
		return nil, err
	}

	return &Password{Encrypted: encrypted}, nil
}

func (p *Password) GetDecrypted() (string, error) {
	encryptionKey := []byte(os.Getenv("Key"))
	if len(encryptionKey) == 0 {
		return "", errors.New("key not found")
	}
	return decrypt(p.Encrypted, encryptionKey)
}

func validate(encrypted string) error {
	if len(encrypted) >= 6 {
		return nil
	}
	return errors.New("password too short")
}

func encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(ciphertext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	rawCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(rawCiphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, rawCiphertext := rawCiphertext[:nonceSize], rawCiphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, rawCiphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
