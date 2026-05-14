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

func TestGitApplierAddsGPGConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitconfig")
	profile := config.Profile{
		User: config.UserConfig{Name: "GPG User", Email: "gpg@example.com"},
		GPG:  &config.GPGConfig{SigningKey: "ABC123", SignCommits: true},
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

	if !strings.Contains(content, "signingkey = ABC123") {
		t.Fatal("GPG signingkey was not applied")
	}
	if !strings.Contains(content, "gpgsign = true") {
		t.Fatal("commit.gpgsign was not applied")
	}
	if !strings.Contains(content, "[commit]") {
		t.Fatal("expected [commit] section")
	}
}

func TestGitApplierRevert(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitconfig")

	existing := "[user]\n\tname = Original\n\temail = original@example.com\n"
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatalf("write existing gitconfig: %v", err)
	}

	profile := config.Profile{
		User: config.UserConfig{Name: "Changed", Email: "changed@example.com"},
	}

	applier := NewGitApplier(path)
	if err := applier.ApplyGitConfig(context.Background(), profile); err != nil {
		t.Fatalf("apply git config: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "Changed") {
		t.Fatal("config was not changed")
	}

	if err := applier.Revert(); err != nil {
		t.Fatalf("revert: %v", err)
	}

	data, _ = os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "Original") {
		t.Fatalf("revert did not restore original:\n%s", content)
	}
	if strings.Contains(content, "Changed") {
		t.Fatal("revert did not remove changes")
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

func TestSSHApplierMultiHost(t *testing.T) {
	sshDir := t.TempDir()
	profile := config.Profile{
		Profile: config.ProfileMeta{Name: "multi"},
		SSH: &config.SSHConfig{
			KeyFile:   "~/.ssh/id_rsa",
			Hosts:     []string{"github.com", "gitlab.com"},
			HostAlias: "git-switch",
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
	content := string(generated)

	if !strings.Contains(content, "HostName github.com") {
		t.Fatal("missing github.com host entry")
	}
	if !strings.Contains(content, "HostName gitlab.com") {
		t.Fatal("missing gitlab.com host entry")
	}
}

func TestSSHApplierRevert(t *testing.T) {
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

	// Verify file was created
	if _, err := os.Stat(filepath.Join(sshDir, "config.d", "git-switch")); err != nil {
		t.Fatal("generated file was not created")
	}

	if err := applier.Revert(); err != nil {
		t.Fatalf("revert: %v", err)
	}

	if _, err := os.Stat(filepath.Join(sshDir, "config.d", "git-switch")); !os.IsNotExist(err) {
		t.Fatal("generated file was not removed on revert")
	}
}

func TestGitApplierRejectsSecondApplyBeforeRevert(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitconfig")
	profile := config.Profile{User: config.UserConfig{Name: "First", Email: "first@example.com"}}

	applier := NewGitApplier(path)
	if err := applier.ApplyGitConfig(context.Background(), profile); err != nil {
		t.Fatalf("first apply git config: %v", err)
	}
	if err := applier.ApplyGitConfig(context.Background(), profile); err == nil {
		t.Fatal("expected second apply before revert to fail")
	}
}

func TestGitApplierDoesNotBlockAfterFailedApply(t *testing.T) {
	dir := t.TempDir()
	profile := config.Profile{User: config.UserConfig{Name: "First", Email: "first@example.com"}}

	applier := NewGitApplier(dir)
	if err := applier.ApplyGitConfig(context.Background(), profile); err == nil {
		t.Fatal("expected apply to fail when config path is a directory")
	}

	applier.configPath = filepath.Join(dir, ".gitconfig")
	if err := applier.ApplyGitConfig(context.Background(), profile); err != nil {
		t.Fatalf("expected apply after failed apply to succeed: %v", err)
	}
}

func TestSSHApplierPreservesExistingConfigPermissions(t *testing.T) {
	sshDir := t.TempDir()
	configPath := filepath.Join(sshDir, "config")
	if err := os.WriteFile(configPath, []byte("Host github.com\n"), 0o644); err != nil {
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
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("stat ssh config: %v", err)
	}
	if info.Mode().Perm() != 0o644 {
		t.Fatalf("expected permissions 0644, got %o", info.Mode().Perm())
	}
}

func TestSSHApplierExpandsIdentityFileHome(t *testing.T) {
	sshDir := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)
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
	if strings.Contains(string(generated), "IdentityFile ~/") {
		t.Fatalf("IdentityFile home was not expanded:\n%s", generated)
	}
	if !strings.Contains(string(generated), filepath.Join(home, ".ssh", "id_rsa")) {
		t.Fatalf("IdentityFile missing expanded home:\n%s", generated)
	}
}

func TestSSHApplierRejectsInvalidHostAlias(t *testing.T) {
	sshDir := t.TempDir()
	profile := config.Profile{
		Profile: config.ProfileMeta{Name: "personal"},
		SSH: &config.SSHConfig{
			KeyFile:   "~/.ssh/id_rsa",
			Hosts:     []string{"github.com"},
			HostAlias: "bad alias",
		},
	}

	applier := NewSSHApplier(sshDir)
	if err := applier.ApplySSHConfig(context.Background(), profile); err == nil {
		t.Fatal("expected invalid host alias to fail")
	}
}

func TestSSHApplierRejectsInvalidHost(t *testing.T) {
	sshDir := t.TempDir()
	profile := config.Profile{
		Profile: config.ProfileMeta{Name: "personal"},
		SSH: &config.SSHConfig{
			KeyFile: "~/.ssh/id_rsa",
			Hosts:   []string{"bad host"},
		},
	}

	applier := NewSSHApplier(sshDir)
	if err := applier.ApplySSHConfig(context.Background(), profile); err == nil {
		t.Fatal("expected invalid host to fail")
	}
}

func TestSSHApplierRejectsInvalidKeyFile(t *testing.T) {
	sshDir := t.TempDir()
	profile := config.Profile{
		Profile: config.ProfileMeta{Name: "personal"},
		SSH: &config.SSHConfig{
			KeyFile: "~/.ssh/id_rsa\n    ProxyCommand evil",
			Hosts:   []string{"github.com"},
		},
	}

	applier := NewSSHApplier(sshDir)
	if err := applier.ApplySSHConfig(context.Background(), profile); err == nil {
		t.Fatal("expected invalid key file to fail")
	}
}

func TestSSHApplierReturnsExpandHomeErrorForKeyFile(t *testing.T) {
	sshDir := t.TempDir()
	t.Setenv("HOME", "")
	profile := config.Profile{
		Profile: config.ProfileMeta{Name: "personal"},
		SSH: &config.SSHConfig{
			KeyFile: "~/.ssh/id_rsa",
			Hosts:   []string{"github.com"},
		},
	}

	applier := NewSSHApplier(sshDir)
	if err := applier.ApplySSHConfig(context.Background(), profile); err == nil {
		t.Fatal("expected key file home expansion error")
	}
}

func TestSSHApplierSkipsCommentedInclude(t *testing.T) {
	sshDir := t.TempDir()

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
