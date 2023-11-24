package handlers

import (
	"GoGin-API-CuentasClaras/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParseUserFromContext(ctx *gin.Context) dao.User {
	user, ctxError := ctx.Get("user")
	userStruct, ok := user.(dao.User)
	if !ctxError || !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return dao.User{}
	}
	return userStruct
}
