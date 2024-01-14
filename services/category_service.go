package services

import (
	"GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/dto"
	"GoGin-API-CuentasClaras/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryService interface {
	Index(user dao.User) (int, []dto.TransformedIndexCategory)
	Create(user dao.User, categoryCreateRequest dto.CategoryCreateRequest) (int, interface{})
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

func (u CategoryServiceImpl) Create(user dao.User, categoryCreateRequest dto.CategoryCreateRequest) (int, interface{}) {
	categoryDao := dao.Category{
		Name:        categoryCreateRequest.Name,
		Color:       categoryCreateRequest.Color,
		Description: categoryCreateRequest.Description,
		UserID:      uint(user.ID),
	}

	_, recordError := u.categoryRepository.Save(&categoryDao)
	if recordError != nil {
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred in the creation of the category."}
	}

	return http.StatusCreated, gin.H{"message": "Category successfully created."}
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

func CategoryServiceInit(categoryRepository repository.CategoryRepository) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		categoryRepository: categoryRepository,
	}
}
