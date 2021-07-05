package user

import (
	"fmt"
	"github.com/dbielecki97/bookstore-users-api/domain/users"
	"github.com/dbielecki97/bookstore-users-api/services"
	"github.com/dbielecki97/bookstore-users-api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Create(c *gin.Context) {
	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println(err)
		restErr := errors.NewBadRequest("invalid JSON body")
		c.JSON(restErr.StatusCode, restErr)
		return
	}

	result, err := services.CreateUser(user)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func Get(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		restErr := errors.NewBadRequest("invalid user id")
		fmt.Println(restErr)
		c.JSON(restErr.StatusCode, restErr)
		return
	}

	result, restErr := services.GetUser(userId)
	if restErr != nil {
		fmt.Println(restErr)
		c.JSON(restErr.StatusCode, restErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func Search(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me!")
}
