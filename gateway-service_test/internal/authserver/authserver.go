package authserver

import (
	"fmt"
	"gateway-service-test/internal/utils/jsonutil"
	"log"
)

type JWTGenerator interface {
	GenerateJWT(username string) (string, error)
	Validate(tokenString string) bool
}

type AuthServer struct {
	login        string
	password     string
	jwtGenerator JWTGenerator
}

func New(login, password string, jwtGenerator *jsonutil.JWTGenerator) *AuthServer {
	return &AuthServer{
		jwtGenerator: jwtGenerator,
		password:     password,
		login:        login,
	}
}

func (authServer *AuthServer) DoLogin(login, password string) (string, error) {
	if login != authServer.login || password != authServer.password {
		return "", fmt.Errorf("incorrect creditians")
	}

	token, _ := authServer.jwtGenerator.GenerateJWT(login)

	log.Println("Generated", token)

	return token, nil
}

func (authServer *AuthServer) Authenticate(token string) bool {
	return authServer.jwtGenerator.Validate(token)
}
