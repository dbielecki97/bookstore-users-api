package errs

import (
	"database/sql"
	"github.com/dbielecki97/bookstore-users-api/utils/errors"
	"github.com/go-sql-driver/mysql"
)

const errDupEntry = 1062

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if err == sql.ErrNoRows {
			return errors.NewInternalServerError("no records matching given id")
		}

		return errors.NewInternalServerError("error processing request")
	}

	switch sqlErr.Number {
	case errDupEntry:
		return errors.NewBadRequestError("invalid data")
	}

	return errors.NewInternalServerError("error processing request")
}
