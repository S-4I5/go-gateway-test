package authserver

type Provider interface {
	Authenticate(token string) bool
}

type AuthServer struct {
	Provider
}

func New(provider Provider) *AuthServer {
	return &AuthServer{Provider: provider}
}

func (au *AuthServer) Authenticate(token string) bool {
	return au.Provider.Authenticate(token)
}
