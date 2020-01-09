package jwt_process

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("stupid dog")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GeneratedToken(username string) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = username
	claims["exp"] = time.Now().Unix()

	expireTime := time.Now()
	signedToken, err := token.SignedString(signingKey)

	return signedToken, expireTime, err
}

func DecodeToken(tokenString string) (bool, int, string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, http.StatusUnauthorized, "", err
		}

		return false, http.StatusBadRequest, "", err
	}

	if !token.Valid {
		return false, http.StatusUnauthorized, "", err
	}

	return true, 200, claims.Username, nil
}
