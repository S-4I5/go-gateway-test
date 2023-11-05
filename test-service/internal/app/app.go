package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"test-service/internal/config"
	pingHandler "test-service/internal/controller/http/handler/ping"
	onwChiLogger "test-service/internal/controller/http/middleware"
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
