package app

import (
	"example.com/discovery-service/pkg/proto/discovery"
	"fmt"
	gprcAuth "gateway-service/internal/authprovider/gprc"
	"gateway-service/internal/config"
	discoveryClient "gateway-service/internal/discovery-client"
	"gateway-service/pkg/proto/authenticator"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

const (
	AUTH_REQ       = "AuthReq"
	RELAY_USERNAME = "RelayUsername"
)

type Filters struct {
	AuthReq       bool
	RelayUsername bool
}

func New(filtersStrings []string) *Filters {
	filter := Filters{
		AuthReq:       false,
		RelayUsername: false,
	}

	for i := range filtersStrings {
		switch filtersStrings[i] {
		case AUTH_REQ:
			filter.AuthReq = true
		case RELAY_USERNAME:
			filter.RelayUsername = true
		}

	}

	return &filter
}

func StartServer(cfg *config.Config) {

	myAuth := setupAuthService(cfg)

	setupCustomRoutes(cfg, myAuth)

	err := http.ListenAndServe(cfg.HTTPServer.Host+":"+cfg.HTTPServer.Port, nil)
	if err != nil {
		log.Panic(err)
	}
}

func setupDiscoveryClient(config *config.Config) *discoveryClient.DiscoveryGatewayClient {
	conn, err := grpc.Dial(config.Discovery.Uri, grpc.WithInsecure())
	if err != nil {
		fmt.Println("cannot setup con with discovery")
		return nil
	}

	return discoveryClient.New(discovery.NewDiscoveryClient(conn))
}

func setupCustomRoutes(config *config.Config, myAuth *gprcAuth.Provider) {
	routes := config.Routes

	disc := setupDiscoveryClient(config)

	for i := range routes {
		filters := New(routes[i].Filters)
		id := routes[i].Id
		endPoint := routes[i].Uri

		http.HandleFunc(routes[i].Uri, func(res http.ResponseWriter, req *http.Request) {

			log.Printf("incoming request: %s %s", req.Host, req.URL.String())

			username := ""

			if filters.AuthReq {
				var err error
				username, err = myAuth.Authenticate(req.Header.Get("Authorization"))
				if err != nil {
					res.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			if filters.RelayUsername {
				req.Header.Add("X-Username", username)
			}

			redir, err := disc.GetService(id)
			if err != nil {
				log.Println("cannot reach discovery")
			}

			proxy("http://"+redir+endPoint, res, req)
		})
	}
}

func setupAuthService(config *config.Config) *gprcAuth.Provider {

	http.HandleFunc("/login", func(res http.ResponseWriter, req *http.Request) {
		proxy("http://"+config.AuthServer.Uri+config.AuthServer.Login, res, req)
	})

	http.HandleFunc("/register", func(res http.ResponseWriter, req *http.Request) {
		proxy("http://"+config.AuthServer.Uri+config.AuthServer.Register, res, req)
	})

	conn, err := grpc.Dial(config.Grpc, grpc.WithInsecure())
	if err != nil {
		return nil
	}

	return gprcAuth.New(authenticator.NewAuthenticatorClient(conn))
}

func proxy(targetURL string, res http.ResponseWriter, req *http.Request) {
	target, err := url.Parse(targetURL)

	if err != nil {
		http.Error(res, "Invalid URL", 500)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(request *http.Request) {
		request.URL.Scheme = target.Scheme
		request.URL.Host = target.Host
		request.URL.Path = target.Path
	}

	log.Printf("Forwarding request to %v", target)
	proxy.ServeHTTP(res, req)
}
