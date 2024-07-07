package services_repository

import (
	"fmt"
	"sync"
	"time"
)

type Service struct {
	lastBeat time.Time
	address  string
}

type ServiceGroup struct {
	services []*Service
	current  int
}

type Repository struct {
	sync.Map
	waitTime time.Duration
}

func New() *Repository {
	return &Repository{waitTime: 4 * time.Second}
}

func (s *ServiceGroup) DeleteInactive(waitTime time.Duration) {

	for i, service := range s.services {

		if service.lastBeat.Add(waitTime).Before(time.Now()) {
			s.current = s.current % len(s.services)
			s.services = append(s.services[:i], s.services[i+1:]...)
		}

	}
}

func (r *Repository) SaveHeartbeat(topic, address string) error {

	if topicAddresses, ok := r.Map.Load(topic); ok {
		groupService := topicAddresses.(*ServiceGroup)

		for _, service := range groupService.services {
			if service.address == address {
				service.lastBeat = time.Now()
				break
			}
		}

	} else {
		return fmt.Errorf("service not found")
	}
	return nil
}

func (r *Repository) GetServiceForTopic(topic string) (error, string) {
	if topicAddresses, ok := r.Map.Load(topic); ok {
		group := topicAddresses.(*ServiceGroup).services
		current := topicAddresses.(*ServiceGroup).current

		fmt.Println(current)

		for group[current].lastBeat.Add(r.waitTime).Before(time.Now()) {
			fmt.Println("to delete" + group[current].address)
			topicAddresses.(*ServiceGroup).services = append(group[:current], group[current+1:]...)
			group = topicAddresses.(*ServiceGroup).services

			if len(group) == 0 {
				return fmt.Errorf("topic not found"), ""
			}

			current = (current + 1) % len(group)
		}

		topicAddresses.(*ServiceGroup).current = (current + 1) % len(topicAddresses.(*ServiceGroup).services)

		return nil, group[current].address
	} else {
		return fmt.Errorf("topic not found"), ""
	}

}

func (r *Repository) AddService(topic, address string) {

	newService := Service{
		address:  address,
		lastBeat: time.Now(),
	}

	var group []*Service
	current := 0

	if topicAddresses, ok := r.Map.Load(topic); ok {
		group = append(topicAddresses.(*ServiceGroup).services, &newService)
		current = topicAddresses.(*ServiceGroup).current
		current++
	} else {
		group = []*Service{&newService}
	}

	r.Map.Store(topic, &ServiceGroup{services: group, current: current})
}
