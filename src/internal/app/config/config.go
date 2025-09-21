package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string        `yaml:"env" env-default:"local"`
	Db             Db            `yaml:"db"`
	GRPC           GRPCConfig    `yaml:"grpc"`
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
	MigrationsPath string
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Db struct {
	Url string `yaml:"url"`
}

var GlobalInstance Config

func MustLoad() *Config {
	configPath, secretConfigPath := getConfigsPath()
	if configPath == "" {
		panic("config path is empty")
	}
	if secretConfigPath == "" {
		panic("secret config path is empty")
	}

	mustExistConfig(configPath)

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("can not read config: " + err.Error())
	}

	if existConfig(secretConfigPath) {
		if err := cleanenv.ReadConfig(secretConfigPath, &cfg); err != nil {
			panic("can not read secret configFatalf: " + err.Error())
		}
	}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic("Error reading env vars: " + err.Error())
	}

	GlobalInstance = cfg

	return &cfg
}

func getConfigsPath() (configPath string, secretConfigPath string) {
	var path string

	flag.StringVar(&path, "config-path", "", "Path to config")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	ext := filepath.Ext(path)
	nameWithoutExt := strings.TrimSuffix(path, ext)
	secretPath := fmt.Sprintf("%s.secret%s", nameWithoutExt, ext)

	return path, secretPath
}

func mustExistConfig(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}
}

func existConfig(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
