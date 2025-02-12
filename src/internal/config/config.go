package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	FirstConsumer Consumer `yaml:"consumer_1" env-required:"true"`
	HTTPServer    `yaml:"http_server" env-required:"true"`
}
type Consumer struct {
	Brokers string `yaml:"env" env-required:"true"`
	GroupId string `yaml:"consumer_group" env-required:"true"`
	Topic   string `yaml:"message_topic" env-required:"true"`
}
type HTTPServer struct {
	Host        string        `yaml:"host" env-required:"true"`
	Port        int           `yaml:"port" env-required:"true"`
	Protocol    string        `yaml:"protocol" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

func EnvLoad() *Config {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("failed to load environment file, error: ", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	cfg := EnvLoadInPath("../../" + configPath)
	return cfg
}

func EnvLoadInPath(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file not found: ", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("failed to read config, error: ", err)
	}

	return &cfg
}
