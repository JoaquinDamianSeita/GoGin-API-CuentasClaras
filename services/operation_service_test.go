package services

import (
	dao "GoGin-API-CuentasClaras/dao"
	testhelpers "GoGin-API-CuentasClaras/test_helpers"
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type MockUserRepositoryOperations struct{}

func (m *MockUserRepositoryOperations) Save(user *dao.User) (dao.User, error) {
	return dao.User{}, nil
}

func (m *MockUserRepositoryOperations) FindUserById(id int) (dao.User, error) {
	if id == 1 || id == 2 {
		return dao.User{
			ID:       id,
			Username: "test.user",
			Email:    "test.user@example.com",
		}, nil
	} else {
		return dao.User{}, errors.New("User not found.")
	}
}

func (m *MockUserRepositoryOperations) FindUserByEmail(email string) (dao.User, error) {
	return dao.User{}, nil
}

type MockOperationRepositoryOperations struct{}

func (u MockOperationRepositoryOperations) FindOperationsByUser(user dao.User) ([]dao.Operation, error) {
	date, _ := time.Parse(time.RFC3339, "2023-10-23T21:33:03.73297-03:00")

	if user.ID == 1 {
		operations := []dao.Operation{}
		user.Operations = append(operations, dao.Operation{
			ID:     1,
			Type:   "income",
			Amount: 1200.5,
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

func (u MockOperationRepositoryOperations) FindOperationByUserAndId(user dao.User, operationID int) (dao.Operation, error) {
	date, _ := time.Parse(time.RFC3339, "2023-10-23T21:33:03.73297-03:00")

	if operationID == 1 {
		return dao.Operation{
			ID:     1,
			Type:   "income",
			Amount: 1200.5,
			Date:   date,
			Category: dao.Category{
				Name:        "Work",
				Color:       "#fdg123",
				Description: "Work",
			},
			Description: "Salario",
		}, nil
	} else {
		return dao.Operation{}, errors.New("Operation not found.")
	}
}

func (u MockOperationRepositoryOperations) Save(operation *dao.Operation) (dao.Operation, error) {
	if operation.Description == "Payment for work" {
		return dao.Operation{}, errors.New("Invalid operation.")
	}
	return dao.Operation{}, nil
}

type MockCategoryRepositoryOperations struct{}

func (u MockCategoryRepositoryOperations) FindCategoryByOperation(operation dao.Operation) (dao.Category, error) {
	return operation.Category, nil
}

func (u MockCategoryRepositoryOperations) Save(category *dao.Category) (dao.Category, error) {
	return dao.Category{}, nil
}

func (u MockCategoryRepositoryOperations) FindCategoryById(id int) (dao.Category, error) {
	if id == 1 {
		return dao.Category{}, nil
	}
	return dao.Category{}, errors.New("Category not found.")
}

func TestOperationServiceImpl_Index(t *testing.T) {
	userRepository := &MockUserRepositoryOperations{}
	operationRepository := &MockOperationRepositoryOperations{}
	categoryRepository := &MockCategoryRepositoryOperations{}
	operationService := OperationServiceInit(userRepository, operationRepository, categoryRepository)
	serviceUri := "/api/operations"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the user has operations",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "[{\"id\":1,\"type\":\"income\",\"amount\":1200.5,\"date\":\"2023-10-24T00:33:03.73297Z\",\"category\":{\"name\":\"Work\",\"color\":\"#fdg123\"}}]",
		},
		{
			Name:         "when the user has no operations",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "[]",
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
			ctx, responseRecorder := testhelpers.MockGetRequest(serviceUri)

			if tt.Name == "when the user has operations" {
				ctx.Set("user_id", "1")
			} else if tt.Name == "when the user has no operations" {
				ctx.Set("user_id", "2")
			} else {
				ctx.Set("user_id", "3")
			}

			operationService.Index(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestOperationServiceImpl_Show(t *testing.T) {
	userRepository := &MockUserRepositoryOperations{}
	operationRepository := &MockOperationRepositoryOperations{}
	categoryRepository := &MockCategoryRepositoryOperations{}
	operationService := OperationServiceInit(userRepository, operationRepository, categoryRepository)
	serviceUri := "/api/operations"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the operation is found",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"id\":1,\"type\":\"income\",\"amount\":1200.5,\"date\":\"2023-10-24T00:33:03.73297Z\",\"category\":{\"name\":\"Work\",\"color\":\"#fdg123\",\"description\":\"Work\"},\"description\":\"Salario\"}",
		},
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
			user_id := "1"
			operation_id := 1

			if tt.Name == "when the user does not exist" {
				user_id = "3"
			} else if tt.Name == "when the operation is not found" {
				operation_id = 2
			}

			ctx, responseRecorder := testhelpers.MockGetRequest(serviceUri)

			ctx.Set("user_id", user_id)

			ctx.Params = []gin.Param{
				{
					Key:   "id",
					Value: strconv.Itoa(operation_id),
				},
			}

			operationService.Show(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestOperationServiceImpl_Create(t *testing.T) {
	userRepository := &MockUserRepositoryOperations{}
	operationRepository := &MockOperationRepositoryOperations{}
	categoryRepository := &MockCategoryRepositoryOperations{}
	operationService := OperationServiceInit(userRepository, operationRepository, categoryRepository)
	serviceUri := "/api/operations"
	invalidDate := time.Now().Add(time.Hour).Format(time.RFC3339)
	validDate := time.Now().Add(-time.Hour).Format(time.RFC3339)

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the operation is created successfully",
			Params:       `{"type": "income", "amount": 200.50, "date": "` + validDate + `", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "{\"message\":\"Operation successfully created.\"}",
		},
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
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when there is an error in the creation of the operation",
			Params:       `{"type": "expense", "amount": 200.50, "date": "` + validDate + `", "description": "Payment for work", "category_id": "1"}`,
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"An error occurred in the creation of the operation.\"}",
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
			user_id := "1"

			if tt.Name == "when the user does not exist" {
				user_id = "3"
			}

			ctx, responseRecorder := testhelpers.MockPostRequest(tt.Params, serviceUri)

			ctx.Set("user_id", user_id)

			operationService.Create(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}
