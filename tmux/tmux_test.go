package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

/*
Taking inspiration from bubbletea/exec_test.go of making a small tea program to see if the error gets returned properly
*/

type mockOption func(*exec.Cmd) *exec.Cmd
type execCommand func(string, ...string) *exec.Cmd
type model struct {
	testCmd func() tea.Cmd
	Err     error
}
type tmuxTest struct {
	name         string
	insideTmux   bool
	cmdRunner    execCommand
	expectedArgs []string
	expectedRes  any
	expectedErr  error
}

/*
important: clear out previous captured args since they will only update after
mockCommand is ran, which it wont if an error is expected prior to running the
cmdRunner
*/
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

// Unless another option that modifies the exit code is passed it will default to exiting with code 0
func mockCommand(mockOpts ...mockOption) execCommand {
	return func(command string, args ...string) *exec.Cmd {
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

func withExitCodeone(cmd *exec.Cmd) *exec.Cmd {
	cmd.Env = append(cmd.Env, "MOCK_CMD_EXIT_CODE=1")
	return cmd
}

func withNonExistingSession(cmd *exec.Cmd) *exec.Cmd {
	cmd.Err = fmt.Errorf("can't find session")
	return cmd
}

func withDuplicateSession(cmd *exec.Cmd) *exec.Cmd {
	cmd.Err = fmt.Errorf("duplicate session")
	return cmd
}

func withNoClient(cmd *exec.Cmd) *exec.Cmd {
	cmd.Err = fmt.Errorf("no client found")
	return cmd
}

func failArgsDoNotMatch(t *testing.T, expectedArgs, capturedArgs []string) {
	fmt.Println("arguments passed to tmux do not match")
	fmt.Printf("expected: %v, got:%v\n", expectedArgs, capturedArgs)
	t.FailNow()
}
