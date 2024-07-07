package discovery_client

import (
	"context"
	"example.com/discovery-service/pkg/proto/discovery"
	"fmt"
	"time"
)

type DiscoveryClient struct {
	discovery.DiscoveryClient
	Topic   string
	Address string
}

func New(client discovery.DiscoveryClient, topic, address string) *DiscoveryClient {
	fmt.Println(topic, address)
	return &DiscoveryClient{
		DiscoveryClient: client,
		Topic:           topic,
		Address:         address,
	}
}

func (c *DiscoveryClient) Start() error {
	_, err := c.DiscoveryClient.Register(context.Background(), &discovery.RegisterRequest{
		Topic:   c.Topic,
		Address: c.Address,
	})
	if err != nil {
		return err
	}

	for {
		fmt.Println("ping")

		_, err = c.DiscoveryClient.Beat(context.Background(), &discovery.BeatRequest{
			Topic:   c.Topic,
			Address: c.Address,
		})
		if err != nil {
			return err
		}

		time.Sleep(time.Second * 2)
	}

}
