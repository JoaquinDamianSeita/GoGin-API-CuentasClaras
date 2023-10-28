package integration_tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	testhelpers "GoGin-API-CuentasClaras/test_helpers"
)

func TestOperationsIntegration_Index_ValidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the user has operations",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "[{\"id\":1,\"type\":\"income\",\"amount\":1200.5,\"date\":\"2023-10-23T21:33:03.73297Z\",\"category\":{\"name\":\"Work\",\"color\":\"#fdg123\"}}]",
		},
		{
			Name:         "when the user has no operations",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "[]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/api/operations", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")

			if tt.Name == "when the user has operations" {
				request.Header.Set("Authorization", "Bearer "+token)
			} else {
				_, anotherToken, _ := authService.GenerateJWT("2")
				request.Header.Set("Authorization", "Bearer "+anotherToken)
			}

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestOperationsIntegration_Index_InvalidRequest(t *testing.T) {
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
			request, _ := http.NewRequest("GET", "/api/operations", strings.NewReader(tt.Params))
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

func TestOperationsIntegration_Show_ValidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the operation is found",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"id\":1,\"type\":\"income\",\"amount\":1200.5,\"date\":\"2023-10-23T21:33:03.73297Z\",\"category\":{\"name\":\"Work\",\"color\":\"#fdg123\",\"description\":\"Work\"},\"description\":\"Salario\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/api/operations/1", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+token)

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestOperationsIntegration_Show_InvalidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the operation is not found",
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
			request, _ := http.NewRequest("GET", "/api/operations/2", strings.NewReader(tt.Params))
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
