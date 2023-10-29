package api

import (
	"GoGin-API-CuentasClaras/api/middleware"
	"GoGin-API-CuentasClaras/api/routes"
	"GoGin-API-CuentasClaras/config"
	"os"

	"github.com/gin-gonic/gin"
)

func Init(init *config.Initialization) *gin.Engine {
	apiMode := os.Getenv("ENVIRONMENT")
	if apiMode == "production" || apiMode == "test" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler())

	api := router.Group("/api")
	routes.HealthRoutes(api, init)
	routes.UserRoutes(api, init)
	routes.OperationRoutes(api, init)

	return router
}
