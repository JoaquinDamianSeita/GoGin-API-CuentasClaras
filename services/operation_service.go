package services

import (
	"GoGin-API-CuentasClaras/dao"
	dto "GoGin-API-CuentasClaras/dto"
	"GoGin-API-CuentasClaras/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OperationService interface {
	Index(user dao.User) (int, []dto.TransformedOperation)
	Show(user dao.User, operationID int) (int, interface{})
	Create(user dao.User, createOperationRequest dto.CreateOperationRequest) (int, map[string]any)
}

type OperationServiceImpl struct {
	userRepository      repository.UserRepository
	operationRepository repository.OperationRepository
	categoryRepository  repository.CategoryRepository
}

var createCategoryOperation dao.Category

func (u OperationServiceImpl) Index(user dao.User) (int, []dto.TransformedOperation) {
	operations, _ := u.operationRepository.FindOperationsByUser(user)
	transformedResponse := []dto.TransformedOperation{}
	for _, operation := range operations {
		category, _ := u.categoryRepository.FindCategoryByOperation(operation)
		transformed := dto.TransformedOperation{
			ID:     operation.ID,
			Type:   operation.Type,
			Amount: operation.Amount,
			Date:   operation.Date.In(utcLocation),
			Category: dto.TransformedCategory{
				Name:  category.Name,
				Color: category.Color,
			},
		}
		transformedResponse = append(transformedResponse, transformed)
	}

	return http.StatusOK, transformedResponse
}

func (u OperationServiceImpl) Show(user dao.User, operationID int) (int, interface{}) {
	operation, recordErrorOperation := u.operationRepository.FindOperationByUserAndId(user, operationID)

	if recordErrorOperation != nil {
		return http.StatusNotFound, gin.H{"error": "Not found."}
	}

	TransformedOperation := dto.TransformedShowOperation{
		ID:          operation.ID,
		Type:        operation.Type,
		Amount:      operation.Amount,
		Date:        operation.Date.In(utcLocation),
		Description: operation.Description,
		Category: dto.TransformedShowCategory{
			Name:        operation.Category.Name,
			Color:       operation.Category.Color,
			Description: operation.Category.Description,
		},
	}

	return http.StatusOK, TransformedOperation
}

func (u OperationServiceImpl) Create(user dao.User, createOperationRequest dto.CreateOperationRequest) (int, map[string]any) {
	if invalidCategoryID(createOperationRequest.CategoryID, u.categoryRepository) {
		return http.StatusUnprocessableEntity, gin.H{"error": "Invalid category."}
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
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred in the creation of the operation."}
	}

	return http.StatusCreated, gin.H{"message": "Operation successfully created."}
}

func invalidCategoryID(categoryID string, categoryRepository repository.CategoryRepository) bool {
	if categoryID == "" {
		return true
	}
	categoryIdInt, errParseInt := strconv.Atoi(categoryID)
	if errParseInt != nil {
		return true
	}
	category, errFindCategory := categoryRepository.FindCategoryById(categoryIdInt)
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
