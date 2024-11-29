package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type HTTPServer struct {
	Addr string `yaml:"address" env:"address" env-required:"true"`
}
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
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
