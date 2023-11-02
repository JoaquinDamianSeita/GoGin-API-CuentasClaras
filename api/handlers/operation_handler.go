package handlers

import (
	"GoGin-API-CuentasClaras/services"

	"github.com/gin-gonic/gin"
)

type OperationHandler interface {
	Index(c *gin.Context)
	Show(c *gin.Context)
	Create(c *gin.Context)
}

type OperationHandlerImpl struct {
	svc services.OperationService
}

func (u OperationHandlerImpl) Index(c *gin.Context) {
	u.svc.Index(c)
}

func (u OperationHandlerImpl) Show(c *gin.Context) {
	u.svc.Show(c)
}

func (u OperationHandlerImpl) Create(c *gin.Context) {
	u.svc.Create(c)
}

func OperationHandlerInit(operationService services.OperationService) *OperationHandlerImpl {
	return &OperationHandlerImpl{
		svc: operationService,
	}
}
