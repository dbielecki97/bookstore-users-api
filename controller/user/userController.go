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
		restErr := errors.NewBadRequestError("invalid JSON body")
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
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		fmt.Println(idErr)
		c.JSON(idErr.StatusCode, idErr)
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

func Update(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		fmt.Println(idErr)
		c.JSON(idErr.StatusCode, idErr)
		return
	}

	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println(err)
		restErr := errors.NewBadRequestError("invalid JSON body")
		c.JSON(restErr.StatusCode, restErr)
		return
	}

	user.ID = userId

	var restErr *errors.RestErr
	var result *users.User
	if c.Request.Method == http.MethodPatch {
		result, restErr = services.PatchUser(user)
		if restErr != nil {
			c.JSON(restErr.StatusCode, restErr)
			return
		}
	} else {
		result, restErr = services.UpdateUser(user)
		if restErr != nil {
			c.JSON(restErr.StatusCode, restErr)
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		fmt.Println(idErr)
		c.JSON(idErr.StatusCode, idErr)
		return
	}

	if err := services.DeleteUser(userId); err != nil {
		fmt.Println(err)
		c.JSON(err.StatusCode, err)
		return

	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")

	results, err := services.Search(status)
	if err != nil {
		fmt.Println(err)
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, results)
}

func getUserId(param string) (int64, *errors.RestErr) {
	userId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		restErr := errors.NewBadRequestError("invalid user id")

		return 0, restErr
	}
	return userId, nil
}
