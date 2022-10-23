package app

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ItemDelegate struct{}

// Height implements list.ItemDelegate
func (ItemDelegate) Height() int {
	return 1
}

// Render implements list.ItemDelegate
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	worktree, ok := item.(Worktree)
	if !ok {
		return
	}
	slNo := index + 1
	stylizeItem := itemStyle.Render
	dir := dirStyle.Render(worktree.Dir())
	if index == m.Index() {
		stylizeItem = func(str string) string {
			return highlightedItemStyle.Render("> " + str)
		}
	}
	itemLabel := stylizeItem(fmt.Sprintf("%d. %s", slNo, worktree.FilterValue()))
	label := itemLabel + dir
	fmt.Fprint(w, label)
}

// Spacing implements list.ItemDelegate
func (ItemDelegate) Spacing() int {
	return 0
}

// Update implements list.ItemDelegate
func (ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
