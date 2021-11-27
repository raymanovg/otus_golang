package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("env manipulation", func(t *testing.T) {
		os.Clearenv()

		os.Setenv("FOO", "foo")
		os.Setenv("BAR", "bar")

		newEnv := Environment{
			"FOO": EnvValue{"reset_foo", false}, // resetting old value of env var
			"BAR": EnvValue{"bar", true},        // deleting exist env var
			"BAZ": EnvValue{"baz", true},        // trying to delete not exist env var
			"BAM": EnvValue{"bam", false},       // setting new env var
		}

		NormalizeEnv(newEnv)

		expectedEnv := []string{"BAM=bam", "FOO=reset_foo"}
		actualEnv := os.Environ()
		require.ElementsMatch(t, expectedEnv, actualEnv)
	})

	t.Run("exit code check", func(t *testing.T) {
		returnCodes := []int{SuccessReturnCode, 1, 2, 3, 4, 5}
		for _, expectedReturnCode := range returnCodes {
			cmd := []string{"./testdata/exit_code.sh", strconv.Itoa(expectedReturnCode)}
			os.Clearenv()
			actualReturnCode := RunCmd(cmd, nil)
			require.Equal(t, expectedReturnCode, actualReturnCode)
		}

		notExistCmd := []string{"./testdata/not-exist-cmd.sh"}
		require.Equal(t, CantRunReturnCode, RunCmd(notExistCmd, nil))
	})
}
