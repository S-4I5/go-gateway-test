package gprc

import (
	"context"
	"gateway-service/internal/proto/authenticator"
)

type Provider struct {
	authenticator.AuthenticatorClient
}

func New(auth authenticator.AuthenticatorClient) *Provider {
	return &Provider{AuthenticatorClient: auth}
}

func (provider *Provider) Authenticate(token string) (string, error) {
	res, err := provider.AuthenticatorClient.Authenticate(context.Background(), &authenticator.AuthenticationRequest{Token: token})
	if err != nil {
		return "", err
	}

	return res.Username, nil
}
