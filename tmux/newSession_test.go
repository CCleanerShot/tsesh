package tmux

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewSession(t *testing.T) {
	tt := map[string]tmuxTest{
		"non-duplicate session outside tmux": {
			insideTmux: false,
			cmdRunner:  mockCommand(),
			expectedArgs: []string{
				"new-session", "-s", "session-name", "-c", "/path/to/dir",
			},
		},
		"non-duplicate session inside tmux": {
			insideTmux: true,
			cmdRunner:  mockCommand(),
			expectedArgs: []string{
				"new-session", "-ds", "session-name", "-c", "/path/to/dir",
			},
		},
		"duplicate session outside tmux": {
			insideTmux: false,
			cmdRunner: mockCommand(
				withDuplicateSession,
			),
			expectedArgs: []string{
				"new-session", "-s", "session-name", "-c", "/path/to/dir",
			},
			expectedErr: ErrDuplicateSession,
		},
		"duplicate session inside tmux": {
			insideTmux: true,
			cmdRunner: mockCommand(
				withDuplicateSession,
			),
			expectedArgs: []string{
				"new-session", "-ds", "session-name", "-c", "/path/to/dir",
			},
			expectedErr: ErrDuplicateSession,
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
					return NewSession("session-name", "/path/to/dir")
				},
			}

			var in, out bytes.Buffer
			p := tea.NewProgram(initModel, tea.WithInput(&in), tea.WithOutput(&out))
			outModel, _ := p.Run()
			finalModel := outModel.(model)

			if !errors.Is(tc.expectedErr, finalModel.Err) {
				fmt.Printf("tea.Cmd returned something unexpected\n")
				fmt.Printf("expected: %v, got: %v\n", tc.expectedErr, finalModel.Err)
				t.FailNow()
			}

			if !reflect.DeepEqual(tc.expectedArgs, capturedArgs) {
				failArgsDoNotMatch(t, tc.expectedArgs, capturedArgs)
			}
		})
	}
}
