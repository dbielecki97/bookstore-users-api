package app

import (
	"github.com/dbielecki97/bookstore-users-api/controller/ping"
	"github.com/dbielecki97/bookstore-users-api/controller/user"
)

func createUrlMappings() {
	router.GET("/ping", ping.Ping)

	router.GET("/users/:user_id", user.Get)
	router.POST("/users", user.Create)
	router.PUT("/users/:user_id", user.Update)
	router.PATCH("/users/:user_id", user.Update)
	router.DELETE("/users/:user_id", user.Delete)
	router.GET("/internal/users/search", user.Search)
	router.POST("/users/login", user.Login)
}
