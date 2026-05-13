package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildCLI(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "git-switch")
	cmd := exec.Command("go", "build", "-o", bin, "./cli")
	cmd.Dir = findProjectRoot(t)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build cli: %v\n%s", err, out)
	}
	return bin
}

func findProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root with go.work")
		}
		dir = parent
	}
}

func TestCLIInitAndProfileList(t *testing.T) {
	bin := buildCLI(t)
	home := t.TempDir()
	t.Setenv("HOME", home)

	// Run init
	out, err := exec.Command(bin, "init").CombinedOutput()
	if err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Initialized") {
		t.Fatalf("expected 'Initialized' in output, got: %s", out)
	}

	// Verify config was created
	configPath := filepath.Join(home, ".config", "git-switch", "config.yaml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config.yaml not created: %v", err)
	}

	// Run profile list
	out, err = exec.Command(bin, "profile", "list").CombinedOutput()
	if err != nil {
		t.Fatalf("profile list: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "personal") {
		t.Fatalf("expected 'personal' in output, got: %s", out)
	}
}

func TestCLIStatus(t *testing.T) {
	bin := buildCLI(t)
	home := t.TempDir()
	t.Setenv("HOME", home)

	// Init first
	out, err := exec.Command(bin, "init").CombinedOutput()
	if err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}

	// Run status
	out, err = exec.Command(bin, "status").CombinedOutput()
	if err != nil {
		t.Fatalf("status: %v\n%s", err, out)
	}
	output := string(out)
	if !strings.Contains(output, "Active profile") {
		t.Fatalf("expected 'Active profile' in output, got: %s", output)
	}
}

func TestCLIHelp(t *testing.T) {
	bin := buildCLI(t)
	out, err := exec.Command(bin, "--help").CombinedOutput()
	if err != nil {
		t.Fatalf("help: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "git-switch") {
		t.Fatalf("expected 'git-switch' in output, got: %s", out)
	}
}
