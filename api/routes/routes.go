package routes

import (
	"GoGin-API-Base/api/middleware"
	"GoGin-API-Base/config"

	"github.com/gin-gonic/gin"
)

var initConfig *config.Initialization = config.Init()

func UserRoutes(router *gin.RouterGroup) {
	user := router.Group("/users")
	{
		user.POST("", initConfig.UserHdler.RegisterUser)
		user.POST("/login", initConfig.UserHdler.LoginUser)
		user.GET("/current", middleware.AuthMiddleware(), initConfig.UserHdler.CurrentUser)
	}
}

func OperationRoutes(router *gin.RouterGroup) {
	operation := router.Group("/operations")
	{
		operation.GET("", middleware.AuthMiddleware(), initConfig.OperationHdler.Index)
	}
}
