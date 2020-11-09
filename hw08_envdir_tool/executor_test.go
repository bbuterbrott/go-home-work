package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cmd := make([]string, 0)
		cmd = append(cmd, "testdata/echo.sh", "arg1", "arg2")
		env := make(map[string]string, 0)
		env["BAR"] = "V1"
		env["FOO"] = "V2"
		env["HELLO"] = "V3"
		env["HOME"] = ""
		returnCode := RunCmd(cmd, env)
		require.Equal(t, 0, returnCode, "return code is not 0")
	})

	t.Run("error", func(t *testing.T) {
		cmd := make([]string, 0)
		cmd = append(cmd, "testdata/error.sh", "arg1", "arg2")
		env := make(map[string]string, 0)
		returnCode := RunCmd(cmd, env)
		require.Equal(t, 3, returnCode, "return code is not 3")
	})
}
