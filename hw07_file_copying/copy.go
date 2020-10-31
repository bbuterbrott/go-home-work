package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

const (
	BufSize = 1024
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fromFi, err := os.Stat(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	fromFileSize := fromFi.Size()
	if fromFileSize == 0 {
		return ErrUnsupportedFile
	}

	if offset > fromFileSize {
		return ErrOffsetExceedsFileSize
	}

	progressBar := pb.StartNew(int(limit))
	progressBar.Set(pb.Bytes, true)
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("couldn't open from file '%v'", fromPath)
	}
	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("couldn't create to file '%v'", toPath)
	}

	if _, err := fromFile.Seek(offset, io.SeekCurrent); err != nil {
		return fmt.Errorf("couldn't seek from file '%v'", fromPath)
	}

	remaining := limit
	if limit == 0 {
		remaining = fromFileSize
	}
	var currentLimit int64 = BufSize
	for {
		if remaining < BufSize {
			currentLimit = remaining
		}
		written, err := io.CopyN(toFile, fromFile, currentLimit)
		progressBar.Add(int(written))
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Completed copying of file")
				break
			}
			return errors.New("couldn't copy bytes")
		}
		if written < BufSize {
			fmt.Println("Completed copying of file")
			break
		}
	}
	progressBar.Finish()

	return nil
}
