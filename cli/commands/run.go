package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sylphy/git-switch/core/applier"
	"github.com/sylphy/git-switch/core/matcher"
)

func runCommand(args []string) error {
	ctx := context.Background()

	loader, err := defaultLoader()
	if err != nil {
		return err
	}

	cfg, err := loader.LoadMain(ctx)
	if err != nil {
		return err
	}

	profiles, err := loader.LoadProfiles(ctx)
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	resolver := matcher.NewResolver()
	result, err := resolver.Resolve(ctx, matcher.ResolveInput{
		CurrentDir:     dir,
		RemoteURL:      getRemoteURL(),
		Profiles:       profiles,
		DefaultProfile: cfg.General.DefaultProfile,
	})
	if err != nil {
		return err
	}

	// Find and apply the matched profile
	for _, p := range profiles {
		if p.Profile.Name == result.ProfileName {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			profileApplier := applier.NewProfileApplier(
				filepath.Join(home, ".gitconfig"),
				filepath.Join(home, ".ssh"),
			)
			if err := profileApplier.ApplyGitConfig(ctx, p); err != nil {
				return fmt.Errorf("apply git config: %w", err)
			}
			if p.SSH != nil {
				if err := profileApplier.ApplySSHConfig(ctx, p); err != nil {
					return fmt.Errorf("apply ssh config: %w", err)
				}
			}
			break
		}
	}

	// Execute git with remaining args
	if len(args) > 0 {
		gitCmd := exec.Command("git", args...)
		gitCmd.Stdin = os.Stdin
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		return gitCmd.Run()
	}
	return nil
}
