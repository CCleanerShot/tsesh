package tmux

import (
	"reflect"
	"testing"
)


/*
I should really think about how I handle being outside of tmux vs inside tmux.
I feel that that logic should be here and depending of where I am I check 
different arguments (the 2nd one to see if it is -s or -ds)

NewSession:
	✓ TestNewSessionArgs
	- TestNewSessionOutsideTmux
	- TestNewSessionInsideTmux
	- TestNewSessionExists
	- TestNewSessionDoesNotExist
*/

func TestNewSessionArgs(t *testing.T) {
	cmdRunner = mockCommandWithExitCode(0)
	NewSession("test1", "/path/to/test1", false)
	expectedArgs := []string{"new-session", "-s", "test1", "-c", "/path/to/test1"}
	if !reflect.DeepEqual(expectedArgs, capturedArgs) {
		argsNotMatching(t, expectedArgs, capturedArgs)
	}
}
