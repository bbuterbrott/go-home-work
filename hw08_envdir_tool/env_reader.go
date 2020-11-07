package main

import (
	"bufio"
	"errors"
	"fmt"
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

	dirStat, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !dirStat.IsDir() {
		return nil, fmt.Errorf("'%v' is not a directory", dir)
	}

	fileStats, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, fileStat := range fileStats {
		fileName := strings.TrimSpace(fileStat.Name())

		if strings.Contains(fileName, "=") {
			return nil, fmt.Errorf("env file name (%v) shouldn't contain '='", fileName)
		}

		file, err := os.Open(path.Join(dir, fileName))
		if err != nil {
			return nil, err
		}

		runes := make([]rune, 0)
		reader := bufio.NewReader(file)
		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return nil, err
			}

			if r == '\n' {
				break
			}

			if r == 0 {
				r = '\n'
			}

			runes = append(runes, r)
		}

		firstLine := string(runes)
		env[fileName] = strings.TrimRight(firstLine, "\t \n")
	}

	return env, nil
}
