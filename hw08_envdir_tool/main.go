package main

import "os"

func main() {
	envDir := os.Args[1]
	cmd := os.Args[2:]
	env, err := ReadDir(envDir)
	rc := 1
	if err == nil {
		rc = RunCmd(cmd, env)
	}

	os.Exit(rc)
}
