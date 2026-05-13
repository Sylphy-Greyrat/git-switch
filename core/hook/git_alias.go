package hook

import (
	"fmt"
	"os/exec"
)

func InstallGitAlias() error {
	out, err := exec.Command("git", "config", "--global", "alias.sw", "!git-switch").CombinedOutput()
	if err != nil {
		return fmt.Errorf("git config: %s: %w", string(out), err)
	}
	return nil
}

func UninstallGitAlias() error {
	out, err := exec.Command("git", "config", "--global", "--unset", "alias.sw").CombinedOutput()
	if err != nil {
		return fmt.Errorf("git config: %s: %w", string(out), err)
	}
	return nil
}
