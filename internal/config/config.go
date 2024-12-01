package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type HTTPServer struct {
	Addr string `yaml:"address" env:"address" env-required:"true"`
}
type Database struct {
	Host     string `yaml:"host" env:"host" env-required:"true"`
	Port     string `yaml:"port" env:"port" env-required:"true"`
	User     string `yaml:"user" env:"user" env-required:"true"`
	Password string `yaml:"password" env:"password" env-required:"true"`
	Name     string `yaml:"name" env:"name" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env:"ssl_mode" env-required:"true"`
}

type Jwt struct {
	Secret string `yaml:"secret" env:"secret" env-required:"true"`
}
type Config struct {
	Env        string `yaml:"env" env:"ENV" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	Database   `yaml:"database"`
	Jwt        `yaml:"jwt"`
}

func MustLoad() *Config {
	var configPath string
	configPath = "config/local.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found %s", configPath)
	}

	var cgf Config
	err := cleanenv.ReadConfig(configPath, &cgf)

	if err != nil {
		log.Fatalf("can't read config file %s", err.Error())
	}

	return &cgf
}
