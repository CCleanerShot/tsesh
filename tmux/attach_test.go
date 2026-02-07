package tmux

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestAttach(t *testing.T) {
	tt := map[string]tmuxTest {
		"existing session outside tmux": {
			insideTmux: false,
			cmdRunner: mockCommand(),
			expectedArgs: []string{"attach-session", "-t", "session-name"},
			expectedErr: nil,
		},
		"existing session inside tmux": {
			insideTmux: true,
			cmdRunner: mockCommand(),
			expectedArgs: nil,
			expectedErr: ErrNestedSession,
		},
		"non-existing session outside tmux": {
			insideTmux: false,
			cmdRunner: mockCommand(
				withNonExistingSession,
			),
			expectedArgs: []string{"attach-session", "-t", "session-name"},
			expectedErr: ErrSessionNotFound,
		},
		"non-existing session inside tmux": {
			insideTmux: true,
			cmdRunner: mockCommand(
				withNonExistingSession,
			),
			expectedArgs: nil,
			expectedErr: ErrNestedSession,
		},
	}

	for name, tc := range tt {
		capturedArgs = nil 
		t.Run(name, func(t *testing.T) {

			tmuxEnvStr := ""
			if tc.insideTmux {
				tmuxEnvStr = "inside"
			}
			t.Setenv("TMUX", tmuxEnvStr)

			cmdRunner = tc.cmdRunner
			initModel := model {
				testCmd: func () tea.Cmd {
					return Attach("session-name")
				},
			}
			
			var in, out bytes.Buffer 
			p := tea.NewProgram(initModel, tea.WithInput(&in), tea.WithOutput(&out))
			outModel, _ := p.Run()
			finalModel := outModel.(model)
			
			if !errors.Is(tc.expectedErr, finalModel.Err) {
				fmt.Printf("tea.Cmd returned something unexpected\n")
				fmt.Printf("expected: %v, got: %v\n", ErrNestedSession, finalModel.Err)
				t.FailNow()
			}

			if !reflect.DeepEqual(tc.expectedArgs, capturedArgs) {
				failArgsDoNotMatch(t, tc.expectedArgs, capturedArgs)
			}
		})
	}
}
