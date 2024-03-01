package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Module struct {
	Name   string
	Type   ResultType
	Method string
	Params []string
	Labels map[string]string
}

type ResultType string

const (
	BoolResult ResultType = "bool"
	HexResult  ResultType = "hex"
	//NumberResult ResultType = "number"
)

type Config struct {
	Modules map[string]Module `yaml:"modules"`
}

func LoadConfig(configPath string) (Config, error) {
	var config Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}
