package main

import (
	"os"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name         string
		cmd          []string
		env          Environment
		expectedEnv  map[string]string
		expectedCode int
	}{
		{
			name: "Simple echo command",
			cmd:  []string{"echo", "Hello, World!"},
			env: Environment{
				"FOO": EnvValue{Value: "123", NeedRemove: false},
				"BAR": EnvValue{Value: "value", NeedRemove: false},
			},
			expectedEnv: map[string]string{
				"FOO": "123",
				"BAR": "value",
			},
			expectedCode: 0,
		},
		{
			name: "Command with non-zero exit code",
			cmd:  []string{"sh", "-c", "exit 1"},
			env: Environment{
				"FOO": EnvValue{Value: "123", NeedRemove: false},
			},
			expectedEnv: map[string]string{
				"FOO": "123",
			},
			expectedCode: 1,
		},
		{
			name: "Command with environment variable removal",
			cmd:  []string{"sh", "-c", "echo $REMOVE_ME"},
			env: Environment{
				"REMOVE_ME": EnvValue{Value: "", NeedRemove: true},
			},
			expectedEnv:  map[string]string{},
			expectedCode: 0,
		},
		{
			name:         "Command with no environment variables",
			cmd:          []string{"echo", "Hello, World!"},
			env:          Environment{},
			expectedEnv:  map[string]string{},
			expectedCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalEnv := os.Environ()
			returnCode := RunCmd(tt.cmd, tt.env)

			if returnCode != tt.expectedCode {
				t.Errorf("Expected return code %d, got %d", tt.expectedCode, returnCode)
			}

			for key, expectedValue := range tt.expectedEnv {
				actualValue := os.Getenv(key)
				if actualValue != expectedValue {
					t.Errorf("Expected %s=%s, got %s", key, expectedValue, actualValue)
				}
			}

			os.Clearenv()
			for _, envVar := range originalEnv {
				parts := strings.SplitN(envVar, "=", 2)
				os.Setenv(parts[0], parts[1])
			}
		})
	}
}
