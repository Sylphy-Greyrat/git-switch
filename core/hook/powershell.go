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

const (
	psBlockBegin = "\n# ------ git-switch BLOCK BEGIN ------"
	psBlockEnd   = "# ------ git-switch BLOCK END ------"
	psOldMarker  = "# git-switch hook"
)

func PowerShellHookScript() string {
	return psBlockBegin + "\n" +
		"$global:GitSwitchOriginalPrompt = $function:prompt\n" +
		"function git_switch_prompt {\n" +
		"    $realLASTEXITCODE = $LASTEXITCODE\n" +
		"    git-switch status --quiet 2>$null\n" +
		"    $LASTEXITCODE = $realLASTEXITCODE\n" +
		"    if ($global:GitSwitchOriginalPrompt) { & $global:GitSwitchOriginalPrompt }\n" +
		"}\n" +
		"Set-Item -Path function:prompt -Value ${function:git_switch_prompt}\n" +
		psBlockEnd + "\n"
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
	start, end := findPSBlockRange(content)
	if start == -1 {
		return nil
	}
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
	s := string(data)
	return strings.Contains(s, psBlockBegin) || strings.Contains(s, psOldMarker), nil
}

func findPSBlockRange(content string) (int, int) {
	start := strings.Index(content, psBlockBegin)
	if start != -1 {
		end := strings.Index(content[start:], psBlockEnd)
		if end != -1 {
			return start, start + end + len(psBlockEnd)
		}
	}
	start = strings.Index(content, psOldMarker)
	if start != -1 {
		end := strings.Index(content[start:], "Set-Item -Path function:prompt -Value ${function:git_switch_prompt}\n")
		if end != -1 {
			return start, start + end + len("Set-Item -Path function:prompt -Value ${function:git_switch_prompt}\n")
		}
	}
	return -1, -1
}
