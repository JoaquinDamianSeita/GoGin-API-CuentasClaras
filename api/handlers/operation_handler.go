package handlers

import (
	"GoGin-API-CuentasClaras/dto"
	"GoGin-API-CuentasClaras/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var createOperationRequest dto.CreateOperationRequest

type OperationHandler interface {
	Index(c *gin.Context)
	Show(c *gin.Context)
	Create(c *gin.Context)
}

type OperationHandlerImpl struct {
	svc services.OperationService
}

func (u OperationHandlerImpl) Index(ctx *gin.Context) {
	code, response := u.svc.Index(ParseUserFromContext(ctx))
	ctx.JSON(code, response)
}

func (u OperationHandlerImpl) Show(ctx *gin.Context) {
	operationID, _ := strconv.Atoi(ctx.Param("id"))
	code, response := u.svc.Show(ParseUserFromContext(ctx), operationID)
	ctx.JSON(code, response)
}

func (u OperationHandlerImpl) Create(ctx *gin.Context) {
	validationError := ctx.ShouldBindJSON(&createOperationRequest)
	if validationError != nil || invalidType() || invalidAmount() || invalidDate() {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters."})
		return
	}
	code, response := u.svc.Create(ParseUserFromContext(ctx), createOperationRequest)
	ctx.JSON(code, response)
}

func invalidType() bool {
	if createOperationRequest.Type == "" {
		return true
	}
	return createOperationRequest.Type != "income" && createOperationRequest.Type != "expense"
}

func invalidAmount() bool {
	return createOperationRequest.Amount <= 0.0
}

func invalidDate() bool {
	if createOperationRequest.Date == "" {
		return true
	}
	parsedDate, err := time.Parse(time.RFC3339, createOperationRequest.Date)
	return err != nil || parsedDate.After(time.Now())
}

func OperationHandlerInit(operationService services.OperationService) *OperationHandlerImpl {
	return &OperationHandlerImpl{
		svc: operationService,
	}
}
