package tmux

import (
	"reflect"
	"testing"
)

/*
Attach:
	✓ TestAttachArgs
	- TestAttachInsideTmux
	- TestAttachOutsideTmux
	- TestAttachExists
	- TestAttachDoesNotExist
*/

func TestAttachArgs(t *testing.T) {
	cmdRunner = mockCommand()
	Attach("test-attach")
	expectedArgs := []string {"attach-session", "-t", "test-attach"}
	if !reflect.DeepEqual(expectedArgs, capturedArgs) {
		failArgsDoNotMatch(t, expectedArgs, capturedArgs)
	}
}

func TestAttachInsideTmux(t *testing.T) {
	cmdRunner = mockCommand(
		insideTmux,
	)
}


