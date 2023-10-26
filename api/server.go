package api

import (
	"GoGin-API-CuentasClaras/api/middleware"
	"GoGin-API-CuentasClaras/api/routes"
	"GoGin-API-CuentasClaras/config"

	"github.com/gin-gonic/gin"
)

func Init(init *config.Initialization) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler())

	api := router.Group("/api")
	routes.UserRoutes(api, init)
	routes.OperationRoutes(api, init)

	return router
}
