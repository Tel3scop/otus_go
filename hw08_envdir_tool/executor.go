package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	// #nosec G204
	command := exec.Command(cmd[0], cmd[1:]...)

	for name, envValue := range env {
		if envValue.NeedRemove {
			err := os.Unsetenv(name)
			if err != nil {
				return 0
			}
		} else {
			err := os.Setenv(name, envValue.Value)
			if err != nil {
				return 0
			}
		}
	}

	command.Env = os.Environ()
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		return 1
	}

	return 0
}
