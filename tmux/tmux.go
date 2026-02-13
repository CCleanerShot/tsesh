package tmux

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type TmuxMsg struct{ Err error }

var cmdRunner = exec.Command
var ErrNestedSession = errors.New("sessions should be nested with care, unset $TMUX to force")
var ErrSessionNotFound = errors.New("session was not found")
var ErrDuplicateSession = errors.New("duplicateSession")
var ErrNoClientFound = errors.New("no client found")

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
func NewSession(sessionName, workingDirectory string) tea.Cmd {
	sessionFlags := "-s"
	if Inside() {
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

func tmux(args ...string) *exec.Cmd {
	cmd := cmdRunner("tmux", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func execCallback(err error) tea.Msg {
	if err == nil {
		return TmuxMsg{
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
	if strings.Contains(err.Error(), "no client found") {
		tmuxErr = ErrNoClientFound
	}

	return TmuxMsg{
		Err: tmuxErr,
	}
}

func nestedSessionsNotAllowed() tea.Msg {
	return TmuxMsg{
		Err: ErrNestedSession,
	}
}
