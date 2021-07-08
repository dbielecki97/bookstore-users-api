package user

import (
	"github.com/dbielecki97/bookstore-oauth-go/oauth"
	"github.com/dbielecki97/bookstore-users-api/domain/users"
	"github.com/dbielecki97/bookstore-users-api/services"
	"github.com/dbielecki97/bookstore-utils-go/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Create(c *gin.Context) {
	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errs.NewBadRequestErr("invalid JSON body")
		c.JSON(restErr.StatusCode, restErr)
		return
	}

	result, err := services.UserService.CreateUser(user)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}
	result.Marshall(c.GetHeader("X-Public") == "true")
	c.JSON(http.StatusCreated, result)
}

func Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.StatusCode, idErr)
		return
	}

	user, restErr := services.UserService.GetUser(userId)
	if restErr != nil {
		c.JSON(restErr.StatusCode, restErr)
		return
	}

	if oauth.GetCallerId(c.Request) == user.ID {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}

	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func Update(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.StatusCode, idErr)
		return
	}

	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errs.NewBadRequestErr("invalid JSON body")
		c.JSON(restErr.StatusCode, restErr)
		return
	}

	user.ID = userId

	var restErr *errs.RestErr
	var result *users.User
	if c.Request.Method == http.MethodPatch {
		result, restErr = services.UserService.PatchUser(user)
		if restErr != nil {
			c.JSON(restErr.StatusCode, restErr)
			return
		}
	} else {
		result, restErr = services.UserService.UpdateUser(user)
		if restErr != nil {
			c.JSON(restErr.StatusCode, restErr)
			return
		}
	}

	c.JSON(http.StatusOK, result.Marshall(oauth.IsPublic(c.Request)))
}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.StatusCode, idErr)
		return
	}

	if err := services.UserService.DeleteUser(userId); err != nil {
		c.JSON(err.StatusCode, err)
		return

	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")

	results, err := services.UserService.Search(status)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, results.Marshall(oauth.IsPublic(c.Request)))
}

func Login(c *gin.Context) {
	var r users.LoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		restErr := errs.NewBadRequestErr("invalid json body")
		c.JSON(restErr.StatusCode, restErr)
		return
	}

	var u *users.User
	var err *errs.RestErr
	if u, err = services.UserService.FindByEmail(r); err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, u.Marshall(oauth.IsPublic(c.Request)))
}

func getUserId(param string) (int64, *errs.RestErr) {
	userId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		restErr := errs.NewBadRequestErr("invalid user id")

		return 0, restErr
	}
	return userId, nil
}
