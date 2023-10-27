package services

import (
	dao "GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OperationService interface {
	Index(c *gin.Context)
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

func (u OperationServiceImpl) Index(c *gin.Context) {
	var user dao.User
	userID, _ := strconv.Atoi(c.GetString("user_id"))
	user, recordError := u.userRepository.FindUserById(userID)

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
			Date:   operation.Date,
			Category: TransformedCategory{
				Name:  category.Name,
				Color: category.Color,
			},
		}
		transformedResponse = append(transformedResponse, transformed)
	}

	c.JSON(http.StatusOK, transformedResponse)
}

func OperationServiceInit(userRepository repository.UserRepository, operationRepository repository.OperationRepository, categoryRepository repository.CategoryRepository) *OperationServiceImpl {
	return &OperationServiceImpl{
		userRepository:      userRepository,
		operationRepository: operationRepository,
		categoryRepository:  categoryRepository,
	}
}
