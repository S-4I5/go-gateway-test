package register_case

import "log/slog"

type Registrar interface {
	AddService(topic, address string)
}

type RegisterCase struct {
	Registrar Registrar
	Logger    *slog.Logger
}

func New(saver Registrar, logger *slog.Logger) *RegisterCase {
	return &RegisterCase{
		Registrar: saver,
		Logger:    logger,
	}
}

func (c RegisterCase) RegisterService(topic, address string) {
	c.Registrar.AddService(topic, address)
}
