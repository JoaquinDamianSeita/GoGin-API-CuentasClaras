package handlers

import (
	"GoGin-API-CuentasClaras/services"

	"github.com/gin-gonic/gin"
)

type OperationHandler interface {
	Index(c *gin.Context)
}

type OperationHandlerImpl struct {
	svc services.OperationService
}

func (u OperationHandlerImpl) Index(c *gin.Context) {
	u.svc.Index(c)
}

func OperationHandlerInit(operationService services.OperationService) *OperationHandlerImpl {
	return &OperationHandlerImpl{
		svc: operationService,
	}
}
