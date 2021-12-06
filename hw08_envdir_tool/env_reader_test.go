package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("successful read dir with env", func(t *testing.T) {
		envDir := "./testdata/env"

		expectEnv := Environment{
			"BAR":   EnvValue{"bar", false},
			"EMPTY": EnvValue{"", false},
			"FOO":   EnvValue{"   foo\nwith new line", false},
			"HELLO": EnvValue{"\"hello\"", false},
			"UNSET": EnvValue{"", true},
		}

		actualEnv, err := ReadDir(envDir)

		require.Nil(t, err, "err should be nil")
		require.Len(t, actualEnv, len(expectEnv), "not all env files has been read")

		for n, ev := range expectEnv {
			av, ok := actualEnv[n]
			require.Truef(t, ok, "there is no env %s", n)
			require.Equalf(t, ev.Value, av.Value, "values of env %s is not equal", n)
			require.Equalf(t, ev.NeedRemove, av.NeedRemove, "need remove flags of env %s is not equal", n)
		}

		_, err = ReadDir("./testdata/invalid_env")
		require.ErrorIs(t, err, ErrInvalidEnvFile, "error should be return if exists invalid env file")
	})

	t.Run("unable to read dir with env files", func(t *testing.T) {
		envDir := "./testdata/not-exist-env-files-dir"

		_, err := ReadDir(envDir)
		require.NotNil(t, err, "ReadDir should return err if env dir does not exist")
		require.ErrorIs(t, err, os.ErrNotExist, "ReadDir should return wrapped os.ErrNotExist error")
	})
}
