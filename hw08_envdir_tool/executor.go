package main

import (
	"errors"
	"os"
	"os/exec"
)

const (
	SuccessReturnCode = iota
	CantRunReturnCode
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(commands []string, env Environment) (returnCode int) {
	c := commands[0]
	args := commands[1:]
	cmd := exec.Command(c, args...)

	NormalizeEnv(env)

	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return CantRunReturnCode
	}
	return SuccessReturnCode
}

func NormalizeEnv(env Environment) {
	for n, e := range env {
		if _, ok := os.LookupEnv(n); ok {
			os.Unsetenv(n)
		}
		if !e.NeedRemove {
			os.Setenv(n, e.Value)
		}
	}
}
