package usecase

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	salt = 10
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), salt)

	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func IsMatchHashAndPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
