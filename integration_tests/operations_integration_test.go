package integration_tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestOperationsIntegration_Create_ValidRequest(t *testing.T) {
	router := setupTest()
	validDate := time.Now().Add(-time.Hour).Format(time.RFC3339)
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the operation is created successfully",
			Params:       `{"type": "income", "amount": 200.50, "date": "` + validDate + `", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "{\"message\":\"Operation successfully created.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/api/operations", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+token)

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestOperationsIntegration_Create_InvalidRequest(t *testing.T) {
	router := setupTest()
	validDate := time.Now().Add(-time.Hour).Format(time.RFC3339)
	invalidDate := time.Now().Add(time.Hour).Format(time.RFC3339)
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the operation has invalid type",
			Params:       `{"type": "invalid", "amount": 200.50, "date": "` + validDate + `", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the operation has invalid amount",
			Params:       `{"type": "income", "amount": 0.0, "date": "` + validDate + `", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the operation has invalid date",
			Params:       `{"type": "income", "amount": 200.50, "date": "` + invalidDate + `", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the operation has invalid category ID",
			Params:       `{"type": "income", "amount": 200.50, "date": "` + validDate + `", "description": "Payment for services", "category_id": "2"}`,
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"Invalid category.\"}",
		},
		{
			Name:         "when the user does not exist",
			Params:       `{"type": "income", "amount": 200.50, "date": "` + validDate + `", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "{\"error\":\"Not authorized\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/api/operations", strings.NewReader(tt.Params))
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

func TestOperationsIntegration_Update_ValidRequest(t *testing.T) {
	router := setupTest()
	validDate := time.Now().Add(-time.Hour).Format(time.RFC3339)
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the operation is updated successfully",
			Params:       `{"type": "income", "amount": 200.50, "date": "` + validDate + `", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"message\":\"Operation successfully updated.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("PUT", "/api/operations/1", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+token)

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestOperationsIntegration_Update_InvalidRequest(t *testing.T) {
	router := setupTest()
	validDate := time.Now().Add(-time.Hour).Format(time.RFC3339)
	invalidDate := time.Now().Add(time.Hour).Format(time.RFC3339)
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the operation has invalid type",
			Params:       `{"type": "invalid", "amount": 200.50, "date": "` + validDate + `", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the operation is not found",
			Params:       `{"type": "income", "amount": 200.50, "date": "` + validDate + `", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: "{\"error\":\"Not found.\"}",
		},
		{
			Name:         "when the operation has invalid amount",
			Params:       `{"type": "income", "amount": 0.0, "date": "` + validDate + `", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the operation has invalid date",
			Params:       `{"type": "income", "amount": 200.50, "date": "` + invalidDate + `", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the operation has invalid category ID",
			Params:       `{"type": "income", "amount": 200.50, "date": "` + validDate + `", "description": "Payment for services", "category_id": "2"}`,
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"Invalid category.\"}",
		},
		{
			Name:         "when the user does not exist",
			Params:       `{"type": "income", "amount": 200.50, "date": "` + validDate + `", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "{\"error\":\"Not authorized\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			operation_id := "1"

			if tt.Name == "when the operation is not found" {
				operation_id = "22"
			}

			request, _ := http.NewRequest("PUT", "/api/operations/"+operation_id, strings.NewReader(tt.Params))
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
