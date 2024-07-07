package get_service_case

import "log/slog"

type ServiceProvider interface {
	GetServiceForTopic(topic string) (error, string)
}

type RegisterCase struct {
	ServiceProvider ServiceProvider
	Logger          *slog.Logger
}

func New(saver ServiceProvider, logger *slog.Logger) *RegisterCase {
	return &RegisterCase{
		ServiceProvider: saver,
		Logger:          logger,
	}
}

func (c RegisterCase) GetServiceForTopic(topic string) (error, string) {
	return c.ServiceProvider.GetServiceForTopic(topic)
}
