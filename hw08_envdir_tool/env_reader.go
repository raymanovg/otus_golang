package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir with env files: %w", err)
	}

	env := make(Environment, len(files))
	for _, envFile := range files {
		envName := envFile.Name()
		fileName := path.Join(dir, envName)
		content, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("unable to read env file: %w", err)
		}

		replaced := bytes.Replace(content, []byte("\x00"), []byte("\n"), -1)
		lines := bytes.Split(replaced, []byte("\n"))
		val := strings.TrimRight(string(lines[0]), "\n \t")

		env[envName] = EnvValue{val, len(content) == 0}
	}

	return env, nil
}
