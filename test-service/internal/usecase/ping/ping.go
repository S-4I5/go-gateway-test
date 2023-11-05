package ping

import "fmt"

type UseCase struct {
	ServerName string
}

func New(serverName string) *UseCase {
	return &UseCase{serverName}
}

func (useCase *UseCase) Ping(userName string) (string, error) {
	return fmt.Sprintf("Hello %s from %s!", userName, useCase.ServerName), nil
}
