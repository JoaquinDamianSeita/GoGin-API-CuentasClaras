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
	BalanceUser(ctx *gin.Context)
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

func (u UserHandlerImpl) LoginUser(ctx *gin.Context) {
	var loginUserRequest dto.LoginRequest
	validationError := ctx.ShouldBindJSON(&loginUserRequest)
	if validationError != nil || loginUserRequest.Email == "" || loginUserRequest.Password == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters."})
		return
	}
	code, response := u.svc.LoginUser(loginUserRequest)
	ctx.JSON(code, response)
}

func (u UserHandlerImpl) CurrentUser(ctx *gin.Context) {
	code, response := u.svc.CurrentUser(ParseUserFromContext(ctx))
	ctx.JSON(code, response)
}

func (u UserHandlerImpl) BalanceUser(ctx *gin.Context) {
	code, response := u.svc.BalanceUser(ParseUserFromContext(ctx))
	ctx.JSON(code, response)
}

func UserHandlerInit(userService services.UserService) *UserHandlerImpl {
	return &UserHandlerImpl{
		svc: userService,
	}
}
