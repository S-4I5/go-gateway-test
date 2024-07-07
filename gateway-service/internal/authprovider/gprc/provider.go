package gprc

import (
	"context"
	authenticator2 "gateway-service/pkg/proto/authenticator"
)

type Provider struct {
	authenticator2.AuthenticatorClient
}

func New(auth authenticator2.AuthenticatorClient) *Provider {
	return &Provider{AuthenticatorClient: auth}
}

func (provider *Provider) Authenticate(token string) (string, error) {
	res, err := provider.AuthenticatorClient.Authenticate(context.Background(), &authenticator2.AuthenticationRequest{Token: token})
	if err != nil {
		return "", err
	}

	return res.Username, nil
}
