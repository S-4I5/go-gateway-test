package main

import (
	"gateway-service/internal/app"
	"gateway-service/internal/config"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")

	app.StartServer(cfg)
}
