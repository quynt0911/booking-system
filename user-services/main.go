package main

import (
	"user-services/config"
	"user-services/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	router := gin.Default()

	routes.AuthRoutes(router)

	router.Run(":8080")
}
