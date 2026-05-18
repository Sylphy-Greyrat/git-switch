package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHookInstallAutoDetectsCurrentShell(t *testing.T) {
	home := configureIsolatedHookInstall(t)
	t.Setenv("SHELL", "/bin/zsh")
	t.Setenv("ComSpec", "")

	output, err := executeHookInstall()
	if err != nil {
		t.Fatalf("hook install error = %v", err)
	}

	assertContains(t, output, "Installed git alias: git sw -> git-switch")
	assertContains(t, output, "Installed shell hook for zsh")
	assertFileContains(t, filepath.Join(home, ".zshrc"), "git_switch_cd()")
	assertFileContains(t, filepath.Join(home, ".zshrc"), "# git-switch completion BEGIN")
	assertContains(t, output, "source")
	assertContains(t, output, "for completion to take effect")
	assertGlobalGitAlias(t)
}

func TestHookInstallShellFlagOverridesAutoDetection(t *testing.T) {
	home := configureIsolatedHookInstall(t)
	t.Setenv("SHELL", "/usr/local/bin/fish")
	t.Setenv("ComSpec", "")

	output, err := executeHookInstall("--shell", "bash")
	if err != nil {
		t.Fatalf("hook install error = %v", err)
	}

	assertContains(t, output, "Installed shell hook for bash")
	assertFileContains(t, filepath.Join(home, ".bashrc"), "git_switch_cd()")
	assertGlobalGitAlias(t)
}

func TestHookInstallDetectFailureReturnsActionableError(t *testing.T) {
	configureIsolatedHookInstall(t)
	t.Setenv("SHELL", "/usr/local/bin/fish")
	t.Setenv("ComSpec", "")

	_, err := executeHookInstall()
	if err == nil {
		t.Fatal("hook install error = nil, want shell detection error")
	}
	assertContains(t, err.Error(), "could not detect current shell")
	assertContains(t, err.Error(), "--shell pwsh")
	assertGlobalGitAlias(t)
}

func TestHookInstallHelpListsSupportedShells(t *testing.T) {
	output, err := executeHookInstall("--help")
	if err != nil {
		t.Fatalf("hook install --help error = %v", err)
	}

	assertContains(t, output, "--shell string")
	assertContains(t, output, "bash, zsh, powershell, pwsh")
}

func TestHookInstallCompletionIdempotent(t *testing.T) {
	home := configureIsolatedHookInstall(t)
	t.Setenv("SHELL", "/bin/bash")

	// First install
	_, err := executeHookInstall()
	if err != nil {
		t.Fatalf("first hook install error = %v", err)
	}
	rcPath := filepath.Join(home, ".bashrc")
	data1, _ := os.ReadFile(rcPath)
	firstCount := strings.Count(string(data1), "# git-switch completion BEGIN")

	// Second install
	_, err = executeHookInstall()
	if err != nil {
		t.Fatalf("second hook install error = %v", err)
	}
	data2, _ := os.ReadFile(rcPath)
	secondCount := strings.Count(string(data2), "# git-switch completion BEGIN")
	if secondCount > firstCount {
		t.Fatal("completion block should be idempotent")
	}
}

func configureIsolatedHookInstall(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(home, ".gitconfig"))
	return home
}

func executeHookInstall(args ...string) (string, error) {
	cmd := hookInstallCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), err
}

func executeHookStatus(args ...string) (string, error) {
	cmd := hookStatusCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), err
}

func TestHookStatusOutput(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(home, ".gitconfig"))

	output, err := executeHookStatus()
	if err != nil {
		t.Fatalf("hook status error = %v", err)
	}

	assertContains(t, output, "git alias 'sw': not installed")
	assertContains(t, output, "shell hook (bash): not installed")
	assertContains(t, output, "shell hook (zsh): not installed")
	assertContains(t, output, "powershell hook: not installed")
	assertContains(t, output, "completion (bash): not installed")
	assertContains(t, output, "completion (zsh): not installed")
	assertContains(t, output, "completion (pwsh): not installed")
}

func TestHookStatusShowsInstalledCompletion(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(home, ".gitconfig"))
	t.Setenv("SHELL", "/bin/bash")

	if _, err := executeHookInstall(); err != nil {
		t.Fatalf("hook install error = %v", err)
	}

	output, err := executeHookStatus()
	if err != nil {
		t.Fatalf("hook status error = %v", err)
	}

	assertContains(t, output, "git alias 'sw': installed")
	assertContains(t, output, "shell hook (bash): installed")
	assertContains(t, output, "completion (bash): installed")
}

func assertGlobalGitAlias(t *testing.T) {
	t.Helper()
	data, err := os.ReadFile(os.Getenv("GIT_CONFIG_GLOBAL"))
	if err != nil {
		t.Fatalf("read global git config: %v", err)
	}
	content := string(data)
	assertContains(t, content, "[alias]")
	assertContains(t, content, "sw = !git-switch")
}

func assertFileContains(t *testing.T, path string, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	assertContains(t, string(data), want)
}

func assertFileNotContains(t *testing.T, path string, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		return // file not existing is fine
	}
	if strings.Contains(string(data), want) {
		t.Fatalf("expected %s to NOT contain %q", path, want)
	}
}

func assertContains(t *testing.T, got string, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("expected %q to contain %q", got, want)
	}
}
