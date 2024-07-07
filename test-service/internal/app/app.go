package app

import (
	"example.com/discovery-service/pkg/proto/discovery"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"log/slog"
	"net/http"
	"test-service/internal/config"
	pingHandler "test-service/internal/controller/http/handler/ping"
	onwChiLogger "test-service/internal/controller/http/middleware"
	discovery_client "test-service/internal/discovery-client"
	"test-service/internal/usecase/ping"
)

func StartServer(config *config.Config, logger *slog.Logger) {
	router := setupRouter(config, logger)

	server := &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}

	go func() {
		err := setupAndStartDiscoveryClient(config, logger)
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	logger.Info("Starting server", slog.String("address", config.Address))

	if err := server.ListenAndServe(); err != nil {
		logger.Error("Failed to start server")
	}
}

func setupRouter(config *config.Config, logger *slog.Logger) chi.Router {

	router := chi.NewRouter()

	router.Use(onwChiLogger.New(logger))
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/ping", pingHandler.New(logger, ping.New(config.Name)))

	return router
}

func setupAndStartDiscoveryClient(config *config.Config, logger *slog.Logger) error {

	if !config.Enable {
		return nil
	}

	logger.Info("setting up discovery")

	conn, err := grpc.Dial(config.Uri, grpc.WithInsecure())
	if err != nil {
		return nil
	}

	discoverClient := discovery_client.New(discovery.NewDiscoveryClient(conn), config.Topic, config.Address)

	if err = discoverClient.Start(); err != nil {
		return fmt.Errorf("cannot start discovery client")
	}

	return fmt.Errorf("discovery error")
}
