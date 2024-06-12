package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config file
type Config struct {
	Port        string            `yaml:"port"`
	HTTPS       bool              `yaml:"https"`
	PostgresDB  PostgresConfig    `yaml:"postgres_db"`
	AuthService AuthServiceConfig `yaml:"auth_service"`
}

// Postgres config
type PostgresConfig struct {
	URL string `yaml:"url"`
}

// Auth service config
type AuthServiceConfig struct {
	URL     string `yaml:"url"`
	Timeout int    `yaml:"timeout"`
	JWTKey  string `yaml:"jwt_key"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
