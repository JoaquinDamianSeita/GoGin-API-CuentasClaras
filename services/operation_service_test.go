package services

import (
	dao "GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/dto"
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

func (u MockOperationRepositoryOperations) Update(operation *dao.Operation) (dao.Operation, error) {
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

	var tests = []testhelpers.TestInterfaceStructure{
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
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			user := dao.User{ID: 1}

			if tt.Name == "when the user has no operations" {
				user = dao.User{ID: 2}
			}

			code, response := operationService.Index(user)

			testhelpers.AssertExpectedCodeAndResponseServiceDto(t, tt, code, response)
		})
	}
}

func TestOperationServiceImpl_Show(t *testing.T) {
	userRepository := &MockUserRepositoryOperations{}
	operationRepository := &MockOperationRepositoryOperations{}
	categoryRepository := &MockCategoryRepositoryOperations{}
	operationService := OperationServiceInit(userRepository, operationRepository, categoryRepository)

	var tests = []testhelpers.TestInterfaceStructure{
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
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			operation_id := 1

			if tt.Name == "when the operation is not found" {
				operation_id = 2
			}

			code, response := operationService.Show(dao.User{ID: 1}, operation_id)

			testhelpers.AssertExpectedCodeAndResponseServiceDto(t, tt, code, response)
		})
	}
}

func TestOperationServiceImpl_Create(t *testing.T) {
	userRepository := &MockUserRepositoryOperations{}
	operationRepository := &MockOperationRepositoryOperations{}
	categoryRepository := &MockCategoryRepositoryOperations{}
	operationService := OperationServiceInit(userRepository, operationRepository, categoryRepository)
	validDate := time.Now().Add(-time.Hour).Format(time.RFC3339)

	var tests = []testhelpers.TestInterfaceStructure{
		{
			Name:         "when the operation is created successfully",
			Params:       dto.OperationRequest{Type: "income", Amount: 200.50, Date: validDate, Description: "Payment for services", CategoryID: "1"},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "{\"message\":\"Operation successfully created.\"}",
		},
		{
			Name:         "when the operation has invalid category ID",
			Params:       dto.OperationRequest{Type: "income", Amount: 200.50, Date: validDate, Description: "Payment for services", CategoryID: "2"},
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"Invalid category.\"}",
		},
		{
			Name:         "when there is an error in the creation of the operation",
			Params:       dto.OperationRequest{Type: "expense", Amount: 200.50, Date: validDate, Description: "Payment for work", CategoryID: "1"},
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"An error occurred in the creation of the operation.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			code, response := operationService.Create(dao.User{ID: 1}, tt.Params.(dto.OperationRequest))

			testhelpers.AssertExpectedCodeAndResponseServiceDto(t, tt, code, response)
		})
	}
}
