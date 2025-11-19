package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInputFileDoesNotExist = errors.New("input file does not exist")
	ErrCreateOutputFile      = errors.New("failed to create output file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return ErrInputFileDoesNotExist
	}
	defer fromFile.Close()

	fi, _ := fromFile.Stat()
	size := fi.Size()
	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return ErrCreateOutputFile
	}
	defer toFile.Close()

	fromFile.Seek(offset, io.SeekStart)

	var bytesToCopy int64
	if limit == 0 || limit > size-offset {
		bytesToCopy = size - offset
	} else {
		bytesToCopy = limit
	}

	bar := pb.Full.Start64(bytesToCopy)
	defer bar.Finish()

	barReader := bar.NewProxyReader(fromFile)

	_, err = io.CopyN(toFile, barReader, bytesToCopy)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}
