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
		return "", fmt.Errorf("unsupported shell: %s (supported: bash, zsh, powershell, pwsh)", shell)
	}
}

const (
	completionBlockBegin = "\n# git-switch completion BEGIN"
	completionBlockEnd   = "# git-switch completion END"
)

func completionSourceLine(shell string) (string, error) {
	path, err := completionFilePath(shell)
	if err != nil {
		return "", err
	}
	switch shell {
	case "bash":
		return fmt.Sprintf("%s\nsource %s\n%s\n", completionBlockBegin, path, completionBlockEnd), nil
	case "zsh":
		return fmt.Sprintf("%s\nfpath=(%s $fpath)\n%s\n",
			completionBlockBegin, filepath.Dir(path), completionBlockEnd), nil
	case "powershell", "pwsh":
		return fmt.Sprintf("%s\n. %s\n%s\n", completionBlockBegin, path, completionBlockEnd), nil
	default:
		return "", fmt.Errorf("unsupported shell: %s (supported: bash, zsh, powershell, pwsh)", shell)
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
	if entries, err := os.ReadDir(dir); err == nil && len(entries) == 0 {
		_ = os.Remove(dir)
	}
	return nil
}

func InjectCompletionBlock(rcContent, shell string) (string, error) {
	srcLine, err := completionSourceLine(shell)
	if err != nil {
		return "", err
	}
	if strings.Contains(rcContent, completionBlockBegin) {
		return rcContent, nil
	}
	return strings.TrimRight(rcContent, "\n") + srcLine, nil
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
