package main

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("trim whitespaces", func(t *testing.T) {
		err := os.MkdirAll("testdata/testenv", 0o777)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer os.Remove("testdata/testenv")

		f, err := os.OpenFile("testdata/testenv/VAR", os.O_RDWR|os.O_CREATE, 0o644)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer os.Remove("testdata/testenv/VAR")

		if _, err := f.Write([]byte("var  \t  ")); err != nil {
			f.Close()
			log.Fatal(err)
			return
		}
		defer f.Close()

		got, _ := ReadDir("testdata/testenv")
		want := Environment{"VAR": {"var", false}}
		eq := reflect.DeepEqual(want, got)
		if !eq {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("no = in file name", func(t *testing.T) {
		err := os.MkdirAll("testdata/testenv", 0o777)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer os.Remove("testdata/testenv")

		f1, err := os.OpenFile("testdata/testenv/VAR1", os.O_RDWR|os.O_CREATE, 0o644)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer os.Remove("testdata/testenv/VAR1")

		if _, err := f1.Write([]byte("var1")); err != nil {
			f1.Close()
			log.Fatal(err)
			return
		}
		defer f1.Close()
		f2, err := os.OpenFile("testdata/testenv/VAR=2", os.O_RDWR|os.O_CREATE, 0o644)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer os.Remove("testdata/testenv/VAR=2")

		if _, err := f2.Write([]byte("var2")); err != nil {
			f2.Close()
			log.Fatal(err)
			return
		}
		defer f2.Close()
		f3, err := os.OpenFile("testdata/testenv/VAR3", os.O_RDWR|os.O_CREATE, 0o644)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer os.Remove("testdata/testenv/VAR3")

		if _, err := f3.Write([]byte("var3")); err != nil {
			f3.Close()
			log.Fatal(err)
			return
		}
		defer f3.Close()

		got, _ := ReadDir("testdata/testenv")
		want := Environment{
			"VAR1": {"var1", false},
			"VAR3": {"var3", false},
		}
		eq := reflect.DeepEqual(want, got)
		if !eq {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("empty file", func(t *testing.T) {
		err := os.MkdirAll("testdata/testenv", 0o777)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer os.Remove("testdata/testenv")

		f, err := os.OpenFile("testdata/testenv/VAR", os.O_RDWR|os.O_CREATE, 0o644)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer os.Remove("testdata/testenv/VAR")
		defer f.Close()

		got, _ := ReadDir("testdata/testenv")
		want := Environment{"VAR": {"", true}}
		eq := reflect.DeepEqual(want, got)
		if !eq {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
