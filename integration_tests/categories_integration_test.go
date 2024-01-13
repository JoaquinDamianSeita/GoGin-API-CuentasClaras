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
				"{\"id\":2,\"name\":\"Custom\",\"color\":\"#123fge\",\"description\":\"Custom\",\"is_default\":false}]",
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
