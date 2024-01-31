package handlers

import (
	"GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/dto"
	testhelpers "GoGin-API-CuentasClaras/test_helpers"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockUserService struct{}

func (m *MockUserService) RegisterUser(registerUserRequest dto.RegisterUserRequest) (int, map[string]any) {
	if registerUserRequest.Username == "" || registerUserRequest.Email == "" || registerUserRequest.Password == "" {
		return http.StatusBadRequest, gin.H{"error": "Invalid parameters."}
	}

	if registerUserRequest.Username == "invalid.user" || registerUserRequest.Email == "invalid.user@example.com" {
		return http.StatusBadRequest, gin.H{"error": "the email or the user is already in use"}
	}

	return http.StatusOK, gin.H{"message": "User successfully created."}
}

func (m *MockUserService) LoginUser(loginUserRequest dto.LoginRequest) (int, map[string]any) {
	if loginUserRequest.Email == "test.user@example.com" && loginUserRequest.Password == "password123" {
		return http.StatusOK, gin.H{"token": "token", "expires_in": "3600"}
	}

	return http.StatusUnauthorized, gin.H{"error": "invalid credentials"}
}

func (m *MockUserService) CurrentUser(user dao.User) (int, map[string]any) {
	if user.ID != 1 {
		return http.StatusUnauthorized, gin.H{"error": "Not authorized"}
	}

	return http.StatusOK, gin.H{"username": "test.user", "email": "test@example.com"}
}

func (m *MockUserService) BalanceUser(user dao.User) (int, interface{}) {
	if user.ID != 1 {
		return http.StatusUnauthorized, gin.H{"error": "Not authorized"}
	}

	return http.StatusOK, gin.H{"total_balance": "100.50"}
}

func TestUserHandlerImpl_RegisterUser(t *testing.T) {
	userService := &MockUserService{}
	userHandler := UserHandlerInit(userService)
	serviceUri := "/api/users"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the request is successful",
			Params:       `{"username": "test.user", "email": "test@example.com", "password": "password123"}`,
			ExpectedCode: http.StatusOK,
			ExpectedBody: `{"message":"User successfully created."}`,
		},
		{
			Name:         "when email is not present",
			Params:       `{"username": "test.user", "password": "password123"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when email is empty",
			Params:       `{"email": "", "username": "test.user", "password": "password123"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when password is not present",
			Params:       `{"email": "test.user@example.com", "username": "user.test"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when password is empty",
			Params:       `{"email": "test.user@example.com", "username": "test.user", "password": ""}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when username is not present",
			Params:       `{"email": "test.user@example.com", "password": "password123"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when username is empty",
			Params:       `{"email": "test.user@example.com", "username": "", "password": "password123"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when email is not available",
			Params:       `{"email": "invalid.user@example.com", "username": "test.user", "password": "password123"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"the email or the user is already in use\"}",
		},
		{
			Name:         "when username is not available",
			Params:       `{"email": "test.user@example.com", "username": "invalid.user", "password": "password123"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"the email or the user is already in use\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, responseRecorder := testhelpers.MockPostRequest(tt.Params, serviceUri)

			userHandler.RegisterUser(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestUserHandlerImpl_LoginUser(t *testing.T) {
	userService := &MockUserService{}
	userHandler := UserHandlerInit(userService)
	serviceUri := "/api/users/login"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the request is successful",
			Params:       `{"email": "test.user@example.com", "password": "password123"}`,
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"expires_in\":\"3600\",\"token\":\"token\"}",
		},
		{
			Name:         "when email is not present",
			Params:       `{"password": "password123"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when email is empty",
			Params:       `{"email": "", "password": "password123"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when password is not present",
			Params:       `{"email": "test.user@example.com"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when password is empty",
			Params:       `{"email": "test.user@example.com", "password": ""}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "with invalid email",
			Params:       `{"email": "invalid.user@example.com", "password": "password123"}`,
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "{\"error\":\"invalid credentials\"}",
		},
		{
			Name:         "with invalid password",
			Params:       `{"email": "test.user@example.com", "password": "invalidpassword"}`,
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "{\"error\":\"invalid credentials\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, responseRecorder := testhelpers.MockPostRequest(tt.Params, serviceUri)

			userHandler.LoginUser(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestUserHandlerImpl_CurrentUser(t *testing.T) {
	userService := &MockUserService{}
	userHandler := UserHandlerInit(userService)
	serviceUri := "/api/users/current"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the request is successful",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"email\":\"test@example.com\",\"username\":\"test.user\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, responseRecorder := testhelpers.MockGetRequest(serviceUri)
			ctx.Set("user", dao.User{ID: 1})

			userHandler.CurrentUser(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestUserHandlerImpl_BalanceUser(t *testing.T) {
	userService := &MockUserService{}
	userHandler := UserHandlerInit(userService)
	serviceUri := "/api/users/balance"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the request is successful",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"total_balance\":\"100.50\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, responseRecorder := testhelpers.MockGetRequest(serviceUri)
			ctx.Set("user", dao.User{ID: 1})

			userHandler.BalanceUser(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}
