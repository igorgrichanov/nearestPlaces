package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	DataPath   string  `yaml:"data_path"`
	SchemaPath string  `yaml:"schema_path"`
	Elastic    Elastic `yaml:"elastic"`
	Server     Server  `yaml:"server"`
	Token      Token   `yaml:"token"`
}

type Elastic struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Server struct {
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
}

type Token struct {
	Secret string        `yaml:"secret"`
	TTL    time.Duration `yaml:"ttl"`
	Skew   time.Duration `yaml:"skew"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("CONFIG_PATH does not exist: ", configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatal("cannot read config: ", err)
	}

	return &config
}
