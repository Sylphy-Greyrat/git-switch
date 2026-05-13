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

func TestGitApplierPreservesExistingConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitconfig")

	existing := "[core]\n\teditor = vim\n\tautocrlf = input\n"
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatalf("write existing gitconfig: %v", err)
	}

	profile := config.Profile{
		User: config.UserConfig{Name: "New User", Email: "new@example.com"},
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

	if !strings.Contains(content, "[core]") {
		t.Fatal("existing [core] section was removed")
	}
	if !strings.Contains(content, "editor = vim") {
		t.Fatal("existing core.editor was removed")
	}
	if !strings.Contains(content, "name = New User") {
		t.Fatal("user.name was not added")
	}
	if !strings.Contains(content, "email = new@example.com") {
		t.Fatal("user.email was not added")
	}
}

func TestGitApplierUpdatesExistingUserKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitconfig")

	existing := "[user]\n\tname = Old Name\n\temail = old@example.com\n"
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatalf("write existing gitconfig: %v", err)
	}

	profile := config.Profile{
		User: config.UserConfig{Name: "Updated Name", Email: "updated@example.com"},
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

	if strings.Contains(content, "Old Name") {
		t.Fatal("old user.name was not replaced")
	}
	if strings.Contains(content, "old@example.com") {
		t.Fatal("old user.email was not replaced")
	}
	if !strings.Contains(content, "name = Updated Name") {
		t.Fatal("new user.name not found")
	}
	nameCount := strings.Count(content, "name = ")
	if nameCount != 1 {
		t.Fatalf("expected 1 name entry, got %d\n%s", nameCount, content)
	}
}

func TestGitApplierAddsURLInsteadOfForSSHProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitconfig")
	profile := config.Profile{
		User: config.UserConfig{Name: "SSH User", Email: "ssh@example.com"},
		SSH: &config.SSHConfig{
			KeyFile:   "~/.ssh/id_rsa",
			Hosts:     []string{"github.com"},
			HostAlias: "github.com-personal",
		},
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

	if !strings.Contains(content, `[url "git@github.com-personal:"]`) {
		t.Fatalf("expected URL insteadOf section:\n%s", content)
	}
	if !strings.Contains(content, "insteadOf = git@github.com:") {
		t.Fatalf("expected insteadOf:\n%s", content)
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

func TestSSHApplierSkipsCommentedInclude(t *testing.T) {
	sshDir := t.TempDir()

	// Write SSH config with a commented-out Include directive
	commentedConfig := "# Include config.d/*\nHost github.com\n    IdentityFile ~/.ssh/id_rsa\n"
	if err := os.WriteFile(filepath.Join(sshDir, "config"), []byte(commentedConfig), 0o600); err != nil {
		t.Fatalf("write ssh config: %v", err)
	}

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

	mainConfig, err := os.ReadFile(filepath.Join(sshDir, "config"))
	if err != nil {
		t.Fatalf("read ssh config: %v", err)
	}
	content := string(mainConfig)
	if !strings.Contains(content, "Include config.d/*") {
		t.Fatal("Include directive was not added")
	}
	// Should have exactly one uncommented Include line
	uncommentedCount := 0
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) == "Include config.d/*" {
			uncommentedCount++
		}
	}
	if uncommentedCount != 1 {
		t.Fatalf("expected 1 uncommented Include line, got %d\n%s", uncommentedCount, content)
	}
}
