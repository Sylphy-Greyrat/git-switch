package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func rcFile(shell string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch shell {
	case "bash":
		return filepath.Join(home, ".bashrc"), nil
	case "zsh":
		return filepath.Join(home, ".zshrc"), nil
	default:
		return "", fmt.Errorf("unsupported shell: %s (supported: bash, zsh)", shell)
	}
}

const hookMarker = "# git-switch hook"

func ShellHookScript() string {
	return hookMarker + "\n" +
		"git_switch_cd() {\n" +
		"    \\cd \"$@\" || return\n" +
		"    if git rev-parse --git-dir >/dev/null 2>&1; then\n" +
		"        git-switch status --quiet\n" +
		"    fi\n" +
		"}\n" +
		"alias cd=git_switch_cd\n"
}

func InstallShellHook(shell string) error {
	path, err := rcFile(shell)
	if err != nil {
		return err
	}
	if installed, err := IsShellHookInstalled(shell); err != nil {
		return err
	} else if installed {
		return nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	if _, err := f.WriteString("\n" + ShellHookScript()); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func UninstallShellHook(shell string) error {
	path, err := rcFile(shell)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}
	content := string(data)
	start := strings.Index(content, hookMarker)
	if start == -1 {
		return nil
	}
	end := strings.Index(content[start:], "alias cd=git_switch_cd\n")
	if end == -1 {
		return nil
	}
	end = start + end + len("alias cd=git_switch_cd\n")
	newContent := strings.TrimRight(content[:start], "\n") + content[end:]
	if err := os.WriteFile(path, []byte(newContent), 0o600); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func IsShellHookInstalled(shell string) (bool, error) {
	path, err := rcFile(shell)
	if err != nil {
		return false, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("read %s: %w", path, err)
	}
	return strings.Contains(string(data), hookMarker), nil
}
