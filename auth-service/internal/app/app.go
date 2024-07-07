package app

import (
	"auth-service/internal/authenticator"
	"auth-service/internal/config"
	"auth-service/internal/controllers/http/handler/user/login"
	"auth-service/internal/controllers/http/handler/user/register"
	ownChiLogger "auth-service/internal/controllers/http/middleware"
	"auth-service/internal/infrastructure/postgresdb"
	"auth-service/internal/lib/logger/sl"
	authenticatorService "auth-service/internal/proto/authenticator"
	"auth-service/internal/repositories/user"
	"auth-service/internal/usecases/user/authenticate"
	loginUseCase "auth-service/internal/usecases/user/login"
	registerUseCase "auth-service/internal/usecases/user/register"
	"auth-service/internal/util/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
)

func StartServer(config *config.Config, logger *slog.Logger) {
	dataBase := setupDB(config, logger)
	jwtUtil := setupJWTUtil(config)
	router := setupRouter(logger, dataBase, jwtUtil)

	addr := config.HTTPServer.Host + ":" + config.HTTPServer.Port

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}

	logger.Info("Starting ", config.Name, slog.String("address", addr))

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("Failed to start server")
		}
	}()

	lis, err := net.Listen("tcp", ":"+config.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := setupGRPC(logger, jwtUtil)

	if err == grpcServer.Serve(lis) {
		log.Fatalf("failed to listen: %v", err)
	}
}

func setupGRPC(logger *slog.Logger, jwtUtil *jwt.JWTGenerator) *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)
	authenticatorService.RegisterAuthenticatorServer(s, &authenticator.GRPCServer{Authenticator: authenticate.New(logger, jwtUtil)})
	return s
}

func setupDB(cfg *config.Config, logger *slog.Logger) *user.Repository {
	logger.Info("Connecting to database",
		slog.String("host", cfg.DB.Host+":"+cfg.DB.Port),
		slog.String("mode", cfg.Mode))

	dataBase, err := postgresdb.New(cfg.DB)
	if err != nil {
		logger.Error("Failed to init storage", sl.Error(err))
		os.Exit(1)
	}

	logger.Info("Connection ready")

	return &user.Repository{Holder: dataBase}
}

func setupJWTUtil(cfg *config.Config) *jwt.JWTGenerator {
	return jwt.New(cfg.Secret)
}

func setupRouter(logger *slog.Logger, rep *user.Repository, jwtUtil *jwt.JWTGenerator) chi.Router {

	router := chi.NewRouter()

	router.Use(ownChiLogger.New(logger))
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/register", register.New(logger, registerUseCase.New(rep, logger)))
	router.Put("/login", login.New(logger, loginUseCase.New(logger, rep, jwtUtil)))

	return router
}
