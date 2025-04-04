package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var bcryptCost = 13

type ICryptoService interface {
	HashPassword(password string) (*string, error)
	ValidatePassword(password string, hash string) (bool, error)
}

type CryptoService struct {
	pepper string
}

func NewCryptoService(pepper string) ICryptoService {
	return &CryptoService{
		pepper: pepper,
	}
}

func (s *CryptoService) HashPassword(password string) (*string, error) {

	if len(password) == 0 {
		return nil, errors.New("password cannot be empty")
	} else if len(password) < 10 {
		return nil, errors.New("password must be at least 10 characters long")
	}

	var bytesToHash = []byte(password + s.pepper)
	var hashBytes, err = bcrypt.GenerateFromPassword(bytesToHash, bcryptCost)
	if err != nil {
		return nil, err
	}
	var hashedString = string(hashBytes)

	return &hashedString, nil
}

func (s *CryptoService) ValidatePassword(password string, hash string) (bool, error) {

	var bytesToHash = []byte(password + s.pepper)
	var err = bcrypt.CompareHashAndPassword([]byte(hash), bytesToHash)
	if err != nil {
		return false, err
	}
	return true, nil
}
