package handlers

import (
	"GoGin-API-CuentasClaras/services"
	testhelpers "GoGin-API-CuentasClaras/test_helpers"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type MockOperationService struct{}

func (m *MockOperationService) Index(c *gin.Context) {
	userID, _ := strconv.Atoi(c.GetString("user_id"))

	date, _ := time.Parse(time.RFC3339, "2023-10-23T21:33:03.73297-03:00")

	transformedResponse := []services.TransformedOperation{}
	transformed := services.TransformedOperation{
		ID:     1,
		Type:   "income",
		Amount: 1200.5,
		Date:   date,
		Category: services.TransformedCategory{
			Name:  "Work",
			Color: "#fdg123",
		},
	}

	transformedResponse = append(transformedResponse, transformed)

	if userID == 1 {
		c.JSON(http.StatusOK, transformedResponse)
		return
	} else if userID == 2 {
		c.JSON(http.StatusOK, []services.TransformedOperation{})
		return
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
}

func (m *MockOperationService) Show(c *gin.Context) {
	userID, _ := strconv.Atoi(c.GetString("user_id"))

	if userID == 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	operationID, _ := strconv.Atoi(c.Param("id"))
	date, _ := time.Parse(time.RFC3339, "2023-10-23T21:33:03.73297-03:00")

	if operationID == 1 {
		c.JSON(http.StatusOK, services.TransformedShowOperation{
			ID:     1,
			Type:   "income",
			Amount: 1200.5,
			Date:   date,
			Category: services.TransformedShowCategory{
				Name:        "Work",
				Color:       "#fdg123",
				Description: "Work",
			},
			Description: "Salario",
		})
		return
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found."})
	}
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
			ExpectedBody: "{\"id\":1,\"type\":\"income\",\"amount\":1200.5,\"date\":\"2023-10-23T21:33:03.73297-03:00\",\"category\":{\"name\":\"Work\",\"color\":\"#fdg123\",\"description\":\"Work\"},\"description\":\"Salario\"}",
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
				user_id = "2"
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

			operationHandler.Show(ctx)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
}
