package authenticator

import (
	"auth-service/internal/proto/authenticator"
	"context"
	"fmt"
)

type Authenticator interface {
	Authenticate(token string) (string, error)
}

type GRPCServer struct {
	Authenticator
	authenticator.UnimplementedAuthenticatorServer
}

func (server *GRPCServer) Authenticate(ctx context.Context, req *authenticator.AuthenticationRequest) (*authenticator.AuthenticationResponse, error) {
	response, err := server.Authenticator.Authenticate(req.Token)
	if err != nil {
		fmt.Println("bg! ", response)
		return nil, err
	}

	fmt.Println("XDD ", response)

	return &authenticator.AuthenticationResponse{
		Username: response,
	}, err
}
