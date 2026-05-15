package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func completionDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "git-switch", "completions"), nil
}

func completionFilePath(shell string) (string, error) {
	dir, err := completionDir()
	if err != nil {
		return "", err
	}
	switch shell {
	case "bash":
		return filepath.Join(dir, "git-switch.bash"), nil
	case "zsh":
		return filepath.Join(dir, "_git-switch"), nil
	case "powershell", "pwsh":
		return filepath.Join(dir, "git-switch.ps1"), nil
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}

const (
	completionBlockBegin = "# git-switch completion BEGIN"
	completionBlockEnd   = "# git-switch completion END"
)

func completionSourceLine(shell string) string {
	path, _ := completionFilePath(shell)
	switch shell {
	case "bash":
		return fmt.Sprintf("\n%s\nsource %s\n%s\n", completionBlockBegin, path, completionBlockEnd)
	case "zsh":
		return fmt.Sprintf("\n%[1]s\nfpath=(%[2]s $fpath)\nautoload -Uz compinit && compinit -i\n%[3]s\n",
			completionBlockBegin, filepath.Dir(path), completionBlockEnd)
	case "powershell", "pwsh":
		return fmt.Sprintf("\n%s\n. %s\n%s\n", completionBlockBegin, path, completionBlockEnd)
	default:
		return ""
	}
}

func WriteCompletionScript(shell, script string) error {
	path, err := completionFilePath(shell)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create completion dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(script), 0o644); err != nil {
		return fmt.Errorf("write completion script: %w", err)
	}
	return nil
}

func RemoveCompletionScript(shell string) error {
	path, err := completionFilePath(shell)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove completion script: %w", err)
	}
	dir := filepath.Dir(path)
	entries, _ := os.ReadDir(dir)
	if len(entries) == 0 {
		_ = os.Remove(dir)
	}
	return nil
}

func InjectCompletionBlock(rcContent, shell string) string {
	srcLine := completionSourceLine(shell)
	if strings.Contains(rcContent, completionBlockBegin) {
		return rcContent
	}
	return strings.TrimRight(rcContent, "\n") + srcLine
}

func RemoveCompletionBlock(rcContent, shell string) string {
	begin := strings.Index(rcContent, completionBlockBegin)
	if begin == -1 {
		return rcContent
	}
	end := strings.Index(rcContent[begin:], completionBlockEnd)
	if end == -1 {
		return rcContent
	}
	end += begin + len(completionBlockEnd)
	return strings.TrimRight(rcContent[:begin], "\n") + rcContent[end:]
}

func IsCompletionInstalled(rcContent string) bool {
	return strings.Contains(rcContent, completionBlockBegin)
}
