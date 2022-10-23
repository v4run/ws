package app

import (
	"bytes"
	"testing"
)

var parseCommandOutputTestCases = []struct {
	In       bytes.Buffer
	Expected []Worktree
}{
	{
		In: *bytes.NewBufferString(
			`worktree /path/to/bare-source
bare

worktree /path/to/linked-worktree
HEAD abcd1234abcd1234abcd1234abcd1234abcd1234
branch refs/heads/master

worktree /path/to/other-linked-worktree
HEAD 1234abc1234abc1234abc1234abc1234abc1234a
detached

worktree /path/to/linked-worktree-locked-no-reason
HEAD 5678abc5678abc5678abc5678abc5678abc5678c
branch refs/heads/locked-no-reason
locked

worktree /path/to/linked-worktree-locked-with-reason
HEAD 3456def3456def3456def3456def3456def3456b
branch refs/heads/locked-with-reason
locked reason why is locked

worktree /path/to/linked-worktree-prunable
HEAD 1233def1234def1234def1234def1234def1234b
detached
prunable gitdir file points to non-existent location

`,
		),
		Expected: []Worktree{
			{
				RootPath: "/path/to/bare-source",
				Bare:     true,
			},
			{
				RootPath: "/path/to/linked-worktree",
				Head:     "abcd1234abcd1234abcd1234abcd1234abcd1234",
				Branch:   "master",
			},
			{
				RootPath: "/path/to/other-linked-worktree",
				Head:     "1234abc1234abc1234abc1234abc1234abc1234a",
				Detached: true,
			},
			{
				RootPath: "/path/to/linked-worktree-locked-no-reason",
				Head:     "5678abc5678abc5678abc5678abc5678abc5678c",
				Branch:   "locked-no-reason",
				Locked:   true,
			},
			{
				RootPath: "/path/to/linked-worktree-locked-with-reason",
				Head:     "3456def3456def3456def3456def3456def3456b",
				Branch:   "locked-with-reason",
				Locked:   true,
				Reason:   "reason why is locked",
			},
			{
				RootPath: "/path/to/linked-worktree-prunable",
				Head:     "1233def1234def1234def1234def1234def1234b",
				Detached: true,
				Prunable: true,
				Reason:   "gitdir file points to non-existent location",
			},
		},
	},
}

func TestParseCommandOutput(t *testing.T) {
	for i, v := range parseCommandOutputTestCases {
		a := App{}
		worktrees := a.ParseCommandOutput(v.In)
		if len(worktrees) != len(v.Expected) {
			t.Errorf("Case: %d. Wrong number of workspaces found. Expected: %d, Got: %d", i+1, len(v.Expected), len(worktrees))
		}
		for j, w := range worktrees {
			if w != v.Expected[j] {
				t.Errorf("Case %d:%d. Worktree mismatch. Expected: %+v, Got: %+v", i+1, j+1, v.Expected[j], w)
			}
		}
	}
}
