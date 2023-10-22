package integration_tests

import (
	"GoGin-API-Base/api"
	"GoGin-API-Base/api/auth"
	"GoGin-API-Base/config"
	"GoGin-API-Base/dao"
	"GoGin-API-Base/repository"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	testhelpers "GoGin-API-Base/test_helpers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var db *gorm.DB
var token string

func TestMain(m *testing.M) {
	db = InitTest()
	defer cleanDB()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func InitTest() *gorm.DB {
	godotenv.Load("../.env")
	os.Setenv("DB_DSN", os.Getenv("DB_DSN_TEST"))

	return config.ConnectToDB()
}

func cleanDB() {
	db.Exec("DELETE FROM users")
	fmt.Println("Database cleaned.")
}

func setupTest() *gin.Engine {
	fmt.Println("Before Test Execution.")
	userRepositoryImpl := repository.UserRepositoryInit(db)
	user, _ := userRepositoryImpl.Save(&dao.User{
		Username: "pedro.fuentes",
		Email:    "pedro.fuentes@gmail.com",
		Password: "password123",
	})
	authService := auth.AuthInit()
	_, token, _ = authService.GenerateJWT(fmt.Sprintf(strconv.Itoa(user.ID)))
	init := config.Init()
	return api.Init(init)
}

func teardownTest() {
	cleanDB()
}

func TestUsersIntegration_Login_ValidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the request is successful",
			Params:       `{"email": "pedro.fuentes@gmail.com", "password": "password123"}`,
			ExpectedCode: http.StatusOK,
			ExpectedBody: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/api/users/login", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestUsersIntegration_Login_InvalidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when email is not present",
			Params:       `{"email": "", "password": "password123"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when password is not present",
			Params:       `{"email": "pedro.fuentes@gmail.com", "password": ""}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"Invalid parameters.\"}",
		},
		{
			Name:         "when the email is wrong",
			Params:       `{"email": "pedro@gmail.com", "password": "password123"}`,
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "{\"error\":\"invalid credentials\"}",
		},
		{
			Name:         "when the password is wrong",
			Params:       `{"email": "pedro.fuente@gmail.com", "password": "password12"}`,
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "{\"error\":\"invalid credentials\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/api/users/login", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestUsersIntegration_Register_ValidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the request is successful",
			Params:       `{"username": "test.user", "email": "test@example.com", "password": "password123"}`,
			ExpectedCode: http.StatusOK,
			ExpectedBody: `{"message":"User successfully created."}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/api/users", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestUsersIntegration_Register_InvalidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
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
			Params:       `{"email": "pedro.fuentes@gmail.com", "username": "test.user", "password": "password123"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"the email or the user is already in use\"}",
		},
		{
			Name:         "when username is not available",
			Params:       `{"email": "test.user@example.com", "username": "pedro.fuentes", "password": "password123"}`,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "{\"error\":\"the email or the user is already in use\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/api/users", strings.NewReader(tt.Params))
			request.Header.Set("Content-Type", "application/json")

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestUsersIntegration_Current_ValidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when the request is successful",
			Params:       "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{\"email\":\"pedro.fuentes@gmail.com\",\"username\":\"pedro.fuentes\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/api/users/current", strings.NewReader(tt.Params))
			request.Header.Set("Authorization", "Bearer " + token)

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}

func TestUsersIntegration_Current_InvalidRequest(t *testing.T) {
	router := setupTest()
	var tests = []testhelpers.TestStructure{
		{
			Name:         "when user does not exists",
			Params:       "",
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "{\"error\":\"Not authorized\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/api/users/current", strings.NewReader(tt.Params))

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			testhelpers.AssertExpectedCodeAndBodyResponse(t, tt, responseRecorder)
		})
	}
	teardownTest()
}
