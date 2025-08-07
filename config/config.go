package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Address         string `yaml:"address"`
	JWTSecret       string `yaml:"jwt_secret"`
	Debug           bool   `yaml:"debug"`
	InsightsService struct {
		BaseUrl string `yaml:"base_url"`
		ApiKey  string `yaml:"api_key"`
	} `yaml:"insights_service"`
	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
		Schema   string `yaml:"schema"`
	} `yaml:"postgres"`
}

func New() *Config {
	cfg := &Config{}

	// Try to read from config.yaml
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Warning: Could not read config.yaml: %v. Using default configuration.", err)
		return cfg
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		log.Printf("Error: Could not unmarshal config.yaml: %v. Using default configuration.", err)
		return cfg
	}

	return cfg
}
