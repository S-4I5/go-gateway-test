package heartbeat_case

import "log/slog"

type Saver interface {
	SaveHeartbeat(topic, address string) error
}

type RegisterCase struct {
	Saver  Saver
	Logger *slog.Logger
}

func New(saver Saver, logger *slog.Logger) *RegisterCase {
	return &RegisterCase{
		Saver:  saver,
		Logger: logger,
	}
}

func (c RegisterCase) SaveHeartbeat(topic, address string) error {
	return c.Saver.SaveHeartbeat(topic, address)
}
