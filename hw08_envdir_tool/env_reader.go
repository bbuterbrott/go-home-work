package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(map[string]string)

	fileStats, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, fileStat := range fileStats {
		fileName := strings.TrimSpace(fileStat.Name())

		if fileStat.IsDir() {
			continue
		}

		if strings.Contains(fileName, "=") {
			continue
		}

		file, err := os.Open(path.Join(dir, fileName))
		if err != nil {
			return nil, err
		}
		defer file.Close()

		reader := bufio.NewReader(file)

		firstLine, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, err
			}
		}

		firstLine = strings.ReplaceAll(firstLine, "\000", "\n")

		env[fileName] = strings.TrimRight(firstLine, "\t \n")
	}

	return env, nil
}
