package testhelpers

import (
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
