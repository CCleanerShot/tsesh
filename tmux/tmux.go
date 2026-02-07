package tmux

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type TmuxMsg struct { Err error }

var cmdRunner = exec.Command
var ErrNestedSession = errors.New("sessions should be nested with care, unset $TMUX to force")
var ErrSessionNotFound = errors.New("session was not found")
var ErrDuplicateSession = errors.New("duplicateSession")

// Checks environment variable $TMUX to determine if user is currently inside a tmux session
func Inside() bool {
	val, _ := os.LookupEnv("TMUX")
	return val != ""
}

// Checks if a session already exists in tmux server
// Uses the same lookup rules, see tmux docs for rules.
func HasSession(targetSession string) bool {
	cmd := tmux("has-session", "-t", targetSession)
	return cmd.Run() == nil
}

// Attach session if outside of tmux
func Attach(sessionName string) tea.Cmd {
	if Inside() {
		return nestedSessionsNotAllowed
	}

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
func NewSession(sessionName, workingDirectory string) tea.Cmd {
	// use Inside() instead of detached to determine how it is ran
	sessionFlags := "-s"
	if Inside() {
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
	if err == nil {
		return TmuxMsg {
			Err: nil,
		}
	}

	var tmuxErr error
	if strings.Contains(err.Error(), "can't find session") {
		tmuxErr = ErrSessionNotFound
	} 
	if strings.Contains(err.Error(), "duplicate session") {
		tmuxErr = ErrDuplicateSession
	}

	return TmuxMsg {
		Err: tmuxErr,
	}
}

func nestedSessionsNotAllowed() tea.Msg {
	return TmuxMsg {
		Err: ErrNestedSession,
	}
}
