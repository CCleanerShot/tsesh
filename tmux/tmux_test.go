package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

/* 

Attach:
	- TestAttachArgs
	- TestAttachOutsideTmux
	- TestAttachInsideTmux
	- TestAttachExists
	- TestAttachDoesNotExist

SwitchClient:
	- TestSwitchClientArgs
	- TestSwitchClientOutsideTmux
	- TestSwitchClientInsideTmux
	- TestSwitchClientExist
	- TestSwitchClientDoesNotExist
*/

var capturedArgs []string

func TestMain(m *testing.M) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		exitCode, _ := strconv.Atoi(os.Getenv("MOCK_CMD_EXIT_CODE"))
		os.Exit(exitCode)
	}

	os.Exit(m.Run())
}

func mockCommandWithExitCode(exitCode int) (func(string, ...string) *exec.Cmd) {
	return func(command string, args...string) *exec.Cmd {
		testBinary := os.Args[0]
		cs := []string{"--", command}
		cs = append(cs, args...) 
		capturedArgs = args
		cmd := exec.Command(testBinary, cs...)
		cmd.Env = []string{
			"GO_WANT_HELPER_PROCESS=1",
			fmt.Sprintf("MOCK_CMD_EXIT_CODE=%d", exitCode),
		}
		return cmd
	}
}

func argsNotMatching(t *testing.T, expectedArgs, capturedArgs []string) {
	fmt.Println("arguments passes to tmux do not match")
	fmt.Printf("expected: %v, got:%v\n", expectedArgs, capturedArgs)
	t.FailNow()
}
