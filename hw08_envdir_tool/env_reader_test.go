package main

import (
	"os"
	"testing"
)

func createTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeFile(t *testing.T, dir, filename, content string) {
	t.Helper()
	err := os.WriteFile(dir+"/"+filename, []byte(content), 0o644)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadDir(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T, dir string)
		expectedEnv Environment
		expectError bool
	}{
		{
			name: "Simple case",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				writeFile(t, dir, "FOO", "123\n")
				writeFile(t, dir, "BAR", "value")
			},
			expectedEnv: Environment{
				"FOO": EnvValue{Value: "123", NeedRemove: false},
				"BAR": EnvValue{Value: "value", NeedRemove: false},
			},
			expectError: false,
		},
		{
			name: "Empty file",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				writeFile(t, dir, "EMPTY", "")
			},
			expectedEnv: Environment{
				"EMPTY": EnvValue{Value: "", NeedRemove: true},
			},
			expectError: false,
		},
		{
			name: "File with spaces and tabs",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				writeFile(t, dir, "SPACES", "  value  \t")
			},
			expectedEnv: Environment{
				"SPACES": EnvValue{Value: "  value", NeedRemove: false},
			},
			expectError: false,
		},
		{
			name: "File with null bytes",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				writeFile(t, dir, "NULLS", "value\x00with\x00nulls")
			},
			expectedEnv: Environment{
				"NULLS": EnvValue{Value: "value\nwith\nnulls", NeedRemove: false},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTempDir(t)
			defer os.RemoveAll(dir)

			tt.setup(t, dir)

			env, err := ReadDir(dir)
			if (err != nil) != tt.expectError {
				t.Fatalf("ReadDir() error = %v, expectError = %v", err, tt.expectError)
			}

			if !tt.expectError {
				for key, expectedValue := range tt.expectedEnv {
					actualValue, exists := env[key]
					if !exists {
						t.Errorf("Expected key %s not found", key)
						continue
					}
					if actualValue.Value != expectedValue.Value {
						t.Errorf("Expected %s=%s, got %s", key, expectedValue.Value, actualValue.Value)
					}
					if actualValue.NeedRemove != expectedValue.NeedRemove {
						t.Errorf("Expected NeedRemove for %s to be %v, got %v", key, expectedValue.NeedRemove, actualValue.NeedRemove)
					}
				}
			}
		})
	}
}
