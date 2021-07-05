package users

import (
	"fmt"
	"github.com/dbielecki97/bookstore-users-api/datasource/mysql/userdb"
	"github.com/dbielecki97/bookstore-users-api/utils/date"
	"github.com/dbielecki97/bookstore-users-api/utils/errors"
	"strings"
)

const (
	insertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?)"

	emailUniqueEmail = "email_unique"
)

var (
	usersDB = make(map[int64]*User)
)

func (u *User) Get() *errors.RestErr {
	if err := userdb.Client.Ping(); err != nil {
		panic(err)
	}

	result := usersDB[u.ID]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", u.ID))
	}

	u.ID = result.ID
	u.Email = result.Email
	u.LastName = result.LastName
	u.FirstName = result.FirstName
	u.DateCreated = result.DateCreated

	return nil
}
func (u *User) Save() *errors.RestErr {
	stmt, err := userdb.Client.Prepare(insertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	u.DateCreated = date.GetNowString()

	result, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.DateCreated)
	if err != nil {
		if strings.Contains(err.Error(), emailUniqueEmail) {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", u.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying save user: %s", err))
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err))
	}

	u.ID = userId

	return nil
}
