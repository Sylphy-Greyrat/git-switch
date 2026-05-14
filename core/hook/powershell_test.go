package hook

import (
	"strings"
	"testing"
)

func TestPowerShellHookScriptPreservesExistingPrompt(t *testing.T) {
	script := PowerShellHookScript()
	if strings.Contains(script, "function prompt") {
		t.Fatalf("hook must not replace an existing PowerShell prompt function:\n%s", script)
	}
	if !strings.Contains(script, "git-switch status --quiet") {
		t.Fatalf("hook should run git-switch status:\n%s", script)
	}
}

func TestPowerShellHookScriptHasBlockMarkers(t *testing.T) {
	script := PowerShellHookScript()
	if !strings.Contains(script, psBlockBegin) {
		t.Fatalf("hook missing block begin marker:\n%s", script)
	}
	if !strings.Contains(script, psBlockEnd) {
		t.Fatalf("hook missing block end marker:\n%s", script)
	}
}

func TestPowerShellHookInstallUninstallWithBlockMarkers(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := InstallPowerShellHook(); err != nil {
		t.Fatalf("install: %v", err)
	}
	if err := InstallPowerShellHook(); err != nil {
		t.Fatalf("second install should be no-op: %v", err)
	}
	ok, err := IsPowerShellHookInstalled()
	if err != nil {
		t.Fatalf("is installed: %v", err)
	}
	if !ok {
		t.Fatal("expected hook to be installed")
	}
	if err := UninstallPowerShellHook(); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	ok, err = IsPowerShellHookInstalled()
	if err != nil {
		t.Fatalf("is installed after uninstall: %v", err)
	}
	if ok {
		t.Fatal("expected hook to be uninstalled")
	}
}
