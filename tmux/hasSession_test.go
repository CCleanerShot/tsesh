package tmux

import (
	"fmt"
	"reflect"
	"testing"
)

func TestHasSessionArgs(t *testing.T) {
	cmdRunner = mockCommandWithExitCode(0)
	HasSession("test")
	expectedArgs := []string{"has-session", "-t", "test"}
	if !reflect.DeepEqual(capturedArgs, expectedArgs) {
		argsNotMatching(t, expectedArgs, capturedArgs)
	}
}

func TestHasSessionExists(t *testing.T) {
	cmdRunner = mockCommandWithExitCode(0)
	found := HasSession("existing_session")
	if !found {
		fmt.Printf("expected session to be found if return code is 0\n")
		t.FailNow()
	}
}

func TestHasSessionDoesNotExist(t *testing.T) {
	cmdRunner = mockCommandWithExitCode(1)
	found := HasSession("non_existing_session")
	if found {
		fmt.Printf("expected session to not be found if return code is 1\n")
		t.FailNow()
	}
}

