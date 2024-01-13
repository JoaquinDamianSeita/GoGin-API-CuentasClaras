package handlers

import (
	"GoGin-API-CuentasClaras/services"

	"github.com/gin-gonic/gin"
)

type CategoryHandler interface {
	Index(c *gin.Context)
}

type CategoryHandlerImpl struct {
	svc services.CategoryService
}

func (u CategoryHandlerImpl) Index(ctx *gin.Context) {
	code, response := u.svc.Index(ParseUserFromContext(ctx))
	ctx.JSON(code, response)
}

func CategoryHandlerInit(categoryService services.CategoryService) *CategoryHandlerImpl {
	return &CategoryHandlerImpl{
		svc: categoryService,
	}
}
