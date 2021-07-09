package app

import (
	"github.com/dbielecki97/bookstore-utils-go/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	createUrlMappings()
	logger.Info("about to start the application...")
	logger.Fatal("could not start server", router.Run(":8080"))
}
