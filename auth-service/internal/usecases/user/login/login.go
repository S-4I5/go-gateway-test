package login

import (
	user2 "auth-service/internal/repositories/user"
	"fmt"
	"log"
	"log/slog"
)

type UseCase struct {
	repository Repository
	logger     *slog.Logger
	jwtUtil    JWTUtil
}

type JWTUtil interface {
	GenerateJWT(username string) (string, error)
}

type user struct {
	UserId   int
	Username string
	Password string
	Role     string
}

type Repository interface {
	GetUser(username string) (*user2.User, error)
}

func New(logger *slog.Logger, rep Repository, generator JWTUtil) *UseCase {
	return &UseCase{
		repository: rep,
		logger:     logger,
		jwtUtil:    generator,
	}
}

func (useCase *UseCase) Login(username, password string) (string, error) {
	userFormDb, err := useCase.repository.GetUser(username)
	if err != nil || userFormDb == nil {
		return "", err
	}

	if password != userFormDb.Password {
		return "", fmt.Errorf("incorrect cred")
	}

	log.Println("login: ", username)

	token, err := useCase.jwtUtil.GenerateJWT(username)
	if err != nil {
		return "", fmt.Errorf("cannot generate jwt")
	}

	return token, nil
}
