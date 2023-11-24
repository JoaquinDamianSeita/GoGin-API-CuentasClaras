package middleware

import (
	"GoGin-API-CuentasClaras/config"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(initConfig *config.Initialization) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
			return
		}

		tokenString := authHeaderParts[1]
		claims, err := initConfig.Auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
			return
		}

		intUserID, _ := strconv.Atoi(claims.UserID)
		user, recordError := initConfig.UserRepo.FindUserById(intUserID)

		if recordError != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
