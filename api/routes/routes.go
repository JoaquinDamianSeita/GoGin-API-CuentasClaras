package routes

import (
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

func UserRoutes(router *gin.RouterGroup, initConfig *config.Initialization, middleware gin.HandlerFunc) {
	user := router.Group("/users")
	{
		user.POST("", initConfig.UserHdler.RegisterUser)
		user.POST("/login", initConfig.UserHdler.LoginUser)
		user.GET("/current", middleware, initConfig.UserHdler.CurrentUser)
	}
}

func OperationRoutes(router *gin.RouterGroup, initConfig *config.Initialization, middleware gin.HandlerFunc) {
	operation := router.Group("/operations")
	{
		operation.GET("", middleware, initConfig.OperationHdler.Index)
		operation.GET("/:id", middleware, initConfig.OperationHdler.Show)
		operation.POST("", middleware, initConfig.OperationHdler.Create)
		operation.PUT("/:id", middleware, initConfig.OperationHdler.Update)
		operation.DELETE("/:id", middleware, initConfig.OperationHdler.Delete)
	}
}

func CategoriesRoutes(router *gin.RouterGroup, initConfig *config.Initialization, middleware gin.HandlerFunc) {
	category := router.Group("/categories")
	{
		category.GET("", middleware, initConfig.CategoryHdler.Index)
		category.POST("", middleware, initConfig.CategoryHdler.Create)
		category.PUT("/:id", middleware, initConfig.CategoryHdler.Update)
		category.DELETE("/:id", middleware, initConfig.CategoryHdler.Delete)
	}
}
