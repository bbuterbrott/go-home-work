package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
}

func testCopy(t *testing.T, offset, limit int64) {
	err := Copy("testdata/input.txt", "/tmp/out.txt", offset, limit)
	require.NoError(t, err)
	result, err := ioutil.ReadFile("/tmp/out.txt")
	if err != nil {
		log.Fatal("read test file error")
	}
	expected, err := ioutil.ReadFile(fmt.Sprintf("testdata/out_offset%v_limit%v.txt", offset, limit))
	require.Equal(t, expected, result, "File contents do not match\nExpected:\n%v\n\n\nActual:\n%v", string(expected), string(result))
}
