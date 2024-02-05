package model

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func NewUser(email, password string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: hashedPassword,
	}, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
