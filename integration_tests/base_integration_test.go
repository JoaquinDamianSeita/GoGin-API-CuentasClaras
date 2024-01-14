package integration_tests

import (
	"GoGin-API-CuentasClaras/api"
	"GoGin-API-CuentasClaras/api/auth"
	"GoGin-API-CuentasClaras/config"
	"GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/repository"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var db *gorm.DB
var token string
var anotherToken string
var authService *auth.AuthImpl

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
	db.Exec("DROP TABLE users CASCADE;")
	db.Exec("DROP TABLE operations CASCADE;")
	db.Exec("DROP TABLE categories CASCADE;")
	fmt.Println("Database cleaned.")
}

func setupTest() *gin.Engine {
	userRepositoryImpl := repository.UserRepositoryInit(db)
	categoryRepositoryImpl := repository.CategoryRepositoryInit(db)
	operationRepositoryImpl := repository.OperationRepositoryInit(db)
	date, _ := time.Parse(time.RFC3339, "2023-10-23T21:33:03.73297Z")
	utcLocation, _ := time.LoadLocation("UTC")
	dateInUTC := date.In(utcLocation)
	user, _ := userRepositoryImpl.Save(&dao.User{
		Username: "pedro.fuentes",
		Email:    "pedro.fuentes@gmail.com",
		Password: "password123",
	})
	anotherUser, _ := userRepositoryImpl.Save(&dao.User{
		Username: "jose.marin",
		Email:    "jose.marin@gmail.com",
		Password: "password123",
	})
	category, _ := categoryRepositoryImpl.Save(&dao.Category{
		Name:        "Work",
		Color:       "#fdg123",
		Description: "Work",
		IsDefault:   true,
	})
	categoryRepositoryImpl.Save(&dao.Category{
		Name:        "Custom",
		Color:       "#6495ed",
		Description: "Custom",
		IsDefault:   false,
		UserID:      uint(anotherUser.ID),
	})
	operationRepositoryImpl.Save(&dao.Operation{
		UserID:      uint(user.ID),
		Category:    category,
		Type:        "income",
		Amount:      1200.5,
		Date:        dateInUTC,
		Description: "Salario",
	})
	authService = auth.AuthInit()
	_, token, _ = authService.GenerateJWT(strconv.Itoa(user.ID))
	_, anotherToken, _ = authService.GenerateJWT(strconv.Itoa(anotherUser.ID))

	init := config.Init()
	return api.Init(init)
}

func teardownTest() {
	cleanDB()
}
