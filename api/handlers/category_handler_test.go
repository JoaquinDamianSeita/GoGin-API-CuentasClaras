package handlers

import (
	"GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/dto"
	testhelpers "GoGin-API-CuentasClaras/test_helpers"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockCategoryService struct{}

func (m *MockCategoryService) Index(user dao.User) (int, []dto.TransformedIndexCategory) {
	transformedResponse := []dto.TransformedIndexCategory{}
	transformed := dto.TransformedIndexCategory{
		Id:          1,
		Name:        "Work",
		Color:       "#fdg123",
		Description: "Work",
		IsDefault:   true,
	}

	transformedResponse = append(transformedResponse, transformed)

	return http.StatusOK, transformedResponse
}

func (m *MockCategoryService) Create(user dao.User, categoryCreateRequest dto.CategoryRequest) (int, interface{}) {
	if categoryCreateRequest.Description == "Payment for work" {
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred in the creation of the category."}
	}

	return http.StatusCreated, gin.H{"message": "Category successfully created."}
}

func (m *MockCategoryService) Update(user dao.User, categoryRequest dto.CategoryRequest, categoryID int) (int, interface{}) {
	if categoryID == 2 {
		return http.StatusNotFound, gin.H{"error": "Not found."}
	}

	if categoryID == 3 {
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred in the update of the category."}
	}

	return http.StatusOK, gin.H{"message": "Category successfully updated."}
}

func (m *MockCategoryService) Delete(user dao.User, categoryID int) (int, interface{}) {
	if categoryID == 2 {
		return http.StatusNotFound, gin.H{"error": "Not found."}
	}

	if categoryID == 3 {
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred while deleting the category."}
	}

	return http.StatusOK, gin.H{"message": "Category successfully deleted."}
}

func TestCategoryHandlerImpl_Index(t *testing.T) {
	categoryService := &MockCategoryService{}
	categoryHandler := CategoryHandlerInit(categoryService)
	serviceUri := "/api/categories"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the user has categories",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "[{\"id\":1,\"name\":\"Work\",\"color\":\"#fdg123\",\"description\":\"Work\",\"is_default\":true}]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, responseRecorder := testhelpers.MockGetRequest(serviceUri)

			ctx.Set("user", dao.User{ID: 1})

			categoryHandler.Index(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestCategoryHandlerImpl_Create(t *testing.T) {
	categoryService := &MockCategoryService{}
	categoryHandler := CategoryHandlerInit(categoryService)
	serviceUri := "/api/categories"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the category is created successfully",
			Params:       `{"name": "Custom", "color": "#6495ed", "description": "Custom"}`,
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "{\"message\":\"Category successfully created.\"}",
		},
		{
			Name:         "when the category has invalid name",
			Params:       `{"name": "", "color": "#6495ed", "description": "Custom"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the category has invalid color",
			Params:       `{"name": "Custom", "color": "193zge", "description": "Custom"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when there is an error in the creation of the category",
			Params:       `{"name": "Custom", "color": "#193zge", "description": "Payment for work"}`,
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"An error occurred in the creation of the category.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, responseRecorder := testhelpers.MockPostRequest(tt.Params, serviceUri)

			ctx.Set("user", dao.User{ID: 1})

			categoryHandler.Create(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestCategoryHandlerImpl_Update(t *testing.T) {
	categoryService := &MockCategoryService{}
	categoryHandler := CategoryHandlerInit(categoryService)
	serviceUri := "/api/categories"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the category is category successfully",
			Params:       `{"name": "Custom", "color": "#6495ed", "description": "Custom"}`,
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"message\":\"Category successfully updated.\"}",
		},
		{
			Name:         "when the category is not found",
			Params:       `{"name": "Custom", "color": "#6495ed", "description": "Custom"}`,
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: "{\"error\":\"Not found.\"}",
		},
		{
			Name:         "when the category has invalid name",
			Params:       `{"name": "", "color": "#6495ed", "description": "Custom"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the category has invalid color",
			Params:       `{"name": "Custom", "color": "193zge", "description": "Custom"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when there is an error in the update of the category",
			Params:       `{"name": "Custom", "color": "#193zge", "description": "Payment for work"}`,
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"An error occurred in the update of the category.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			categoryId := 1

			if tt.Name == "when the category is not found" {
				categoryId = 2
			} else if tt.Name == "when there is an error in the update of the category" {
				categoryId = 3
			}

			ctx, responseRecorder := testhelpers.MockPutRequest(tt.Params, serviceUri)

			ctx.Set("user", dao.User{ID: 1})

			ctx.Params = []gin.Param{
				{
					Key:   "id",
					Value: strconv.Itoa(categoryId),
				},
			}

			categoryHandler.Update(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestCategoryHandlerImpl_Delete(t *testing.T) {
	categoryService := &MockCategoryService{}
	categoryHandler := CategoryHandlerInit(categoryService)
	serviceUri := "/api/categories"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the category is found",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"message\":\"Category successfully deleted.\"}",
		},
		{
			Name:         "when the category is not found",
			Params:       "",
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: "{\"error\":\"Not found.\"}",
		},
		{
			Name:         "when there is an error while deleting the category",
			Params:       "",
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"An error occurred while deleting the category.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			categoryId := 1

			if tt.Name == "when the category is not found" {
				categoryId = 2
			} else if tt.Name == "when there is an error while deleting the category" {
				categoryId = 3
			}

			ctx, responseRecorder := testhelpers.MockDeleteRequest(serviceUri)

			ctx.Set("user", dao.User{ID: 1})

			ctx.Params = []gin.Param{
				{
					Key:   "id",
					Value: strconv.Itoa(categoryId),
				},
			}

			categoryHandler.Delete(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}
