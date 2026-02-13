package picker

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/e-mar404/tsesh/tmux"
)

// Wrapper for list.Item to add extra fields
type Item struct {
	SessionName string
	Path        string
	// TODO: at some point I should also add a configuration var for code that should be executed before and after entering the tmux session
}

func (i Item) FilterValue() string {
	return i.SessionName
}

func (i Item) Title() string {
	return i.SessionName
}

func (i Item) Description() string {
	return i.Path
}

type Picker struct {
	List list.Model
	Err  error
}

func (p Picker) Init() tea.Cmd {
	return nil
}

func (p Picker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.List.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return p, tea.Quit

		case "enter":
			if p.List.SelectedItem() == nil {
				return p, tea.Quit
			}

			choice := p.List.SelectedItem().(Item)

			switch tmux.HasSession(choice.SessionName) {
			case true:
				if tmux.Inside() {
					return p, tmux.SwitchClient(choice.SessionName)
				}
				return p, tmux.Attach(choice.SessionName)

			case false:
				var cmds []tea.Cmd
				cmds = append(cmds, tmux.NewSession(choice.SessionName, choice.Path))
				if tmux.Inside() {
					cmds = append(cmds, tmux.SwitchClient(choice.SessionName))
				}
				return p, tea.Sequence(cmds...)
			}
		}

	case tmux.TmuxMsg:
		p.Err = msg.Err
		return p, tea.Quit
	}

	p.List, cmd = p.List.Update(msg)

	return p, cmd
}

func (p Picker) View() string {
	return p.List.View()
}

func New(searchPaths []string) Picker {
	return Picker{
		List: list.New(
			findDirectories(searchPaths),
			list.NewDefaultDelegate(),
			0,
			0,
		),
	}
}
