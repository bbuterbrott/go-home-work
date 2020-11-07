package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)

	for name, value := range env {
		if value == "" {
			os.Unsetenv(name)
		} else {
			os.Setenv(name, value)
		}
	}
	command.Env = os.Environ()

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Got error while running command: %v\n", err)
			return exitError.ExitCode()
		}
		return 1
	}
	return 0
}
