package handlers

import (
	"GoGin-API-CuentasClaras/dto"
	"GoGin-API-CuentasClaras/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var categoryCreateRequest dto.CategoryCreateRequest

type CategoryHandler interface {
	Index(c *gin.Context)
	Create(ctx *gin.Context)
}

type CategoryHandlerImpl struct {
	svc services.CategoryService
}

func (u CategoryHandlerImpl) Index(ctx *gin.Context) {
	code, response := u.svc.Index(ParseUserFromContext(ctx))
	ctx.JSON(code, response)
}

func (u CategoryHandlerImpl) Create(ctx *gin.Context) {
	validationError := ctx.ShouldBindJSON(&categoryCreateRequest)
	if validationError != nil || invalidName() || invalidColor() {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters."})
		return
	}
	code, response := u.svc.Create(ParseUserFromContext(ctx), categoryCreateRequest)
	ctx.JSON(code, response)
}

func invalidName() bool {
	return categoryCreateRequest.Name == ""
}

func invalidColor() bool {
	if len(categoryCreateRequest.Color) != 7 || categoryCreateRequest.Color[0] != '#' {
		return true
	}
	return false
}

func CategoryHandlerInit(categoryService services.CategoryService) *CategoryHandlerImpl {
	return &CategoryHandlerImpl{
		svc: categoryService,
	}
}
