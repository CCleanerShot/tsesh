package main

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		h, err := os.UserHomeDir()
		if err != nil {
			return h, err
		}
		return filepath.Join(h, path[1:]), err
	}
	return path, nil
}

type model struct {
	loaded bool
	searchPaths []string
	pathOptions []string
}

func newModel() model {
	return model{
		loaded: false,
		searchPaths: []string{
			"~/projects",
			"~/code",
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.loaded {
		for _, path := range m.searchPaths {
			p, _ := expandPath(path)
			dirEntries, _:= os.ReadDir(p)

			for _, entry := range dirEntries {
				m.pathOptions = append(m.pathOptions, entry.Name())
			}
		}
		m.loaded = true
	}

	switch msg := msg.(type)  {
		case tea.KeyMsg :
			switch msg.String() {
			case "q":
				return m, tea.Quit
			}
	}
	return m, nil
}

func (m model) View() string {
	list := strings.Builder{}
	for _, option := range m.pathOptions {
		list.WriteString(option + "\n")
	}
	return list.String()
}
