package services

import (
	"GoGin-API-CuentasClaras/api/auth"
	dao "GoGin-API-CuentasClaras/dao"
	dto "GoGin-API-CuentasClaras/dto"
	"GoGin-API-CuentasClaras/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	RegisterUser(registerUserRequest dto.RegisterUserRequest) (int, map[string]any)
	LoginUser(loginUserRequest dto.LoginRequest) (int, map[string]any)
	CurrentUser(user dao.User) (int, map[string]any)
	BalanceUser(user dao.User) (int, interface{})
}

type UserServiceImpl struct {
	userRepository      repository.UserRepository
	auth                auth.Auth
	operationRepository repository.OperationRepository
}

func (u UserServiceImpl) RegisterUser(registerUserRequest dto.RegisterUserRequest) (int, map[string]any) {
	_, recordError := u.userRepository.Save(&dao.User{
		Username: registerUserRequest.Username,
		Password: registerUserRequest.Password,
		Email:    registerUserRequest.Email,
	})

	if recordError != nil {
		return http.StatusBadRequest, gin.H{"error": recordError.Error()}
	}

	return http.StatusOK, gin.H{"message": "User successfully created."}
}

func (u UserServiceImpl) LoginUser(loginUserRequest dto.LoginRequest) (int, map[string]any) {
	var user dao.User

	user, recordError := u.userRepository.FindUserByEmail(loginUserRequest.Email)
	if recordError != nil {
		return http.StatusUnauthorized, gin.H{"error": "invalid credentials"}
	}

	credentialError := user.CheckPassword(loginUserRequest.Password)
	if credentialError != nil {
		return http.StatusUnauthorized, gin.H{"error": "invalid credentials"}
	}

	expiresIn, tokenString, err := u.auth.GenerateJWT(fmt.Sprint(user.ID))
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}
	}
	return http.StatusOK, gin.H{"token": tokenString, "expires_in": expiresIn}
}

func (u UserServiceImpl) CurrentUser(user dao.User) (int, map[string]any) {
	return http.StatusOK, gin.H{"email": user.Email, "username": user.Username}
}

func (u UserServiceImpl) BalanceUser(user dao.User) (int, interface{}) {
	println("ACAAA ID DE OPERACION: ")
	operations, _ := u.operationRepository.FindOperationsByUser(user)
	var balance float32 = 0.0
	for _, operation := range operations {
		if operation.Type == INCOME_TYPE {
			balance += float32(operation.Amount)
		} else {
			balance -= float32(operation.Amount)
		}
	}

	return http.StatusOK, gin.H{"total_balance": fmt.Sprintf("%.2f", balance)}
}

func UserServiceInit(userRepository repository.UserRepository, auth auth.Auth, operationRepository repository.OperationRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository:      userRepository,
		auth:                auth,
		operationRepository: operationRepository,
	}
}
