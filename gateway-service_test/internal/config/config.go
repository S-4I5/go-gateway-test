package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Route struct {
	Id         string   `yaml:"id"`
	Uri        string   `yaml:"uri"`
	Predicates string   `yaml:"predicates"`
	Filters    []string `yaml:"filters"`
}

type HTTPServer struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle-timeout" env-default:"60s"`
}

type AuthServer struct {
	Uri string `yaml:"uri"`
}

type Discovery struct {
	Uri string `yaml:"uri"`
}

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	Routes     []Route `yaml:"routes"`
	AuthServer `yaml:"auth_server"`
	Discovery  `yaml:"discovery"`
}

func MustLoad(configPath string) *Config {
	/*configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH not set")
	}*/

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Cannot find config file")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Error while reading config")
	}

	return &cfg
}
