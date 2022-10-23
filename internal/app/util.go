package app

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func GetSubtreePath() string {
	var output bytes.Buffer
	cmd := exec.Command("git", "rev-parse", "--show-prefix")
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		return ""
	}
	p := strings.TrimSpace(output.String())
	if p == "" {
		return p
	}
	return "/" + strings.TrimSuffix(p, "/")
}

func GetCWD() string {
	dir, _ := os.Getwd()
	return dir
}
