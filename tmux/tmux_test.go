package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

/* 

SwitchClient:
	- TestSwitchClientArgs
	- TestSwitchClientOutsideTmux
	- TestSwitchClientInsideTmux
	- TestSwitchClientExist
	- TestSwitchClientDoesNotExist
*/

type mockOption func(*exec.Cmd) *exec.Cmd
type execCommand func(string, ...string) *exec.Cmd

var capturedArgs []string

func TestMain(m *testing.M) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		var exitCode int

		codeStr := os.Getenv("MOCK_CMD_EXIT_CODE")
		if codeStr == "" {
			exitCode = 0
		} else {
			exitCode, _ = strconv.Atoi(codeStr)
		}

		os.Exit(exitCode)
	}

	os.Exit(m.Run())
}

// Unless another option that modifies the exit code is passed it will default to exiting with code 0
func mockCommand(mockOpts...mockOption) execCommand {
	return func(command string, args...string) *exec.Cmd {
		capturedArgs = args
		testBinary := os.Args[0]

		cs := []string{"--", command}
		cs = append(cs, args...) 

		cmd := exec.Command(testBinary, cs...)
		cmd.Env = []string{
			"GO_WANT_HELPER_PROCESS=1",
		}
		
		for _, f := range mockOpts {
			cmd = f(cmd)
		}

		return cmd
	}
}

func withExitCodeOne(cmd *exec.Cmd) *exec.Cmd {
	cmd.Env = append(cmd.Env, "MOCK_CMD_EXIT_CODE=1")
	return cmd
}

func insideTmux(cmd *exec.Cmd) *exec.Cmd {
	cmd.Env = append(cmd.Env, "TMUX=/run/user/100/tmux-1000/default,0000,0")
	return cmd
}

func failArgsDoNotMatch(t *testing.T, expectedArgs, capturedArgs []string) {
	fmt.Println("arguments passed to tmux do not match")
	fmt.Printf("expected: %v, got:%v\n", expectedArgs, capturedArgs)
	t.FailNow()
}
