package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func createTempFile(t *testing.T, name string, content []byte) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tmpFile.Write(content)
	if err != nil {
		tmpFile.Close()
		t.Fatal(err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func readFileContent(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestCopy(t *testing.T) {
	content := []byte("Hello, this is test data for copying.")

	fromFile := createTempFile(t, "from", content)
	defer os.Remove(fromFile)

	toFile := filepath.Join(os.TempDir(), "to_test")
	defer os.Remove(toFile)

	tests := []struct {
		name    string
		offset  int64
		limit   int64
		want    []byte
		wantErr bool
	}{
		{"full copy", 0, 0, content, false},
		{"with offset", 7, 0, content[7:], false},
		{"with limit", 0, 5, content[:5], false},
		{"offset and limit", 7, 10, content[7:17], false},
		{"offset exceeds", int64(len(content)) + 1, 0, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean before run
			os.Remove(toFile)

			err := Copy(fromFile, toFile, tt.offset, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Copy() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			got := readFileContent(t, toFile)
			if string(got) != string(tt.want) {
				t.Errorf("Content mismatch. Got %q, want %q", got, tt.want)
			}
		})
	}

	t.Run("self copy", func(t *testing.T) {
		err := Copy(fromFile, fromFile, 0, 0)
		if !errors.Is(err, ErrSelfCopy) {
			t.Errorf("The case for the same 'from' and 'to' is not handled")
		}
	})

	t.Run("from /dev/urandom", func(t *testing.T) {
		err := Copy("/dev/urandom", toFile, 0, 0)
		if !errors.Is(err, ErrUnsupportedFile) {
			t.Errorf("Failed to handle copying from /dev/urandom")
		}
	})

	t.Run("from /dev/zero", func(t *testing.T) {
		err := Copy("/dev/zero", toFile, 0, 0)
		if !errors.Is(err, ErrUnsupportedFile) {
			t.Errorf("Failed to handle copying from /dev/zero")
		}
	})
}
