package tmux

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSwitchClient(t *testing.T) {
	tt := map[string]tmuxTest{
		"existing session inside tmux": {
			insideTmux:   true,
			cmdRunner:    mockCommand(),
			expectedArgs: []string{"switch-client", "-t", "session-name"},
			expectedErr:  nil,
		},
		"existing session outside tmux": {
			insideTmux: false,
			cmdRunner: mockCommand(
				withNoClient,
			),
			expectedArgs: []string{"switch-client", "-t", "session-name"},
			expectedErr:  ErrNoClientFound,
		},
		"non-existing session inside tmux": {
			insideTmux: true,
			cmdRunner: mockCommand(
				withNonExistingSession,
			),
			expectedArgs: []string{"switch-client", "-t", "session-name"},
			expectedErr:  ErrSessionNotFound,
		},
		"non-existing session outside tmux": {
			insideTmux: false,
			cmdRunner: mockCommand(
				withNoClient,
			),
			expectedArgs: []string{"switch-client", "-t", "session-name"},
			expectedErr:  ErrNoClientFound,
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
			initModel := model{
				testCmd: func() tea.Cmd {
					return SwitchClient("session-name")
				},
			}

			var in, out bytes.Buffer
			p := tea.NewProgram(initModel, tea.WithInput(&in), tea.WithOutput(&out))
			outModel, _ := p.Run()
			finalModel := outModel.(model)

			if !errors.Is(tc.expectedErr, finalModel.Err) {
				fmt.Printf("expected: %v, got: %v\n", tc.expectedErr, finalModel.Err)
				t.FailNow()
			}

			if !reflect.DeepEqual(tc.expectedArgs, capturedArgs) {
				failArgsDoNotMatch(t, tc.expectedArgs, capturedArgs)
			}
		})
	}
}
