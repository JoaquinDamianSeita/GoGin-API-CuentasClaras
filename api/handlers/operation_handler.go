package handlers

import (
	"GoGin-API-Base/services"
)

type OperationHandler interface {
}

type OperationHandlerImpl struct {
	svc services.OperationService
}

func OperationHandlerInit(operationService services.OperationService) *OperationHandlerImpl {
	return &OperationHandlerImpl{
		svc: operationService,
	}
}
