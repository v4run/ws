package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultWidth  = 20
	defaultHeight = 16
)

var (
	titleStyle           = lipgloss.NewStyle().MarginLeft(2)
	itemStyle            = lipgloss.NewStyle().PaddingLeft(4)
	highlightedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	helpStyle            = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	dirStyle             = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).PaddingLeft(2)
)

// Keys
const (
	CtrlR = "ctrl+r"
	CtrlS = "ctrl+s"
	CtrlC = "ctrl+c"
	Enter = "enter"
)

func InitList(items []list.Item, delegate ItemDelegate) list.Model {
	l := list.New(items, delegate, defaultWidth, defaultHeight)
	l.Title = "Select worktree"
	l.SetShowFilter(true)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.HelpStyle = helpStyle
	l.AdditionalShortHelpKeys = additionalBindHelps
	l.AdditionalFullHelpKeys = additionalBindHelps
	return l
}

func additionalBindHelps() []key.Binding {
	bindings := []key.Binding{
		key.NewBinding(
			key.WithKeys(CtrlR),
			key.WithHelp(CtrlR, "toggle relative path"),
		),
		key.NewBinding(
			key.WithKeys(CtrlS),
			key.WithHelp(CtrlS, "toggle subtree"),
		),
	}
	return bindings
}
