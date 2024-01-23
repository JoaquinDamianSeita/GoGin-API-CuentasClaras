package handlers

import (
	"GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/dto"
	testhelpers "GoGin-API-CuentasClaras/test_helpers"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type MockOperationService struct{}

func (m *MockOperationService) Index(user dao.User) (int, []dto.TransformedOperation) {
	date, _ := time.Parse(time.RFC3339, "2023-10-23T21:33:03.73297-03:00")

	transformedResponse := []dto.TransformedOperation{}
	transformed := dto.TransformedOperation{
		ID:     1,
		Type:   "income",
		Amount: 1200.5,
		Date:   date,
		Category: dto.TransformedCategory{
			Name:  "Work",
			Color: "#fdg123",
		},
	}

	transformedResponse = append(transformedResponse, transformed)

	if user.ID == 1 {
		return http.StatusOK, transformedResponse
	} else {
		return http.StatusOK, []dto.TransformedOperation{}
	}
}

func (m *MockOperationService) Show(user dao.User, operationID int) (int, interface{}) {
	date, _ := time.Parse(time.RFC3339, "2023-10-23T21:33:03.73297-03:00")

	if operationID == 1 {
		return http.StatusOK, dto.TransformedShowOperation{
			ID:     1,
			Type:   "income",
			Amount: 1200.5,
			Date:   date,
			Category: dto.TransformedShowCategory{
				Name:        "Work",
				Color:       "#fdg123",
				Description: "Work",
				IsDefault:   true,
			},
			Description: "Salario",
		}
	} else {
		return http.StatusNotFound, gin.H{"error": "Not found."}
	}
}

func (m *MockOperationService) Create(user dao.User, operationRequest dto.OperationRequest) (int, interface{}) {
	if operationRequest.CategoryID == "2" {
		return http.StatusUnprocessableEntity, gin.H{"error": "Invalid Category."}
	}

	if operationRequest.Description == "Payment for work" {
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred in the creation of the operation."}
	}

	return http.StatusCreated, gin.H{"message": "Operation successfully created."}
}

func (m *MockOperationService) Update(user dao.User, operationRequest dto.OperationRequest, operationID int) (int, interface{}) {
	if operationID == 2 {
		return http.StatusNotFound, gin.H{"error": "Not found."}
	}

	if operationRequest.CategoryID == "2" {
		return http.StatusUnprocessableEntity, gin.H{"error": "Invalid Category."}
	}

	if operationID == 3 {
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred in the update of the operation."}
	}

	return http.StatusOK, gin.H{"message": "Operation successfully updated."}
}

func (m *MockOperationService) Delete(user dao.User, operationID int) (int, interface{}) {
	if operationID == 2 {
		return http.StatusNotFound, gin.H{"error": "Not found."}
	}

	if operationID == 3 {
		return http.StatusUnprocessableEntity, gin.H{"error": "An error occurred while deleting the operation."}
	}

	return http.StatusOK, gin.H{"message": "Operation successfully deleted."}
}

func TestOperationHandlerImpl_Index(t *testing.T) {
	operationService := &MockOperationService{}
	operationHandler := OperationHandlerInit(operationService)
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
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, responseRecorder := testhelpers.MockGetRequest(serviceUri)

			if tt.Name == "when the user has operations" {
				ctx.Set("user", dao.User{ID: 1})
			} else {
				ctx.Set("user", dao.User{ID: 2})
			}

			operationHandler.Index(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestOperationHandlerImpl_Show(t *testing.T) {
	operationService := &MockOperationService{}
	operationHandler := OperationHandlerInit(operationService)
	serviceUri := "/api/operations"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the operation is found",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"id\":1,\"type\":\"income\",\"amount\":1200.5,\"date\":\"2023-10-23T21:33:03.73297-03:00\",\"category\":{\"name\":\"Work\",\"color\":\"#fdg123\",\"description\":\"Work\",\"is_default\":true},\"description\":\"Salario\"}",
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

			ctx, responseRecorder := testhelpers.MockGetRequest(serviceUri)

			ctx.Set("user", dao.User{ID: 1})

			ctx.Params = []gin.Param{
				{
					Key:   "id",
					Value: strconv.Itoa(operation_id),
				},
			}

			operationHandler.Show(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestOperationHandlerImpl_Create(t *testing.T) {
	operationService := &MockOperationService{}
	operationHandler := OperationHandlerInit(operationService)
	serviceUri := "/api/operations"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the operation is created successfully",
			Params:       `{"type": "income", "amount": 200.50, "date": "2023-11-02T23:07:00Z", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "{\"message\":\"Operation successfully created.\"}",
		},
		{
			Name:         "when the operation has invalid type",
			Params:       `{"type": "invalid", "amount": 200.50, "date": "2023-11-02T23:07:00Z", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the operation has invalid amount",
			Params:       `{"type": "income", "amount": 0.0, "date": "2023-11-02T23:07:00Z", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the operation has invalid date",
			Params:       `{"type": "income", "amount": 200.50, "date": "", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when there is an error in the creation of the operation",
			Params:       `{"type": "expense", "amount": 200.50, "date": "2023-11-02T23:07:00Z", "description": "Payment for work", "category_id": "1"}`,
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"An error occurred in the creation of the operation.\"}",
		},
		{
			Name:         "when the operation has invalid category ID",
			Params:       `{"type": "income", "amount": 200.50, "date": "2023-11-02T23:07:00Z", "description": "Payment for services", "category_id": "2"}`,
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"Invalid Category.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, responseRecorder := testhelpers.MockPostRequest(tt.Params, serviceUri)

			ctx.Set("user", dao.User{ID: 1})

			operationHandler.Create(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestOperationHandlerImpl_Update(t *testing.T) {
	operationService := &MockOperationService{}
	operationHandler := OperationHandlerInit(operationService)
	serviceUri := "/api/operations"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the operation is updated successfully",
			Params:       `{"type": "income", "amount": 200.50, "date": "2023-11-02T23:07:00Z", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"message\":\"Operation successfully updated.\"}",
		},
		{
			Name:         "when the operation is not found",
			Params:       `{"type": "income", "amount": 200.50, "date": "2023-11-02T23:07:00Z", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: "{\"error\":\"Not found.\"}",
		},
		{
			Name:         "when the operation has invalid type",
			Params:       `{"type": "invalid", "amount": 200.50, "date": "2023-11-02T23:07:00Z", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the operation has invalid amount",
			Params:       `{"type": "income", "amount": 0.0, "date": "2023-11-02T23:07:00Z", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the operation has invalid date",
			Params:       `{"type": "income", "amount": 200.50, "date": "", "description": "Payment for services", "category_id": "1"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when there is an error in the update of the operation",
			Params:       `{"type": "expense", "amount": 200.50, "date": "2023-11-02T23:07:00Z", "description": "Payment for work", "category_id": "1"}`,
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"An error occurred in the update of the operation.\"}",
		},
		{
			Name:         "when the operation has invalid category ID",
			Params:       `{"type": "income", "amount": 200.50, "date": "2023-11-02T23:07:00Z", "description": "Payment for services", "category_id": "2"}`,
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"Invalid Category.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			operation_id := 1

			if tt.Name == "when the operation is not found" {
				operation_id = 2
			} else if tt.Name == "when there is an error in the update of the operation" {
				operation_id = 3
			}

			ctx, responseRecorder := testhelpers.MockPutRequest(tt.Params, serviceUri)

			ctx.Set("user", dao.User{ID: 1})

			ctx.Params = []gin.Param{
				{
					Key:   "id",
					Value: strconv.Itoa(operation_id),
				},
			}

			operationHandler.Update(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}

func TestOperationHandlerImpl_Delete(t *testing.T) {
	operationService := &MockOperationService{}
	operationHandler := OperationHandlerInit(operationService)
	serviceUri := "/api/operations"

	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the operation is found",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"message\":\"Operation successfully deleted.\"}",
		},
		{
			Name:         "when the operation is not found",
			Params:       "",
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: "{\"error\":\"Not found.\"}",
		},
		{
			Name:         "when there is an error while deleting the operation",
			Params:       "",
			ExpectedCode: http.StatusUnprocessableEntity,
			ExpectedBody: "{\"error\":\"An error occurred while deleting the operation.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			operation_id := 1

			if tt.Name == "when the operation is not found" {
				operation_id = 2
			} else if tt.Name == "when there is an error while deleting the operation" {
				operation_id = 3
			}

			ctx, responseRecorder := testhelpers.MockDeleteRequest(serviceUri)

			ctx.Set("user", dao.User{ID: 1})

			ctx.Params = []gin.Param{
				{
					Key:   "id",
					Value: strconv.Itoa(operation_id),
				},
			}

			operationHandler.Delete(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}
