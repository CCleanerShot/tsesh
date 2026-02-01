package tmux

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: ideas
// check out session groups and see if I can implement a shortcut to add a new session to a group

var cmdRunner = exec.Command

type TmuxMsg struct { Err error }

// Checked environment variable $TMUX to determine if user is currently inside a tmux session
func Inside() bool {
	_, inside := os.LookupEnv("TMUX")
	return inside 
}

// Checks if a session already exists in tmux server
// Uses the same lookup rules, see tmux docs for rules.
func HasSession(targetSession string) bool {
	cmd := tmux("has-session", "-t", targetSession)
	return cmd.Run() == nil
}

// Attach session if outside of tmux
func Attach(sessionName string) tea.Cmd {
	return tea.ExecProcess(
		tmux("attach-session", "-t", sessionName),
		execCallback,
	)
}

// Move active session to existing session if inside tmux 
func SwitchClient(sessionName string) tea.Cmd {
	return tea.ExecProcess(
		tmux("switch-client", "-t", sessionName),
		execCallback,
	)
}

// New session will be created with specified name at provided workingDirectory
func NewSession(sessionName, workingDirectory string, detached bool) tea.Cmd {
	sessionFlags := "-s"
	if detached {
		sessionFlags = "-ds"
	}

	return tea.ExecProcess(
		tmux(
			"new-session",
			sessionFlags, sessionName, 
			"-c", workingDirectory,
		),
		execCallback,
	) 
}

func tmux(args... string) *exec.Cmd {
	cmd := cmdRunner("tmux", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func execCallback(err error) tea.Msg {
	return TmuxMsg {
		Err: err,
	}
}
