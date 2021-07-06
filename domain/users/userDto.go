package users

import (
	"github.com/dbielecki97/bookstore-users-api/utils/errors"
	"strings"
)

type User struct {
	ID          int64  `json:"id,omitempty"`
	FirstName   string `json:"first_name,omitempty" db:"first_name"`
	LastName    string `json:"last_name,omitempty" db:"last_name"`
	Email       string `json:"email,omitempty"`
	DateCreated string `json:"date_created,omitempty" db:"date_created"`
}

func (u *User) Validate() *errors.RestErr {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)

	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if u.Email == "" {
		return errors.NewBadRequestError("invalid email address")
	}

	return nil
}
