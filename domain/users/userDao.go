package users

import (
	"fmt"
	"github.com/dbielecki97/bookstore-users-api/utils/errors"
)

var (
	usersDB = make(map[int64]*User)
)

func (u *User) Save() *errors.RestErr {
	cur := usersDB[u.ID]
	if cur != nil {
		if cur.Email == u.Email {
			return errors.NewBadRequest(fmt.Sprintf("email %s already registered", u.Email))
		}

		return errors.NewBadRequest("user already exists")
	}

	usersDB[u.ID] = u

	return nil
}

func (u *User) Get() *errors.RestErr {
	result := usersDB[u.ID]
	if result == nil {
		return errors.NewNotFound(fmt.Sprintf("user %d not found", u.ID))
	}

	u.ID = result.ID
	u.Email = result.Email
	u.LastName = result.LastName
	u.FirstName = result.FirstName
	u.DateCreated = result.DateCreated

	return nil
}
