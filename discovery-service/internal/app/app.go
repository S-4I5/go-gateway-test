package app

import (
	"discovery-service/internal/config"
	discoveryHandler "discovery-service/internal/handlers/grpc/discovery"
	"discovery-service/internal/repository/services_repository"
	"discovery-service/internal/use_case/get_service_case"
	"discovery-service/internal/use_case/heartbeat_case"
	"discovery-service/internal/use_case/register_case"
	"discovery-service/pkg/proto/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"log/slog"
	"net"
)

func StartServer(config *config.Config, logger *slog.Logger) {

	lis, err := net.Listen("tcp", ":"+config.GRPCServer.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repository := services_repository.New()

	grpcServer := setupGRPC(logger, repository)

	if err == grpcServer.Serve(lis) {
		log.Fatalf("failed to listen: %v", err)
	}
}

func setupGRPC(logger *slog.Logger, rep *services_repository.Repository) *grpc.Server {

	s := grpc.NewServer()
	reflection.Register(s)

	discovery.RegisterDiscoveryServer(s, &discoveryHandler.Service{
		Registrar:       register_case.New(rep, logger),
		ServiceProvider: get_service_case.New(rep, logger),
		BeatHolder:      heartbeat_case.New(rep, logger),
	})

	return s
}
