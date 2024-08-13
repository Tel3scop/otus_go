package main

import "testing"

func TestRunCmd(t *testing.T) {
	env := Environment{
		"FOO": EnvValue{Value: "123", NeedRemove: false},
		"BAR": EnvValue{Value: "value", NeedRemove: false},
	}

	cmd := []string{"echo", "Hello, World!"}
	returnCode := RunCmd(cmd, env)

	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}
}
