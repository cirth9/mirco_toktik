package utils

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// UserBasicClaims
// jwt claims
type UserBasicClaims struct {
	UserInfoID int64
	Username   string
	Password   string
	jwt.StandardClaims
}

// TokenExpireDuration  过期时间
const TokenExpireDuration = time.Hour * 2

// Secret 密钥
var Secret = []byte("secret")

//1694321493
//1694328693

// GetToken 获取token
func GetToken(username, password string, userID int64) (string, error) {
	expirationTime := time.Now().Add(TokenExpireDuration)
	log.Println(">>>>>>>>>>>>>", time.Now().Unix(), expirationTime.Unix())
	claims := &UserBasicClaims{
		UserInfoID: userID,
		Username:   username,
		Password:   password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(Secret)
}

// ParseToken 解析token
func ParseToken(tokenString string) (*UserBasicClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString,
		&UserBasicClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return Secret, nil
		})

	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if value, ok := tokenClaims.Claims.(*UserBasicClaims); ok && tokenClaims.Valid {
			return value, nil
		}
	}

	return nil, nil
}
