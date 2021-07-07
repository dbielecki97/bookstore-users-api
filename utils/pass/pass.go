package pass

import (
	"github.com/dbielecki97/bookstore-users-api/logger"
	"golang.org/x/crypto/bcrypt"
)

func Generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		logger.Error("error when hashing password", err)
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func Compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
