package services

import (
	"GoGin-API-Base/api/auth"
	dao "GoGin-API-Base/dao"
	"GoGin-API-Base/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	RegisterUser(c *gin.Context)
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

func (u UserServiceImpl) RegisterUser(c *gin.Context) {
	var request dao.User

	validationError := c.ShouldBindJSON(&request)
	if validationError != nil || request.Username == "" || request.Email == "" || request.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters."})
		return
	}

	_, recordError := u.userRepository.Save(&request)
	if recordError != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": recordError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully created."})
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
	var user dao.User
	userID, _ := strconv.Atoi(c.GetString("user_id"))
	user, recordError := u.userRepository.FindUserById(userID)

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
