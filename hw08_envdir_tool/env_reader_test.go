package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		env, err := ReadDir("testdata/env")

		require.NoError(t, err)

		expectedEnv := Environment(make(map[string]string, 0))
		expectedEnv["BAR"] = "bar"
		expectedEnv["FOO"] = "foo"
		expectedEnv["HELLO"] = "\"hello\""
		expectedEnv["UNSET"] = ""
		require.Equal(t, expectedEnv, env, "Environments doesn't match")
	})

	t.Run("= in name", func(t *testing.T) {
		env, err := ReadDir("testdata/env-error")
		require.EqualError(t, err, "env file name (BAR=) shouldn't contain '='")
		require.Nil(t, env)
	})

	t.Run("not a folder", func(t *testing.T) {
		env, err := ReadDir("testdata/error.sh")
		require.EqualError(t, err, "'testdata/error.sh' is not a directory")
		require.Nil(t, env)
	})
}
