package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func isEmpty(f *os.File) bool {
	fInfo, _ := f.Stat()

	return fInfo.Size() == 0
}

func processFile(fname string) (EnvValue, error) {
	var res EnvValue
	file, err := os.Open(fname)
	if err != nil {
		return res, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if isEmpty(file) {
		return EnvValue{"", true}, nil
	}

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		firstLine := scanner.Text()
		cleaned := strings.ReplaceAll(firstLine, "\x00", "\n")
		res.Value = strings.TrimRight(cleaned, " \t")
	}

	if err := scanner.Err(); err != nil {
		return res, fmt.Errorf("error reading file: %w", err)
	}
	return res, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	res := make(map[string]EnvValue)

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
		return res, err
	}

	for _, file := range files {
		name := file.Name()
		if strings.Contains(name, "=") {
			continue
		}
		fullName := dir + "/" + name
		envVar, err := processFile(fullName)
		if err != nil {
			return res, err
		}
		res[name] = envVar
	}
	return res, nil
}
