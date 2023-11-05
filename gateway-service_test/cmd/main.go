package main

import (
	"fmt"
	"gateway-service-test/internal/authserver"
	"gateway-service-test/internal/config"
	"gateway-service-test/internal/utils/jsonutil"
	"github.com/go-chi/render"
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

func main() {

	cfg := config.MustLoad("./config/config.yaml")

	fmt.Println(cfg)

	jwtGenerator := jsonutil.New("204a2a5a7d652931bf8b33e87578f1d655fbed400527fa9c08f4925e82009fc6")

	myAuth := authserver.New("amogus", "sus", jwtGenerator)

	routes := cfg.Routes

	for i := range routes {
		log.Println("XD", i)
		redir := routes[i].Predicates
		filters := New(routes[i].Filters)

		http.HandleFunc(routes[i].Uri, func(res http.ResponseWriter, req *http.Request) {
			/*if req.Method != "POST" {
				log.Println("XD1")
				http.NotFound(res, req)
				return
			}*/

			//log.Println(req.Method)

			log.Printf("incoming request: %s %s", req.Host, req.URL.String())

			if filters.AuthReq && !myAuth.Authenticate(req.Header.Get("Authorization")) {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}

			if filters.RelayUsername {
				fmt.Println("XDD")
				req.Header.Add("X-Kys", "kiss yourself")
			}

			proxy("http://"+redir, res, req)
		})
	}

	http.HandleFunc("/login", func(res http.ResponseWriter, req *http.Request) {
		log.Println("XD")
		if req.Method != "POST" {
			log.Println("XD1")
			http.NotFound(res, req)
			return
		}

		var request LoginRequest

		err := render.DecodeJSON(req.Body, &request)
		if err != nil {
			log.Println("XD2")
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		log.Println(request.Login, request.Password)

		token, err := myAuth.DoLogin(request.Login, request.Password)
		if err == nil {
			res.WriteHeader(http.StatusOK)
			render.JSON(res, req, LoginResponse{Token: token})
		} else {
			res.WriteHeader(http.StatusUnauthorized)
			render.JSON(res, req, LoginResponse{Token: "fuck u!"})
		}
	})

	http.HandleFunc("/logout", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.NotFound(res, req)
			return
		}

		if !myAuth.Authenticate(req.Header.Get("Authorization")) {
			res.WriteHeader(http.StatusUnauthorized)
			//...
			return
		}
	})

	log.Println("Started shit")
	err := http.ListenAndServe(":8880", nil)
	if err != nil {
		log.Panic(err)
	}
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
