package services

import (
	dao "GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/dto"
	testhelpers "GoGin-API-CuentasClaras/test_helpers"
	"errors"
	"net/http"
	"testing"
)

type MockCategoryRepositoryCategories struct{}

func (u MockCategoryRepositoryCategories) FindCategoryByOperation(operation dao.Operation) (dao.Category, error) {
	return operation.Category, nil
}

func (u MockCategoryRepositoryCategories) Save(category *dao.Category) (dao.Category, error) {
	if category.Description == "Payment for work" {
		return dao.Category{}, errors.New("Invalid category.")
	}
	return dao.Category{}, nil
}

func (u MockCategoryRepositoryCategories) FindCategoryByUserAndId(user dao.User, categoryID int) (dao.Category, error) {
	if categoryID == 2 {
		return dao.Category{}, errors.New("Invalid category.")
	}
	return dao.Category{}, nil
}

func (u MockCategoryRepositoryCategories) Update(category *dao.Category) (dao.Category, error) {
	if category.Description == "Payment for work" {
		return dao.Category{}, errors.New("Invalid category.")
	}
	return dao.Category{}, nil
}

func (u MockCategoryRepositoryCategories) FindCategoryById(id int) (dao.Category, error) {
	if id == 1 {
		return dao.Category{}, nil
	}
	return dao.Category{}, errors.New("Category not found.")
}

func (u MockCategoryRepositoryCategories) FindCategoriesByUser(user dao.User) ([]dao.Category, error) {
	if user.ID == 2 {
		categories := []dao.Category{}
		categoriesResponse := append(categories, dao.Category{
			ID:          2,
			Name:        "Custom",
			Color:       "#6495ed",
			Description: "Custom",
			IsDefault:   false,
		})
		return categoriesResponse, nil
	}
	return []dao.Category{}, nil
}

func (u MockCategoryRepositoryCategories) FindDefaultCategories() ([]dao.Category, error) {
	categories := []dao.Category{}
	categoriesResponse := append(categories, dao.Category{
		ID:          1,
		Name:        "Work",
		Color:       "#fdg123",
		Description: "Work",
		IsDefault:   true,
	})
	return categoriesResponse, nil
}

func TestCategoryServiceImpl_Index(t *testing.T) {
	categoryRepository := &MockCategoryRepositoryCategories{}
	categoryService := CategoryServiceInit(categoryRepository)

	var tests = []testhelpers.TestInterfaceStructure{
		{
			Name:         "when the user has default categories",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "[{\"id\":1,\"name\":\"Work\",\"color\":\"#fdg123\",\"description\":\"Work\",\"is_default\":true}]",
		},
		{
			Name:         "when the user has default and custom categories",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "[{\"id\":1,\"name\":\"Work\",\"color\":\"#fdg123\",\"description\":\"Work\",\"is_default\":true}," +
				"{\"id\":2,\"name\":\"Custom\",\"color\":\"#6495ed\",\"description\":\"Custom\",\"is_default\":false}]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			user := dao.User{ID: 1}

			if tt.Name == "when the user has default and custom categories" {
				user = dao.User{ID: 2}
			}

			code, response := categoryService.Index(user)

			testhelpers.AssertExpectedCodeAndResponseServiceDto(t, tt, code, response)
		})
	}
}

func TestCategoryServiceImpl_Create(t *testing.T) {
	categoryRepository := &MockCategoryRepositoryCategories{}
	categoryService := CategoryServiceInit(categoryRepository)

	var tests = []testhelpers.TestInterfaceStructure{
		{
			Name:         "when the category is created successfully",
			Params:       dto.CategoryRequest{Name: "Custom", Color: "#6495ed", Description: "Custom"},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "{\"message\":\"Category successfully created.\"}",
		},
		{
			Name:         "when there is an error in the creation of the category",
			Params:       dto.CategoryRequest{Name: "Custom", Color: "#6495ed", Description: "Payment for work"},
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"An error occurred in the creation of the category.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			code, response := categoryService.Create(dao.User{ID: 1}, tt.Params.(dto.CategoryRequest))

			testhelpers.AssertExpectedCodeAndResponseServiceDto(t, tt, code, response)
		})
	}
}

func TestCategoryServiceImpl_Update(t *testing.T) {
	categoryRepository := &MockCategoryRepositoryCategories{}
	categoryService := CategoryServiceInit(categoryRepository)

	var tests = []testhelpers.TestInterfaceStructure{
		{
			Name:         "when the category is updated successfully",
			Params:       dto.CategoryRequest{Name: "Custom", Color: "#6495ed", Description: "Custom"},
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"message\":\"Category successfully updated.\"}",
		},
		{
			Name:         "when the category is not found",
			Params:       dto.CategoryRequest{Name: "Custom", Color: "#6495ed", Description: "Custom"},
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: "{\"error\":\"Not found.\"}",
		},
		{
			Name:         "when there is an error in the update of the category",
			Params:       dto.CategoryRequest{Name: "Custom", Color: "#6495ed", Description: "Payment for work"},
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"An error occurred in the update of the category.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			categoryId := 1

			if tt.Name == "when the category is not found" {
				categoryId = 2
			}

			if tt.Name == "when there is an error in the update of the category" {
				categoryId = 3
			}

			code, response := categoryService.Update(dao.User{ID: 1}, tt.Params.(dto.CategoryRequest), categoryId)

			testhelpers.AssertExpectedCodeAndResponseServiceDto(t, tt, code, response)
		})
	}
}
