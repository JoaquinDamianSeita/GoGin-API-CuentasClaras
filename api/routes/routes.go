package routes

import (
	"GoGin-API-CuentasClaras/api/middleware"
	"GoGin-API-CuentasClaras/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthRoutes(router *gin.RouterGroup, initConfig *config.Initialization) {
	health_check := router.Group("/health_check")
	{
		health_check.GET("", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "OK"}) })
	}
}

func UserRoutes(router *gin.RouterGroup, initConfig *config.Initialization) {
	user := router.Group("/users")
	{
		user.POST("", initConfig.UserHdler.RegisterUser)
		user.POST("/login", initConfig.UserHdler.LoginUser)
		user.GET("/current", middleware.AuthMiddleware(initConfig), initConfig.UserHdler.CurrentUser)
	}
}

func OperationRoutes(router *gin.RouterGroup, initConfig *config.Initialization) {
	operation := router.Group("/operations")
	{
		operation.GET("", middleware.AuthMiddleware(initConfig), initConfig.OperationHdler.Index)
		operation.GET("/:id", middleware.AuthMiddleware(initConfig), initConfig.OperationHdler.Show)
		operation.POST("", middleware.AuthMiddleware(initConfig), initConfig.OperationHdler.Create)
		operation.PUT("/:id", middleware.AuthMiddleware(initConfig), initConfig.OperationHdler.Update)
	}
}
