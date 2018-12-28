package helper

import (
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct{}

func (h *Bcrypt) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (h *Bcrypt) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewBcryptHelper() *Bcrypt {
	return &Bcrypt{}
}
