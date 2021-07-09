package userdb

import (
	"fmt"
	"github.com/dbielecki97/bookstore-utils-go/logger"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
)

const (
	mysqlUsersUsername = "MYSQL_USERS_USERNAME"
	mysqlUsersPassword = "MYSQL_USERS_PASSWORD"
	mysqlUsersHost     = "MYSQL_USERS_HOST"
	mysqlUsersSchema   = "MYSQL_USERS_SCHEMA"
)

var (
	Client   *sqlx.DB
	username = os.Getenv(mysqlUsersUsername)
	password = os.Getenv(mysqlUsersPassword)
	host     = os.Getenv(mysqlUsersHost)
	schema   = os.Getenv(mysqlUsersSchema)
)

func init() {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username,
		password,
		host,
		schema)

	Client, err = sqlx.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}

	err = mysql.SetLogger(logger.GetLogger())
	if err != nil {
		logger.Error("could not set up logger for mysql", err)
	}
	logger.Info("database successfully configured")
}
