package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string        `yaml:"env" env-default:"local"`
	StoragePath    string        `yaml:"storage_path" env-required:"true"`
	GRPC           GRPCConfig    `yaml:"grpc"`
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
	MigrationsPath string
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := GetConfigPath()
	if configPath == "" {
		panic("config path is emptry")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

func GetConfigPath() string {
	var path string

	flag.StringVar(&path, "config-path", "", "Path to config")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
