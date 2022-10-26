# Switch worktrees with ease

## ws - Prints a menu to choose a git worktree

[![asciicast](https://asciinema.org/a/531887.svg)](https://asciinema.org/a/531887)

## Usage

### Installation

```bash
go install github.com/v4run/ws@latest
```

### Add the following alias to go to the selected worktree directory

```bash
alias wt='cd $(ws || echo .)'
```
