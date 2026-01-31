package picker

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/e-mar404/tsesh/tmux"
)

type TmuxMsg error

func TmuxError(err error) (func() tea.Msg) {
	return func() tea.Msg {
		return TmuxMsg(err)
	}
}

// Wrapper for list.Item to add extra fields
type Item struct {
	Name string
	Path string
	// TODO: at some point I should also add a configuration var for code that should be executed before and after entering the tmux session
}

func (i Item) FilterValue() string {
	return i.Name
}

func (i Item) Title() string {
	return i.Name
}

func (i Item) Description() string {
	return i.Path
}

type Picker struct {
	list list.Model
	info string
}

func (p Picker) Init() tea.Cmd {
	return nil
}

func (p Picker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return p, tea.Quit

		case "enter":
			if p.list.SelectedItem() == nil {
				return p, tea.Quit
			}

			choice := p.list.SelectedItem().(Item)
			switch tmux.HasSession(choice.Name) {
			case true:
				if tmux.Inside() {
					err := tmux.SwitchClient(choice.Name)
					cmds = append(cmds, TmuxError(err))
				}
				err := tmux.Attach(choice.Name)
				cmds = append(cmds, TmuxError(err))

			case false:
				if tmux.Inside() {
					err := tmux.NewSession(choice.Name, choice.Path, true)
					cmds = append(cmds, TmuxError(err))
					err = tmux.SwitchClient(choice.Name)
					cmds = append(cmds, TmuxError(err))
				}
				err := tmux.NewSession(choice.Name, choice.Path, false)
				cmds = append(cmds, TmuxError(err))
			}
		}
	case TmuxMsg:
		p.info = string(msg.Error())
	}

	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	cmds = append(cmds, cmd)
	
	return p, tea.Batch(cmds...)
}

func (p Picker) View() string {
	if p.info != "" {
		return p.info
	}

	return p.list.View()
}

func New() Picker {
	// TODO: this is were I should load the files from the config
	// for now they will be hardcoded
	searchPaths := []string{
		"~/",
		"~/code",
		"~/projects",
	}

	return Picker{
		list: list.New(
			findDirectories(searchPaths),
			list.NewDefaultDelegate(),
			0,
			0,
		),
	}
}
