package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	// Selected worktree
	Selected Worktree

	// App state
	AppState *AppState

	// List of worktrees
	WorktreeList list.Model
}

// Init implements tea.Model
func (Model) Init() tea.Cmd { return nil }

// View implements tea.Model
func (m Model) View() string {
	return m.WorktreeList.View()
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.WorktreeList.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case CtrlC:
			return m, tea.Quit
		case CtrlS:
			m.AppState.SwitchToWorktreeRoot = !m.AppState.SwitchToWorktreeRoot
		case CtrlR:
			m.AppState.ShowRelativePath = !m.AppState.ShowRelativePath
		case Enter:
			worktree, ok := m.WorktreeList.SelectedItem().(Worktree)
			if ok {
				m.Selected = worktree
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.WorktreeList, cmd = m.WorktreeList.Update(msg)
	return m, cmd
}
