package main

import (
	"bufio"
	"os"
	"log"
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

	if fInfo.Size() == 0 {
		return true
	}
	return false
}

func processFile(fname string) (EnvValue, error) {
	var res EnvValue
    file, err := os.Open(fname)
    if err != nil {
        log.Fatalf("failed to open file: %s", err)
		return res, err
    }
    defer file.Close()

	if isEmpty(file) {
		return EnvValue{"", true}, nil
	}

    scanner := bufio.NewScanner(file)

    if scanner.Scan() {
        firstLine := scanner.Text()
		cleaned := strings.ReplaceAll(firstLine, "\x00", "\n")
		res.Value = strings.TrimRight(cleaned," \t")
    }

    if err := scanner.Err(); err != nil {
        log.Fatalf("error reading file: %s", err)
		return res, err
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
		full_name := dir + "/" + name
		envVar, err := processFile(full_name)
		if err != nil {
			return res, err
		}
		res[name] = envVar
	}
	return res, nil
}
