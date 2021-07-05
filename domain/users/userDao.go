package users

import (
	"database/sql"
	"fmt"
	"github.com/dbielecki97/bookstore-users-api/datasource/mysql/userdb"
	"github.com/dbielecki97/bookstore-users-api/utils/date"
	"github.com/dbielecki97/bookstore-users-api/utils/errors"
	"strings"
)

const (
	insertUser       = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?)"
	getUser          = "SELECT id, first_name, last_name, email, date_created FROM users where id = ? "
	emailUniqueEmail = "email_unique"
)

func (u *User) Get() *errors.RestErr {
	if err := userdb.Client.Ping(); err != nil {
		panic(err)
	}

	stmt, err := userdb.Client.Preparex(getUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRowx(u.ID)
	if err := row.StructScan(u); err != nil {
		if err == sql.ErrNoRows {
			return errors.NewNotFoundError(fmt.Sprintf("user with id %d not found", u.ID))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when tring to get uset %d: %s", u.ID, err))
	}

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
