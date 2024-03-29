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
	Create(user dao.User, operationRequest dto.OperationRequest) (int, interface{})
	Update(user dao.User, operationRequest dto.OperationRequest, operationID int) (int, interface{})
	Delete(user dao.User, operationID int) (int, interface{})
}

type OperationServiceImpl struct {
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
	invalidOperationID, operation := validateOperationID(operationID, user, u.operationRepository)

	if invalidOperationID {
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
			IsDefault:   operation.Category.IsDefault,
		},
	}

	return http.StatusOK, TransformedOperation
}

func (u OperationServiceImpl) Create(user dao.User, operationRequest dto.OperationRequest) (int, interface{}) {
	if invalidCategoryID(operationRequest.CategoryID, u.categoryRepository) {
		return http.StatusUnprocessableEntity, gin.H{"error": "Invalid category."}
	}

	dateOperation, _ := time.Parse(time.RFC3339, operationRequest.Date)

	operationDao := dao.Operation{
		Type:        operationRequest.Type,
		Amount:      operationRequest.Amount,
		Date:        dateOperation,
		Category:    createCategoryOperation,
		Description: operationRequest.Description,
		UserID:      uint(user.ID),
	}

	_, recordError := u.operationRepository.Save(&operationDao)
	if recordError != nil {
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred in the creation of the operation."}
	}

	return http.StatusCreated, gin.H{"message": "Operation successfully created."}
}

func (u OperationServiceImpl) Update(user dao.User, operationRequest dto.OperationRequest, operationID int) (int, interface{}) {
	invalidOperationID, operation := validateOperationID(operationID, user, u.operationRepository)
	if invalidOperationID {
		return http.StatusNotFound, gin.H{"error": "Not found."}
	}

	if invalidCategoryID(operationRequest.CategoryID, u.categoryRepository) {
		return http.StatusUnprocessableEntity, gin.H{"error": "Invalid category."}
	}

	dateOperation, _ := time.Parse(time.RFC3339, operationRequest.Date)

	operationDao := dao.Operation{
		ID:          operation.ID,
		Type:        operationRequest.Type,
		Amount:      operationRequest.Amount,
		Date:        dateOperation,
		Category:    createCategoryOperation,
		Description: operationRequest.Description,
		UserID:      uint(user.ID),
	}

	_, recordError := u.operationRepository.Update(&operationDao)
	if recordError != nil {
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred in the update of the operation."}
	}

	return http.StatusOK, gin.H{"message": "Operation successfully updated."}
}

func (u OperationServiceImpl) Delete(user dao.User, operationID int) (int, interface{}) {
	invalidOperationID, operation := validateOperationID(operationID, user, u.operationRepository)

	if invalidOperationID {
		return http.StatusNotFound, gin.H{"error": "Not found."}
	}

	_, recordError := u.operationRepository.Delete(&operation)
	if recordError != nil {
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred while deleting the operation."}
	}

	return http.StatusOK, gin.H{"message": "Operation successfully deleted."}
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

func validateOperationID(operationID int, user dao.User, operationRepository repository.OperationRepository) (bool, dao.Operation) {
	operation, errFindOperation := operationRepository.FindOperationByUserAndId(user, operationID)
	return errFindOperation != nil, operation
}

func OperationServiceInit(operationRepository repository.OperationRepository, categoryRepository repository.CategoryRepository) *OperationServiceImpl {
	return &OperationServiceImpl{
		operationRepository: operationRepository,
		categoryRepository:  categoryRepository,
	}
}
