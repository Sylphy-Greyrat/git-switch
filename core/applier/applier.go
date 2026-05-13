package applier

import (
	"context"
	"errors"

	"github.com/sylphy/git-switch/core/config"
)

type Applier interface {
	ApplyGitConfig(ctx context.Context, profile config.Profile) error
	ApplySSHConfig(ctx context.Context, profile config.Profile) error
	Revert(ctx context.Context) error
}

type ProfileApplier struct {
	git GitApplier
	ssh SSHApplier
}

func NewProfileApplier(gitConfigPath, sshDir string) *ProfileApplier {
	return &ProfileApplier{
		git: NewGitApplier(gitConfigPath),
		ssh: NewSSHApplier(sshDir),
	}
}

func (a *ProfileApplier) ApplyGitConfig(ctx context.Context, profile config.Profile) error {
	return a.git.ApplyGitConfig(ctx, profile)
}

func (a *ProfileApplier) ApplySSHConfig(ctx context.Context, profile config.Profile) error {
	return a.ssh.ApplySSHConfig(ctx, profile)
}

func (a *ProfileApplier) Revert(ctx context.Context) error {
	return errors.New("revert not yet implemented")
}
