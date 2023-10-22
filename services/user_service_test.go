package services

import (
	auth "GoGin-API-Base/api/auth"
	dao "GoGin-API-Base/dao"
	testhelpers "GoGin-API-Base/test_helpers"
	"errors"
	"net/http"
	"testing"
)

type MockUserRepository struct{}

func (m *MockUserRepository) Save(user *dao.User) (dao.User, error) {
	if user.Email == "invalid.user@example.com" || user.Username == "invalid.user" {
		return dao.User{}, errors.New("the email or the user is already in use")
	}

	return dao.User{}, nil
}

func (m *MockUserRepository) FindUserById(id int) (dao.User, error) {
	if id == 1 {
		return dao.User{
			ID:       1,
			Username: "test.user",
			Email:    "test.user@example.com",
		}, nil
	} else {
		return dao.User{}, errors.New("User not found.")
	}
}

func (m *MockUserRepository) FindUserByEmail(email string) (dao.User, error) {
	if email == "invalid.user@example.com" {
		return dao.User{}, errors.New("User not found.")
	}

	if email == "test.user@example.com" {
		return dao.User{
			Email:    "test.user@example.com",
			Password: "$2a$12$tDdX/jDY.JEoFMfk6bbuROMkJnxvDFV7VuQIqT88GaJI.auLEp.iq", // Contraseña válida en formato bcrypt
		}, nil
	}

	return dao.User{}, nil
}

type MockAuth struct{}

func (auth *MockAuth) GenerateJWT(userId string) (expiresIn int64, tokenString string, err error) {
	return 3600, "token", nil
}

func (auth *MockAuth) ValidateToken(signedToken string) (claims *auth.JWTClaim, err error) {
	return nil, nil
}

func TestUserServiceImpl_RegisterUser(t *testing.T) {
	userRepository := &MockUserRepository{}
	auth := &MockAuth{}
	userService := UserServiceInit(userRepository, auth)
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

			userService.RegisterUser(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestUserHandlerImpl_LoginUser(t *testing.T) {
	userRepository := &MockUserRepository{}
	auth := &MockAuth{}
	userService := UserServiceInit(userRepository, auth)
	serviceUri := "/api/users/login"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the request is successful",
			Params:       `{"email": "test.user@example.com", "password": "password123"}`,
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"expires_in\":3600,\"token\":\"token\"}",
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

			userService.LoginUser(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestUserHandlerImpl_CurrentUser(t *testing.T) {
	userRepository := &MockUserRepository{}
	auth := &MockAuth{}
	userService := UserServiceInit(userRepository, auth)
	serviceUri := "/api/users/current"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the request is successful",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"email\":\"test.user@example.com\",\"username\":\"test.user\"}",
		},
		{
			Name:         "when user does not exists",
			Params:       "",
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "{\"error\":\"Not authorized\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, responseRecorder := testhelpers.MockGetRequest(serviceUri)

			if tt.Name == "when the request is successful" {
				ctx.Set("user_id", "1")
			} else {
				ctx.Set("user_id", "2")
			}

			userService.CurrentUser(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}
