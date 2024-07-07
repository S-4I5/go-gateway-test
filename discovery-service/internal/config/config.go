package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type GRPCServer struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port string `yaml:"port" env-default:"8080"`
	Name string `yaml:"name"`
}

type ServerRegistry struct {
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
}

type Config struct {
	Env            string `yaml:"env" env-required:"true"`
	ServerRegistry `yaml:"server-registry"`
	GRPCServer     `yaml:"grpc-server"`
}

func MustLoad(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Cannot find config file")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Error while reading config")
	}

	return &cfg
}
