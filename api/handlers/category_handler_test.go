package handlers

import (
	"GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/dto"
	testhelpers "GoGin-API-CuentasClaras/test_helpers"
	"net/http"
	"testing"
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
