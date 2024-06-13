package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

// Config file
type Config struct {
	Port        string            `yaml:"port"`
	HTTPS       bool              `yaml:"https"`
	URLAUTH     string            `yaml:"url_auth"`
	URLDB       string            `yaml:"url_db"`
	AuthService AuthServiceConfig `yaml:"auth_service"`
}

// Auth service config
type AuthServiceConfig struct {
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
	dbURL := os.Getenv("DB_URL")
	if dbURL != "" {
		config.URLDB = dbURL
	}
	flag.StringVar(&dbURL, "url_db", config.URLDB, "database url")
	if dbURL != "" {
		config.URLDB = dbURL
	}
	authURL := os.Getenv("AUTH_URL")
	if authURL != "" {
		config.URLAUTH = authURL
	}
	flag.StringVar(&authURL, "url_auth", config.URLAUTH, "auth url")
	if authURL != "" {
		config.URLAUTH = authURL
	}
	return config, nil
}
