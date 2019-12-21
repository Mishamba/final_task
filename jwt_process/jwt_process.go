package jwt

import (
	"time"

	"github.com/Mishamba/FinalTaks/model"
	"github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("stupid dog")

func GeneratedToken(user model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = user.Name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return token.SignedString(signingKey)
}

func DecodeToken() int {

	return 0
}
