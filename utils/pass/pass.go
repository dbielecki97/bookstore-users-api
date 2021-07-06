package pass

import (
	"fmt"
	"github.com/dbielecki97/bookstore-users-api/utils/errors"
	"golang.org/x/crypto/bcrypt"
)

func Generate(s string) (string, *errors.RestErr) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return "", errors.NewInternalServerError("error processing request")
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func Compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
