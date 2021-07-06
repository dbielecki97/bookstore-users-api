package users

import (
	"fmt"
	"github.com/dbielecki97/bookstore-users-api/datasource/mysql/userdb"
	"github.com/dbielecki97/bookstore-users-api/utils/errors"
	"github.com/dbielecki97/bookstore-users-api/utils/mysql/errs"
)

const (
	insertUser       = "INSERT INTO users(first_name, last_name, email, date_created, password, status) VALUES(?, ?, ?, ?, ?, ?);"
	getUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?;"
	updateUser       = "UPDATE users SET first_name = ?, last_name = ?, email = ?, password = ? WHERE id = ?;"
	deleteUser       = "DELETE FROM users where id = ?;"
	findUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status = ?;"
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

	result, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.DateCreated, u.Password, u.Status)
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

func (u *User) Delete() *errors.RestErr {
	stmt, err := userdb.Client.Prepare(deleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.ID)
	if err != nil {
		return errs.ParseError(err)
	}

	return nil
}

func (u *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := userdb.Client.Preparex(findUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.StructScan(&user)
		if err != nil {
			return nil, errs.ParseError(err)
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}

	return results, nil
}
