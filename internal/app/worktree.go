package app

import (
	"os"
	"path/filepath"
)

// Worktree defines a git worktree
type Worktree struct {
	// Absolute path of the git worktree directory
	RootPath string
	// Relative path of RootPath with respect
	// to the current working directory
	RootPathRelative string
	// Absolute path of the directory which can be reached
	// from appending the path of current directory relative
	// to the top-level directory to the RootPath. This path
	// may or may not exist.
	// Eg.
	//	CWD = $BRANCH-1/a/b/c
	//	Here,
	//		subtree is a/b/c
	//		ValidSubtreePath would be <current_worktree_root>/a/b/c
	//			if it exist
	ValidSubtreePath string
	// Relative path of ValidSubtreePath with respect
	// to the current working directory
	ValidSubtreePathRelative string
	// Worktree branch
	Branch string
	// Head
	Head string
	// Reason for lock
	Reason string
	// Whether the worktree is locked
	Locked bool
	// Whether the head is detached
	Detached bool
	// Whether the worktree can be pruned
	Prunable bool
	// Whether the repo is a bare repo
	Bare bool

	appState *AppState
}

func (w *Worktree) Enrich(app App) {
	// Set the relative root path
	relativeRoot, err := filepath.Rel(app.cwd, w.RootPath)
	if err != nil {
		relativeRoot = w.RootPath
	}
	w.RootPathRelative = relativeRoot

	// Set the absolute subtree path
	subtreePath := filepath.Join(w.RootPath, app.subtree)
	if _, err := os.Stat(subtreePath); err != nil {
		subtreePath = w.RootPath
	}
	w.ValidSubtreePath = subtreePath

	// Set the relative subtree path
	relativeSubtreePath := filepath.Join(w.RootPathRelative, app.subtree)
	if _, err := os.Stat(relativeSubtreePath); err != nil {
		relativeSubtreePath = w.RootPathRelative
	}
	w.ValidSubtreePathRelative = relativeSubtreePath
	w.appState = app.state
}

func (w Worktree) FilterValue() string {
	if w.Bare {
		return "bare"
	}
	if w.Detached {
		return w.Head[:7]
	}
	return w.Branch
}

func (w Worktree) Dir() string {
	if w.appState.SwitchToWorktreeRoot {
		if w.appState.ShowRelativePath {
			return w.RootPathRelative
		}
		return w.RootPath
	}
	if w.appState.ShowRelativePath {
		return w.ValidSubtreePathRelative
	}
	return w.ValidSubtreePath
}
