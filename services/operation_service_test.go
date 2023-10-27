package services

import (
	dao "GoGin-API-CuentasClaras/dao"
	testhelpers "GoGin-API-CuentasClaras/test_helpers"
	"errors"
	"net/http"
	"testing"
	"time"
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

func (u MockOperationRepositoryOperations) Save(operation *dao.Operation) (dao.Operation, error) {
	return dao.Operation{}, nil
}

type MockCategoryRepositoryOperations struct{}

func (u MockCategoryRepositoryOperations) FindCategoryByOperation(operation dao.Operation) (dao.Category, error) {
	return operation.Category, nil
}

func (u MockCategoryRepositoryOperations) Save(category *dao.Category) (dao.Category, error) {
	return dao.Category{}, nil
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
			ExpectedBody: "[{\"id\":1,\"type\":\"income\",\"amount\":1200.5,\"date\":\"2023-10-23T21:33:03.73297-03:00\",\"category\":{\"name\":\"Work\",\"color\":\"#fdg123\"}}]",
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
