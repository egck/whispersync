package config

import (
	"errors"
)

// Agent config part

type AgentConfig struct {
	ServerURI          string `yaml:"server_uri" env:"WS_SERVER_URI"`
	BufferFlushPeriod  int    `yaml:"buffer_flush_period" env:"WS_BUFFER_FLUSH_PERIOD"`
	WatchDirectoryPath string `yaml:"watch_directory_path" env:"WS_WATCH_DIRECTORY_PATH"`
	Secret             string `yaml:"secret" env:"WS_SECRET"`
}

func (ac *AgentConfig) Validate() error {
	if ac.ServerURI == "" {
		return errors.New("missing 'server_uri' field")
	}
	if ac.WatchDirectoryPath == "" {
		return errors.New("missing 'watch_directory_path' field")
	}
	return nil
}

// Server config part

type ServerConfig struct {
	APIEndpoint     string `yaml:"api_endpoint" env:"WS_API_ENDPOINT"`
	DataStoragePath string `yaml:"data_storage_path" env:"WS_DATA_STORAGE_PATH"`
	Secret          string `yaml:"secret" env:"WS_SECRET"`
}

func (sc *ServerConfig) Validate() error {
	if sc.APIEndpoint == "" {
		return errors.New("missing 'api_endpoint' field")
	}
	if sc.DataStoragePath == "" {
		return errors.New("missing 'data_storage_path' field")
	}
	return nil
}
