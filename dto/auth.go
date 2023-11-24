package dto

import "github.com/dgrijalva/jwt-go"

type JWTClaim struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}
