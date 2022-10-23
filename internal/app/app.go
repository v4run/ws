package app

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/pkg/errors"
)

type AppState struct {
	ShowRelativePath     bool
	SwitchToWorktreeRoot bool
}
type App struct {
	activeWorktreeRootPath string
	cwd                    string
	// activeWorktreeRootPath + subtree = cwd
	// eg:
	//	activeWorktreeRootPath = /abs/path/to/project
	//	cwd = /abs/path/to/project/some/dir
	//	subtree = /some/dir
	subtree string
	state   *AppState
}

func InitApp(state AppState) App {
	app := App{}
	app.state = &state
	app.cwd = GetCWD()
	app.subtree = GetSubtreePath()
	app.activeWorktreeRootPath = strings.TrimSuffix(app.cwd, app.subtree)
	return app
}

// GetWorktrees returns the worktree list for the repo
// of which cwd is part of
func (a App) GetWorktrees() ([]Worktree, error) {
	var outputBuf, errorBuf bytes.Buffer
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Stdout = &outputBuf
	cmd.Stderr = &errorBuf
	cmd.Dir = a.cwd
	if err := cmd.Run(); err != nil {
		return nil, errors.Wrap(err, errorBuf.String())
	}
	worktrees := a.ParseCommandOutput(outputBuf)
	return worktrees, nil
}

// ParseCommandOutput parses the "git worktree list --porcelain"
// command output and returns a list of Worktree
func (a App) ParseCommandOutput(buf bytes.Buffer) []Worktree {
	scanner := bufio.NewScanner(&buf)
	var (
		workspaceBlockIdx int
		workTreeIdx       int
		worktrees         []Worktree
		currentWorktree   Worktree
	)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			workTreeIdx++
			worktrees = append(worktrees, currentWorktree)
			workspaceBlockIdx = 0
			currentWorktree = Worktree{}
		}
		label, value, _ := strings.Cut(text, " ")
		switch label {
		case "worktree":
			currentWorktree.RootPath = value
		case "HEAD":
			currentWorktree.Head = value
		case "bare":
			currentWorktree.Bare = true
		case "branch":
			currentWorktree.Branch = strings.TrimPrefix(value, "refs/heads/")
		case "locked":
			currentWorktree.Locked = true
			currentWorktree.Reason = value
		case "detached":
			currentWorktree.Detached = true
		case "prunable":
			currentWorktree.Prunable = true
			currentWorktree.Reason = value
		}
		workspaceBlockIdx++
	}
	return worktrees
}

func (a App) NewModel() (Model, error) {
	workspaces, err := a.GetWorktrees()
	if err != nil {
		return Model{}, err
	}
	var (
		items             []list.Item
		activeWorktreeIdx int
	)
	var i int
	for _, w := range workspaces {
		// enrich the worktree with more details
		w.Enrich(a)
		if w.Prunable {
			continue
		}
		if w.RootPath == a.activeWorktreeRootPath {
			activeWorktreeIdx = i
		}
		items = append(items, w)
		i++
	}
	l := InitList(items, ItemDelegate{})
	l.Select(activeWorktreeIdx)
	m := Model{
		WorktreeList: l,
		AppState:     a.state,
	}
	return m, nil
}

func (a App) Run() error {
	model, err := a.NewModel()
	if err != nil {
		return err
	}
	// force colors
	lipgloss.SetColorProfile(termenv.TrueColor)
	options := []tea.ProgramOption{
		// hack to allow running in subshell
		tea.WithOutput(os.Stderr),
		tea.WithAltScreen(),
	}
	pgm := tea.NewProgram(model, options...)
	finalModel, err := pgm.StartReturningModel()
	if err != nil {
		return err
	}

	appModel, _ := finalModel.(Model)
	if appModel.Selected.RootPath == "" {
		return ErrCancelled
	}
	if a.state.SwitchToWorktreeRoot {
		fmt.Println(appModel.Selected.RootPath)
	} else {
		fmt.Println(appModel.Selected.ValidSubtreePath)
	}
	return nil
}
