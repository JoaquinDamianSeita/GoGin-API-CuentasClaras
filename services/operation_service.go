package services

import (
	"GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/repository"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OperationService interface {
	Index(c *gin.Context)
	Show(c *gin.Context)
	Create(c *gin.Context)
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

type CreateOperationRequest struct {
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id"`
}

var createOperationRequest CreateOperationRequest
var createCategoryOperation dao.Category

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

func (u OperationServiceImpl) Create(c *gin.Context) {
	user, recordErrorUser := RetrieveCurrentUser(u.userRepository, c.GetString("user_id"))

	if recordErrorUser != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	validationError := c.ShouldBindJSON(&createOperationRequest)
	if validationError != nil || invalidType() || invalidAmount() || invalidDate() || invalidCategoryID(u.categoryRepository) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters."})
		return
	}

	dateOperation, _ := time.Parse(time.RFC3339, createOperationRequest.Date)

	operationDao := dao.Operation{
		Type:        createOperationRequest.Type,
		Amount:      createOperationRequest.Amount,
		Date:        dateOperation,
		Category:    createCategoryOperation,
		Description: createOperationRequest.Description,
		UserID:      uint(user.ID),
	}

	_, recordError := u.operationRepository.Save(&operationDao)
	if recordError != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity, gin.H{"error": "An error occurred in the creation of the operation."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Operation successfully created."})
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
	return err != nil || parsedDate.Before(time.Now())
}

func invalidCategoryID(categoryRepository repository.CategoryRepository) bool {
	if createOperationRequest.CategoryID == "" {
		return true
	}
	categoryIdInt, errParseInt := strconv.Atoi(createOperationRequest.CategoryID)
	if errParseInt != nil {
		return true
	}
	category, errFindCategory := categoryRepository.FindCategoryById(categoryIdInt)
	log.Println(errFindCategory)
	createCategoryOperation = category
	return errFindCategory != nil
}

func OperationServiceInit(userRepository repository.UserRepository, operationRepository repository.OperationRepository, categoryRepository repository.CategoryRepository) *OperationServiceImpl {
	return &OperationServiceImpl{
		userRepository:      userRepository,
		operationRepository: operationRepository,
		categoryRepository:  categoryRepository,
	}
}
