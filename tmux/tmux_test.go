package tmux

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

/*

SwitchClient:
	- TestSwitchClientArgs
	- TestSwitchClientOutsideTmux
	- TestSwitchClientInsideTmux
	- TestSwitchClientExist
	- TestSwitchClientDoesNotExist
*/

/*
Taking inspiration from bubbletea/exec_test.go of making a small tea program to see if the error gets returned properly
*/

type mockOption func(*exec.Cmd) *exec.Cmd
type execCommand func(string, ...string) *exec.Cmd
type model struct {
	testCmd func () tea.Cmd 
	Err error
}
type tmuxTest struct {
	cmdRunner execCommand
	sessionName string
	testCmd func () tea.Cmd
	expectedArgs []string
	expectedErr error
}

var capturedArgs []string

func (m model) Init() tea.Cmd {
	return m.testCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TmuxMsg:
		m.Err = msg.Err
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	return "\n"
}

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

func TestArgs(t *testing.T) {
	t.Setenv("TMUX", "")
	tt := []tmuxTest {
		{
			cmdRunner: mockCommand(),
			testCmd: func() tea.Cmd {
				return Attach("test-attach-args")
			},
			expectedArgs: []string {"attach-session", "-t", "test-attach-args"},
		},
		{
			cmdRunner: mockCommand(),
			testCmd: func() tea.Cmd {
				return NewSession("test1", "/path/to/test1")
			},
			expectedArgs: []string{"new-session", "-s", "test1", "-c", "/path/to/test1"},
		},
	}

	for _, tc := range tt {
		cmdRunner = tc.cmdRunner
		tc.testCmd()

		if !reflect.DeepEqual(tc.expectedArgs, capturedArgs) {
			failArgsDoNotMatch(t, tc.expectedArgs, capturedArgs)
		}
	}
}

func TestInsideTmux(t *testing.T) {
	t.Setenv("TMUX", "inside")
	var in, out bytes.Buffer 

	tt := []tmuxTest {
		{
			cmdRunner: mockCommand(),
			sessionName: "attach-inside-tmux",
			testCmd: func() tea.Cmd {
				return Attach("attach-inside-tmux")
			},
			expectedErr: ErrNestedSession,
		},
		{
			cmdRunner: mockCommand(),
			testCmd: func() tea.Cmd {
				return NewSession("new-session-unique-session", "/path/to/dir")
			},
			expectedErr: nil,
		},
		{
			cmdRunner: mockCommand(
				withDuplicateSession,
			),
			testCmd: func() tea.Cmd {
				return NewSession("new-session-duplicate-session", "/path/to/dir")
			},
			expectedErr: ErrDuplicateSession,
		},
	}

	for _, tc := range tt {
		cmdRunner = tc.cmdRunner 
		initModel := model {
			testCmd: tc.testCmd,
		}

		p := tea.NewProgram(initModel, tea.WithInput(&in), tea.WithOutput(&out))
		outModel, _ := p.Run()
		finalModel := outModel.(model)

		if !errors.Is(finalModel.Err, tc.expectedErr) {
			fmt.Printf("tea.Cmd returned something unexpected\n")
			fmt.Printf("expected: %v, got: %v\n", ErrNestedSession, finalModel.Err)
			t.FailNow()
		}
	}
}

func TestOutsideTmux(t *testing.T) {
	t.Setenv("TMUX", "")

	tt := []tmuxTest {
		{
			cmdRunner: mockCommand(),
			testCmd: func() tea.Cmd {
				return Attach("attach-existing-session")
			},
			expectedErr: nil,
		},
		{
			cmdRunner: mockCommand(
				withNonExistingSession,
			),
			testCmd: func() tea.Cmd {
				return Attach("atttach-non-existing-session")
			},
			expectedErr: ErrSessionNotFound, 
		},
		{
			cmdRunner: mockCommand(),
			testCmd: func() tea.Cmd {
				return NewSession("new-session-unique-session", "/path/to/dir")
			},
			expectedErr: nil,
		},
		{
			cmdRunner: mockCommand(
				withDuplicateSession,
			),
			testCmd: func() tea.Cmd {
				return NewSession("new-session-duplicate-session", "/path/to/dir")
			},
			expectedErr: ErrDuplicateSession,
		},
	}

	var in, out bytes.Buffer
	for _, tc := range tt {
		cmdRunner = tc.cmdRunner
		initModel := model {
			testCmd: tc.testCmd,
		}

		p := tea.NewProgram(initModel, tea.WithInput(&in), tea.WithOutput(&out))
		outModel, err := p.Run()
		if err != nil {
			t.Fatalf("something went wrong with the tea program: %v\n", err)
		}
		finalModel := outModel.(model)

		if !errors.Is(finalModel.Err, tc.expectedErr) {
			fmt.Printf("expected: %v, got: %v\n", tc.expectedErr, finalModel.Err)
			t.FailNow()
		}
	}
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

func withNonExistingSession(cmd *exec.Cmd) *exec.Cmd {
	cmd.Err = fmt.Errorf("can't find session") 
	return cmd
}

func withDuplicateSession(cmd *exec.Cmd) *exec.Cmd {
	cmd.Err = fmt.Errorf("duplicate session")
	return cmd
}

func failArgsDoNotMatch(t *testing.T, expectedArgs, capturedArgs []string) {
	fmt.Println("arguments passed to tmux do not match")
	fmt.Printf("expected: %v, got:%v\n", expectedArgs, capturedArgs)
	t.FailNow()
}
