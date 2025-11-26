package main

import (
	"os/exec"
	"os"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cVec := []string{}
	for _, cw := range cmd {
		cVec = append(cVec, cw)
	}
	c := exec.Command(cVec[0], cVec[1:]...)
	cEnv := os.Environ()
	for k,v := range env {
		varStr := k + "=" + v.Value
		cEnv = append(cEnv, varStr)
	}

	c.Env = cEnv
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()

    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            return exitErr.ExitCode()
        }
    }
    return 0
}
