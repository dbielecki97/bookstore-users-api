package users

import (
	"database/sql"
	"fmt"
	"github.com/dbielecki97/bookstore-users-api/datasource/mysql/userdb"
	"github.com/dbielecki97/bookstore-users-api/logger"
	"github.com/dbielecki97/bookstore-users-api/utils/errors"
)

const (
	insertUser             = "INSERT INTO users(first_name, last_name, email, date_created, password, status) VALUES(?, ?, ?, ?, ?, ?);"
	getUser                = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?;"
	updateUser             = "UPDATE users SET first_name = ?, last_name = ?, email = ?, password = ? WHERE id = ?;"
	deleteUser             = "DELETE FROM users where id = ?;"
	findByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status = ?;"
	findByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status from users where email = ? and password = ?;"
	findByEmail            = "SELECT id, first_name, last_name, email, date_created, status, password from users where email = ?;"
)

func (u *User) Get() *errors.RestErr {
	if err := userdb.Client.Ping(); err != nil {
		panic(err)
	}

	stmt, err := userdb.Client.Preparex(getUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	row := stmt.QueryRowx(u.ID)
	if err := row.StructScan(u); err != nil {
		if err == sql.ErrNoRows {
			return errors.NewNotFoundError(fmt.Sprintf("user with id %d not found", u.ID))
		}
		logger.Error("error when trying to get user by id", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (u *User) Save() *errors.RestErr {
	stmt, err := userdb.Client.Prepare(insertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.DateCreated, u.Password, u.Status)
	if err != nil {
		logger.Error("error when trying to save user", err)
		return errors.NewInternalServerError("database error")
	}

	userId, err := result.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return errors.NewInternalServerError("database error")
	}

	u.ID = userId
	u.Password = ""

	return nil
}

func (u *User) Update() *errors.RestErr {
	stmt, err := userdb.Client.Prepare(updateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("database error")

	}
	defer stmt.Close()

	_, err = stmt.Exec(u.FirstName, u.LastName, u.Email, u.ID)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (u *User) Delete() *errors.RestErr {
	stmt, err := userdb.Client.Prepare(deleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.ID)
	if err != nil {
		logger.Error("error when trying to delete user", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (u *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := userdb.Client.Preparex(findByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.StructScan(&user)
		if err != nil {
			logger.Error("error when trying to struct scan user struct", err)
			return nil, errors.NewInternalServerError("database error")
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}

	return results, nil
}

func (u *User) FindByEmailAndPassword() *errors.RestErr {
	if err := userdb.Client.Ping(); err != nil {
		panic(err)
	}

	stmt, err := userdb.Client.Preparex(findByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	row := stmt.QueryRowx(u.Email, u.Password)
	if err := row.StructScan(u); err != nil {
		if err == sql.ErrNoRows {
			return errors.NewNotFoundError(fmt.Sprintf("invalid email or password"))
		}
		logger.Error("error when trying to get user by email and password", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (u *User) FindByEmail() *errors.RestErr {
	if err := userdb.Client.Ping(); err != nil {
		panic(err)
	}

	stmt, err := userdb.Client.Preparex(findByEmail)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	row := stmt.QueryRowx(u.Email)
	if err := row.StructScan(u); err != nil {
		if err == sql.ErrNoRows {
			return errors.NewAuthenticationError("invalid credentials")
		}
		logger.Error("error when trying to get user by email", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}
