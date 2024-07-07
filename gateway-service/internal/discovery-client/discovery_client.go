package discovery_client

import (
	"context"
	"example.com/discovery-service/pkg/proto/discovery"
	"fmt"
)

type DiscoveryGatewayClient struct {
	discovery.DiscoveryClient
}

func New(client discovery.DiscoveryClient) *DiscoveryGatewayClient {
	return &DiscoveryGatewayClient{client}
}

func (c DiscoveryGatewayClient) GetService(topic string) (string, error) {
	result, err := c.DiscoveryClient.GetService(context.Background(), &discovery.GetServiceRequest{Topic: topic})
	if err != nil {

		fmt.Println(err.Error())
		return "", err
	}
	fmt.Println(result.Address, 2)
	return result.Address, nil
}
