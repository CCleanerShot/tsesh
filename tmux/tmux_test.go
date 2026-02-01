package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

var capturedArgs []string

func TestMain(m *testing.M) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		os.Exit(0)
	}

	os.Exit(m.Run())
}

func TestHasSessionArgs(t *testing.T) {
	cmdRunner = mockCommand
	HasSession("test")
	expectedArgs := []string{"has-session", "-t", "test"}
	if !reflect.DeepEqual(capturedArgs, expectedArgs) {
		fmt.Println("arguments passes to tmux do not match")
		fmt.Printf("expected: %v, got:%v\n", expectedArgs, capturedArgs)
		t.FailNow()
	}
}

/* 
HasSession:
	- TestHasSessionExists
	- TestHasSessionDoesNotExist

NewSession:
	- TestNewSessionArgs 
	- TestNewSessionOutsideTmux
	- TestNewSessionInsideTmux
	- TestNewSessionExists
	- TestNewSessionDoesNotExist

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

func mockCommand(command string, args...string) *exec.Cmd {
	testBinary := os.Args[0]
	cs := []string{"--", command}
	cs = append(cs, args...) 
	capturedArgs = args
	cmd := exec.Command(testBinary, cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}
