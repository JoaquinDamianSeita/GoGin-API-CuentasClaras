package middleware

import (
	"GoGin-API-CuentasClaras/config"
	"GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/dto"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type JWTClaim struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type MockAuthValidator struct{}

func (auth *MockAuthValidator) GenerateJWT(userId string) (expiresIn int64, tokenString string, err error) {
	return
}

func (auth *MockAuthValidator) ValidateToken(signedToken string) (claims *dto.JWTClaim, err error) {
	log.Println(signedToken)
	if signedToken == "valid_token" {
		claims = &dto.JWTClaim{UserID: "1"}
		return claims, nil
	} else if signedToken == "invalid_user" {
		claims = &dto.JWTClaim{UserID: "2"}
		return claims, nil
	}
	return nil, errors.New("Invalid token")
}

type MockUserRepository struct{}

func (u MockUserRepository) FindUserById(id int) (dao.User, error) {
	if id == 1 {
		return dao.User{ID: 1}, nil
	}
	return dao.User{}, errors.New("User not found")
}
func (u MockUserRepository) FindUserByEmail(email string) (dao.User, error) { return dao.User{}, nil }
func (u MockUserRepository) Save(user *dao.User) (dao.User, error)          { return dao.User{}, nil }

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &config.Initialization{
		Auth:     &MockAuthValidator{},
		UserRepo: &MockUserRepository{},
	}

	middleware := AuthMiddleware(config)

	t.Run("Valid Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer valid_token")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		middleware(c)

		user, _ := c.Get("user")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, dao.User{ID: 1}, user)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Empty(t, c.GetString("user_id"))
	})

	t.Run("User Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalid_user")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		config.UserRepo = &MockUserRepository{}

		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Empty(t, c.GetString("user_id"))
	})
}
