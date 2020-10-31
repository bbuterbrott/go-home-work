package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testCopy(t, 0, 0)
		testCopy(t, 0, 10)
		testCopy(t, 0, 1000)
		testCopy(t, 100, 1000)
		testCopy(t, 6000, 1000)
	})

	t.Run("from file doesn't exist", func(t *testing.T) {
		tmpDir, tmpFile := createTmpFilePath()
		defer os.RemoveAll(tmpDir)

		err := Copy("testdata/input.txt1111111", tmpFile, 0, 100)

		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("empty file", func(t *testing.T) {
		tmpDir, tmpFile := createTmpFilePath()
		defer os.RemoveAll(tmpDir)

		err := Copy("testdata/empty.txt", tmpFile, 0, 100)

		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("offset greater then file size", func(t *testing.T) {
		tmpDir, tmpFile := createTmpFilePath()
		defer os.RemoveAll(tmpDir)

		err := Copy("testdata/input.txt", tmpFile, 50000000000, 50000000001)

		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})
}

func testCopy(t *testing.T, offset, limit int64) {
	tmpDir, tmpFile := createTmpFilePath()
	defer os.RemoveAll(tmpDir)

	err := Copy("testdata/input.txt", tmpFile, offset, limit)
	require.NoError(t, err)

	result, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		log.Fatal("read test file error")
	}
	expected, err := ioutil.ReadFile(fmt.Sprintf("testdata/out_offset%v_limit%v.txt", offset, limit))
	require.Equal(t, expected, result, "File contents do not match\nExpected:\n%v\n\n\nActual:\n%v", string(expected), string(result))
}

func createTmpFilePath() (tmpDir string, tmpFile string) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}
	return tmpDir, filepath.Join(tmpDir, "out.txt")
}
