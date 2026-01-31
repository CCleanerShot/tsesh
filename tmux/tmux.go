package tmux

import (
	"os"
	"os/exec"
)

// TODO: ideas
// check out session groups and see if I can implement a shortcut to add a new session to a group

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
func Attach(sessionName string) error {
	cmd := tmux("attach-session", "-t", sessionName)
	return cmd.Run()
}

// Move active session to existing session if inside tmux 
func SwitchClient(sessionName string) error {
	cmd:= tmux("switch-client", "-t", sessionName)
	return cmd.Run()
}

// New session will be created with specified name at provided workingDirectory
func NewSession(sessionName, workingDirectory string, detached bool) error {
	conditionalDetached := ""
	if detached {
		conditionalDetached = "-d"
	}

	cmd := tmux(
		"new-session",
		"-s", sessionName, 
		"-c", workingDirectory,
		conditionalDetached,
	)

	return cmd.Run()
}

func tmux(args... string) *exec.Cmd {
	cmd := exec.Command("tmux", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
