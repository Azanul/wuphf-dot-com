package model

import (
	"encoding/json"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Receivers string `json:"receivers"`
}

func NewUser(email, password string) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	receivers, err := json.Marshal(map[string]string{"email": email})
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  hashedPassword,
		Receivers: string(receivers),
	}, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
