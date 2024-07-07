package discovery

import (
	"context"
	discovery2 "discovery-service/pkg/proto/discovery"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Registrar interface {
	RegisterService(topic, address string)
}

type ServiceProvider interface {
	GetServiceForTopic(topic string) (error, string)
}

type BeatHolder interface {
	SaveHeartbeat(topic, address string) error
}

type Service struct {
	Registrar       Registrar
	ServiceProvider ServiceProvider
	BeatHolder      BeatHolder
	discovery2.UnimplementedDiscoveryServer
}

func (s *Service) Register(ctx context.Context, req *discovery2.RegisterRequest) (*emptypb.Empty, error) {

	s.Registrar.RegisterService(req.Topic, req.Address)

	return &emptypb.Empty{}, nil
}

func (s *Service) Beat(ctx context.Context, req *discovery2.BeatRequest) (*emptypb.Empty, error) {

	err := s.BeatHolder.SaveHeartbeat(req.Topic, req.Address)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) GetService(ctx context.Context, req *discovery2.GetServiceRequest) (*discovery2.GetServiceResponse, error) {
	fmt.Println("getReq", req.Topic)
	err, result := s.ServiceProvider.GetServiceForTopic(req.Topic)
	if err != nil {
		return nil, err
	}
	fmt.Println("res", result)
	return &discovery2.GetServiceResponse{Address: result}, nil
}
