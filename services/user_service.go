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
	RegisterUser(registerUserRequest dto.RegisterUserRequest, c *gin.Context) (int, map[string]any)
	LoginUser(c *gin.Context)
	CurrentUser(c *gin.Context)
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
	auth           auth.Auth
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u UserServiceImpl) RegisterUser(registerUserRequest dto.RegisterUserRequest, c *gin.Context) (int, map[string]any) {

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

func (u UserServiceImpl) LoginUser(c *gin.Context) {
	var request LoginRequest
	var user dao.User
	validationError := c.ShouldBindJSON(&request)
	if validationError != nil || request.Email == "" || request.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters."})
		return
	}

	user, recordError := u.userRepository.FindUserByEmail(request.Email)
	if recordError != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	expiresIn, tokenString, err := u.auth.GenerateJWT(fmt.Sprint(user.ID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "expires_in": expiresIn})
}

func (u UserServiceImpl) CurrentUser(c *gin.Context) {
	user, recordError := RetrieveCurrentUser(u.userRepository, c.GetString("user_id"))

	if recordError != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": user.Email, "username": user.Username})
}

func UserServiceInit(userRepository repository.UserRepository, auth auth.Auth) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
		auth:           auth,
	}
}
