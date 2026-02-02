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
/*
TODO:
needs to be re written to account if it is inside tmux or not
if it is inside tmux then return a sequence of new-session (detached) and then switch-client but if it is outside tmux then just new-session (not detached)

This way on picker.go if there is no existing session it can just call tmux.NewSession()

This will change how it is tested so I will leave that till the end, the other ones are free to create tests for since those do not need this much conditional branching
*/
func NewSession(sessionName, workingDirectory string, detached bool) tea.Cmd {
	// use Inside() instead of detached to determine how it is ran
	sessionFlags := "-s"
	if detached {
		sessionFlags = "-ds"
	}

	// use a cmds []tea.Cmd that will append appropriate cmds based on Inside() 
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
