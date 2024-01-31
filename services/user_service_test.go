package services

import (
	dao "GoGin-API-CuentasClaras/dao"
	dto "GoGin-API-CuentasClaras/dto"
	testhelpers "GoGin-API-CuentasClaras/test_helpers"
	"errors"
	"net/http"
	"testing"
	"time"
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
			Password: "$2a$12$tDdX/jDY.JEoFMfk6bbuROMkJnxvDFV7VuQIqT88GaJI.auLEp.iq",
		}, nil
	}

	return dao.User{}, nil
}

type MockAuth struct{}

func (auth *MockAuth) GenerateJWT(userId string) (expiresIn int64, tokenString string, err error) {
	return 3600, "token", nil
}

func (auth *MockAuth) ValidateToken(signedToken string) (claims *dto.JWTClaim, err error) {
	if signedToken == "valid_token" {
		claims = &dto.JWTClaim{}
		return claims, nil
	}
	return nil, errors.New("Invalid token")
}

type MockOperationRepositoryUser struct{}

func (u MockOperationRepositoryUser) FindOperationsByUser(user dao.User) ([]dao.Operation, error) {
	date, _ := time.Parse(time.RFC3339, "2023-10-23T21:33:03.73297-03:00")

	if user.ID == 1 {
		operations := []dao.Operation{}
		user.Operations = append(operations, dao.Operation{
			ID:     1,
			Type:   "income",
			Amount: 100.5,
			Date:   date,
			Category: dao.Category{
				Name:  "Work",
				Color: "#fdg123",
			},
		})
	} else if user.ID == 2 {
		user.Operations = []dao.Operation{}
	}
	return user.Operations, nil
}

func (u MockOperationRepositoryUser) FindOperationByUserAndId(user dao.User, operationID int) (dao.Operation, error) {
	return dao.Operation{}, nil
}

func (u MockOperationRepositoryUser) Save(operation *dao.Operation) (dao.Operation, error) {
	return dao.Operation{}, nil
}

func (u MockOperationRepositoryUser) Update(operation *dao.Operation) (dao.Operation, error) {
	return dao.Operation{}, nil
}

func (u MockOperationRepositoryUser) Delete(operation *dao.Operation) (dao.Operation, error) {
	return dao.Operation{}, nil
}

func TestUserServiceImpl_RegisterUser(t *testing.T) {
	userRepository := &MockUserRepository{}
	auth := &MockAuth{}
	operationRepository := &MockOperationRepositoryUser{}
	userService := UserServiceInit(userRepository, auth, operationRepository)
	serviceUri := "/api/users"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the request is successful",
			Params:       `{"username": "test.user", "email": "test@example.com", "password": "password123"}`,
			ExpectedCode: http.StatusOK,
			ExpectedBody: `{"message":"User successfully created."}`,
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
			ctx, _ := testhelpers.MockPostRequest(tt.Params, serviceUri)
			var registerUserRequest dto.RegisterUserRequest
			ctx.ShouldBindJSON(&registerUserRequest)

			code, response := userService.RegisterUser(registerUserRequest)

			testhelpers.AssertExpectedCodeAndResponseService(t, tt, code, response)
		})
	}
}

func TestUserServiceImpl_LoginUser(t *testing.T) {
	userRepository := &MockUserRepository{}
	auth := &MockAuth{}
	operationRepository := &MockOperationRepositoryUser{}
	userService := UserServiceInit(userRepository, auth, operationRepository)
	serviceUri := "/api/users/login"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the request is successful",
			Params:       `{"email": "test.user@example.com", "password": "password123"}`,
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"expires_in\":3600,\"token\":\"token\"}",
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
			ctx, _ := testhelpers.MockPostRequest(tt.Params, serviceUri)
			var loginUserRequest dto.LoginRequest
			ctx.ShouldBindJSON(&loginUserRequest)

			code, response := userService.LoginUser(loginUserRequest)

			testhelpers.AssertExpectedCodeAndResponseService(t, tt, code, response)
		})
	}
}

func TestUserServiceImpl_CurrentUser(t *testing.T) {
	userRepository := &MockUserRepository{}
	auth := &MockAuth{}
	operationRepository := &MockOperationRepositoryUser{}
	userService := UserServiceInit(userRepository, auth, operationRepository)

	var tests = []testhelpers.TestInterfaceStructure{
		{
			Name:         "when the request is successful",
			Params:       dao.User{ID: 1, Username: "test.user", Email: "test.user@example.com"},
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"email\":\"test.user@example.com\",\"username\":\"test.user\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			code, response := userService.CurrentUser(tt.Params.(dao.User))

			testhelpers.AssertExpectedCodeAndResponseServiceDto(t, tt, code, response)
		})
	}
}

func TestUserServiceImpl_BalanceUser(t *testing.T) {
	userRepository := &MockUserRepository{}
	auth := &MockAuth{}
	operationRepository := &MockOperationRepositoryUser{}
	userService := UserServiceInit(userRepository, auth, operationRepository)

	var tests = []testhelpers.TestInterfaceStructure{
		{
			Name:         "when the request is successful",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"total_balance\":\"100.50\"}",
		},
		{
			Name:         "when the user has no registered operations",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"total_balance\":\"0.00\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			user := dao.User{ID: 1}

			if tt.Name == "when the user has no registered operations" {
				user = dao.User{ID: 2}
			}

			code, response := userService.BalanceUser(user)

			testhelpers.AssertExpectedCodeAndResponseServiceDto(t, tt, code, response)
		})
	}
}
