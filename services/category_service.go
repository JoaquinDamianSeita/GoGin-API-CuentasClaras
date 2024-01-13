package services

import (
	"GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/dto"
	"GoGin-API-CuentasClaras/repository"
	"net/http"
)

type CategoryService interface {
	Index(user dao.User) (int, []dto.TransformedIndexCategory)
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
