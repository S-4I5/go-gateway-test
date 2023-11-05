package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"strconv"
	"time"
)

type JWTGenerator struct {
	secret string
}

type Payload struct {
	User string
	Exp  string
}

func New(secret string) *JWTGenerator {
	return &JWTGenerator{secret: secret}
}

func (gen *JWTGenerator) GenerateJWT(username string) (string, error) {

	log.Println("gen" + username)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	log.Println("Token", token, time.Now().Add(time.Hour*24).Unix())

	singedToken, _ := token.SignedString([]byte(gen.secret))

	log.Println("SToken", singedToken)

	return singedToken, nil
}

func (gen *JWTGenerator) Validate(tokenString string) (string, error) {

	tokenString = tokenString[7:len(tokenString)]

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(gen.secret), nil
	}

	token, err := jwt.Parse(tokenString, keyFunc, jwt.WithJSONNumber())
	if err != nil {
		return "", fmt.Errorf("incorrect jwt")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		tokenExpTime, err := strconv.Atoi(fmt.Sprint(claims["exp"]))
		if err != nil {
			return "", fmt.Errorf("jwt expiered")
		}

		log.Println("time", tokenExpTime)
		log.Println("user: ", fmt.Sprint(claims["user"]))

		if int64(tokenExpTime) > time.Now().Unix() {
			return fmt.Sprint(claims["user"].(string)), nil
		}
	}

	return "", fmt.Errorf("incorrect jwt")
}
