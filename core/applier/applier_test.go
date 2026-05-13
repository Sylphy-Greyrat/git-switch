package applier

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sylphy/git-switch/core/config"
)

func TestGitApplierWritesUserConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitconfig")
	profile := config.Profile{
		User: config.UserConfig{Name: "Test User", Email: "test@example.com"},
	}

	applier := NewGitApplier(path)
	if err := applier.ApplyGitConfig(context.Background(), profile); err != nil {
		t.Fatalf("apply git config: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read gitconfig: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "name = Test User") || !strings.Contains(content, "email = test@example.com") {
		t.Fatalf("unexpected gitconfig content:\n%s", content)
	}
}

func TestSSHApplierWritesConfigAndInclude(t *testing.T) {
	sshDir := t.TempDir()
	profile := config.Profile{
		Profile: config.ProfileMeta{Name: "personal"},
		SSH: &config.SSHConfig{
			KeyFile:   "~/.ssh/id_rsa",
			Hosts:     []string{"github.com"},
			HostAlias: "github.com-personal",
		},
	}

	applier := NewSSHApplier(sshDir)
	if err := applier.ApplySSHConfig(context.Background(), profile); err != nil {
		t.Fatalf("apply ssh config: %v", err)
	}

	generated, err := os.ReadFile(filepath.Join(sshDir, "config.d", "git-switch"))
	if err != nil {
		t.Fatalf("read generated ssh config: %v", err)
	}
	if !strings.Contains(string(generated), "Host github.com-personal") {
		t.Fatalf("generated config missing host alias:\n%s", generated)
	}

	mainConfig, err := os.ReadFile(filepath.Join(sshDir, "config"))
	if err != nil {
		t.Fatalf("read ssh config: %v", err)
	}
	if strings.Count(string(mainConfig), "Include config.d/*") != 1 {
		t.Fatalf("expected one include line, got:\n%s", mainConfig)
	}
}
