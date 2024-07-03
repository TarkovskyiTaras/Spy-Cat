package main

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

type ServerConfig struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type DatabaseConfig struct {
	ConnectionString string `json:"connectionString"`
}

func LoadConfig(path string) (*Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	bytes, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
