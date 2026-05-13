package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func psProfilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	docDir := filepath.Join(home, "Documents")
	// Try multiple possible PowerShell profile paths
	paths := []string{
		filepath.Join(docDir, "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1"),
		filepath.Join(docDir, "PowerShell", "Microsoft.PowerShell_profile.ps1"),
		filepath.Join(home, ".config", "powershell", "Microsoft.PowerShell_profile.ps1"),
	}
	for _, p := range paths {
		if _, err := os.Stat(filepath.Dir(p)); err == nil {
			return p, nil
		}
	}
	return paths[0], nil
}

const psHookMarker = "# git-switch hook"

func PowerShellHookScript() string {
	return psHookMarker + "\n" +
		"function prompt {\n" +
		"    $realLASTEXITCODE = $LASTEXITCODE\n" +
		"    git-switch status --quiet 2>$null\n" +
		"    $LASTEXITCODE = $realLASTEXITCODE\n" +
		"    \"PS $($executionContext.SessionState.Path.CurrentLocation)$('>' * ($nestedPromptLevel + 1)) \"\n" +
		"}\n"
}

func InstallPowerShellHook() error {
	path, err := psProfilePath()
	if err != nil {
		return err
	}
	if installed, err := IsPowerShellHookInstalled(); err != nil {
		return err
	} else if installed {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	if _, err := f.WriteString("\n" + PowerShellHookScript()); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func UninstallPowerShellHook() error {
	path, err := psProfilePath()
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
	start := strings.Index(content, psHookMarker)
	if start == -1 {
		return nil
	}
	// Find end of hook block (closing brace of prompt function)
	end := strings.Index(content[start:], "}\n")
	if end == -1 {
		return nil
	}
	end = start + end + 2
	newContent := strings.TrimRight(content[:start], "\n") + content[end:]
	if err := os.WriteFile(path, []byte(newContent), 0o600); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func IsPowerShellHookInstalled() (bool, error) {
	path, err := psProfilePath()
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
	return strings.Contains(string(data), psHookMarker), nil
}
