package internal

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger      LoggerConfig      `yaml:"logger"`
	GRPCServer  GRPCServerConfig  `yaml:"grpcserver"`
	GRPCGateway GRPCGatewayConfig `yaml:"grpcgateway"`
	Swagger     SwaggerConfig     `yaml:"swagger"`
}

type LoggerConfig struct {
	Level    string `yaml:"level"`
	Filename string `yaml:"filename"`
}

type GRPCServerConfig struct {
	Addr string `yaml:"address"`
	Port string `yaml:"port"`
}
type GRPCGatewayConfig struct {
	Addr string `yaml:"address"`
	Port string `yaml:"port"`
}

type SwaggerConfig struct {
	Addr string `yaml:"address"`
	Port string `yaml:"port"`
}

func NewConfig() (*Config, error) {
	conf := new(Config)

	file, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	err = conf.load(file)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (c *Config) load(configFile []byte) error {

	if configFile == nil {
		return errors.New("config file not found")
	}

	if err := yaml.Unmarshal(configFile, c); err != nil {
		return err
	}

	return nil
}
