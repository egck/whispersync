package config

import "testing"

func TestAgentConfigValidate(t *testing.T) {
	tests := []struct {
		name          string
		config        AgentConfig
		expectedError string
	}{
		{
			name: "valid configuration",
			config: AgentConfig{
				ServerURI:          "localhost:8080",
				WatchDirectoryPath: "/tmp/watch",
			},
		},
		{
			name:          "missing server URI",
			config:        AgentConfig{WatchDirectoryPath: "/tmp/watch"},
			expectedError: "missing 'server_uri' field",
		},
		{
			name:          "missing watch directory",
			config:        AgentConfig{ServerURI: "localhost:8080"},
			expectedError: "missing 'watch_directory_path' field",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.config.Validate()
			if test.expectedError == "" {
				if err != nil {
					t.Fatalf("expected valid configuration, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error %q, got nil", test.expectedError)
			}
			if err.Error() != test.expectedError {
				t.Fatalf("expected error %q, got %q", test.expectedError, err.Error())
			}
		})
	}
}

func TestServerConfigValidate(t *testing.T) {
	tests := []struct {
		name          string
		config        ServerConfig
		expectedError string
	}{
		{
			name: "valid configuration",
			config: ServerConfig{
				APIEndpoint:     ":8080",
				DataStoragePath: "/tmp/storage",
			},
		},
		{
			name:          "missing API endpoint",
			config:        ServerConfig{DataStoragePath: "/tmp/storage"},
			expectedError: "missing 'api_endpoint' field",
		},
		{
			name:          "missing storage path",
			config:        ServerConfig{APIEndpoint: ":8080"},
			expectedError: "missing 'data_storage_path' field",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.config.Validate()
			if test.expectedError == "" {
				if err != nil {
					t.Fatalf("expected valid configuration, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error %q, got nil", test.expectedError)
			}
			if err.Error() != test.expectedError {
				t.Fatalf("expected error %q, got %q", test.expectedError, err.Error())
			}
		})
	}
}
