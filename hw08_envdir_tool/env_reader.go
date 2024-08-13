package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if strings.Contains(fileName, "=") {
			return nil, errors.New("file name contains '='")
		}

		filePath := fmt.Sprintf("%s/%s", dir, fileName)
		fileContent, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(fileContent)
		if scanner.Scan() {
			value := scanner.Text()
			value = strings.TrimRight(value, " \t")
			value = strings.ReplaceAll(value, "\x00", "\n")
			env[fileName] = EnvValue{Value: value, NeedRemove: false}
		} else {
			env[fileName] = EnvValue{Value: "", NeedRemove: true}
		}
		err = fileContent.Close()
		if err != nil {
			return nil, err
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return env, nil
}
