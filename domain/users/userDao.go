package users

import (
	"fmt"
	"github.com/dbielecki97/bookstore-users-api/datasource/mysql/userdb"
	"github.com/dbielecki97/bookstore-users-api/utils/date"
	"github.com/dbielecki97/bookstore-users-api/utils/errors"
	"github.com/dbielecki97/bookstore-users-api/utils/mysql/errs"
)

const (
	insertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	getUser    = "SELECT id, first_name, last_name, email, date_created FROM users where id = ?;"
	updateUser = "UPDATE users SET first_name = ?, last_name = ?, email = ? where id = ?;"
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
		return errs.ParseError(err)
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
		return errs.ParseError(err)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err))
	}

	u.ID = userId

	return nil
}

func (u *User) Update() *errors.RestErr {
	stmt, err := userdb.Client.Prepare(updateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.FirstName, u.LastName, u.Email, u.ID)
	if err != nil {
		return errs.ParseError(err)
	}

	return nil
}
