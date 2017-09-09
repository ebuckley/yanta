package sync

import (
	"os"
	"os/exec"

	"github.com/ebuckley/yanta/context"
)

// Commit will preform a git commit on the page
func Commit(p *context.Page, msg string) (string, error) {
	addCmd := exec.Command("git", "add", p.Path)
	addCmd.Env = os.Environ()
	commitCmd := exec.Command("git", "commit", "-m", msg)
	commitCmd.Env = os.Environ()

	printCommand(addCmd)
	addOutput, err := addCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	printCommand(commitCmd)
	commitOutput, err := commitCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(addOutput) + "\r\n" + string(commitOutput), nil
}
