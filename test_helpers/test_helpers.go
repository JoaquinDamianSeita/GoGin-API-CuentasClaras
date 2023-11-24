package testhelpers

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestStructure struct {
	Name         string
	Params       string
	ExpectedCode int
	ExpectedBody string
}

type TestInterfaceStructure struct {
	Name         string
	Params       interface{}
	ExpectedCode int
	ExpectedBody string
}

func MockPostRequest(request_body string, uri string) (*gin.Context, *httptest.ResponseRecorder) {
	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.Request = httptest.NewRequest("POST", uri, strings.NewReader(request_body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	return ctx, responseRecorder
}

func MockGetRequest(uri string) (*gin.Context, *httptest.ResponseRecorder) {
	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.Request = httptest.NewRequest("GET", uri, nil)
	return ctx, responseRecorder
}

func AssertExpectedCodeAndBodyResponse(t *testing.T, tt TestStructure, responseRecorder *httptest.ResponseRecorder) {
	assert.Equal(t, tt.ExpectedCode, responseRecorder.Code)
	if tt.ExpectedBody != "" {
		assert.Equal(t, tt.ExpectedBody, responseRecorder.Body.String())
	}
}

func AssertExpectedCodeAndResponseService(t *testing.T, tt TestStructure, code int, response map[string]any) {
	assert.Equal(t, tt.ExpectedCode, code)
	if tt.ExpectedBody != "" {
		responseString := mapToString(response)

		assert.Equal(t, tt.ExpectedBody, responseString)
	}
}

func AssertExpectedCodeAndResponseServiceDto(t *testing.T, tt TestInterfaceStructure, code int, response interface{}) {
	assert.Equal(t, tt.ExpectedCode, code)
	if tt.ExpectedBody != "" {
		responseString := mapStructToString(response)

		assert.Equal(t, tt.ExpectedBody, responseString)
	}
}

func mapToString(m map[string]interface{}) string {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func mapStructToString(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
