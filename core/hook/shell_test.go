package hook

import "testing"

func TestDetectCurrentShellFromShellEnv(t *testing.T) {
	tests := []struct {
		name  string
		shell string
		want  string
	}{
		{name: "zsh", shell: "/bin/zsh", want: "zsh"},
		{name: "bash", shell: "/usr/bin/bash", want: "bash"},
		{name: "pwsh", shell: "/opt/homebrew/bin/pwsh", want: "pwsh"},
		{name: "powershell", shell: `C:\\Program Files\\PowerShell\\7\\powershell.exe`, want: "powershell"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SHELL", tt.shell)
			t.Setenv("ComSpec", "")
			got, err := DetectCurrentShell()
			if err != nil {
				t.Fatalf("DetectCurrentShell() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("DetectCurrentShell() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDetectCurrentShellFromComSpec(t *testing.T) {
	t.Setenv("SHELL", "")
	t.Setenv("ComSpec", `C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe`)

	got, err := DetectCurrentShell()
	if err != nil {
		t.Fatalf("DetectCurrentShell() error = %v", err)
	}
	if got != "powershell" {
		t.Fatalf("DetectCurrentShell() = %q, want %q", got, "powershell")
	}
}

func TestDetectCurrentShellUnsupported(t *testing.T) {
	t.Setenv("SHELL", "/usr/local/bin/fish")
	t.Setenv("ComSpec", "")

	_, err := DetectCurrentShell()
	if err == nil {
		t.Fatal("DetectCurrentShell() error = nil, want unsupported shell error")
	}
}

func TestDetectCurrentShellEmptyEnvironment(t *testing.T) {
	t.Setenv("SHELL", "")
	t.Setenv("ComSpec", "")

	_, err := DetectCurrentShell()
	if err == nil {
		t.Fatal("DetectCurrentShell() error = nil, want missing shell error")
	}
}
