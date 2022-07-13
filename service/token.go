package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserClaims struct {
	UserId     string
	OpenID     string
	NickName   string
	Role       string
	StoreId    string
	SessionKey string
	// StandardClaims结构体实现了Claims接口(Valid()函数)
	jwt.StandardClaims
}

func Sign(claims *UserClaims, ExpiredAt time.Duration) (string, error) {

	expiredAt := time.Now().Add(time.Duration(time.Second) * ExpiredAt).Unix()

	jwtSecretKey := os.Getenv("JwtSecretKey")

	// metadata for your jwt
	claims.ExpiresAt = expiredAt

	to := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := to.SignedString([]byte(jwtSecretKey))

	if err != nil {
		return accessToken, err
	}

	return accessToken, nil
}

func VerifyTokenHeader(ctx *gin.Context) (*jwt.Token, error) {
	tokenHeader := ctx.GetHeader("Authorization")
	if tokenHeader == "" {
		return nil, fmt.Errorf("Authorization must not empty")
	}

	accessToken := strings.SplitAfter(tokenHeader, "Bearer")[1]
	jwtSecretKey := os.Getenv("JwtSecretKey")

	token, err := jwt.ParseWithClaims(strings.Trim(accessToken, " "), &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func VerifyToken(accessToken string) (*jwt.Token, error) {
	jwtSecretKey := os.Getenv("JwtSecretKey")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
