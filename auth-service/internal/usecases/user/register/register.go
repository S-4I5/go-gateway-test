package register

import (
	"fmt"
	"log/slog"
)

type Repository interface {
	SaveUser(username, password string) (int, error)
}

type UseCase struct {
	Repository Repository
	Logger     *slog.Logger
}

func New(rep Repository, logger *slog.Logger) *UseCase {
	return &UseCase{
		Repository: rep,
		Logger:     logger,
	}
}

func (useCase *UseCase) Register(username, password string) (int, error) {
	fmt.Println("XDD")
	userId, err := useCase.Repository.SaveUser(username, password)
	fmt.Println(userId, err)
	if err != nil {
		return -1, fmt.Errorf("this username is already taken")
	}

	return userId, nil
}
