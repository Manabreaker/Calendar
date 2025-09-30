package configs

import (
	"github.com/Manabreaker/Calendar/internal/app/apiserver"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server apiserver.ServerConfig `yaml:"server"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
