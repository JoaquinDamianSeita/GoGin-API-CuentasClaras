package services

import (
	"GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/dto"
	"GoGin-API-CuentasClaras/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CategoryService interface {
	Index(user dao.User) (int, []dto.TransformedIndexCategory)
	Create(user dao.User, categoryCreateRequest dto.CategoryRequest) (int, interface{})
}

type CategoryServiceImpl struct {
	categoryRepository repository.CategoryRepository
}

func (u CategoryServiceImpl) Index(user dao.User) (int, []dto.TransformedIndexCategory) {
	userCategories, _ := u.categoryRepository.FindCategoriesByUser(user)
	defaultCategories, _ := u.categoryRepository.FindDefaultCategories()
	transformedResponse := FormatCategories(userCategories, defaultCategories)

	return http.StatusOK, transformedResponse
}

func (u CategoryServiceImpl) Create(user dao.User, categoryRequest dto.CategoryRequest) (int, interface{}) {
	categoryDao := dao.Category{
		Name:        categoryRequest.Name,
		Color:       categoryRequest.Color,
		Description: categoryRequest.Description,
		UserID:      uint(user.ID),
	}

	_, recordError := u.categoryRepository.Save(&categoryDao)
	if recordError != nil {
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred in the creation of the category."}
	}

	return http.StatusCreated, gin.H{"message": "Category successfully created."}
}

func (u CategoryServiceImpl) Update(user dao.User, categoryRequest dto.CategoryRequest, categoryID int) (int, interface{}) {
	invalidOperationID, category := validateCategoryID(categoryID, user, u.categoryRepository)
	if invalidOperationID {
		return http.StatusNotFound, gin.H{"error": "Not found."}
	}

	dateOperation, _ := time.Parse(time.RFC3339, operationRequest.Date)

	operationDao := dao.Category{
		ID:          category.ID,
		Type:        operationRequest.Type,
		Amount:      operationRequest.Amount,
		Date:        dateOperation,
		Category:    createCategoryOperation,
		Description: operationRequest.Description,
		UserID:      uint(user.ID),
	}

	_, recordError := u.operationRepository.Update(&operationDao)
	if recordError != nil {
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred in the update of the category."}
	}

	return http.StatusOK, gin.H{"message": "Category successfully updated."}
}

func FormatCategories(userCategories []dao.Category, defaultCategories []dao.Category) []dto.TransformedIndexCategory {
	transformedCategories := []dto.TransformedIndexCategory{}
	userCategoriesAndDefaults := append(defaultCategories, userCategories...)
	for _, category := range userCategoriesAndDefaults {
		transformed := dto.TransformedIndexCategory{
			Id:          category.ID,
			Name:        category.Name,
			Color:       category.Color,
			Description: category.Description,
			IsDefault:   category.IsDefault,
		}
		transformedCategories = append(transformedCategories, transformed)
	}
	return transformedCategories
}

func validateCategoryID(categoryID int, user dao.User, categoryRepository repository.CategoryRepository) (bool, dao.Category) {
	category, errFindOperation := categoryRepository.FindCategoryByUserAndId(user, categoryID)
	return errFindOperation != nil, category
}

func CategoryServiceInit(categoryRepository repository.CategoryRepository) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		categoryRepository: categoryRepository,
	}
}
