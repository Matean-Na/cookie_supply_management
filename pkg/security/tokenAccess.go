package security

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

type TokenDetails struct {
	AccessToken string
	AccessUuid  string
	RefreshUuid string
	AtExpires   int64
	RtExpires   int64
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string, jwtKey string, expiresIn time.Duration) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresIn).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

func ValidateToken(tokenString string, jwtKey string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func ExtractToken(re *http.Request) string {
	bearToken := re.Header.Get("Authorization")

	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
