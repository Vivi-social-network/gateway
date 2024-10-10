package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/Vivi-social-network/core/logger"
)

const (
	envDev  = "develop"
	envProd = "prod"
)

type Config struct {
	Env     string        `yaml:"env"`
	Servers Servers       `yaml:"servers"`
	Logger  logger.Config `yaml:"logger"`
}

func (c Config) IsDev() bool {
	return c.Env == envDev
}

func (c Config) IsProd() bool {
	return c.Env == envProd
}

func Parse(configPath string) (Config, error) {
	cfg := Config{}
	cfgFileContent, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(cfgFileContent, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
