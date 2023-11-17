package auth

import (
	"GoGin-API-CuentasClaras/repository"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("SECRET_JWT_KEY"))

type JWTClaim struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type Auth interface {
	GenerateJWT(userId string) (expiresIn int64, tokenString string, err error)
	ValidateToken(signedToken string) (claims *JWTClaim, err error)
}

type AuthImpl struct {
	userRepository repository.UserRepository
}

func (auth AuthImpl) GenerateJWT(userId string) (expiresIn int64, tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	expiresIn = int64(time.Until(expirationTime).Seconds())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func (auth AuthImpl) ValidateToken(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil, err
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return nil, err
	}

	intUserID, _ := strconv.Atoi(claims.UserID)
	_, recordError := auth.userRepository.FindUserById(intUserID)

	if recordError != nil {
		err = errors.New("invalid token")
		return nil, err
	}

	return claims, nil
}

func AuthInit(userRepository repository.UserRepository) *AuthImpl {
	return &AuthImpl{
		userRepository: userRepository,
	}
}
