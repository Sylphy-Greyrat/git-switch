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
