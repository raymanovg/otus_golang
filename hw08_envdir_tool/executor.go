package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(commands []string, env Environment) (returnCode int) {
	cmd := exec.Command(commands[0], commands[1:]...)

	cmd.Env = append(os.Environ(), getEnvList(env)...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	return
}

func getEnvList(env Environment) []string {
	var envList []string
	for n, e := range env {
		if !e.NeedRemove {
			envList = append(envList, fmt.Sprintf("%s=%s", n, e.Value))
		}
	}
	return envList
}
