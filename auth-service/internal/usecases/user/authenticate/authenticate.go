package authenticate

import (
	"fmt"
	"log/slog"
)

type JWTUtil interface {
	Validate(tokenString string) (string, error)
}

type UseCase struct {
	logger  *slog.Logger
	jwtUtil JWTUtil
}

func New(logger *slog.Logger, generator JWTUtil) *UseCase {
	return &UseCase{
		logger:  logger,
		jwtUtil: generator,
	}
}

func (useCase *UseCase) Authenticate(token string) (string, error) {
	username, err := useCase.jwtUtil.Validate(token)
	fmt.Println("sus")

	if err != nil {
		fmt.Println("susXD")
		return "", err
	}

	return username, nil
}
