package main

import (
	"testing"
	"os"
)

func TestRunCmd(t *testing.T) {
	t.Run("simple ls to file", func(t *testing.T) {
		cmd := []string{"/bin/bash", "t1.sh"}
		ret := RunCmd(cmd, Environment{})
		if ret != 0 {
			t.Errorf("RunCmd return code [%d] != [0]", ret)
		}
		contents, err := os.ReadFile("res")
		if err != nil {
			t.Errorf("Failed to read 'res' file")
		}
		got := string(contents)
		want := "env_reader.go\nenv_reader_test.go\nexecutor.go\nexecutor_test.go\ngo.mod\nmain.go\nREADME.md\nres\nt1.sh\nt2.sh\nt3.sh\ntestdata\ntest.sh\n"
		if got != want {
			t.Errorf("got '%v' want '%v'", got, want)
		}
		os.Remove("res")
	})

	t.Run("return error", func(t *testing.T) {
		cmd := []string{"/bin/bash", "t2.sh"}
		ret := RunCmd(cmd, Environment{})
		if ret != 1 {
			t.Errorf("RunCmd return code [%d] != [1]", ret)
		}
	})

	t.Run("set VAR=var", func(t *testing.T) {
		env := Environment{ "VARTEST": {"var", false}}
		cmd := []string{"/bin/bash", "t3.sh"}
		ret := RunCmd(cmd, env)
		if ret != 0 {
			t.Errorf("RunCmd return code [%d] != [0]", ret)
		}
		contents, err := os.ReadFile("res")
		if err != nil {
			t.Errorf("Failed to read 'res' file")
		}
		got := string(contents)
		want := "VARTEST=var\n"
		if got != want {
			t.Errorf("got '%v' want '%v'", got, want)
		}
		os.Remove("res")
	})


}
