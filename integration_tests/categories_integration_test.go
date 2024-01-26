package integration_tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	testhelpers "GoGin-API-CuentasClaras/test_helpers"
)

func TestCategoriesIntegration_Index_ValidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
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
			request, _ := http.NewRequest("GET", "/api/categories", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")

			if tt.Name == "when the user has default categories" {
				request.Header.Set("Authorization", "Bearer "+token)
			} else {
				request.Header.Set("Authorization", "Bearer "+anotherToken)
			}

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestCategoriesIntegration_Index_InvalidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the user does not exist",
			Params:       "",
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "{\"error\":\"Not authorized\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/api/categories", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")
			_, anotherToken, _ := authService.GenerateJWT("3")
			request.Header.Set("Authorization", "Bearer "+anotherToken)

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestCategoriesIntegration_Create_ValidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the category is created successfully",
			Params:       `{"name": "Custom", "color": "#6495ed", "description": "Custom"}`,
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "{\"message\":\"Category successfully created.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/api/categories", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+token)

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestCategoriesIntegration_Create_InvalidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
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
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/api/categories", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+token)

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestCategoriesIntegration_Update_ValidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the category is update successfully",
			Params:       `{"name": "Custom", "color": "#6495ed", "description": "Custom"}`,
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"message\":\"Category successfully updated.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("PUT", "/api/categories/2", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+anotherToken)

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestCategoriesIntegration_Update_InvalidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the category is not found",
			Params:       `{"name": "Custom", "color": "#193zge", "description": "Custom"}`,
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
			Name:         "when the user does not exist",
			Params:       `{"name": "Custom", "color": "#193zge", "description": "Custom"}`,
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "{\"error\":\"Not authorized\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			categoryId := "1"
			request, _ := http.NewRequest("PUT", "/api/categories/"+categoryId, strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+token)

			if tt.Name == "when the user does not exist" {
				_, anotherToken, _ := authService.GenerateJWT("3")
				request.Header.Set("Authorization", "Bearer "+anotherToken)
			}

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestCategoriesIntegration_Delete_ValidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the category is deleted successfully",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"message\":\"Category successfully deleted.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("DELETE", "/api/categories/2", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+anotherToken)

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestCategoriesIntegration_Delete_InvalidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the category is not found",
			Params:       "",
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: "{\"error\":\"Not found.\"}",
		},
		{
			Name:         "when the user does not exist",
			Params:       "",
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "{\"error\":\"Not authorized\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			categoryId := "1"
			request, _ := http.NewRequest("DELETE", "/api/categories/"+categoryId, strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+token)

			if tt.Name == "when the user does not exist" {
				_, anotherToken, _ := authService.GenerateJWT("3")
				request.Header.Set("Authorization", "Bearer "+anotherToken)
			}

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}
