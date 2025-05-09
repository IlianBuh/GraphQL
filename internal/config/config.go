package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	defaultConfigPath = "./config/config.yml"
)

type Config struct {
	Env       string    `yaml:"env" env-required=true`
	Port      int       `yaml:"port" env-required=true`
	SSOClient SSOConfig `yaml:"sso-client" env-required=true`
}

type SSOConfig struct {
	Host string `yaml:"host" env-default="localhost"`
	Port int    `yaml:"port" env-required=true`
}

func New() *Config {
	const op = "config.New"
	path := fetchPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	cfg := MustLoad(path)

	return cfg
}

func MustLoad(path string) *Config {
	const op = "config.MustLoad"

	cfg, err := Load(path)
	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	return cfg
}

func Load(path string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func fetchPath() string {
	res := ""

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res != "" {
		return res
	}

	res = os.Getenv("CONFIG-PATH")
	if res != "" {
		return res
	}

	return defaultConfigPath
}
