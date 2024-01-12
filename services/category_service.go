package services

import (
	"GoGin-API-CuentasClaras/api/auth"
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
	auth               auth.Auth
}

func (u CategoryServiceImpl) Index(user dao.User) (int, []dto.TransformedIndexCategory) {
	userCategories, _ := u.categoryRepository.FindCategoriesByUser(user)
	defaultCategories, _ := u.categoryRepository.FindDefaultCategories()
	transformedResponse := FormatCategories(userCategories, defaultCategories)

	return http.StatusOK, transformedResponse
}

func FormatCategories(userCategories []dao.Category, defaultCategories []dao.Category) []dto.TransformedIndexCategory {
	transformedCategories := []dto.TransformedIndexCategory{}
	userCategoriesAndDefaults := append(userCategories, defaultCategories...)
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

func CategoryServiceInit(categoryRepository repository.CategoryRepository, auth auth.Auth) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		categoryRepository: categoryRepository,
		auth:               auth,
	}
}
