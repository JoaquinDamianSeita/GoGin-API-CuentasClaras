package handlers

import (
	"GoGin-API-CuentasClaras/dto"
	"GoGin-API-CuentasClaras/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
	CurrentUser(c *gin.Context)
}

type UserHandlerImpl struct {
	svc services.UserService
}

func (u UserHandlerImpl) RegisterUser(ctx *gin.Context) {
	var registerUserRequest dto.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&registerUserRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters."})
		return
	}
	code, response := u.svc.RegisterUser(registerUserRequest)
	ctx.JSON(code, response)
}

func (u UserHandlerImpl) LoginUser(c *gin.Context) {
	u.svc.LoginUser(c)
}

func (u UserHandlerImpl) CurrentUser(c *gin.Context) {
	u.svc.CurrentUser(c)
}

func UserHandlerInit(userService services.UserService) *UserHandlerImpl {
	return &UserHandlerImpl{
		svc: userService,
	}
}
