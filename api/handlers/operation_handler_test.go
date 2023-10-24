package handlers

import (
	testhelpers "GoGin-API-Base/test_helpers"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type MockOperationService struct{}

type TransformedOperation struct {
	ID       int                 `json:"id"`
	Type     string              `json:"type"`
	Amount   float64             `json:"amount"`
	Date     time.Time           `json:"date"`
	Category TransformedCategory `json:"category"`
}

type TransformedCategory struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (m *MockOperationService) Index(c *gin.Context) {
	userID, _ := strconv.Atoi(c.GetString("user_id"))

	date, _ := time.Parse(time.RFC3339, "2023-10-23T21:33:03.73297-03:00")

	transformedResponse := []TransformedOperation{}
	transformed := TransformedOperation{
		ID:     1,
		Type:   "income",
		Amount: 1200.5,
		Date:   date,
		Category: TransformedCategory{
			Name:  "Work",
			Color: "#fdg123",
		},
	}

	transformedResponse = append(transformedResponse, transformed)

	if userID == 1 {
		c.JSON(http.StatusOK, transformedResponse)
		return
	} else if userID == 2 {
		c.JSON(http.StatusOK, []TransformedOperation{})
		return
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
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
			ctx, responseRecorder := testhelpers.MockPostRequest(tt.Params, serviceUri)

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
