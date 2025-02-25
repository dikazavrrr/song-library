package config

import (
	"log"

	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	API struct {
		ExternalURL string `yaml:"external_url"`
	} `yaml:"api"`
}

var AppConfig Config

func LoadConfig() {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
}
