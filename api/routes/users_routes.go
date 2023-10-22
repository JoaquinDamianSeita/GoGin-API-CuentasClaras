package routes

import (
	"GoGin-API-Base/api/middleware"
	"GoGin-API-Base/config"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup) {
	init := config.Init()

	user := router.Group("/users")
	{
		user.POST("", init.UserHdler.RegisterUser)
		user.POST("/login", init.UserHdler.LoginUser)
		user.GET("/current", middleware.AuthMiddleware(), init.UserHdler.CurrentUser)
	}

}
