package tmux

import (
	"fmt"
	"reflect"
	"testing"
)

func TestHasSession(t *testing.T) {
	tt := map[string]tmuxTest {
		"existing session": {
			cmdRunner: mockCommand(),
			expectedArgs: []string{"has-session", "-t", "session-name"},
			expectedRes: true,
			expectedErr: nil,
		},
		"non-existing session": {
			cmdRunner: mockCommand(
				withExitCodeone,
			),
			expectedArgs: []string{"has-session", "-t", "session-name"},
			expectedRes: false,
			expectedErr: nil,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			capturedArgs = nil
			cmdRunner = tc.cmdRunner
			res := HasSession("session-name")

			if !reflect.DeepEqual(tc.expectedArgs, capturedArgs) {
				failArgsDoNotMatch(t, tc.expectedArgs, capturedArgs)
			}

			if tc.expectedRes != res {
				fmt.Printf("expected: %v, got: %v\n", tc.expectedRes, res)
				t.FailNow()
			}
		})
	}
}
