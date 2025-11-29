package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cVec := []string{}
	cVec = append(cVec, cmd...)
	c := exec.Command(cVec[0], cVec[1:]...) //nolint
	cEnv := os.Environ()
	for k, v := range env {
		varStr := k + "=" + v.Value
		cEnv = append(cEnv, varStr)
	}

	c.Env = cEnv
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	err := c.Run()
	if err != nil {
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}
