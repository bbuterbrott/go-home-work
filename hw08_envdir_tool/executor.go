package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec. Да, это небезопасно, но в данном случае это необходимо

	for name, value := range env {
		if value == "" {
			os.Unsetenv(name)
			continue
		}
		os.Setenv(name, value)
	}
	command.Env = os.Environ()

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		fmt.Printf("Got error while running command: %v\n", err)
		if exitError, ok := err.(*exec.ExitError); ok { //nolint:errorlint. Мы знаем, что ошибка будет без вложенностей
			return exitError.ExitCode()
		}
		return 1
	}
	return 0
}
