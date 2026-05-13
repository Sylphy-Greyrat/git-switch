package hook

import (
	"fmt"
	"os/exec"
	"strings"
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

func StatusGitAlias() (bool, string, error) {
	out, err := exec.Command("git", "config", "--global", "alias.sw").Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return false, "", nil
		}
		return false, "", err
	}
	val := strings.TrimSpace(string(out))
	return val == "!git-switch", val, nil
}
