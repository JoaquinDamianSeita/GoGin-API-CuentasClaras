package services

import (
	"GoGin-API-CuentasClaras/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OperationService interface {
	Index(c *gin.Context)
	Show(c *gin.Context)
}

type OperationServiceImpl struct {
	userRepository      repository.UserRepository
	operationRepository repository.OperationRepository
	categoryRepository  repository.CategoryRepository
}

type TransformedOperation struct {
	ID       int                 `json:"id"`
	Type     string              `json:"type"`
	Amount   float64             `json:"amount"`
	Date     time.Time           `json:"date"`
	Category TransformedCategory `json:"category"`
}

type TransformedCategory struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TransformedShowOperation struct {
	ID          int                     `json:"id"`
	Type        string                  `json:"type"`
	Amount      float64                 `json:"amount"`
	Date        time.Time               `json:"date"`
	Category    TransformedShowCategory `json:"category"`
	Description string                  `json:"description"`
}

type TransformedShowCategory struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

func (u OperationServiceImpl) Index(c *gin.Context) {
	user, recordError := RetrieveCurrentUser(u.userRepository, c.GetString("user_id"))

	if recordError != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	operations, _ := u.operationRepository.FindOperationsByUser(user)

	transformedResponse := []TransformedOperation{}
	for _, operation := range operations {
		category, _ := u.categoryRepository.FindCategoryByOperation(operation)
		transformed := TransformedOperation{
			ID:     operation.ID,
			Type:   operation.Type,
			Amount: operation.Amount,
			Date:   operation.Date.In(utcLocation),
			Category: TransformedCategory{
				Name:  category.Name,
				Color: category.Color,
			},
		}
		transformedResponse = append(transformedResponse, transformed)
	}

	c.JSON(http.StatusOK, transformedResponse)
}

func (u OperationServiceImpl) Show(c *gin.Context) {
	user, recordErrorUser := RetrieveCurrentUser(u.userRepository, c.GetString("user_id"))

	if recordErrorUser != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	operationID, _ := strconv.Atoi(c.Param("id"))
	operation, recordErrorOperation := u.operationRepository.FindOperationByUserAndId(user, operationID)

	if recordErrorOperation != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found."})
		return
	}

	TransformedOperation := TransformedShowOperation{
		ID:          operation.ID,
		Type:        operation.Type,
		Amount:      operation.Amount,
		Date:        operation.Date.In(utcLocation),
		Description: operation.Description,
		Category: TransformedShowCategory{
			Name:        operation.Category.Name,
			Color:       operation.Category.Color,
			Description: operation.Category.Description,
		},
	}

	c.JSON(http.StatusOK, TransformedOperation)
}

func OperationServiceInit(userRepository repository.UserRepository, operationRepository repository.OperationRepository, categoryRepository repository.CategoryRepository) *OperationServiceImpl {
	return &OperationServiceImpl{
		userRepository:      userRepository,
		operationRepository: operationRepository,
		categoryRepository:  categoryRepository,
	}
}
