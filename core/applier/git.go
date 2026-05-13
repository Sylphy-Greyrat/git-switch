package applier

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sylphy/git-switch/core/config"
)

type GitApplier struct {
	configPath string
}

func NewGitApplier(configPath string) GitApplier {
	return GitApplier{configPath: configPath}
}

func (a GitApplier) ApplyGitConfig(ctx context.Context, profile config.Profile) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	path, err := config.ExpandHome(a.configPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create git config directory: %w", err)
	}

	content := fmt.Sprintf("[user]\n\tname = %s\n\temail = %s\n", profile.User.Name, profile.User.Email)
	if profile.SSH != nil && profile.SSH.HostAlias != "" && len(profile.SSH.Hosts) > 0 {
		content += fmt.Sprintf("\n[url \"git@%s:\"]\n\tinsteadOf = git@%s:\n", profile.SSH.HostAlias, profile.SSH.Hosts[0])
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write git config %s: %w", path, err)
	}
	return nil
}
